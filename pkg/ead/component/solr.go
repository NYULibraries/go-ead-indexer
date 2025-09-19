package component

import (
	"fmt"
	"github.com/nyulibraries/go-ead-indexer/pkg/ead/eadutil"
	"github.com/nyulibraries/go-ead-indexer/pkg/util"
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

// Note that unlike the struct tags used for `SolrAddMessage` and `AddElement`,
// the struct tags for `DocElement` below are mandatory.  They are used for the
// field name in the Solr HTTP request XML body.
type DocElement struct {
	Address_teim           []string `xml:"address_teim"`
	Appraisal_teim         []string `xml:"appraisal_teim"`
	Author_teim            []string `xml:"author_teim"`
	BiogHist_teim          []string `xml:"bioghist_teim"`
	ChronList_teim         []string `xml:"chronlist_teim"`
	CollectionUnitID_ssm   string   `xml:"collection_unitid_ssm"`
	CollectionUnitID_teim  string   `xml:"collection_unitid_teim"`
	Collection_sim         string   `xml:"collection_sim"`
	Collection_ssm         string   `xml:"collection_ssm"`
	Collection_teim        string   `xml:"collection_teim"`
	ComponentChildren_bsi  string   `xml:"component_children_bsi"`
	ComponentLevel_isim    string   `xml:"component_level_isim"`
	CorpName_ssm           []string `xml:"corpname_ssm"`
	CorpName_teim          []string `xml:"corpname_teim"`
	Creator_sim            []string `xml:"creator_sim"`
	Creator_ssm            []string `xml:"creator_ssm"`
	Creator_teim           []string `xml:"creator_teim"`
	DAO_sim                []string `xml:"dao_sim"`
	DAO_ssm                []string `xml:"dao_ssm"`
	DAO_teim               []string `xml:"dao_teim"`
	DateRange_sim          []string `xml:"date_range_sim"`
	EAD_ssi                string   `xml:"ead_ssi"`
	FamName_ssm            []string `xml:"famname_ssm"`
	FamName_teim           []string `xml:"famname_teim"`
	Format_sim             []string `xml:"format_sim"`
	Format_ssm             []string `xml:"format_ssm"`
	Function_ssm           []string `xml:"function_ssm"`
	Function_teim          []string `xml:"function_teim"`
	GenreForm_ssm          []string `xml:"genreform_ssm"`
	GenreForm_teim         []string `xml:"genreform_teim"`
	GeogName_ssm           []string `xml:"geogname_ssm"`
	GeogName_teim          []string `xml:"geogname_teim"`
	Heading_ssm            []string `xml:"heading_ssm"`
	ID                     string   `xml:"id"`
	Language_sim           string   `xml:"language_sim"`
	Language_ssm           string   `xml:"language_ssm"`
	Level_sim              string   `xml:"level_sim"`
	Location_si            string   `xml:"location_si"`
	Location_ssm           []string `xml:"location_ssm"`
	MaterialType_sim       []string `xml:"material_type_sim"`
	MaterialType_ssm       []string `xml:"material_type_ssm"`
	Name_sim               []string `xml:"name_sim"`
	Name_ssm               []string `xml:"name_ssm"`
	Name_teim              []string `xml:"name_teim"`
	Note_ssm               []string `xml:"note_ssm"`
	Note_teim              []string `xml:"note_teim"`
	Occupation_ssm         []string `xml:"occupation_ssm"`
	Occupation_teim        []string `xml:"occupation_teim"`
	ParentUnitTitles_ssm   []string `xml:"parent_unittitles_ssm"`
	ParentUnitTitles_teim  []string `xml:"parent_unittitles_teim"`
	Parent_ssi             string   `xml:"parent_ssi"`
	Parent_ssm             []string `xml:"parent_ssm"`
	PersName_ssm           []string `xml:"persname_ssm"`
	PersName_teim          []string `xml:"persname_teim"`
	PhysTech_teim          []string `xml:"phystech_teim"`
	Place_sim              []string `xml:"place_sim"`
	Place_ssm              []string `xml:"place_ssm"`
	Ref_ssi                string   `xml:"ref_ssi"`
	Repository_sim         string   `xml:"repository_sim"`
	Repository_ssi         string   `xml:"repository_ssi"`
	Repository_ssm         string   `xml:"repository_ssm"`
	ScopeContent_teim      []string `xml:"scopecontent_teim"`
	Series_si              string   `xml:"series_si"`
	Series_sim             []string `xml:"series_sim"`
	Sort_ii                string   `xml:"sort_ii"`
	Subject_sim            []string `xml:"subject_sim"`
	Subject_ssm            []string `xml:"subject_ssm"`
	Subject_teim           []string `xml:"subject_teim"`
	Title_ssm              []string `xml:"title_ssm"`
	Title_teim             []string `xml:"title_teim"`
	UnitDateBulk_teim      []string `xml:"unitdate_bulk_teim"`
	UnitDateEnd_si         string   `xml:"unitdate_end_si"`
	UnitDateEnd_sim        []string `xml:"unitdate_end_sim"`
	UnitDateEnd_ssm        []string `xml:"unitdate_end_ssm"`
	UnitDateInclusive_teim []string `xml:"unitdate_inclusive_teim"`
	UnitDateNormal_sim     []string `xml:"unitdate_normal_sim"`
	UnitDateNormal_ssm     []string `xml:"unitdate_normal_ssm"`
	UnitDateNormal_teim    []string `xml:"unitdate_normal_teim"`
	UnitDateStart_si       string   `xml:"unitdate_start_si"`
	UnitDateStart_sim      []string `xml:"unitdate_start_sim"`
	UnitDateStart_ssm      []string `xml:"unitdate_start_ssm"`
	UnitDate_ssm           []string `xml:"unitdate_ssm"`
	UnitDate_teim          []string `xml:"unitdate_teim"`
	UnitID_ssm             []string `xml:"unitid_ssm"`
	UnitID_teim            []string `xml:"unitid_teim"`
	UnitTitle_ssm          []string `xml:"unittitle_ssm"`
	UnitTitle_teim         []string `xml:"unittitle_teim"`
}

func (component *Component) setSolrAddMessage() {
	docElement := &component.SolrAddMessage.Add.Doc

	docElement.Address_teim = component.Parts.Address.Values

	docElement.Appraisal_teim = component.Parts.Appraisal.Values

	docElement.Author_teim = component.Parts.Author

	docElement.BiogHist_teim = component.Parts.BiogHist.Values

	docElement.ChronList_teim = component.Parts.ChronListComplex.Values

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
	docElement.Creator_ssm = append(docElement.Creator_ssm,
		util.CompactStringSlicePreserveOrder(component.Parts.CreatorComplex.Values)...)
	docElement.Creator_teim = append(docElement.Creator_teim,
		component.Parts.Creator.Values...)

	docElement.DAO_sim = component.Parts.DAO.Values
	docElement.DAO_ssm = component.Parts.DAODescriptionParagraph.Values
	docElement.DAO_teim = component.Parts.DAODescriptionParagraph.Values

	docElement.DateRange_sim = append(docElement.DateRange_sim,
		util.CompactStringSlicePreserveOrder(component.Parts.DateRange.Values)...)

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

	docElement.Heading_ssm = component.Parts.Heading.Values

	if len(component.Parts.Language.Values) > 0 {
		docElement.Language_sim = component.Parts.Language.Values[0]
		docElement.Language_ssm = component.Parts.Language.Values[0]
	}

	docElement.Level_sim = component.Parts.Level.Values[0]

	if len(component.Parts.Location.Values) > 0 {
		docElement.Location_si = component.Parts.Location.Values[len(component.Parts.Location.Values)-1]
	}
	docElement.Location_ssm = util.CompactStringSlicePreserveOrder(
		component.Parts.Location.Values)

	docElement.MaterialType_sim = append(docElement.MaterialType_sim,
		util.CompactStringSlicePreserveOrder(component.Parts.MaterialType.Values)...)
	docElement.MaterialType_ssm = append(docElement.MaterialType_ssm,
		util.CompactStringSlicePreserveOrder(component.Parts.MaterialType.Values)...)

	docElement.Name_sim = append(docElement.Name_ssm, component.Parts.Name.Values...)
	docElement.Name_ssm = append(docElement.Name_ssm, component.Parts.NameElementAll.Values...)
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

	docElement.Sort_ii = strconv.Itoa(component.Parts.Sort)

	docElement.Subject_sim = component.Parts.SubjectForFacets.Values
	docElement.Subject_ssm = component.Parts.Subject.Values
	docElement.Subject_teim = append(docElement.Subject_teim, component.Parts.Subject.Values...)
	// TODO: is this duplication done in v1 indexer a bug that needs to be added
	// to https://jira.nyu.edu/browse/DLFA-211?
	docElement.Subject_teim = append(docElement.Subject_teim, component.Parts.SubjectForFacets.Values...)

	docElement.Title_ssm = component.Parts.Title.Values
	docElement.Title_teim = component.Parts.Title.Values

	docElement.UnitDateBulk_teim = component.Parts.UnitDateBulk.Values

	docElement.UnitDateInclusive_teim = component.Parts.UnitDateInclusive.Values

	docElement.UnitDateNormal_sim = component.Parts.UnitDateNormal.Values
	docElement.UnitDateNormal_ssm = component.Parts.UnitDateNormal.Values
	docElement.UnitDateNormal_teim = component.Parts.UnitDateNormal.Values

	docElement.UnitDate_ssm = append(docElement.UnitDate_ssm,
		util.CompactStringSlicePreserveOrder(component.Parts.UnitDateDisplay.Values)...)
	docElement.UnitDate_teim = append(docElement.UnitDate_teim, component.Parts.UnitDateNoTypeAttribute.Values...)

	if len(component.Parts.UnitDateEnd.Values) > 0 {
		docElement.UnitDateEnd_si = component.Parts.UnitDateEnd.Values[len(component.Parts.UnitDateEnd.Values)-1]
	}
	docElement.UnitDateEnd_sim = append(docElement.UnitDateEnd_sim,
		util.CompactStringSlicePreserveOrder(component.Parts.UnitDateEnd.Values)...)
	docElement.UnitDateEnd_ssm = append(docElement.UnitDateEnd_ssm,
		util.CompactStringSlicePreserveOrder(component.Parts.UnitDateEnd.Values)...)

	if len(component.Parts.UnitDateStart.Values) > 0 {
		docElement.UnitDateStart_si = component.Parts.UnitDateStart.Values[len(component.Parts.UnitDateStart.Values)-1]
	}
	docElement.UnitDateStart_sim = append(docElement.UnitDateStart_sim,
		util.CompactStringSlicePreserveOrder(component.Parts.UnitDateStart.Values)...)
	docElement.UnitDateStart_ssm = append(docElement.UnitDateStart_ssm,
		util.CompactStringSlicePreserveOrder(component.Parts.UnitDateStart.Values)...)

	docElement.UnitID_ssm = component.Parts.DIDUnitID.Values
	docElement.UnitID_teim = component.Parts.DIDUnitID.Values

	docElement.UnitTitle_ssm = component.Parts.UnitTitleHTML.Values
	docElement.UnitTitle_teim = component.Parts.DIDUnitTitle.Values
}

func (solrAddMessage SolrAddMessage) String() string {
	fields := eadutil.GetDocElementFieldsInAlphabeticalOrder(solrAddMessage.Add.Doc)
	fieldElementStrings := eadutil.MakeSolrAddMessageFieldElementStrings(fields)

	return eadutil.PrettifySolrAddMessageXML(
		fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><add><doc>%s</doc></add>`,
			strings.Join(fieldElementStrings, "")),
	)
}
