package ead

import (
	"github.com/lestrrat-go/libxml2/types"
	"regexp"
	"strconv"
)

type DateRange struct {
	Display   string
	StartDate int
	EndDate   int
}

var datePartsRegexp = regexp.MustCompile(`^\s*(\d{4})\/(\d{4})\s*$`)

var dateRangesCenturies = []DateRange{
	{Display: "1101-1200", StartDate: 1101, EndDate: 1200},
	{Display: "1201-1300", StartDate: 1201, EndDate: 1300},
	{Display: "1301-1400", StartDate: 1301, EndDate: 1400},
	{Display: "1401-1500", StartDate: 1401, EndDate: 1500},
	{Display: "1501-1600", StartDate: 1501, EndDate: 1600},
	{Display: "1601-1700", StartDate: 1601, EndDate: 1700},
	{Display: "1701-1800", StartDate: 1701, EndDate: 1800},
	{Display: "1801-1900", StartDate: 1801, EndDate: 1900},
	{Display: "1901-2000", StartDate: 1901, EndDate: 2000},
	{Display: "2001-2100", StartDate: 2001, EndDate: 2100},
}

const undated = "undated & other"

func getDateRange(unitDates []string) []string {
	dateRange := []string{}

	// Add `dateRangeCentury` display dates for which least one date falls within
	// range.
	for _, dateRangeCentury := range dateRangesCenturies {
		for _, unitDate := range unitDates {
			if isDateInRange(unitDate, dateRangeCentury) {
				dateRange = append(dateRange, dateRangeCentury.Display)
				break
			}
		}
	}

	// Check to see if even a single date couldn't be matched to a date range.
	existsDateWithRangeNotFound := false
	for _, unitDate := range unitDates {
		matchFound := false
		for _, dateRangeCentury := range dateRangesCenturies {
			if isDateInRange(unitDate, dateRangeCentury) {
				matchFound = true
				break
			}
		}

		// No date range found for date
		if !matchFound {
			existsDateWithRangeNotFound = true
			break
		}
	}

	if len(dateRange) == 0 || existsDateWithRangeNotFound {
		dateRange = []string{undated}
	}

	return dateRange
}

func getValuesForXPathQuery(query string, node types.Node) ([]string, error) {
	var values []string

	xpathResult, err := node.Find(query)
	if err != nil {
		return nil, err
	}

	for _, node = range xpathResult.NodeList() {
		values = append(values, node.NodeValue())
	}

	return values, nil
}

// `dateString` should be of the form "YYYY/YYYY", where the left "YYYY" is the
// start date and the right "YYYY" is the end date.
func isDateInRange(dateString string, dateRange DateRange) bool {
	matches := datePartsRegexp.FindStringSubmatch(dateString)

	if len(matches) == 3 {
		startDate := matches[1]
		endDate := matches[2]

		startDateInt, err := strconv.Atoi(startDate)
		if err != nil {
			return false
		}

		endDateInt, err := strconv.Atoi(endDate)
		if err != nil {
			return false
		}

		return (startDateInt >= dateRange.StartDate && startDateInt <= dateRange.EndDate) ||
			(endDateInt >= dateRange.StartDate && endDateInt <= dateRange.EndDate)
	} else {
		return false
	}
}

func replaceMARCSubfieldDemarcatorsInSlice(stringSlice []string) []string {
	newSlice := []string{}
	for _, element := range stringSlice {
		newSlice = append(newSlice, replaceMARCSubfieldDemarcators(element))
	}

	return newSlice
}

// TODO: fix the bug we've intentionally preserved here -- for details, see:
// * https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10154897&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10154897
// * https://jira.nyu.edu/browse/DLFA-229?focusedCommentId=10153922&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10153922
func replaceMARCSubfieldDemarcators(str string) string {
	return marcSubfieldDemarcator.ReplaceAllString(str, "--")
}