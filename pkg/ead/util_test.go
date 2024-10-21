package ead

import (
	"slices"
	"testing"
)

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
		actual := getDateRange(testCase.unitDates)
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
		// TODO DLFA-238
		// For now, preserve v1 indexer bug https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=8378822&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-8378822
		{
			"`unitDateNoTypeAttribute` absent; `unitDateInclusive` absent and `unitDateBulk` present",
			[]string{},
			[]string{},
			[]string{"1930-1990"},
			"Inclusive, ; 1930-1990",
		},
	}

	for _, testCase := range testCases {
		actual := getUnitDateDisplay(testCase.unitDateNoTypeAttribute, testCase.unitDateInclusive,
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
