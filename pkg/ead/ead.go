package ead

import (
	"fmt"
	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/types"
	"go-ead-indexer/pkg/ead/collectiondoc"
	"go-ead-indexer/pkg/ead/component"
	"regexp"
)

// We need to set `xmlns=""` to get the xpath queries working.  See code comment
// in `New()` for more details.  `xmlns=""` is valid according to this post:
// https://stackoverflow.com/questions/1587891/is-xmlns-a-valid-xml-namespace
var namespaceRegexp = regexp.MustCompile(`<((?s)\s*)ead((?s).*)xmlns="(?U).*"`)

type EAD struct {
	CollectionDoc        collectiondoc.CollectionDoc `json:"collection_doc"`
	Components           *[]component.Component      `json:"components"`
	ModifiedFileContents string                      `json:"modified_file_contents"`
	OriginalFileContents string                      `json:"original_file_contents"`
}

// Note that the repository code historically is taken from the name of the
// EAD file's parent directory, not from the anything in the contents of the file
// itself.  For now we are keeping file handling out of this package, so it is
// up to the client pass in the repository code.
func New(repositoryCode string, eadXML string) (EAD, error) {
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

	ead.CollectionDoc, err = collectiondoc.MakeCollectionDoc(repositoryCode, rootNode)
	if err != nil {
		return ead, err
	}

	collectionDocParts := component.ComponentCollectionDocParts{
		// TODO: DLFA-238
		// We collect this Author information but don't include it in `SolrAddMessage`.
		// See: "Solr field `author` in Component Solr doc is never populated"
		// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=8577864&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-8577864
		// After passing DLFA-201 add `author` to `SolrAddMessage`.
		Author:           ead.CollectionDoc.Parts.Author.Values,
		Collection:       ead.CollectionDoc.Parts.Collection.Values[0],
		CollectionUnitID: ead.CollectionDoc.Parts.UnitID.Values[0],
		RepositoryCode:   repositoryCode,
	}
	ead.Components, err = component.MakeComponents(collectionDocParts, rootNode)
	if err != nil {
		return ead, err
	}

	return ead, nil
}

func MakeXMLDoc(eadXML string) (types.Document, error) {
	xmlParser := parser.New()
	xmlDoc, err := xmlParser.ParseString(eadXML)
	if err != nil {
		return xmlDoc, err
	}

	return xmlDoc, nil
}
