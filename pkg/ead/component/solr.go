package component

import (
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
	ID string `xml:"id"`
}

func (component *Component) setSolrAddMessage() {

}

func (solrAddMessage *SolrAddMessage) String() string {
	return "test"
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
