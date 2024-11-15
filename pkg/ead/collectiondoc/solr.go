package collectiondoc

import (
	"fmt"
	eadUtil "go-ead-indexer/pkg/ead/util"
	"go-ead-indexer/pkg/util"
	"reflect"
	"strconv"
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

// TODO DLFA-238:
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
	Author_teim            []string `xml:"author_teim"`
	Author_ssm             []string `xml:"author_ssm"`
	UnitTitle_teim         []string `xml:"unittitle_teim"`
	UnitTitle_ssm          []string `xml:"unittitle_ssm"`
	UnitID_teim            []string `xml:"unitid_teim"`
	UnitID_ssm             []string `xml:"unitid_ssm"`
	Language_ssm           string   `xml:"language_ssm"`
	Language_sim           string   `xml:"language_sim"`
	Abstract_teim          []string `xml:"abstract_teim"`
	Abstract_ssm           []string `xml:"abstract_ssm"`
	Creator_teim           []string `xml:"creator_teim"`
	Creator_ssm            []string `xml:"creator_ssm"`
	UnitDateNormal_ssm     []string `xml:"unitdate_normal_ssm"`
	UnitDateNormal_teim    []string `xml:"unitdate_normal_teim"`
	UnitDateNormal_sim     []string `xml:"unitdate_normal_sim"`
	UnitDateBulk_teim      []string `xml:"unitdate_bulk_teim"`
	UnitDateInclusive_teim []string `xml:"unitdate_inclusive_teim"`
	ScopeContent_teim      []string `xml:"scopecontent_teim"`
	BiogHist_teim          []string `xml:"bioghist_teim"`
	AcqInfo_teim           []string `xml:"acqinfo_teim"`
	CustodHist_teim        []string `xml:"custodhist_teim"`
	Appraisal_teim         []string `xml:"appraisal_teim"`
	PhysTech_teim          []string `xml:"phystech_teim"`
	ChronList_teim         []string `xml:"chronlist_teim"`
	CorpName_teim          []string `xml:"corpname_teim"`
	CorpName_ssm           []string `xml:"corpname_ssm"`
	FamName_teim           []string `xml:"famname_teim"`
	FamName_ssm            []string `xml:"famname_ssm"`
	Function_teim          []string `xml:"function_teim"`
	Function_ssm           []string `xml:"function_ssm"`
	GenreForm_teim         []string `xml:"genreform_teim"`
	GenreForm_ssm          []string `xml:"genreform_ssm"`
	GeogName_teim          []string `xml:"geogname_teim"`
	GeogName_ssm           []string `xml:"geogname_ssm"`
	Name_teim              []string `xml:"name_teim"`
	Name_ssm               []string `xml:"name_ssm"`
	Occupation_teim        []string `xml:"occupation_teim"`
	Occupation_ssm         []string `xml:"occupation_ssm"`
	PersName_teim          []string `xml:"persname_teim"`
	PersName_ssm           []string `xml:"persname_ssm"`
	Subject_teim           []string `xml:"subject_teim"`
	Subject_ssm            []string `xml:"subject_ssm"`
	Title_teim             []string `xml:"title_teim"`
	Title_ssm              []string `xml:"title_ssm"`
	Collection_sim         []string `xml:"collection_sim"`
	Collection_ssm         []string `xml:"collection_ssm"`
	Collection_teim        []string `xml:"collection_teim"`
	ID                     string   `xml:"id"`
	EAD_ssi                string   `xml:"ead_ssi"`
	Repository_ssi         string   `xml:"repository_ssi"`
	Repository_sim         string   `xml:"repository_sim"`
	Repository_ssm         string   `xml:"repository_ssm"`
	Format_sim             string   `xml:"format_sim"`
	Format_ssm             string   `xml:"format_ssm"`
	Format_ii              string   `xml:"format_ii"`
	Creator_sim            []string `xml:"creator_sim"`
	Name_sim               []string `xml:"name_sim"`
	Place_sim              []string `xml:"place_sim"`
	Subject_sim            []string `xml:"subject_sim"`
	DAO_sim                string   `xml:"dao_sim"`
	MaterialType_sim       []string `xml:"material_type_sim"`
	MaterialType_ssm       []string `xml:"material_type_ssm"`
	Heading_ssm            []string `xml:"heading_ssm"`
	UnitDateStart_sim      []string `xml:"unitdate_start_sim"`
	UnitDateStart_ssm      []string `xml:"unitdate_start_ssm"`
	UnitDateStart_si       string   `xml:"unitdate_start_si"`
	UnitDateEnd_sim        []string `xml:"unitdate_end_sim"`
	UnitDateEnd_ssm        []string `xml:"unitdate_end_ssm"`
	UnitDateEnd_si         string   `xml:"unitdate_end_si"`
	UnitDate_ssm           []string `xml:"unitdate_ssm"`
	DateRange_sim          []string `xml:"date_range_sim"`
	// Currently not in Omega golden file, so don't know where to place them.
	Note_ssm      []string `xml:"note_ssm"`
	Note_teim     []string `xml:"note_teim"`
	UnitDate_teim []string `xml:"unitdate_teim"`
}

func (collectionDoc *CollectionDoc) setSolrAddMessage() {
	docElement := &collectionDoc.SolrAddMessage.Add.Doc

	docElement.Abstract_ssm = append(docElement.Abstract_ssm, collectionDoc.Parts.Abstract.Values...)
	docElement.Abstract_teim = append(docElement.Abstract_teim, collectionDoc.Parts.Abstract.Values...)

	docElement.AcqInfo_teim = append(docElement.AcqInfo_teim, collectionDoc.Parts.AcqInfo.Values...)

	docElement.Appraisal_teim = append(docElement.Appraisal_teim, collectionDoc.Parts.Appraisal.Values...)

	docElement.Author_ssm = append(docElement.Author_ssm, collectionDoc.Parts.Author.Values...)
	docElement.Author_teim = append(docElement.Author_teim, collectionDoc.Parts.Author.Values...)

	docElement.BiogHist_teim = append(docElement.BiogHist_teim, collectionDoc.Parts.BiogHist.Values...)

	docElement.ChronList_teim = append(docElement.ChronList_teim, collectionDoc.Parts.ChronListComplex.Values...)

	docElement.Collection_sim = append(docElement.Collection_sim, collectionDoc.Parts.UnitTitle.Values...)
	docElement.Collection_ssm = append(docElement.Collection_ssm, collectionDoc.Parts.UnitTitle.Values...)
	docElement.Collection_teim = append(docElement.Collection_teim, collectionDoc.Parts.UnitTitle.Values...)

	docElement.CorpName_ssm = append(docElement.CorpName_ssm, collectionDoc.Parts.CorpNameNotInDSC.Values...)
	docElement.CorpName_teim = append(docElement.CorpName_teim, collectionDoc.Parts.CorpNameNotInDSC.Values...)

	// See 2nd `Creator_ssm` append below.
	docElement.Creator_sim = append(docElement.Creator_sim,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.CreatorComplex.Values)...)
	docElement.Creator_ssm = append(docElement.Creator_ssm, collectionDoc.Parts.Creator.Values...)
	// TODO: is this duplication done in v1 indexer a bug that needs to be added
	// to https://jira.nyu.edu/browse/DLFA-211?
	docElement.Creator_ssm = append(docElement.Creator_ssm,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.CreatorComplex.Values)...)
	docElement.Creator_teim = append(docElement.Creator_teim, collectionDoc.Parts.Creator.Values...)

	docElement.CustodHist_teim = append(docElement.CustodHist_teim, collectionDoc.Parts.CustodHist.Values...)

	if len(collectionDoc.Parts.OnlineAccess.Values) > 0 {
		docElement.DAO_sim = collectionDoc.Parts.OnlineAccess.Values[0]
	}

	docElement.DateRange_sim = append(docElement.DateRange_sim,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.DateRange.Values)...)

	docElement.EAD_ssi = collectionDoc.Parts.EADID.Values[0]

	docElement.FamName_ssm = append(docElement.FamName_ssm, collectionDoc.Parts.FamNameNotInDSC.Values...)
	docElement.FamName_teim = append(docElement.FamName_teim, collectionDoc.Parts.FamNameNotInDSC.Values...)

	docElement.Format_ii = strconv.Itoa(collectionDoc.Parts.FormatForSort)
	docElement.Format_sim = collectionDoc.Parts.FormatForDisplay
	docElement.Format_ssm = collectionDoc.Parts.FormatForDisplay

	docElement.Function_ssm = append(docElement.Function_ssm, collectionDoc.Parts.Function.Values...)
	docElement.Function_teim = append(docElement.Function_teim, collectionDoc.Parts.Function.Values...)

	docElement.GenreForm_ssm = append(docElement.GenreForm_ssm, collectionDoc.Parts.GenreFormNotInDSC.Values...)
	docElement.GenreForm_teim = append(docElement.GenreForm_teim, collectionDoc.Parts.GenreFormNotInDSC.Values...)

	docElement.GeogName_ssm = append(docElement.GeogName_ssm, collectionDoc.Parts.GeogNameNotInDSC.Values...)
	docElement.GeogName_teim = append(docElement.GeogName_teim, collectionDoc.Parts.GeogNameNotInDSC.Values...)

	docElement.Heading_ssm = append(docElement.Heading_ssm,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.Heading.Values)...)

	docElement.ID = collectionDoc.Parts.EADID.Values[0]

	if len(collectionDoc.Parts.Language.Values) > 0 {
		docElement.Language_sim = collectionDoc.Parts.Language.Values[0]
		docElement.Language_ssm = collectionDoc.Parts.Language.Values[0]
	}

	docElement.MaterialType_sim = append(docElement.MaterialType_sim,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.MaterialType.Values)...)
	docElement.MaterialType_ssm = append(docElement.MaterialType_ssm,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.MaterialType.Values)...)

	docElement.Name_sim = append(docElement.Name_sim,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.Name.Values)...)
	docElement.Name_ssm = append(docElement.Name_ssm, collectionDoc.Parts.NameNotInDSC.Values...)
	docElement.Name_teim = append(docElement.Name_teim, collectionDoc.Parts.NameNotInDSC.Values...)
	docElement.Name_teim = append(docElement.Name_teim,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.Name.Values)...)

	docElement.Note_ssm = append(docElement.Note_ssm, collectionDoc.Parts.NoteNotInDSC.Values...)
	docElement.Note_teim = append(docElement.Note_teim, collectionDoc.Parts.NoteNotInDSC.Values...)

	docElement.Occupation_ssm = append(docElement.Occupation_ssm, collectionDoc.Parts.OccupationNotInDSC.Values...)
	docElement.Occupation_teim = append(docElement.Occupation_teim, collectionDoc.Parts.OccupationNotInDSC.Values...)

	docElement.PersName_ssm = append(docElement.PersName_ssm, collectionDoc.Parts.PersNameNotInDSC.Values...)
	docElement.PersName_teim = append(docElement.PersName_teim, collectionDoc.Parts.PersNameNotInDSC.Values...)

	docElement.PhysTech_teim = append(docElement.PhysTech_teim, collectionDoc.Parts.Phystech.Values...)

	docElement.Place_sim = append(docElement.Place_sim,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.Place.Values)...)

	docElement.Repository_sim = collectionDoc.Parts.RepositoryCode.Values[0]
	docElement.Repository_ssi = collectionDoc.Parts.RepositoryCode.Values[0]
	docElement.Repository_ssm = collectionDoc.Parts.RepositoryCode.Values[0]

	docElement.ScopeContent_teim = append(docElement.ScopeContent_teim, collectionDoc.Parts.ScopeContent.Values...)

	// See 2nd `Subject_teim` append below.
	docElement.Subject_sim = append(docElement.Subject_sim,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.SubjectForFacets.Values)...)
	docElement.Subject_ssm = append(docElement.Subject_ssm, collectionDoc.Parts.SubjectNotInDSC.Values...)
	docElement.Subject_teim = append(docElement.Subject_teim, collectionDoc.Parts.SubjectNotInDSC.Values...)
	// TODO: is this duplication done in v1 indexer a bug that needs to be added
	// to https://jira.nyu.edu/browse/DLFA-211?
	docElement.Subject_teim = append(docElement.Subject_teim,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.SubjectForFacets.Values)...)

	docElement.Title_ssm = append(docElement.Title_ssm, collectionDoc.Parts.TitleNotInDSC.Values...)
	docElement.Title_teim = append(docElement.Title_teim, collectionDoc.Parts.TitleNotInDSC.Values...)

	docElement.UnitDateBulk_teim = append(docElement.UnitDateBulk_teim, collectionDoc.Parts.UnitDateBulk.Values...)

	if len(collectionDoc.Parts.UnitDateEnd.Values) > 0 {
		docElement.UnitDateEnd_si = collectionDoc.Parts.UnitDateEnd.Values[len(collectionDoc.Parts.UnitDateEnd.Values)-1]
	}
	docElement.UnitDateEnd_sim = append(docElement.UnitDateEnd_sim,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.UnitDateEnd.Values)...)
	docElement.UnitDateEnd_ssm = append(docElement.UnitDateEnd_ssm,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.UnitDateEnd.Values)...)

	docElement.UnitDateInclusive_teim = append(docElement.UnitDateInclusive_teim, collectionDoc.Parts.UnitDateInclusive.Values...)

	docElement.UnitDateNormal_sim = append(docElement.UnitDateNormal_sim, collectionDoc.Parts.UnitDateNormal.Values...)
	docElement.UnitDateNormal_ssm = append(docElement.UnitDateNormal_ssm, collectionDoc.Parts.UnitDateNormal.Values...)
	docElement.UnitDateNormal_teim = append(docElement.UnitDateNormal_teim, collectionDoc.Parts.UnitDateNormal.Values...)

	docElement.UnitDate_ssm = append(docElement.UnitDate_ssm,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.UnitDateDisplay.Values)...)
	docElement.UnitDate_teim = append(docElement.UnitDate_teim, collectionDoc.Parts.UnitDateNoTypeAttribute.Values...)

	if len(collectionDoc.Parts.UnitDateStart.Values) > 0 {
		docElement.UnitDateStart_si = collectionDoc.Parts.UnitDateStart.Values[len(collectionDoc.Parts.UnitDateStart.Values)-1]
	}
	docElement.UnitDateStart_sim = append(docElement.UnitDateStart_sim,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.UnitDateStart.Values)...)
	docElement.UnitDateStart_ssm = append(docElement.UnitDateStart_ssm,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.UnitDateStart.Values)...)

	docElement.UnitID_ssm = append(docElement.UnitID_ssm, collectionDoc.Parts.UnitID.Values...)
	docElement.UnitID_teim = append(docElement.UnitID_teim, collectionDoc.Parts.UnitID.Values...)

	docElement.UnitTitle_ssm = append(docElement.UnitTitle_ssm,
		util.CompactStringSlicePreserveOrder(collectionDoc.Parts.UnitTitleHTML.Values)...)
	docElement.UnitTitle_teim = append(docElement.UnitTitle_teim, collectionDoc.Parts.UnitTitle.Values...)
}

// TODO DLFA-238:
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
				if isNonEmptyString(fieldValue) {
					fieldsInV1IndexerInsertionOrder = append(fieldsInV1IndexerInsertionOrder,
						makeSolrAddMessageFieldElementString(fieldName, fieldValue))
				}
			}
		} else if fieldTypeKind == reflect.String {
			fieldValue := field.String()
			if isNonEmptyString(fieldValue) {
				fieldsInV1IndexerInsertionOrder = append(fieldsInV1IndexerInsertionOrder,
					makeSolrAddMessageFieldElementString(fieldName, fieldValue))
			}
		} else {
			// Should never get here!
			panic("Unrecognized `reflect.Type.Kind`: " + fieldTypeKind.String())
		}
	}

	return fieldsInV1IndexerInsertionOrder
}

// Based on: https://stackoverflow.com/questions/18594330/what-is-the-best-way-to-test-for-an-empty-string-in-go
func isNonEmptyString(value string) bool {
	return len(strings.TrimSpace(value)) > 0
}

func makeSolrAddMessageFieldElementString(fieldName string, fieldValue string) string {
	escapedFieldValue := eadUtil.EscapeSolrFieldString(fieldValue)

	return fmt.Sprintf(`    <field name="%s">%s</field>`, fieldName, escapedFieldValue)
}
