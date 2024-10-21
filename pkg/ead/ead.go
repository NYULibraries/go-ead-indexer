package ead

import (
	"fmt"
	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/types"
	"regexp"
)

// TODO DLFA-238: fix the bug we've intentionally preserved in MARC subfield demarcation
// replacement.  For details, see:
//
//   - https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10154897&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10154897
//   - https://jira.nyu.edu/browse/DLFA-229?focusedCommentId=10153922&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10153922
//
// This is the buggy regular expression which replicates the v1 indexer code here:
// https://github.com/NYULibraries/ead_indexer/blob/a367ab8cc791376f0d8a287cbcd5b6ee43d5c04f/lib/ead_indexer/behaviors.rb#L124
var marcSubfieldDemarcator = regexp.MustCompile(`\|\w{1}`)

// We need to set `xmlns=""` to get the xpath queries working.  See code comment
// in `New()` for more details.  `xmlns=""` is valid according to this post:
// https://stackoverflow.com/questions/1587891/is-xmlns-a-valid-xml-namespace
var namespaceRegexp = regexp.MustCompile(`<((?s)\s*)ead((?s).*)xmlns="(?U).*"`)

type EAD struct {
	CollectionDoc        CollectionDoc
	Components           *[]Component
	ModifiedFileContents string
	OriginalFileContents string
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

	ead.CollectionDoc, err = MakeCollectionDoc(repositoryCode, rootNode)
	if err != nil {
		return ead, err
	}

	ead.Components, err = MakeComponents(repositoryCode, rootNode)
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
