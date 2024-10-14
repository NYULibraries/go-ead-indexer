package ead

import (
	"fmt"
	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/types"
	"regexp"
)

// We need to set `xmlns=""` to get the xpath queries working.  See code comment
// in `New()` for more details.  `xmlns=""` is valid according to this post:
// https://stackoverflow.com/questions/1587891/is-xmlns-a-valid-xml-namespace
var namespaceRegexp = regexp.MustCompile(`<((?s)\s*)ead((?s).*)xmlns="(?U).*"`)

type EAD struct {
	Collection           Collection
	Components           *[]Component
	ModifiedFileContents string
	OriginalFileContents string
}

func New(eadXML string) (EAD, error) {
	ead := EAD{}

	// XPath queries fail if we don't set the namespace to empty string.
	// Excepting the `xlink` prefix, the tags in the EAD files don't seem to use
	// namespace prefixes much, and the XPath queries we need for this indexer
	// don't use prefixes at all.  Some brief experimentation suggests that if
	// we don't blank out the `xmlns` attribute, we would have to register a
	// namespace and add the prefix to all tag names in all the xpath queries.
	// We do a string replace instead of using `ctx.RegisterNS("", "")`, because
	// that call fails with the error "cannot register namespace".
	//
	// Note: v1 indexer removed all namespace stuff including prefixes:
	// https://github.com/awead/solr_ead/blob/v0.7.5/lib/solr_ead/om_behaviors.rb#L24
	// There doesn't appear to be a way to do this using `lestrrat-go/libxml2`.
	// There is one method `Element.SetNamespace()` which appears to enable
	// setting the namespace URL and prefix to empty strings, but the methods we
	// use for parsing and traversing/querying the DOM all seem to return `Node`
	// objects, not `Element` objects, and `Node` does not have a `SetNamespace()`
	// method.
	// A quick online search didn't turn up any easy to implement solutions for
	// removing namespace stuff from all nodes using the standard library
	// `encoding/xml` package.
	matchGroups := namespaceRegexp.FindStringSubmatch(eadXML)
	newString := fmt.Sprintf(`<%sead%sxmlns=""`, matchGroups[1], matchGroups[2])
	modifiedEADXML := namespaceRegexp.ReplaceAllString(eadXML, newString)

	ead.OriginalFileContents = eadXML
	ead.ModifiedFileContents = modifiedEADXML

	xmlDoc, err := MakeXMLDoc(modifiedEADXML)
	defer xmlDoc.Free()
	if err != nil {
		return ead, err
	}

	rootNode, err := xmlDoc.DocumentElement()
	if err != nil {
		return ead, err
	}

	ead.Collection, err = MakeCollection(rootNode)
	if err != nil {
		return ead, err
	}

	ead.Components, err = MakeComponents(rootNode)
	if err != nil {
		return ead, err
	}

	return ead, nil
}

func MakeCollection(node types.Node) (Collection, error) {
	newCollection := Collection{
		SolrAddMessage: "",
		XPathParts:     CollectionXPathParts{},
	}

	err := newCollection.populateXPathParts(node)
	if err != nil {
		return newCollection, err
	}

	return newCollection, nil
}

func MakeComponents(node types.Node) (*[]Component, error) {
	// TODO: remove this fake data
	return &[]Component{
		{
			ID:             "mos_2021additional-daos",
			SolrAddMessage: "",
		},
		{
			ID:             "mos_2021dao1",
			SolrAddMessage: "",
		},
		{
			ID:             "mos_2021dao1",
			SolrAddMessage: "",
		},
	}, nil
}

func MakeXMLDoc(eadXML string) (types.Document, error) {
	xmlParser := parser.New()
	xmlDoc, err := xmlParser.ParseString(eadXML)
	if err != nil {
		return xmlDoc, err
	}

	return xmlDoc, nil
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
