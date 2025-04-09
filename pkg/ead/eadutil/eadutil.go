package eadutil

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/lestrrat-go/libxml2/dom"
	"github.com/lestrrat-go/libxml2/types"
	"github.com/nyulibraries/dlts-finding-aids-ead-go-packages/ead/validate"
	languageLib "github.com/nyulibraries/go-ead-indexer/pkg/language"
	"github.com/nyulibraries/go-ead-indexer/pkg/sanitize"
	"github.com/nyulibraries/go-ead-indexer/pkg/util"
	"html"
	"io"
	"maps"
	"path/filepath"
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

// TODO: DLFA-238
// Remove these `consts` for left- and right- padding for matching v1
// indexer bug behavior described here:
// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10849506&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10849506
const daoDescriptionParagraphLeftPadString = "\n          "
const daoDescriptionParagraphRightPadString = "\n        "
const unitTitleLeftPadString = "\n      "
const unitTitleRightPadString = "\n    "

const eadLineBreakTag = "<lb/>"

const undated = "undated & other"

var allowedConvertedEADToHTMLTags = util.CompactStringSlicePreserveOrder(
	slices.Collect(maps.Values(eadTagRenderAttributeToHTMLTagName)))

var datePartsRegexp = regexp.MustCompile(`^\s*(\d{4})\/(\d{4})\s*$`)

// We assume lowercase entity references only because Go escaping uses lowercase.
var multipleEscapedAmpersandRegexp = regexp.MustCompile(`(?:&amp;)+amp;`)

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

var marcSubfieldDemarcator = regexp.MustCompile(`\|\w{1}\s`)

// Go \s metachar is [\t\n\f\r ], and does not include NBSP.
// Source: https://pkg.go.dev/regexp/syntax
var multipleConsecutiveWhitespace = regexp.MustCompile(`[\s ]{2}\s*`)
var leadingWhitespaceInFieldContent = regexp.MustCompile(`^[\s ]+`)
var trailingWhitespaceInFieldContent = regexp.MustCompile(`[\s ]+$`)

// These are not perfect regexps for open and close XML tags, but they are fine
// for our constrained use cases.
var closeTagRegExp = regexp.MustCompile("</[^>]+>$")
var openTagRegExp = regexp.MustCompile("^<[^>]+>")

var validEADIDRegexp = regexp.MustCompile(validate.ValidEADIDRegexpString)

func ConvertEADToHTML(eadString string) (string, error) {
	htmlString, err := convertEADTagsWithRenderAttributesToHTML(eadString)
	if err != nil {
		return htmlString, err
	}

	return sanitize.Clean(htmlString), nil
}

func ConvertToFacetSlice(rawSlice []string) []string {
	return util.CompactStringSlicePreserveOrder(
		replaceMARCSubfieldDemarcatorsInSlice(rawSlice))
}

func EADPathToEADID(path string) (string, error) {
	eadID := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if !IsValidEADID(eadID) {
		return "", fmt.Errorf("invalid EADID: %s", eadID)
	}
	return eadID, nil
}

func GetDateParts(dateString string) DateParts {
	dateParts := DateParts{}

	matches := datePartsRegexp.FindStringSubmatch(dateString)

	if len(matches) == 3 {
		dateParts.Start = matches[1]
		dateParts.End = matches[2]
	}

	return dateParts
}

func GetDateRange(unitDates []string) []string {
	dateRange := []string{}

	// Add `dateRangeCentury` display dates for which at least one date falls within
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

func GetFirstNode(query string, node types.Node) (types.Node, error) {
	xpathResult, err := node.Find(query)
	if err != nil {
		return nil, err
	}
	defer xpathResult.Free()

	nodeList := xpathResult.NodeList()

	if len(nodeList) > 0 {
		return nodeList[0], nil
	} else {
		return nil, nil
	}
}

func GetLanguage(langCodes []string) ([]string, []error) {
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

func GetNodeList(query string, node types.Node) (types.NodeList, error) {
	xpathResult, err := node.Find(query)
	if err != nil {
		return nil, err
	}
	defer xpathResult.Free()

	return xpathResult.NodeList(), nil
}

func GetUnitDateDisplay(unitDateNoTypeAttribute []string, unitDateInclusive []string,
	unitDateBulk []string) string {
	partsUnitDateDisplay := []string{}
	if len(unitDateNoTypeAttribute) > 0 {
		partsUnitDateDisplay = unitDateNoTypeAttribute
	} else if len(unitDateInclusive) == 0 && len(unitDateBulk) == 0 {
		// Do nothing
	} else {
		partsUnitDateDisplay = append(partsUnitDateDisplay, "Inclusive,")
		if len(unitDateInclusive) > 0 {
			partsUnitDateDisplay = append(partsUnitDateDisplay, unitDateInclusive...)
		}
		if len(unitDateBulk) > 0 {
			partsUnitDateDisplay = append(partsUnitDateDisplay, ";")
			partsUnitDateDisplay = append(partsUnitDateDisplay, unitDateBulk...)
		}
	}

	return strings.Join(partsUnitDateDisplay, " ")
}

func GetNodeValuesAndXMLStrings(query string, node types.Node) ([]string, []string, error) {
	var values []string
	var xmlStrings []string

	xpathResult, err := node.Find(query)
	if err != nil {
		return nil, nil, err
	}
	defer xpathResult.Free()

	for _, resultNode := range xpathResult.NodeList() {
		xmlString := resultNode.String()

		var value string
		if resultNode.NodeType() == dom.ElementNode {
			// We were originally using Node.NodeValue() for `values` slice, but
			// it caused problems with element values containing <lb/> tags.
			// We basically want everything we got from Node.NodeValue() but
			// with <lb/> tags replaced with whitespace so that the text on
			// either side of the <lb/> tags don't get fused together.
			// Note that there is downstream whitespace processing that might alter
			// the whitespace replacement choice we make here, but at this stage
			// of processing we just do what seems most natural.
			value, err = parseNodeValue(xmlString)
			if err != nil {
				return values, xmlStrings, err
			}
		} else {
			value = resultNode.NodeValue()
		}

		values = append(values, value)
		xmlStrings = append(xmlStrings, xmlString)
	}

	return values, xmlStrings, nil
}

// This is based on `validateEADID()`:
// https://github.com/NYULibraries/dlts-finding-aids-ead-go-packages/blob/7baee7dfde24a01422ec8e6470fdc8a76d84b3fb/ead/validate/validate.go#L205-L244
func IsValidEADID(eadid string) bool {
	if len(eadid) > validate.MAXIMUM_EADID_LENGTH {
		return false
	}

	return validEADIDRegexp.MatchString(eadid)
}

func MakeSolrAddMessageFieldElementString(fieldName string, fieldValue string) string {
	massagedValue := fieldValue

	massagedValue = html.EscapeString(fieldValue)

	return fmt.Sprintf(`<field name="%s">%s</field>`, fieldName, massagedValue)
}

func MakeTitleHTML(unitTitle string) (string, error) {
	converted, err := ConvertEADToHTML(unitTitle)
	if err != nil {
		return converted, err
	}

	titleHTML, err := StripNonEADToHTMLTags(converted)
	if err != nil {
		return titleHTML, err
	}

	return titleHTML, nil
}

// `node.TextContent()` might contain unescaped characters that would be dangerous
// for XML processing, like "&", ">", or "<".
func ParseEscapedNodeTextContent(node types.Node) (string, error) {
	textContentBytes := []byte(node.TextContent())

	escapedBuffer := new(strings.Builder)
	if err := xml.EscapeText(escapedBuffer, textContentBytes); err != nil {
		return string(textContentBytes), err
	}

	return escapedBuffer.String(), nil
}

// Using `strings.ReplaceAll` instead of full parsing of the XML should be safe
// for `SolrAddMessage` XML strings, which are valid XML and therefore cannot have
// unescaped "<" and ">" characters in text nodes or attribute values.
func PrettifySolrAddMessageXML(xml string) string {
	var xml1 = strings.ReplaceAll(xml, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>",
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	var xml2 = strings.ReplaceAll(xml1, "<add>", "<add>\n")
	var xml3 = strings.ReplaceAll(xml2, "<doc>", "  <doc>\n")
	var xml4 = strings.ReplaceAll(xml3, "<field name=", "    <field name=")
	var xml5 = strings.ReplaceAll(xml4, "</field>", "</field>\n")
	var xml6 = strings.ReplaceAll(xml5, "</doc>", "  </doc>\n")
	var xml7 = strings.ReplaceAll(xml6, "</add>", "</add>\n")

	return xml7
}

// Note that this function only removes child nodes, it does not recursively
// remove all descendant notes which match `elementName`.
//
// This function mutates the `node` arg.  The first version of this function made
// a copy and returned it after removing the appropriate child nodes.  This was
// to prevent surprising the caller with an unwanted mutation, because even though
// the param is `types.Node` and not `*types.Node`, the mutations performed here
// are permanent.
//
// Returning a copy ended up not being as good a choice as it originally seemed.
// There were two undesirable side effects:
//
//  1. The copying process removed this attribute from the root <c> node:
//     `xmlns:xlink="http://www.w3.org/1999/xlink"`.  For this project this is
//     most likely harmless, and in fact it is some ways convenient, it's an
//     unexpected change that can't be opted out of.
//
//  2. The `Node` returned by `node.Copy()` had non of the original node's parent
//     node data, causing all `.ParentNode()` calls to the returned, modified copy
//     to fail.  An attempt was made to attach the modified `node` copy to either
//     the parent node of `node` or a copy of the parent node, but in cases where
//     `node` had no parent, as in the unit test, the return error was cryptic:
//     "unknown node: 9".  This could be from either `go-libxml2` or the `libxml2`
//     C library, but in any case the returned error is a generic string error
//     and not a typed error.  Doing a string match on "unknown node" to differentiate
//     between cases where the error is simply due to `node` being a root node
//     with no parent and an error caused by an actual problem seem too brittle.
//
// So for now, we mutate the arg.  The caller should make a defensive copy to avoid
// any of the risks associated with mutation mentioned in the first list in this
// comment.
func RemoveChildNodesMatchingName(node types.Node, elementName string) error {
	if node == nil {
		return errors.New("`node` arg is `nil`")
	}

	childNodes, err := node.ChildNodes()
	if err != nil {
		return err
	}

	for _, childNode := range childNodes {
		if childNode != nil {
			if childNode.NodeName() == elementName {
				err = node.RemoveChild(childNode)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// This is not a particularly robust solution, but for our constrained use case
// it's fine.  This function is used to remove the open and close tags *only* from
// an XML string returned by `Node.String()`.  It's part of a tricky EAD to HTML
// conversion process wherein certain tags are allowed but converted and other
// tags are not allowed.  It may at first seem redundant given we also have the
// `StripTags()` function below, but these functions run at different times and
// do different things.
func StripOpenAndCloseTags(xmlString string) string {
	return openTagRegExp.ReplaceAllString(
		closeTagRegExp.ReplaceAllString(xmlString, ""), "")
}

func StripNonEADToHTMLTags(xmlString string) (string, error) {
	return StripTags(xmlString, &allowedConvertedEADToHTMLTags)
}

func StripTags(xmlString string, allowedTags *[]string) (string, error) {
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
			if allowedTags == nil || !slices.Contains(*allowedTags, token.Name.Local) {
				continue
			}

			startTagNames = append(startTagNames, token.Name.Local)
			strippedString += stringifyStartElementToken(token)

		case xml.EndElement:
			if allowedTags == nil || !slices.Contains(*allowedTags, token.Name.Local) {
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
			// The XML has been unescaped, need to re-escape it since it needs
			// to go back into XML again.
			buffer := new(strings.Builder)
			if err := xml.EscapeText(buffer, token); err != nil {
				return htmlString, err
			}

			htmlString += buffer.String()
		}
	}

	return htmlString, nil
}

func isDateInRange(dateString string, dateRange DateRange) bool {
	dateParts := GetDateParts(dateString)

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

func parseNodeValue(xmlString string) (string, error) {
	// We can't just strip <lb/> tags because many times the text on either side
	// of the tags have no intervening whitespace, and so simple removal would
	// cause the text on other side of the tags to be fused together.
	value := strings.ReplaceAll(xmlString, eadLineBreakTag, "\n")

	// All other tags must be removed.
	value, err := StripTags(value, nil)
	if err != nil {
		return value, err
	}

	return value, nil
}

// TODO: fix the bug we've intentionally preserved here -- for details, see:
// * https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10154897&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10154897
// * https://jira.nyu.edu/browse/DLFA-229?focusedCommentId=10153922&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10153922
func replaceMARCSubfieldDemarcators(str string) string {
	return marcSubfieldDemarcator.ReplaceAllString(str, "--")
}

func replaceMARCSubfieldDemarcatorsInSlice(stringSlice []string) []string {
	newSlice := []string{}
	for _, element := range stringSlice {
		newSlice = append(newSlice, replaceMARCSubfieldDemarcators(element))
	}

	return newSlice
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

// We assume lowercase entity references only because Go escaping uses lowercase.
func undoMultipleAmpersandEscaping(str string) string {
	return multipleEscapedAmpersandRegexp.ReplaceAllString(str, "&amp;")
}
