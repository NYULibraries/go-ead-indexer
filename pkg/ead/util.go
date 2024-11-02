package ead

import (
	"encoding/xml"
	"fmt"
	"github.com/lestrrat-go/libxml2/types"
	languageLib "go-ead-indexer/pkg/language"
	"go-ead-indexer/pkg/sanitize"
	"go-ead-indexer/pkg/util"
	"io"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type DateParts struct {
	Start string
	End   string
}

type DateRange struct {
	Display   string
	StartDate int
	EndDate   int
}

const undated = "undated & other"

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

var eadTagRenderAttributeToHTMLTagName = map[string]string{
	"altrender":       "em",
	"bold":            "strong",
	"bolddoublequote": "strong",
	"bolditalic":      "strong",
	"boldsinglequote": "strong",
	"boldsmcaps":      "strong",
	"boldunderline":   "strong",
	"doublequote":     "em",
	"italic":          "em",
	"italics":         "em",
	"nonproport":      "em",
	"singlequote":     "em",
	"smcaps":          "em",
	"sub":             "sub",
	"super":           "sup",
	"underline":       "em",
}

var allowedHTMLTags = util.CompactStringSlicePreserveOrder(
	slices.Collect(maps.Values(eadTagRenderAttributeToHTMLTagName)))

func convertToFacetSlice(rawSlice []string) []string {
	return util.CompactStringSlicePreserveOrder(
		replaceMARCSubfieldDemarcatorsInSlice(rawSlice))
}

func convertEADTagsWithRenderAttributesToHTML(eadString string) (string, error) {
	var htmlString string

	var startTagNames []string

	decoder := xml.NewDecoder(strings.NewReader(eadString))

	for {
		token, err := decoder.Token()

		if err == io.EOF {
			break
		} else if err != nil {
			return eadString, err
		}

		switch token := token.(type) {
		case xml.StartElement:
			var renderAttributeValue string
			for i := range token.Attr {
				if token.Attr[i].Name.Local == "render" {
					renderAttributeValue = token.Attr[i].Value
					break
				}
			}

			if renderAttributeValue == "" {
				startTagNames = append(startTagNames, token.Name.Local)

				htmlString += stringifyStartElementToken(token)
			} else {
				if htmlTagName, ok := eadTagRenderAttributeToHTMLTagName[renderAttributeValue]; ok {
					startTagNames = append(startTagNames, htmlTagName)

					token.Name.Local = htmlTagName
					token.Attr = slices.DeleteFunc(token.Attr, func(attribute xml.Attr) bool {
						return attribute.Name.Local == "render"
					})

					htmlString += stringifyStartElementToken(token)
				} else {
					startTagNames = append(startTagNames, token.Name.Local)

					htmlString += stringifyStartElementToken(token)
				}
			}

		case xml.EndElement:
			htmlString += fmt.Sprintf("</%s>", startTagNames[len(startTagNames)-1])
			startTagNames = startTagNames[:len(startTagNames)-1]

		case xml.CharData:
			htmlString += string(token)
		}
	}

	return htmlString, nil
}

func convertEADToHTML(eadString string) (string, error) {
	htmlString, err := convertEADTagsWithRenderAttributesToHTML(eadString)
	if err != nil {
		return htmlString, err
	}

	return sanitize.Clean(htmlString), nil
}

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

func getLanguage(langCodes []string) ([]string, []error) {
	language := []string{}
	errs := []error{}

	for _, langCode := range langCodes {
		languageForLangCode, err := languageLib.GetLanguageForLanguageCode(langCode)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		language = append(language, languageForLangCode)
	}

	return language, errs
}

// TODO DLFA-238: This method preserves probable v1 indexer bug for the purposes of passing
// the DLFA-201 transition acceptance test -- the bug:
// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=8378822&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-8378822
// Fix this bug after we've completed the transition.
func getUnitDateDisplay(unitDateNoTypeAttribute []string, unitDateInclusive []string,
	unitDateBulk []string) string {
	partsUnitDateDisplay := []string{}
	if len(unitDateNoTypeAttribute) > 0 {
		partsUnitDateDisplay = unitDateNoTypeAttribute
	} else if len(unitDateInclusive) == 0 && len(unitDateBulk) == 0 {
		// Do nothing
	} else {
		partsUnitDateDisplay = append(partsUnitDateDisplay, "Inclusive,")
		partsUnitDateDisplay = append(partsUnitDateDisplay, unitDateInclusive...)
		if len(unitDateBulk) > 0 {
			partsUnitDateDisplay = append(partsUnitDateDisplay, ";")
			partsUnitDateDisplay = append(partsUnitDateDisplay, unitDateBulk...)
		}
	}

	return strings.Join(partsUnitDateDisplay, " ")
}

func getDateParts(dateString string) DateParts {
	dateParts := DateParts{}

	matches := datePartsRegexp.FindStringSubmatch(dateString)

	if len(matches) == 3 {
		dateParts.Start = matches[1]
		dateParts.End = matches[2]
	}

	return dateParts
}

func getValuesForXPathQuery(query string, node types.Node) ([]string, []string, error) {
	var values []string
	var xmlStrings []string

	xpathResult, err := node.Find(query)
	if err != nil {
		return nil, nil, err
	}

	for _, node = range xpathResult.NodeList() {
		values = append(values, node.NodeValue())
		xmlStrings = append(xmlStrings, node.String())
	}

	return values, xmlStrings, nil
}

// `dateString` should be of the form "YYYY/YYYY", where the left "YYYY" is the
// start date and the right "YYYY" is the end date.
func isDateInRange(dateString string, dateRange DateRange) bool {
	dateParts := getDateParts(dateString)

	startDateInt, err := strconv.Atoi(dateParts.Start)
	if err != nil {
		return false
	}

	endDateInt, err := strconv.Atoi(dateParts.End)
	if err != nil {
		return false
	}

	return (startDateInt >= dateRange.StartDate && startDateInt <= dateRange.EndDate) ||
		(endDateInt >= dateRange.StartDate && endDateInt <= dateRange.EndDate)
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

func stringifyStartElementToken(token xml.StartElement) string {
	startTag := "<" + token.Name.Local

	// Note that `token.Attr` appears to preserve the order of the attributes as
	// they appear in the XML.
	for _, attribute := range token.Attr {
		startTag += fmt.Sprintf(` %s="%s"`, attribute.Name.Local, attribute.Value)
	}

	startTag += ">"

	return startTag
}

// TODO: If we end up keeping this instead of using a 3rd-party package, make it
// general purpose by adding an `allowedHTMLTags` parameter instead of coupling
// to the package-level `allowedHTMLTags` var.
func stripTags(xmlString string) (string, error) {
	var strippedString string

	var startTagNames []string

	decoder := xml.NewDecoder(strings.NewReader(xmlString))

	for {
		token, err := decoder.Token()

		if err == io.EOF {
			break
		} else if err != nil {
			return xmlString, err
		}

		switch token := token.(type) {
		case xml.StartElement:
			if !slices.Contains(allowedHTMLTags, token.Name.Local) {
				continue
			}

			startTagNames = append(startTagNames, token.Name.Local)
			strippedString += stringifyStartElementToken(token)

		case xml.EndElement:
			if !slices.Contains(allowedHTMLTags, token.Name.Local) {
				continue
			}

			strippedString += fmt.Sprintf("</%s>", startTagNames[len(startTagNames)-1])
			startTagNames = startTagNames[:len(startTagNames)-1]

		case xml.CharData:
			strippedString += string(token)
		}
	}

	return strippedString, nil
}
