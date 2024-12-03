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
	ID                     string   `xml:"id"`
	EAD_ssi                string   `xml:"ead_ssi"`
	Parent_ssi             string   `xml:"parent_ssi"`
	Parent_ssm             []string `xml:"parent_ssm"`
	ParentUnitTitles_ssm   []string `xml:"parent_unittitles_ssm"`
	ParentUnitTitles_teim  []string `xml:"parent_unittitles_teim"`
	ComponentLevel_isim    string   `xml:"component_level_isim"`
	ComponentChildren_bsi  string   `xml:"component_children_bsi"`
	Collection_sim         string   `xml:"collection_sim"`
	Collection_ssm         string   `xml:"collection_ssm"`
	CollectionUnitID_ssm   string   `xml:"collection_unitid_ssm"`
	Level_sim              string   `xml:"level_sim"`
	UnitTitle_ssm          []string `xml:"unittitle_ssm"`
	UnitTitle_teim         []string `xml:"unittitle_teim"`
	UnitID_teim            []string `xml:"unitid_teim"`
	UnitID_ssm             []string `xml:"unitid_ssm"`
	Creator_teim           []string `xml:"creator_teim"`
	Creator_ssm            []string `xml:"creator_ssm"`
	UnitDateNormal_ssm     []string `xml:"unitdate_normal_ssm"`
	UnitDateNormal_teim    []string `xml:"unitdate_normal_teim"`
	UnitDateNormal_sim     []string `xml:"unitdate_normal_sim"`
	UnitDateInclusive_teim []string `xml:"unitdate_inclusive_teim"`
	ScopeContent_teim      []string `xml:"scopecontent_teim"`
	BiogHist_teim          []string `xml:"bioghist_teim"`
	Address_teim           []string `xml:"address_teim"`
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
	Note_teim              []string `xml:"note_teim"`
	Note_ssm               []string `xml:"note_ssm"`
	DAO_teim               []string `xml:"dao_teim"`
	DAO_ssm                []string `xml:"dao_ssm"`
	Ref_ssi                string   `xml:"ref_ssi"`
	Repository_ssi         string   `xml:"repository_ssi"`
	Repository_sim         string   `xml:"repository_sim"`
	Repository_ssm         string   `xml:"repository_ssm"`
	Format_sim             []string `xml:"format_sim"`
	Format_ssm             []string `xml:"format_ssm"`
	Location_ssm           []string `xml:"location_ssm"`
	Location_si            []string `xml:"location_si"`
	Creator_sim            []string `xml:"creator_sim"`
	Name_sim               []string `xml:"name_sim"`
	DAO_sim                []string `xml:"dao_sim"`
	Place_ssm              []string `xml:"place_ssm"`
	Place_sim              []string `xml:"place_sim"`
	Subject_sim            []string `xml:"subject_sim"`
	Collection_teim        string   `xml:"collection_teim"`
	CollectionUnitID_teim  string   `xml:"collection_unitid_teim"`
	Series_sim             []string `xml:"series_sim"`
	Series_si              string   `xml:"series_si"`
}

func (component *Component) setSolrAddMessage() {
	docElement := &component.SolrAddMessage.Add.Doc

	docElement.Address_teim = component.Parts.Address.Values

	docElement.Appraisal_teim = component.Parts.Appraisal.Values

	docElement.BiogHist_teim = component.Parts.BiogHist.Values

	docElement.ChronList_teim = component.Parts.ChronList.Values

	docElement.Collection_sim = component.Parts.Collection
	docElement.Collection_ssm = component.Parts.Collection
	docElement.Collection_teim = component.Parts.Collection

	docElement.CollectionUnitID_ssm = component.Parts.CollectionUnitID
	docElement.CollectionUnitID_teim = component.Parts.CollectionUnitID

	docElement.ComponentChildren_bsi = component.Parts.ComponentChildren
	docElement.ComponentLevel_isim = component.Parts.ComponentLevel

	docElement.CorpName_ssm = component.Parts.CorpName.Values
	docElement.CorpName_teim = component.Parts.CorpName.Values

	docElement.Creator_sim = append(docElement.Creator_sim,
		util.CompactStringSlicePreserveOrder(component.Parts.CreatorComplex.Values)...)
	// See 2nd `Creator_ssm` append below.
	docElement.Creator_ssm = append(docElement.Creator_ssm,
		component.Parts.CreatorCorpName.Values...)
	docElement.Creator_ssm = append(docElement.Creator_ssm,
		component.Parts.CreatorFamName.Values...)
	docElement.Creator_ssm = append(docElement.Creator_ssm,
		component.Parts.CreatorPersName.Values...)
	// TODO: is this duplication done in v1 indexer a bug that needs to be added
	// to https://jira.nyu.edu/browse/DLFA-211?
	docElement.Creator_ssm = append(docElement.Creator_ssm,
		util.CompactStringSlicePreserveOrder(component.Parts.CreatorComplex.Values)...)
	docElement.Creator_teim = append(docElement.Creator_teim,
		component.Parts.CreatorCorpName.Values...)
	docElement.Creator_teim = append(docElement.Creator_teim,
		component.Parts.CreatorFamName.Values...)
	docElement.Creator_teim = append(docElement.Creator_teim,
		component.Parts.CreatorPersName.Values...)

	docElement.DAO_sim = component.Parts.DAO.Values
	docElement.DAO_ssm = component.Parts.DAODescriptionParagraph.Values
	docElement.DAO_teim = component.Parts.DAODescriptionParagraph.Values

	docElement.EAD_ssi = component.Parts.EADID.Values[0]

	docElement.FamName_ssm = component.Parts.FamName.Values
	docElement.FamName_teim = component.Parts.FamName.Values

	docElement.Format_sim = component.Parts.Format.Values
	docElement.Format_ssm = component.Parts.Format.Values

	docElement.Function_ssm = component.Parts.Function.Values
	docElement.Function_teim = component.Parts.Function.Values

	docElement.GenreForm_ssm = component.Parts.GenreForm.Values
	docElement.GenreForm_teim = component.Parts.GenreForm.Values

	docElement.GeogName_ssm = component.Parts.GeogName.Values
	docElement.GeogName_teim = component.Parts.GeogName.Values

	docElement.ID = component.ID

	docElement.Level_sim = component.Parts.Level.Values[0]

	docElement.Location_si = component.Parts.Location.Values
	docElement.Location_ssm = component.Parts.Location.Values

	docElement.Name_sim = append(docElement.Name_ssm, component.Parts.Name.Values...)
	docElement.Name_ssm = append(docElement.Name_ssm, component.Parts.NameElementAll.Values...)
	docElement.Name_ssm = append(docElement.Name_ssm, component.Parts.Name.Values...)
	docElement.Name_teim = append(docElement.Name_teim, component.Parts.NameElementAll.Values...)
	docElement.Name_teim = append(docElement.Name_teim, component.Parts.Name.Values...)

	docElement.Note_ssm = component.Parts.Note.Values
	docElement.Note_teim = component.Parts.Note.Values

	docElement.Occupation_ssm = component.Parts.Occupation.Values
	docElement.Occupation_teim = component.Parts.Occupation.Values

	docElement.Place_sim = component.Parts.Place.Values
	docElement.Place_ssm = component.Parts.Place.Values

	if component.Parts.ParentForSort != "" {
		docElement.Parent_ssi = component.Parts.ParentForSort
	}
	docElement.Parent_ssm = component.Parts.ParentForDisplay.Values

	// This is identical to `docElement.Series_ssm`.
	docElement.ParentUnitTitles_ssm = component.Parts.AncestorUnitTitleList
	docElement.ParentUnitTitles_teim = component.Parts.AncestorUnitTitleList

	docElement.PersName_ssm = component.Parts.PersName.Values
	docElement.PersName_teim = component.Parts.PersName.Values

	docElement.PhysTech_teim = component.Parts.PhysTech.Values

	docElement.Repository_sim = component.Parts.RepositoryCode
	docElement.Repository_ssi = component.Parts.RepositoryCode
	docElement.Repository_ssm = component.Parts.RepositoryCode

	if len(component.Parts.Ref.Values) > 0 {
		docElement.Ref_ssi = component.Parts.Ref.Values[0]
	}

	docElement.ScopeContent_teim = component.Parts.ScopeContent.Values

	docElement.Series_si = component.Parts.SeriesForSort
	// This is identical to `docElement.ParentUnitTitles_ssm`.
	docElement.Series_sim = component.Parts.AncestorUnitTitleList

	docElement.Subject_sim = component.Parts.SubjectForFacets.Values
	docElement.Subject_ssm = component.Parts.Subject.Values
	docElement.Subject_teim = append(docElement.Subject_teim, component.Parts.Subject.Values...)
	// TODO: is this duplication done in v1 indexer a bug that needs to be added
	// to https://jira.nyu.edu/browse/DLFA-211?
	docElement.Subject_teim = append(docElement.Subject_teim, component.Parts.SubjectForFacets.Values...)

	docElement.Title_ssm = component.Parts.Title.Values
	docElement.Title_teim = component.Parts.Title.Values

	docElement.UnitDateInclusive_teim = component.Parts.UnitDateInclusive.Values

	docElement.UnitDateNormal_sim = component.Parts.UnitDateNormal.Values
	docElement.UnitDateNormal_ssm = component.Parts.UnitDateNormal.Values
	docElement.UnitDateNormal_teim = component.Parts.UnitDateNormal.Values

	docElement.UnitID_ssm = component.Parts.DIDUnitID.Values
	docElement.UnitID_teim = component.Parts.DIDUnitID.Values

	docElement.UnitTitle_ssm = component.Parts.UnitTitleHTML.Values
	docElement.UnitTitle_teim = component.Parts.DIDUnitTitle.Values
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
// Note that this function is duplicated in the `collection` and `component` packages.
// Normally we'd find a way to DRY this up (probably by using a `struct` param
// instead of the `CollectionDoc.SolrAddMessage` and `Component.SolrAddMessage`
// types, but since function is ephemeral, we just copy it.
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
				// TODO: DLFA-238
				// Re-enable non-empty string checks.  v1 indexer does not filter
				// out all whitespace values, it only filters out empty strings:
				// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10840271&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10840271
				// This is a DLFA-238 TODO within a function that is itself a DLFA-238 TODO.
				// Putting this here in case any of the code ends up being copy-pasted
				// into permanent functions.
				// if util.IsNonEmptyString(fieldValue) {
				if fieldValue != "" {
					fieldsInV1IndexerInsertionOrder = append(fieldsInV1IndexerInsertionOrder,
						eadutil.MakeSolrAddMessageFieldElementString(fieldName, fieldValue))
				}
			}
		} else if fieldTypeKind == reflect.String {
			fieldValue := field.String()
			/// TODO: DLFA-238
			// Re-enable non-empty string checks.  v1 indexer does not filter
			// out all whitespace values, it only filters out empty strings:
			// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10840271&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10840271
			// This is a DLFA-238 TODO within a function that is itself a DLFA-238 TODO.
			// Putting this here in case any of the code ends up being copy-pasted
			// into permanent functions.
			// if util.IsNonEmptyString(fieldValue) {
			if fieldValue != "" {
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
