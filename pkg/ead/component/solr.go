package component

import (
	"fmt"
	"go-ead-indexer/pkg/ead/eadutil"
	"go-ead-indexer/pkg/util"
	"reflect"
	"strings"
)

// We are currently using `String()` and not marshaling, but for now we are
// structuring as if we are or might later be marshaling.
type SolrAddMessage struct {
	Add AddElement `xml:"add"`
}

type AddElement struct {
	Doc DocElement `xml:"doc"`
}

// TODO: DLFA-238
// This struct definition replicates the order in which the v1 indexer writes
// out the Solr field elements in the HTTP request to Solr.  We are generating
// the XML request body by using the `reflect` package to loop through the
// struct fields in the order they are defined here (at least that's how it
// seems in the current Go version).
// After we pass the DLFA-201 acceptance test, we need to implement the
// permanent `String()` or custom marshaling that will be free of the need to
// match v1 indexer's ordering, and restore the alphabetical ordering of the field
// definitions in this struct.
type DocElement struct {
	ID                    string   `xml:"id"`
	EAD_ssi               string   `xml:"ead_ssi"`
	Parent_ssi            string   `xml:"parent_ssi"`
	Parent_ssm            []string `xml:"parent_ssm"`
	ParentUnitTitles_ssm  []string `xml:"parent_unittitles_ssm"`
	ParentUnitTitles_teim []string `xml:"parent_unittitles_teim"`
	ComponentLevel_isim   string   `xml:"component_level_isim"`
	ComponentChildren_bsi string   `xml:"component_children_bsi"`
	Collection_teim       string   `xml:"collection_teim"`
	Collection_ssm        string   `xml:"collection_ssm"`
	Repository_ssi        string   `xml:"repository_ssi"`
	Repository_sim        string   `xml:"repository_sim"`
	Repository_ssm        string   `xml:"repository_ssm"`
}

func (component *Component) setSolrAddMessage() {
	docElement := &component.SolrAddMessage.Add.Doc

	docElement.ID = component.ID

	docElement.Collection_ssm = component.Parts.Collection
	docElement.Collection_teim = component.Parts.Collection

	docElement.ComponentChildren_bsi = component.Parts.ComponentChildren
	docElement.ComponentLevel_isim = component.Parts.ComponentLevel

	docElement.EAD_ssi = component.ID

	if component.Parts.ParentForSort != "" {
		docElement.Parent_ssi = component.Parts.ParentForSort
	}
	docElement.Parent_ssm = component.Parts.ParentForDisplay.Values

	docElement.ParentUnitTitles_ssm = component.Parts.AncestorUnitTitleList
	docElement.ParentUnitTitles_teim = component.Parts.AncestorUnitTitleList

	docElement.Repository_sim = component.Parts.RepositoryCode
	docElement.Repository_ssi = component.Parts.RepositoryCode
	docElement.Repository_ssm = component.Parts.RepositoryCode
}

// TODO: DLFA-238
// This replicates the order in which the v1 indexer writes out the Solr
// field elements in the HTTP request to Solr.  After we pass the DLFA-201
// acceptance test, we need to implement the permanent `String()` or custom
// marshaling that will be free of the need to match v1 indexer's ordering.
func (solrAddMessage SolrAddMessage) String() string {
	fields := getSolrFieldElementStringsInV1IndexerInsertionOrder(solrAddMessage)

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<add>
  <doc>
%s
  </doc>
</add>
`, strings.Join(fields, "\n"))
}

// TODO: DLFA-238
// This replicates the order in which the v1 indexer writes out the Solr
// field elements in the HTTP request to Solr.  After we pass the DLFA-201
// acceptance test, we need to implement the permanent `String()` or custom
// marshaling that will be free of the need to match v1 indexer's ordering.
// Note that this function is duplicated in the `component` package.  Normally
// we'd find a way to DRY this up (probably by using a `struct` param instead of
// the `CollectionDoc.SolrAddMessage` and `Component.SolrAddMessage` types, but
// since this is ephemeral, we just copy it.
func getSolrFieldElementStringsInV1IndexerInsertionOrder(solrAddMessage SolrAddMessage) []string {
	var fieldsInV1IndexerInsertionOrder []string

	docElementStructType := reflect.TypeOf(solrAddMessage.Add.Doc)
	docElementStructValue := reflect.ValueOf(solrAddMessage.Add.Doc)

	numFields := docElementStructValue.NumField()
	for i := 0; i < numFields; i++ {
		field := docElementStructValue.Field(i)
		fieldName := strings.Split(docElementStructType.Field(i).Tag.Get("xml"), ",")[0]
		fieldTypeKind := field.Type().Kind()
		if fieldTypeKind == reflect.Slice {
			for _, fieldValue := range field.Interface().([]string) {
				if util.IsNonEmptyString(fieldValue) {
					fieldsInV1IndexerInsertionOrder = append(fieldsInV1IndexerInsertionOrder,
						eadutil.MakeSolrAddMessageFieldElementString(fieldName, fieldValue))
				}
			}
		} else if fieldTypeKind == reflect.String {
			fieldValue := field.String()
			if util.IsNonEmptyString(fieldValue) {
				fieldsInV1IndexerInsertionOrder = append(fieldsInV1IndexerInsertionOrder,
					eadutil.MakeSolrAddMessageFieldElementString(fieldName, fieldValue))
			}
		} else {
			// Should never get here!
			panic("Unrecognized `reflect.Type.Kind`: " + fieldTypeKind.String())
		}
	}

	return fieldsInV1IndexerInsertionOrder
}
