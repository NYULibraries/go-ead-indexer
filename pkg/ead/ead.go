package ead

import (
	"errors"
	"fmt"
	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/types"
	"github.com/nyulibraries/go-ead-indexer/pkg/ead/collectiondoc"
	"github.com/nyulibraries/go-ead-indexer/pkg/ead/component"
	"github.com/nyulibraries/go-ead-indexer/pkg/ead/eadutil"
	"regexp"
)

type EAD struct {
	CollectionDoc        collectiondoc.CollectionDoc `json:"collection_doc"`
	Components           *[]component.Component      `json:"components"`
	ModifiedFileContents string                      `json:"modified_file_contents"`
	OriginalFileContents string                      `json:"original_file_contents"`
}

const errorFormatStringInvalidEADID = `"%s" is not a valid EAD ID`
const errorNoEADTagWithExpectedStructureFound = "No <ead> tag with the expected structure was found"

// This must be to the number of match groups in the regexp below.
const numMatchGroupsInNamespaceRegexp = 3

// We need to set `xmlns=""` to get the xpath queries working.  See code comment
// in `New()` for more details.  `xmlns=""` is valid according to this post:
// https://stackoverflow.com/questions/1587891/is-xmlns-a-valid-xml-namespace
var namespaceRegexp = regexp.MustCompile(`<((?s)\s*)ead((?s).*)xmlns="(?U).*"`)

// We don't have an official repository code format, but there is a comprehensive
// list of repository codes:
// https://jira.nyu.edu/browse/FADESIGN-65
var validRepositoryCodeRegex = regexp.MustCompile(`^[a-z]+$`)

// Note that the repository code historically is taken from the name of the
// EAD file's parent directory, not from the anything in the contents of the file
// itself.  For now we are keeping file handling out of this package, so it is
// up to the client to pass in the repository code.
func New(repositoryCode string, eadXML string) (EAD, error) {
	ead := EAD{}

	if !validRepositoryCodeRegex.MatchString(repositoryCode) {
		return ead, errors.New(fmt.Sprintf(`Invalid repository code: "%s"`,
			repositoryCode))
	}

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
	if len(matchGroups) < numMatchGroupsInNamespaceRegexp {
		return ead, errors.New(errorNoEADTagWithExpectedStructureFound)
	}
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

	// At the moment values are in a slice, as is the case with most Parts.
	var eadID = ead.CollectionDoc.Parts.EADID.Values[0]
	if !eadutil.IsValidEADID(eadID) {
		return ead, errors.New(fmt.Sprintf(errorFormatStringInvalidEADID, eadID))
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
