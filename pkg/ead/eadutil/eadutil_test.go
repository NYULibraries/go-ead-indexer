package eadutil

import (
	"flag"
	"fmt"
	"github.com/lestrrat-go/libxml2/parser"
	languageLib "go-ead-indexer/pkg/language"
	"go-ead-indexer/pkg/util"
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
		{
			"Returns empty `DateParts` for ambiguous date string",
			"2016/2020/2024",
			DateParts{},
		},
		{
			"Returns empty `DateParts` for date string with hypen",
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
		{
			"Returns false for too many date years",
			"2016/2017/2018/2019/2020",
			DateRange{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
			false,
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
