package eadutil

import (
	"flag"
	"fmt"
	"github.com/lestrrat-go/libxml2/parser"
	languageLib "github.com/nyulibraries/go-ead-indexer/pkg/language"
	"github.com/nyulibraries/go-ead-indexer/pkg/util"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

var fixturesDirPath string
var goldenFilesDirPath string

var updateGoldenFiles = flag.Bool("update-golden-files", false, "update the golden files")

// We need to get the absolute path to this package in order to enable the function
// for golden file and fixture file retrieval to be called from other packages
// which would not be able to resolve the hardcoded relative paths used here.
func init() {
	// The `filename` string is the absolute path to this source file, which should
	// be located at the root of the package directory.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("ERROR: `runtime.Caller(0)` failed")
	}

	// Get the path to the parent directory of this file.  Again, this is assuming
	// that this `init()` function is defined in a package top level file -- or
	// more precisely, that this file is in the same directory at the `testdata/`
	// directory that is referenced in the relative paths used in the functions
	// defined in this file.
	herePath := filepath.Dir(filename)
	// Get testdata directory paths
	fixturesDirPath = filepath.Join(herePath, "testdata", "fixtures")
	goldenFilesDirPath = filepath.Join(herePath, "testdata", "golden")
}

func TestConvertEADToHTML(t *testing.T) {
	testConvertEADToHTML_EveryCombinationOfTagAndRenderAttributeWithInvalidChars(t)
	testConvertEADToHTML_GracefulHandlingOfInvalidXML(t)
	testConvertEADToHTML_NestedTags(t)
	testConvertEADToHTML_Specificity(t)
}

func testConvertEADToHTML_EveryCombinationOfTagAndRenderAttributeWithInvalidChars(t *testing.T) {
	type testCase struct {
		name               string
		eadString          string
		expectedHTMLString string
	}

	const invalidOSCCharacter = ``
	const eadElementText = "EAD ELEMENT TEXT"
	const eadElementTextWithInvalidOSCCharacters = invalidOSCCharacter +
		"EAD ELEMENT TEXT" + invalidOSCCharacter
	const textAfterEADTag = "AFTER EAD TAG"
	const textBeforeEADTag = "BEFORE EAD TAG"

	// See "RENDER attribute" in this DLFA-229 comment:
	// https://jira.nyu.edu/browse/DLFA-229?focusedCommentId=10283699&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10283699
	eadTagsToTest := []string{
		"emph",
		"title",
		"titleproper",
	}

	sortedRenderAttributes := []string{}
	for renderAttribute, _ := range eadTagRenderAttributeToHTMLTagName {
		sortedRenderAttributes = append(sortedRenderAttributes, renderAttribute)
	}
	slices.Sort(sortedRenderAttributes)

	testCases := []testCase{}
	for _, renderAttribute := range sortedRenderAttributes {
		htmlTag := eadTagRenderAttributeToHTMLTagName[renderAttribute]
		for _, eadTagToTest := range eadTagsToTest {
			testCase := testCase{
				name: fmt.Sprintf(`<%s render="%s">`,
					eadTagToTest, renderAttribute),
				eadString: fmt.Sprintf(`%s<%s render="%s">%s</%s>%s`,
					textBeforeEADTag,
					eadTagToTest, renderAttribute, eadElementTextWithInvalidOSCCharacters, eadTagToTest,
					textAfterEADTag),
				expectedHTMLString: fmt.Sprintf("%s<%s>%s</%s>%s",
					textBeforeEADTag,
					htmlTag, eadElementText, htmlTag,
					textAfterEADTag),
			}
			testCases = append(testCases, testCase)
		}
	}

	for _, testCase := range testCases {
		actual, err := ConvertEADToHTML(testCase.eadString)
		if err != nil {
			t.Errorf(`%s: expected no error, but got error: "%s"`, testCase.name,
				err)
		}

		if actual != testCase.expectedHTMLString {
			t.Errorf(`%s: expected EAD string "%s" to be converted to HTML string "%s", but got "%s"`,
				testCase.name, testCase.eadString, testCase.expectedHTMLString, actual)
		}
	}
}

func testConvertEADToHTML_GracefulHandlingOfInvalidXML(t *testing.T) {
	invalidXML := `<titleproper>This is invalid EAD</emph>`
	result, err := ConvertEADToHTML(invalidXML)
	if err == nil {
		t.Errorf(`Does not return an error for "%s"`, invalidXML)
	}
	if result != invalidXML {
		t.Errorf(`Does not return the original string "%s" on error`, invalidXML)
	}
}

func testConvertEADToHTML_NestedTags(t *testing.T) {
	testCases := []struct {
		name               string
		eadString          string
		expectedHTMLString string
	}{
		{
			// fales/mss_270.xml
			`<title> with nested <emph> -- each has render="underline"`,
			`<title render="underline"><emph render="underline">In Process</emph></title> Volume 12, No. 2, Summer 2005`,
			"<em><em>In Process</em></em> Volume 12, No. 2, Summer 2005",
		},
		{
			// nyhs/pro056_victor_prevost.xml
			`<title> with nested <emph> -- <emph> has render="italic"`,
			`Statuary at Crystal Palace [<title><emph render="italic">Eve</emph></title> by Hiram Powers]`,
			"Statuary at Crystal Palace [<title><em>Eve</em></title> by Hiram Powers]",
		},
		{
			`<titleproper> with nested <emph> -- <titleproper> has render="bolddoublequote"`,
			`<titleproper render="bolddoublequote">This is a <emph>contrived</emph> example.</titleproper>`,
			`<strong>This is a <emph>contrived</emph> example.</strong>`,
		},
		{
			`<titleproper> with nested <emph> -- neither has a render attribute`,
			`<titleproper>This is a <emph>contrived</emph> example.</titleproper>`,
			`<titleproper>This is a <emph>contrived</emph> example.</titleproper>`,
		},
		{
			`<titleproper> with 2 layers of nested <emph> -- innermost <emph> has a render attribute`,
			`<titleproper>This <emph>is <emph render="italic">a</emph> contrived</emph> example.</titleproper>`,
			`<titleproper>This <emph>is <em>a</em> contrived</emph> example.</titleproper>`,
		},
	}

	for _, testCase := range testCases {
		actual, err := ConvertEADToHTML(testCase.eadString)
		if err != nil {
			t.Errorf(`%s: expected no error, but got error: "%s"`, testCase.name,
				err)
		}

		if actual != testCase.expectedHTMLString {
			t.Errorf(`%s: expected EAD string "%s" to be converted to HTML string "%s", but got "%s"`,
				testCase.name, testCase.eadString, testCase.expectedHTMLString, actual)
		}
	}
}

func testConvertEADToHTML_Specificity(t *testing.T) {
	eadStringTokens := []string{
		"0",
		`<date type="acquisition" normal="19880423">April 23, 1988.</date>`,
		"1",
		"<title>TITLE [no attributes]</title>",
		"2",
		`<emph id="underline" render="underline" altrender="bold">EMPH [render="underline"]</emph>`,
		"3",
		`<emph id="underline" altrender="bold">EMPH [id="underline" altrender="bold"]</emph>`,
		"4",
		`Thomson, John.  Arabia, Egypt, Abyssinia, Red Sea &amp;c.`,
		"5",
	}
	eadString := strings.Join(eadStringTokens, "")

	expectedHTMLStringTokens := []string{
		"0",
		`<date type="acquisition" normal="19880423">April 23, 1988.</date>`,
		"1",
		"<title>TITLE [no attributes]</title>",
		"2",
		`<em id="underline" altrender="bold">EMPH [render="underline"]</em>`,
		"3",
		`<emph id="underline" altrender="bold">EMPH [id="underline" altrender="bold"]</emph>`,
		"4",
		`Thomson, John.  Arabia, Egypt, Abyssinia, Red Sea &amp;c.`,
		"5",
	}
	expectedHTMLString := strings.Join(expectedHTMLStringTokens, "")

	testCases := []struct {
		name               string
		eadString          string
		expectedHTMLString string
	}{
		{
			name:               "Only converts EAD tags with `render` attributes",
			eadString:          eadString,
			expectedHTMLString: expectedHTMLString,
		},
	}

	for _, testCase := range testCases {
		actual, err := ConvertEADToHTML(testCase.eadString)
		if err != nil {
			t.Errorf(`%s: expected no error, but got error: "%s"`, testCase.name,
				err)
		}

		if actual != testCase.expectedHTMLString {
			t.Errorf(`%s: expected EAD string "%s" to be converted to HTML string "%s", but got "%s"`,
				testCase.name, testCase.eadString, testCase.expectedHTMLString, actual)
		}
	}
}

func TestGetDateParts(t *testing.T) {
	testCases := []struct {
		name              string
		dateString        string
		expectedDateParts DateParts
	}{
		{
			"Gets start and end date for valid date string",
			"2016/2020",
			DateParts{
				Start: "2016",
				End:   "2020",
			},
		},
		// TODO: DLFA-238
		// Delete this after passing the transition test and resolving this:
		// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=11550822&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-11550822.
		{
			"Gets start and end date for valid date string: allow yyyy-mm-dd",
			// Value "1911/2023-03-27" appears in mc_286aspace_14c7ab764c20a3d6960975f319b33a4e
			"1911/2023-03-27",
			DateParts{
				Start: "1911",
				End:   "2023-03-27",
			},
		},
		// TODO: DLFA-238
		// Re-enable this test after passing transition test and resolving this
		// v1 indexer bug which allows for the invalid date format in this test case:
		// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=11550822&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-11550822.
		//{
		//	"Returns empty `DateParts` for ambiguous date string",
		//	"2016/2020/2024",
		//	DateParts{},
		//},
		{
			"Returns empty `DateParts` for date string with hyphen",
			"2016-2020",
			DateParts{},
		},
		{
			"Returns empty `DateParts` for invalid date string",
			"BAD DATES, INDY!",
			DateParts{},
		},
		{
			"Returns empty `DateParts` for empty date string",
			"",
			DateParts{},
		},
	}

	for _, testCase := range testCases {
		actual := GetDateParts(testCase.dateString)
		if actual.Start != testCase.expectedDateParts.Start || actual.End != testCase.expectedDateParts.End {
			t.Errorf(`%s: expected start="%s" and end="%s" for date string="%s", but got start="%s" and end="%s"`,
				testCase.name, testCase.expectedDateParts.Start, testCase.expectedDateParts.End,
				testCase.dateString, actual.Start, actual.End)
		}
	}
}

func TestGetDateRange(t *testing.T) {
	testCases := []struct {
		name              string
		unitDates         []string
		expectedDateRange []string
	}{
		{
			"Maps multiple in-range dates and returns date ranges in the right order",
			[]string{
				// Wholly within a single range
				"2016/2020",
				// Start date not in any range, but end date within a range
				"0001/2100",
				// Start date within a range, but end date not within any range
				"1101/9999",
				// Start date within one range and end date within another
				"1201/1901",
			},
			[]string{
				"1101-1200",
				"1201-1300",
				"1901-2000",
				"2001-2100",
			},
		},
		{
			"Returns undated for one mappable date and one syntactically valid but unmappable date",
			[]string{
				"2016/2020",
				"0001/0002",
			},
			[]string{
				undated,
			},
		},
		{
			"Returns undated for one mappable date and one syntactically invalid date",
			[]string{
				"BAD DATES, INDY!",
				"2016/2020",
			},
			[]string{
				undated,
			},
		},
		{
			"Returns undated when no dates",
			[]string{},
			[]string{
				undated,
			},
		},
	}

	for _, testCase := range testCases {
		actual := GetDateRange(testCase.unitDates)
		if slices.Compare(actual, testCase.expectedDateRange) != 0 {
			t.Errorf(`%s: expected dates "%v" to map to ranges %v, but got ranges %v`,
				testCase.name, testCase.unitDates, testCase.expectedDateRange, actual)
		}
	}
}

func TestGetUnitDateDisplay(t *testing.T) {
	testCases := []struct {
		name                    string
		unitDateNoTypeAttribute []string
		unitDateInclusive       []string
		unitDateBulk            []string
		expected                string
	}{
		{
			"`unitDateNoTypeAttribute`, `unitDateInclusive`, `unitDateBulk` all absent",
			[]string{},
			[]string{},
			[]string{},
			"",
		},
		{
			"`unitDateNoTypeAttribute`, `unitDateInclusive`, `unitDateBulk` all present",
			[]string{"29 November 1965"},
			[]string{"1910 - 1990"},
			[]string{"1930-1960"},
			"29 November 1965",
		},
		{
			"`unitDateNoTypeAttribute` absent; `unitDateInclusive`, `unitDateBulk` present",
			[]string{},
			[]string{"1910 - 1990"},
			[]string{"1930-1960"},
			"Inclusive, 1910 - 1990 ; 1930-1960",
		},
		{
			"`unitDateNoTypeAttribute` absent; `unitDateInclusive` present and `unitDateBulk` absent, ",
			[]string{},
			[]string{"1910 - 1990"},
			[]string{},
			"Inclusive, 1910 - 1990",
		},
		// TODO: DLFA-238
		// For now, preserve v1 indexer bug https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=8378822&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-8378822
		{
			"`unitDateNoTypeAttribute` absent; `unitDateInclusive` absent and `unitDateBulk` present",
			[]string{},
			[]string{},
			[]string{"1930-1990"},
			"Inclusive,  ; 1930-1990",
		},
	}

	for _, testCase := range testCases {
		actual := GetUnitDateDisplay(testCase.unitDateNoTypeAttribute, testCase.unitDateInclusive,
			testCase.unitDateBulk)
		if actual != testCase.expected {
			t.Errorf(`%s: expected unitDateNoTypeAttribute=%v, unitDateInclusive=%v, unitDateBulk=%v to return "%s", got "%s"`,
				testCase.name, testCase.unitDateNoTypeAttribute, testCase.unitDateInclusive,
				testCase.unitDateBulk, testCase.expected, actual)
		}
	}
}

func TestIsDateInRange(t *testing.T) {
	testCases := []struct {
		name       string
		dateString string
		dateRange  DateRange
		expected   bool
	}{
		{
			"Returns true for wholly in range date",
			"2016/2020",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			true,
		},
		{
			"Returns true for start date in range but end date not in range",
			"2016/9999",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			true,
		},
		{
			"Returns true for start date not in range but end date in range",
			"0001/2020",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			true,
		},
		{
			"Returns true start and end date on exact borders",
			"2001/2100",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			true,
		},
		{
			"Returns true for start date on border and end date out of range",
			"2001/9999",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			true,
		},
		{
			"Returns true for start date out of range and end on border",
			"0001/2100",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			true,
		},
		{
			"Returns true for start date in one range and end date in another",
			"1200/1900",
			DateRange{Display: "1101-1200", StartDate: 1101, EndDate: 1200},
			true,
		},
		{
			"Returns true for wholly in range, with allowable leading and trailing whitespace",
			" 2016/2020\t\r\n",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			true,
		},
		{
			"Returns false for wholly out of range",
			"0001/9999",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			false,
		},
		{
			"Returns false for empty string",
			"",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			false,
		},
		{
			"Returns false for all whitespace",
			"         ",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			false,
		},
		{
			"Returns false for hyphen instead of slash",
			"2016-2020",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			false,
		},
		{
			"Returns false for single year instead of two years",
			"2016",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			false,
		},
		{
			`Returns false for YYYY-MM-DD formats and " to " instead of slash`,
			"2016-01-01 to 2020-12-31",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			false,
		},
		// TODO: DLFA-238
		// Re-enable this and delete `true` expected result test after passing
		// transition test and confirming from stakeholders that they don't want
		// to pass date strings like this. Or, if they do wish this permissiveness
		// then delete this and keep the `true` test.
		//{
		//	"Returns false for too many date years",
		//	"2016/2017/2018/2019/2020",
		//	DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
		//	false,
		//},
		{
			"Returns true for string of years",
			"2016/2017/2018/2019/2020",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			true,
		},
		{
			"Returns false for not a date",
			"BAD DATES, INDY!",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			false,
		},
	}

	for _, testCase := range testCases {
		actual := isDateInRange(testCase.dateString, testCase.dateRange)
		if actual != testCase.expected {
			t.Errorf(`%s: expected "%s" in "%s" to return %t, got %t`,
				testCase.name, testCase.dateString, testCase.dateRange.Display, testCase.expected, actual)
		}
	}
}

// The `language` package has its own test suite, so we don't need to go crazy
// with coverage here.
func TestLanguage(t *testing.T) {
	testCases := []struct {
		name              string
		langCodes         []string
		expectedLanguages []string
		expectedErrors    []error
	}{
		{
			`Simple lookup`,
			[]string{"ara", "eng"},
			[]string{"Arabic", "English"},
			[]error{},
		},
		{
			`Language code not found`,
			[]string{"xxx"},
			[]string{},
			[]error{languageLib.ErrLanguageNotFound},
		},
		{
			`Invalid language codes`,
			[]string{
				"",
				"!!!",
				"abcdefghijklmnopqrstuvwxyz",
				"a",
				"e n g",
			},
			[]string{},
			[]error{
				languageLib.ErrEmptyLanguageCode,
				languageLib.ErrInvalidCharacters,
				languageLib.ErrInvalidLength,
				languageLib.ErrInvalidLength,
				languageLib.ErrInternalWhitespace,
			},
		},
	}

	for _, testCase := range testCases {
		actualLanguages, actualErrors := GetLanguage(testCase.langCodes)
		if slices.Compare(actualLanguages, testCase.expectedLanguages) != 0 {
			t.Errorf(`%s: expected language codes "%v" to map to languages %v, but got %v`,
				testCase.name, testCase.langCodes, testCase.expectedLanguages,
				actualLanguages)
		}
		if len(actualErrors) > 0 || len(testCase.expectedErrors) > 0 {
			// `slices.Compare(actualErrors, testCase.expectedErrors)` because
			// `error` is not type `cmp.Ordered`.
			actualErrorStrings := []string{}
			for _, err := range actualErrors {
				actualErrorStrings = append(actualErrorStrings, err.Error())
			}
			expectedErrorStrings := []string{}
			for _, err := range testCase.expectedErrors {
				expectedErrorStrings = append(expectedErrorStrings, err.Error())
			}

			if slices.Compare(actualErrorStrings, expectedErrorStrings) != 0 {
				t.Errorf(`%s: expected language codes "%v" to generate errors %v, but got %v`,
					testCase.name, testCase.langCodes, expectedErrorStrings,
					actualErrorStrings)
			}
		}
	}
}

// `MakeTitleHTML()` just calls `ConvertEADToHTML()` and `StripTags()` in succession.
// Those function already have their own unit test coverage, so we don't necessarily
// have to do extensive testing here, which would just have to be mechanically
// updated every time we update those two functions.
func TestMakeTitleHTML(t *testing.T) {
	eadStringTokens := []string{
		"0",
		`<date type="acquisition" normal="19880423">April 23, 1988.</date>`,
		"1",
		"<title>TITLE [no attributes]</title>",
		"2",
		`<emph id="underline" render="underline" altrender="bold">EMPH [render="underline"]</emph>`,
		"3",
		`<emph id="underline" altrender="bold">EMPH [id="underline" altrender="bold"]</emph>`,
		"4",
		"<strong>STRONG</strong>",
		"5",
		"<lb/>",
		"6",
		"<br></br>",
		"7",
		`Statuary at Crystal Palace [<title><emph render="italic">Eve</emph></title> by Hiram Powers]`,
		"8",
	}
	eadString := strings.Join(eadStringTokens, "")

	expectedHTMLStringTokens := []string{
		"0",
		`April 23, 1988.`,
		"1",
		"TITLE [no attributes]",
		"2",
		`<em id="underline" altrender="bold">EMPH [render="underline"]</em>`,
		"3",
		`EMPH [id="underline" altrender="bold"]`,
		"4",
		`<strong>STRONG</strong>`,
		"5",
		"",
		"6",
		"",
		"7",
		"Statuary at Crystal Palace [<em>Eve</em> by Hiram Powers]",
		"8",
	}
	expectedHTMLString := strings.Join(expectedHTMLStringTokens, "")

	testCases := []struct {
		name               string
		eadString          string
		expectedHTMLString string
	}{
		{
			name:               "Basic test",
			eadString:          eadString,
			expectedHTMLString: expectedHTMLString,
		},
	}

	for _, testCase := range testCases {
		actual, err := MakeTitleHTML(testCase.eadString)
		if err != nil {
			t.Errorf(`%s: expected no error, but got error: "%s"`, testCase.name,
				err)
		}

		if actual != testCase.expectedHTMLString {
			t.Errorf(`%s: expected EAD string "%s" to be converted to HTML string "%s", but got "%s"`,
				testCase.name, testCase.eadString, testCase.expectedHTMLString, actual)
		}
	}
}

func TestPrettifySolrAddMessageXML(t *testing.T) {
	// https://github.com/NYULibraries/dlfa-188_v1-indexer-http-requests/blob/206386a464e2b1280021571cbd4e73218990c26c/http-requests/fales/mss_420/mss_420-add.txt
	const inputXML = `<?xml version="1.0" encoding="UTF-8"?><add><doc><field name="author_teim">Laurainne Ojo-Ohikuare, 2015, and Stacey Flatt with the assistance of Madison DeLaere, 2023 </field><field name="author_ssm">Laurainne Ojo-Ohikuare, 2015, and Stacey Flatt with the assistance of Madison DeLaere, 2023 </field><field name="unittitle_teim">Hot Peaches Records</field><field name="unittitle_ssm">Hot Peaches Records</field><field name="unitid_teim">MSS.420</field><field name="unitid_ssm">MSS.420</field><field name="abstract_teim">Hot Peaches was a drag theater company founded in 1972 by writer and performer Jimmy Camicia. The group was known for their collaborative interactions with other drag performance groups and queer revolutionaries, including cast member Marsha P. Johnson. The group was based in New York City and performed shows regularly from the 1970s to the early 2000s not only in New York clubs, but also tours in Europe. Hot Peaches was often a platform for cast members' self-expression, as well as a nurturing queer community. The Hot Peaches Records contain material connected to their performances including scripts; press clippings; show production work books; digital and analog performance recordings; cast and performance photographs; and promotional announcements and posters. This collection also contains personal material created by Jimmy Camicia including school ephemera, writings, and journals dating from the late 1950s to 2012.</field><field name="abstract_ssm">Hot Peaches was a drag theater company founded in 1972 by writer and performer Jimmy Camicia. The group was known for their collaborative interactions with other drag performance groups and queer revolutionaries, including cast member Marsha P. Johnson. The group was based in New York City and performed shows regularly from the 1970s to the early 2000s not only in New York clubs, but also tours in Europe. Hot Peaches was often a platform for cast members' self-expression, as well as a nurturing queer community. The Hot Peaches Records contain material connected to their performances including scripts; press clippings; show production work books; digital and analog performance recordings; cast and performance photographs; and promotional announcements and posters. This collection also contains personal material created by Jimmy Camicia including school ephemera, writings, and journals dating from the late 1950s to 2012.</field><field name="creator_teim">Camicia, Jimmy</field><field name="creator_ssm">Camicia, Jimmy</field><field name="creator_ssm">Camicia, Jimmy</field><field name="unitdate_normal_ssm">1950/2016</field><field name="unitdate_normal_ssm">1971/1992</field><field name="unitdate_normal_teim">1950/2016</field><field name="unitdate_normal_teim">1971/1992</field><field name="unitdate_normal_sim">1950/2016</field><field name="unitdate_normal_sim">1971/1992</field><field name="unitdate_bulk_teim">1971-1992</field><field name="unitdate_inclusive_teim">1950s-2014, undated</field><field name="scopecontent_teim">The Hot Peaches Records (1971-2014) include analog and digital recordings, and paper material related to Hot Peaches, a New York City-based gay theater group who based their shows on political camp and were dominated by drag performers. The Hot Peaches Records document the theater company's artistic processes, as well as the broader context of avant-garde musical performance, male and female gay performers, and queer communities in New York City in the late 20th century.
Performance audio and video recordings in this collection date from the 1970s to the 2000s and include not only Hot Peaches shows, but also fundraisers to assist former castmates like International Chrysis and Ian McKay. The collection also includes photographs of the cast, performances, backstage, and images from their European tours. Publicity material includes fliers, postcards, posters, and press releases. Show production work books include notated scripts, cast lists, stage directions, press clippings, and photographs.</field><field name="bioghist_teim">Hot Peaches was a theater company working mostly in drag in New York City during the 1970s to the early 2000s. The company was founded by Jimmy Camicia in 1972, who befriended a group of drag queens in New York and was inspired to write shows for them to perform. Self-defined as as gay theater group, the Hot Peaches created shows that expressed the gay experience with a campy and political twist. Early shows were known for the castmembers' costumes, which often included vibrant, sparkling glam outfits with liberal use of platform boots, glitter, and feather boas. Camicia almost exclusively wrote the scripts, and the shows were put on three to five times per year at a variety of small theaters in Manhattan including Peach Pitts, Theater for the New City, La Mama, and Theater Genesis. The group would also occasionally do European tours, performing in England, Amsterdam, Scotland, Italy, and Germany, often providing their audiences with an experience that was not necessarily available for the gay communities in these countries. Notable cast members include the activist and performer Marsha P. Johnson, Sister Tooey, Wilhelmina Ross, Ian McKay, and Split Britches founder Peggy Shaw.</field><field name="acqinfo_teim">Donated by Jimmy Camicia, 2014 and 2023. The accession numbers associated with these gifts are 2014.420 and 2023.048.</field><field name="appraisal_teim">The following were removed from the collection: 61 DVDs and CDs (either damaged or duplicates); 19 reference copy DVDs of analog recordings; 2 commercial VHS; approximately 100 rolled posters (duplicates/outside the collection's scope); unmarked commercial sheet music; and approximately 10 publications and books that had existing copies within the NYU library.</field><field name="phystech_teim">Advance notice is required for the use of computer records. Original physical digital media is restricted.
An access terminal for born-digital materials in the collection is available by appointment for reading room viewing and listening only. Researchers may view an item's original container and/or carrier, but the physical carriers themselves are not available for use because of preservation concerns.</field><field name="phystech_teim">Some audiovisual materials have not been preserved and may not be available to researchers. Materials not yet digitized will need to have access copies made before they can be used. To request an access copy, or if you are unsure if an item has been digitized, please contact Fales Library and Special Collections, special.collections@nyu.edu, 212-998-2596 with the collection name, collection number, and a description of the item(s) requested. A staff member will respond to you with further information.


Access to some of the audiovisual materials in this collection is available through digitized access copies. Researchers may view an item's original container, but the media themselves are not available for playback because of preservation concerns. Materials that have already been digitized are noted in the collection's finding aid and can be requested in our reading room.</field><field name="corpname_teim">Fales Library and Special Collections</field><field name="corpname_teim">Hot Peaches</field><field name="corpname_ssm">Fales Library and Special Collections</field><field name="corpname_ssm">Hot Peaches</field><field name="genreform_teim">Scripts (documents)</field><field name="genreform_teim">Video recordings.</field><field name="genreform_teim">Audiocassettes.</field><field name="genreform_teim">Color photographs.</field><field name="genreform_teim">Diaries</field><field name="genreform_ssm">Scripts (documents)</field><field name="genreform_ssm">Video recordings.</field><field name="genreform_ssm">Audiocassettes.</field><field name="genreform_ssm">Color photographs.</field><field name="genreform_ssm">Diaries</field><field name="persname_teim">Camicia, Jimmy</field><field name="persname_teim">Camicia, Jimmy</field><field name="persname_teim">Camicia, Jimmy</field><field name="persname_teim">Johnson, Marsha P., 1945-1992</field><field name="persname_teim">International Chrysis</field><field name="persname_ssm">Camicia, Jimmy</field><field name="persname_ssm">Camicia, Jimmy</field><field name="persname_ssm">Camicia, Jimmy</field><field name="persname_ssm">Johnson, Marsha P., 1945-1992</field><field name="persname_ssm">International Chrysis</field><field name="subject_teim">Drag shows</field><field name="subject_teim">Drag queens</field><field name="subject_teim">Artists and theater</field><field name="subject_teim">Gay theater -- United States</field><field name="subject_teim"> Drag community</field><field name="subject_teim">Gender identity in the theater</field><field name="subject_teim">Musical theater -- New York (State) -- New York</field><field name="subject_teim">Gay liberation movement -- United States.</field><field name="subject_teim">Transgender people -- United States</field><field name="subject_teim">Theatrical companies</field><field name="subject_teim">Theatrical managers</field><field name="subject_teim">Musical theater</field><field name="subject_teim">Drag shows</field><field name="subject_teim">Drag queens</field><field name="subject_teim">Artists and theater</field><field name="subject_teim">Gay theater -- United States</field><field name="subject_teim"> Drag community</field><field name="subject_teim">Gender identity in the theater</field><field name="subject_teim">Musical theater -- New York (State) -- New York</field><field name="subject_teim">Gay liberation movement -- United States.</field><field name="subject_teim">Transgender people -- United States</field><field name="subject_teim">Theatrical companies</field><field name="subject_teim">Theatrical managers</field><field name="subject_teim">Musical theater</field><field name="subject_ssm">Drag shows</field><field name="subject_ssm">Drag queens</field><field name="subject_ssm">Artists and theater</field><field name="subject_ssm">Gay theater -- United States</field><field name="subject_ssm"> Drag community</field><field name="subject_ssm">Gender identity in the theater</field><field name="subject_ssm">Musical theater -- New York (State) -- New York</field><field name="subject_ssm">Gay liberation movement -- United States.</field><field name="subject_ssm">Transgender people -- United States</field><field name="subject_ssm">Theatrical companies</field><field name="subject_ssm">Theatrical managers</field><field name="subject_ssm">Musical theater</field><field name="collection_sim">Hot Peaches Records</field><field name="collection_ssm">Hot Peaches Records</field><field name="collection_teim">Hot Peaches Records</field><field name="id">mss_420</field><field name="ead_ssi">mss_420</field><field name="repository_ssi">fales</field><field name="repository_sim">fales</field><field name="repository_ssm">fales</field><field name="format_sim">Archival Collection</field><field name="format_ssm">Archival Collection</field><field name="format_ii">0</field><field name="creator_sim">Camicia, Jimmy</field><field name="name_sim">Hot Peaches</field><field name="name_sim">Camicia, Jimmy</field><field name="name_sim">Johnson, Marsha P., 1945-1992</field><field name="name_sim">International Chrysis</field><field name="name_teim">Hot Peaches</field><field name="name_teim">Camicia, Jimmy</field><field name="name_teim">Johnson, Marsha P., 1945-1992</field><field name="name_teim">International Chrysis</field><field name="subject_sim">Drag shows</field><field name="subject_sim">Drag queens</field><field name="subject_sim">Artists and theater</field><field name="subject_sim">Gay theater -- United States</field><field name="subject_sim"> Drag community</field><field name="subject_sim">Gender identity in the theater</field><field name="subject_sim">Musical theater -- New York (State) -- New York</field><field name="subject_sim">Gay liberation movement -- United States.</field><field name="subject_sim">Transgender people -- United States</field><field name="subject_sim">Theatrical companies</field><field name="subject_sim">Theatrical managers</field><field name="subject_sim">Musical theater</field><field name="dao_sim">Online Access</field><field name="material_type_sim">Scripts (documents)</field><field name="material_type_sim">Video recordings.</field><field name="material_type_sim">Audiocassettes.</field><field name="material_type_sim">Color photographs.</field><field name="material_type_sim">Diaries</field><field name="material_type_ssm">Scripts (documents)</field><field name="material_type_ssm">Video recordings.</field><field name="material_type_ssm">Audiocassettes.</field><field name="material_type_ssm">Color photographs.</field><field name="material_type_ssm">Diaries</field><field name="heading_ssm">Hot Peaches Records</field><field name="unitdate_start_sim">1950</field><field name="unitdate_start_sim">1971</field><field name="unitdate_start_ssm">1950</field><field name="unitdate_start_ssm">1971</field><field name="unitdate_start_si">1971</field><field name="unitdate_end_sim">2016</field><field name="unitdate_end_sim">1992</field><field name="unitdate_end_ssm">2016</field><field name="unitdate_end_ssm">1992</field><field name="unitdate_end_si">1992</field><field name="unitdate_ssm">Inclusive, 1950s-2014, undated ; 1971-1992</field><field name="date_range_sim">1901-2000</field><field name="date_range_sim">2001-2100</field></doc></add>`

	// Generated by `xmllint --format`
	const expectedXML = `<?xml version="1.0" encoding="UTF-8"?>
<add>
  <doc>
    <field name="author_teim">Laurainne Ojo-Ohikuare, 2015, and Stacey Flatt with the assistance of Madison DeLaere, 2023 </field>
    <field name="author_ssm">Laurainne Ojo-Ohikuare, 2015, and Stacey Flatt with the assistance of Madison DeLaere, 2023 </field>
    <field name="unittitle_teim">Hot Peaches Records</field>
    <field name="unittitle_ssm">Hot Peaches Records</field>
    <field name="unitid_teim">MSS.420</field>
    <field name="unitid_ssm">MSS.420</field>
    <field name="abstract_teim">Hot Peaches was a drag theater company founded in 1972 by writer and performer Jimmy Camicia. The group was known for their collaborative interactions with other drag performance groups and queer revolutionaries, including cast member Marsha P. Johnson. The group was based in New York City and performed shows regularly from the 1970s to the early 2000s not only in New York clubs, but also tours in Europe. Hot Peaches was often a platform for cast members' self-expression, as well as a nurturing queer community. The Hot Peaches Records contain material connected to their performances including scripts; press clippings; show production work books; digital and analog performance recordings; cast and performance photographs; and promotional announcements and posters. This collection also contains personal material created by Jimmy Camicia including school ephemera, writings, and journals dating from the late 1950s to 2012.</field>
    <field name="abstract_ssm">Hot Peaches was a drag theater company founded in 1972 by writer and performer Jimmy Camicia. The group was known for their collaborative interactions with other drag performance groups and queer revolutionaries, including cast member Marsha P. Johnson. The group was based in New York City and performed shows regularly from the 1970s to the early 2000s not only in New York clubs, but also tours in Europe. Hot Peaches was often a platform for cast members' self-expression, as well as a nurturing queer community. The Hot Peaches Records contain material connected to their performances including scripts; press clippings; show production work books; digital and analog performance recordings; cast and performance photographs; and promotional announcements and posters. This collection also contains personal material created by Jimmy Camicia including school ephemera, writings, and journals dating from the late 1950s to 2012.</field>
    <field name="creator_teim">Camicia, Jimmy</field>
    <field name="creator_ssm">Camicia, Jimmy</field>
    <field name="creator_ssm">Camicia, Jimmy</field>
    <field name="unitdate_normal_ssm">1950/2016</field>
    <field name="unitdate_normal_ssm">1971/1992</field>
    <field name="unitdate_normal_teim">1950/2016</field>
    <field name="unitdate_normal_teim">1971/1992</field>
    <field name="unitdate_normal_sim">1950/2016</field>
    <field name="unitdate_normal_sim">1971/1992</field>
    <field name="unitdate_bulk_teim">1971-1992</field>
    <field name="unitdate_inclusive_teim">1950s-2014, undated</field>
    <field name="scopecontent_teim">The Hot Peaches Records (1971-2014) include analog and digital recordings, and paper material related to Hot Peaches, a New York City-based gay theater group who based their shows on political camp and were dominated by drag performers. The Hot Peaches Records document the theater company's artistic processes, as well as the broader context of avant-garde musical performance, male and female gay performers, and queer communities in New York City in the late 20th century.
Performance audio and video recordings in this collection date from the 1970s to the 2000s and include not only Hot Peaches shows, but also fundraisers to assist former castmates like International Chrysis and Ian McKay. The collection also includes photographs of the cast, performances, backstage, and images from their European tours. Publicity material includes fliers, postcards, posters, and press releases. Show production work books include notated scripts, cast lists, stage directions, press clippings, and photographs.</field>
    <field name="bioghist_teim">Hot Peaches was a theater company working mostly in drag in New York City during the 1970s to the early 2000s. The company was founded by Jimmy Camicia in 1972, who befriended a group of drag queens in New York and was inspired to write shows for them to perform. Self-defined as as gay theater group, the Hot Peaches created shows that expressed the gay experience with a campy and political twist. Early shows were known for the castmembers' costumes, which often included vibrant, sparkling glam outfits with liberal use of platform boots, glitter, and feather boas. Camicia almost exclusively wrote the scripts, and the shows were put on three to five times per year at a variety of small theaters in Manhattan including Peach Pitts, Theater for the New City, La Mama, and Theater Genesis. The group would also occasionally do European tours, performing in England, Amsterdam, Scotland, Italy, and Germany, often providing their audiences with an experience that was not necessarily available for the gay communities in these countries. Notable cast members include the activist and performer Marsha P. Johnson, Sister Tooey, Wilhelmina Ross, Ian McKay, and Split Britches founder Peggy Shaw.</field>
    <field name="acqinfo_teim">Donated by Jimmy Camicia, 2014 and 2023. The accession numbers associated with these gifts are 2014.420 and 2023.048.</field>
    <field name="appraisal_teim">The following were removed from the collection: 61 DVDs and CDs (either damaged or duplicates); 19 reference copy DVDs of analog recordings; 2 commercial VHS; approximately 100 rolled posters (duplicates/outside the collection's scope); unmarked commercial sheet music; and approximately 10 publications and books that had existing copies within the NYU library.</field>
    <field name="phystech_teim">Advance notice is required for the use of computer records. Original physical digital media is restricted.
An access terminal for born-digital materials in the collection is available by appointment for reading room viewing and listening only. Researchers may view an item's original container and/or carrier, but the physical carriers themselves are not available for use because of preservation concerns.</field>
    <field name="phystech_teim">Some audiovisual materials have not been preserved and may not be available to researchers. Materials not yet digitized will need to have access copies made before they can be used. To request an access copy, or if you are unsure if an item has been digitized, please contact Fales Library and Special Collections, special.collections@nyu.edu, 212-998-2596 with the collection name, collection number, and a description of the item(s) requested. A staff member will respond to you with further information.


Access to some of the audiovisual materials in this collection is available through digitized access copies. Researchers may view an item's original container, but the media themselves are not available for playback because of preservation concerns. Materials that have already been digitized are noted in the collection's finding aid and can be requested in our reading room.</field>
    <field name="corpname_teim">Fales Library and Special Collections</field>
    <field name="corpname_teim">Hot Peaches</field>
    <field name="corpname_ssm">Fales Library and Special Collections</field>
    <field name="corpname_ssm">Hot Peaches</field>
    <field name="genreform_teim">Scripts (documents)</field>
    <field name="genreform_teim">Video recordings.</field>
    <field name="genreform_teim">Audiocassettes.</field>
    <field name="genreform_teim">Color photographs.</field>
    <field name="genreform_teim">Diaries</field>
    <field name="genreform_ssm">Scripts (documents)</field>
    <field name="genreform_ssm">Video recordings.</field>
    <field name="genreform_ssm">Audiocassettes.</field>
    <field name="genreform_ssm">Color photographs.</field>
    <field name="genreform_ssm">Diaries</field>
    <field name="persname_teim">Camicia, Jimmy</field>
    <field name="persname_teim">Camicia, Jimmy</field>
    <field name="persname_teim">Camicia, Jimmy</field>
    <field name="persname_teim">Johnson, Marsha P., 1945-1992</field>
    <field name="persname_teim">International Chrysis</field>
    <field name="persname_ssm">Camicia, Jimmy</field>
    <field name="persname_ssm">Camicia, Jimmy</field>
    <field name="persname_ssm">Camicia, Jimmy</field>
    <field name="persname_ssm">Johnson, Marsha P., 1945-1992</field>
    <field name="persname_ssm">International Chrysis</field>
    <field name="subject_teim">Drag shows</field>
    <field name="subject_teim">Drag queens</field>
    <field name="subject_teim">Artists and theater</field>
    <field name="subject_teim">Gay theater -- United States</field>
    <field name="subject_teim"> Drag community</field>
    <field name="subject_teim">Gender identity in the theater</field>
    <field name="subject_teim">Musical theater -- New York (State) -- New York</field>
    <field name="subject_teim">Gay liberation movement -- United States.</field>
    <field name="subject_teim">Transgender people -- United States</field>
    <field name="subject_teim">Theatrical companies</field>
    <field name="subject_teim">Theatrical managers</field>
    <field name="subject_teim">Musical theater</field>
    <field name="subject_teim">Drag shows</field>
    <field name="subject_teim">Drag queens</field>
    <field name="subject_teim">Artists and theater</field>
    <field name="subject_teim">Gay theater -- United States</field>
    <field name="subject_teim"> Drag community</field>
    <field name="subject_teim">Gender identity in the theater</field>
    <field name="subject_teim">Musical theater -- New York (State) -- New York</field>
    <field name="subject_teim">Gay liberation movement -- United States.</field>
    <field name="subject_teim">Transgender people -- United States</field>
    <field name="subject_teim">Theatrical companies</field>
    <field name="subject_teim">Theatrical managers</field>
    <field name="subject_teim">Musical theater</field>
    <field name="subject_ssm">Drag shows</field>
    <field name="subject_ssm">Drag queens</field>
    <field name="subject_ssm">Artists and theater</field>
    <field name="subject_ssm">Gay theater -- United States</field>
    <field name="subject_ssm"> Drag community</field>
    <field name="subject_ssm">Gender identity in the theater</field>
    <field name="subject_ssm">Musical theater -- New York (State) -- New York</field>
    <field name="subject_ssm">Gay liberation movement -- United States.</field>
    <field name="subject_ssm">Transgender people -- United States</field>
    <field name="subject_ssm">Theatrical companies</field>
    <field name="subject_ssm">Theatrical managers</field>
    <field name="subject_ssm">Musical theater</field>
    <field name="collection_sim">Hot Peaches Records</field>
    <field name="collection_ssm">Hot Peaches Records</field>
    <field name="collection_teim">Hot Peaches Records</field>
    <field name="id">mss_420</field>
    <field name="ead_ssi">mss_420</field>
    <field name="repository_ssi">fales</field>
    <field name="repository_sim">fales</field>
    <field name="repository_ssm">fales</field>
    <field name="format_sim">Archival Collection</field>
    <field name="format_ssm">Archival Collection</field>
    <field name="format_ii">0</field>
    <field name="creator_sim">Camicia, Jimmy</field>
    <field name="name_sim">Hot Peaches</field>
    <field name="name_sim">Camicia, Jimmy</field>
    <field name="name_sim">Johnson, Marsha P., 1945-1992</field>
    <field name="name_sim">International Chrysis</field>
    <field name="name_teim">Hot Peaches</field>
    <field name="name_teim">Camicia, Jimmy</field>
    <field name="name_teim">Johnson, Marsha P., 1945-1992</field>
    <field name="name_teim">International Chrysis</field>
    <field name="subject_sim">Drag shows</field>
    <field name="subject_sim">Drag queens</field>
    <field name="subject_sim">Artists and theater</field>
    <field name="subject_sim">Gay theater -- United States</field>
    <field name="subject_sim"> Drag community</field>
    <field name="subject_sim">Gender identity in the theater</field>
    <field name="subject_sim">Musical theater -- New York (State) -- New York</field>
    <field name="subject_sim">Gay liberation movement -- United States.</field>
    <field name="subject_sim">Transgender people -- United States</field>
    <field name="subject_sim">Theatrical companies</field>
    <field name="subject_sim">Theatrical managers</field>
    <field name="subject_sim">Musical theater</field>
    <field name="dao_sim">Online Access</field>
    <field name="material_type_sim">Scripts (documents)</field>
    <field name="material_type_sim">Video recordings.</field>
    <field name="material_type_sim">Audiocassettes.</field>
    <field name="material_type_sim">Color photographs.</field>
    <field name="material_type_sim">Diaries</field>
    <field name="material_type_ssm">Scripts (documents)</field>
    <field name="material_type_ssm">Video recordings.</field>
    <field name="material_type_ssm">Audiocassettes.</field>
    <field name="material_type_ssm">Color photographs.</field>
    <field name="material_type_ssm">Diaries</field>
    <field name="heading_ssm">Hot Peaches Records</field>
    <field name="unitdate_start_sim">1950</field>
    <field name="unitdate_start_sim">1971</field>
    <field name="unitdate_start_ssm">1950</field>
    <field name="unitdate_start_ssm">1971</field>
    <field name="unitdate_start_si">1971</field>
    <field name="unitdate_end_sim">2016</field>
    <field name="unitdate_end_sim">1992</field>
    <field name="unitdate_end_ssm">2016</field>
    <field name="unitdate_end_ssm">1992</field>
    <field name="unitdate_end_si">1992</field>
    <field name="unitdate_ssm">Inclusive, 1950s-2014, undated ; 1971-1992</field>
    <field name="date_range_sim">1901-2000</field>
    <field name="date_range_sim">2001-2100</field>
  </doc>
</add>
`
	var actualXML = PrettifySolrAddMessageXML(inputXML)
	diff := util.DiffStrings("expected XML", expectedXML,
		"actual XML", actualXML)
	if diff != "" {
		t.Errorf("Prettified XML does not match expected:\n%s", diff)
	}

}

func TestRemoveChildNodes(t *testing.T) {
	testRemoveChildNodes(t)
	testRemoveChildNodes_errors(t)
}

func testRemoveChildNodes(t *testing.T) {
	nodeWithChildNodesXMLBytes, err := os.ReadFile(path.Join(fixturesDirPath, "test.xml"))
	if err != nil {
		t.Errorf("Error reading fixture file: %s", err)
	}
	nodeWithChildNodes := string(nodeWithChildNodesXMLBytes)

	testCases := []struct {
		name            string
		nodeXML         string
		elementToRemove string
		goldenName      string
	}{
		{
			name:            "Node arg has no child nodes",
			nodeXML:         "<root></root>",
			elementToRemove: "doesnotmatter",
			goldenName:      "node-arg-has-no-child-nodes",
		},
		{
			name:            "Node arg has only a text node",
			nodeXML:         "<root>TEXT NODE</root>",
			elementToRemove: "doesnotmatter",
			goldenName:      "node-arg-has-only-a-text-node",
		},
		{
			name:            "Remove nothing",
			nodeXML:         nodeWithChildNodes,
			elementToRemove: "",
			goldenName:      "remove-nothing",
		},
		{
			name:            "Remove an element that is not present in the XML",
			nodeXML:         nodeWithChildNodes,
			elementToRemove: "omega",
			goldenName:      "remove-an-element-that-is-not-present",
		},
		{
			name:            "Remove all top-level <alpha> elements, but not <Alpha> elements",
			nodeXML:         nodeWithChildNodes,
			elementToRemove: "alpha",
			goldenName:      "remove-top-level-alpha-lowercase",
		},
		{
			name:            "Remove all top-level <Alpha> elements, but not <alpha> elements",
			nodeXML:         nodeWithChildNodes,
			elementToRemove: "Alpha",
			goldenName:      "remove-top-level-alpha-titlecase",
		},
		{
			name:            "Remove all top-level <zulu/>",
			nodeXML:         nodeWithChildNodes,
			elementToRemove: "zulu",
			goldenName:      "remove-top-level-zulu",
		},
	}

	for _, testCase := range testCases {
		xmlParser := parser.New()
		testDoc, err := xmlParser.ParseString(testCase.nodeXML)
		defer testDoc.Free()
		if err != nil {
			t.Errorf("Failed to parse test XML: %s", err)
		}

		testNode, err := testDoc.DocumentElement()
		if err != nil {
			t.Errorf("Failed to get `testNode` from `testDoc`: %s", err)
		}

		actualNode, err := testNode.Copy()
		if err != nil {
			t.Errorf("Failed to copy test node: %s", err)
		}

		err = RemoveChildNodesMatchingName(actualNode, testCase.elementToRemove)
		if err != nil {
			t.Errorf(`%s: expected no error, but got error: "%s"`, testCase.name,
				err)
		}

		actualXML := actualNode.String()

		if *updateGoldenFiles {
			err := updateGoldenFile(testCase.goldenName, actualXML)
			if err != nil {
				t.Fatalf("Error updating golden file: %s", err)
			}
		}

		expectedXML, err := getGoldenFileValue(testCase.goldenName)
		if err != nil {
			t.Errorf("Failed to get `expectedXML`: %s", err)
		}

		if actualXML != expectedXML {
			diff := util.DiffStrings("expectedXML", expectedXML,
				"actualXML", actualXML)

			t.Errorf(`%s: actual XML does not match expected XML: "%s",`,
				testCase.name, diff)
		}
	}
}

func testRemoveChildNodes_errors(t *testing.T) {
	err := RemoveChildNodesMatchingName(nil, "doesnotmatter")
	if err == nil {
		t.Errorf("Expected an error return for `nil` node arg, but didn't get one")
	}
}

func TestReplaceMARCSubfieldDemarcators(t *testing.T) {
	// To see where some of these real life examples came from:
	// https://jira.nyu.edu/browse/DLFA-229?focusedCommentId=10153922&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10153922
	testCases := []struct {
		in  string
		out string
	}{
		{
			"",
			"",
		},
		{
			"Laundry industry |z New York (State) |z New York.",
			"Laundry industry -- New York (State) -- New York.",
		},
		{
			"China |x History |x Tiananmen Square Incident, 1989",
			"China -- History -- Tiananmen Square Incident, 1989",
		},
		{
			"Labor Unions |z United States |y 1980-1990.",
			"Labor Unions -- United States -- 1980-1990.",
		},
		{
			"Elections |z United States |x History |y 20th century |v Statistics.",
			"Elections -- United States -- History -- 20th century -- Statistics.",
		},
		{
			"Randall, Margaret, |d 1936-",
			"Randall, Margaret, -- 1936-",
		},
		{
			"General strikes |Z New York (State) |z Kings County",
			"General strikes -- New York (State) -- Kings County",
		},
		{
			"Theaters |x Employees |X Labor unions |z United States.",
			"Theaters -- Employees -- Labor unions -- United States.",
		},
		{
			"France. |t Constitution (1958).",
			"France. -- Constitution (1958).",
		},
		{
			"United States. Congress. House. |b Committee on Education and Labor. |b Select Subcommittee on Education",
			"United States. Congress. House. -- Committee on Education and Labor. -- Select Subcommittee on Education",
		},
		{
			"Wagner, Richard, 1813-1883. |t Operas. |k Selections",
			"Wagner, Richard, 1813-1883. -- Operas. -- Selections",
		},
		{
			"Germany. |t Treaties, etc. |g Soviet Union, |d 1939 Aug. 23.",
			"Germany. -- Treaties, etc. -- Soviet Union, -- 1939 Aug. 23.",
		},
		{
			"DO | NOT || CHANGE",
			"DO | NOT || CHANGE",
		},
		// TODO: fix the bug we've intentionally preserved in MARC subfield demarcation
		// replacement.  For details, see:
		//
		//   - https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10154897&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10154897
		//   - https://jira.nyu.edu/browse/DLFA-229?focusedCommentId=10153922&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10153922
		//
		// Once that is done, we can uncomment these tests, which currently fail.
		//{
		//	"Violence: Recode / UNDER|STAND",
		//	"Violence: Recode / UNDER|STAND",
		//},
		//{
		//	"85-2126 | John Hans[e|o]n (from Box 4 of 6)",
		//	"85-2126 | John Hans[e|o]n (from Box 4 of 6)",
		//},
	}

	for _, testCase := range testCases {
		actual := replaceMARCSubfieldDemarcators(testCase.in)
		if actual != testCase.out {
			t.Errorf(`Expected output string "%s" for input string "%s", got "%s""`,
				testCase.out, testCase.in, actual)
		}
	}
}

// Testing doesn't have to be extensive, because the function under test is only
// used for strings returned by `types.Node.String()`.
// See function header comment for StripOpenAndCloseTags().
func TestStripOpenAndCloseTags(t *testing.T) {
	testCases := []struct {
		name   string
		before string
		after  string
	}{
		{
			"Basic case",
			`<unittitle>TITLE</unittitle>`,
			"TITLE",
		},
		{
			"Open tag has attributes",
			`<unittitle attr1="1" attr2="2">TITLE</unittitle>`,
			"TITLE",
		},
		{
			"Empty string",
			"",
			"",
		},
		{
			"Not an XML element",
			"TITLE",
			"TITLE",
		},
		{
			"Open but no close tag",
			"<unittitle>TITLE",
			"TITLE",
		},
		{
			"Close but no open tag",
			"TITLE</unittitle>",
			"TITLE",
		},
		{
			"Remove outermost tags only",
			"<unittitle><italic>TITLE</italic></unittitle>",
			"<italic>TITLE</italic>",
		},
	}

	for _, testCase := range testCases {
		actual := StripOpenAndCloseTags(testCase.before)

		if actual != testCase.after {
			t.Errorf(`%s: expected XML string "%s" to be string "%s", but got "%s"`,
				testCase.name, testCase.before, testCase.after, actual)
		}
	}
}

func TestStripTags(t *testing.T) {
	testStripTags_EmptyElements(t)
	testStripTags_Specificity(t)
}

func testStripTags_EmptyElements(t *testing.T) {
	testCases := []struct {
		name               string
		eadString          string
		expectedHTMLString string
	}{
		{
			"Single empty self-closing tag",
			`1<lb/>2`,
			"12",
		},
		{
			"Single empty element with both opening and closing tags",
			`1<lb></lb>2`,
			"12",
		},
		{
			"Single empty self-closing tag with attributes",
			`1<dimensions id="aspace_f01c5dcb7232080131a647dc8b66183b" label="21.1 x 29.7 cm"/>2`,
			"12",
		},
	}

	for _, testCase := range testCases {
		actual, err := StripTags(testCase.eadString)
		if err != nil {
			t.Errorf(`%s: expected no error, but got error: "%s"`, testCase.name,
				err)
		}

		if actual != testCase.expectedHTMLString {
			t.Errorf(`%s: expected XML string "%s" to be converted to HTML string "%s", but got "%s"`,
				testCase.name, testCase.eadString, testCase.expectedHTMLString, actual)
		}
	}
}

func testStripTags_Specificity(t *testing.T) {
	eadStringTokens := []string{
		"0",
		"<title>TITLE</title>",
		"1",
		`<em>EM</em>`,
		"2",
		"<lb/>",
		"3",
		"<br></br>",
		"4",
		`<date type="acquisition" normal="19880423">April 23, 1988.</date>`,
		"5",
		`<strong>STRONG</strong>`,
		"6",
	}
	xmlString := strings.Join(eadStringTokens, "")

	expectedHTMLStringTokens := []string{
		"0",
		"TITLE",
		"1",
		`<em>EM</em>`,
		"2",
		"",
		"3",
		"",
		"4",
		`April 23, 1988.`,
		"5",
		`<strong>STRONG</strong>`,
		"6",
	}
	expectedHTMLString := strings.Join(expectedHTMLStringTokens, "")

	testCases := []struct {
		name               string
		xmlString          string
		expectedHTMLString string
	}{
		{
			name:               "Only strips disallowed XML tags",
			xmlString:          xmlString,
			expectedHTMLString: expectedHTMLString,
		},
	}

	for _, testCase := range testCases {
		actual, err := StripTags(testCase.xmlString)
		if err != nil {
			t.Errorf(`%s: expected no error, but got error: "%s"`, testCase.name,
				err)
		}

		if actual != testCase.expectedHTMLString {
			t.Errorf(`%s: expected XML string "%s" to be converted to HTML string "%s", but got "%s"`,
				testCase.name, testCase.xmlString, testCase.expectedHTMLString, actual)
		}
	}
}

// TODO: Consolidate fixture/helper/diff code with `pkg/ead/testutils/`?
// Note the that latter was designed to be EAD file specific, and also was
// written in anticipation of potentially lifting out of the `ead` package
// entirely so it could be used by multiple packages.
func getGoldenFileValue(goldenName string) (string, error) {
	return getTestdataFileContents(goldenFilePath(goldenName))
}

func getTestdataFileContents(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return filename, err
	}

	return string(bytes), nil
}

func goldenFilePath(goldenName string) string {
	return filepath.Join(goldenFilesDirPath, goldenName+".xml")
}

func updateGoldenFile(goldenName string, data string) error {
	return os.WriteFile(goldenFilePath(goldenName), []byte(data), 0644)
}
