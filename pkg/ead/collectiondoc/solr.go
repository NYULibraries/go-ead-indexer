package collectiondoc

import (
	"strconv"
)

type SolrAddMessage struct {
	Add AddElement `xml:"add"`
}

type AddElement struct {
	Doc DocElement `xml:"doc"`
}

type DocElement struct {
	Abstract_ssm           []string `xml:"abstract_ssm,omitempty"`
	Abstract_teim          []string `xml:"abstract_teim,omitempty"`
	AcqInfo_teim           []string `xml:"acqinfo_teim,omitempty"`
	Appraisal_teim         []string `xml:"appraisal_teim,omitempty"`
	Author_ssm             []string `xml:"author_ssm,omitempty"`
	Author_teim            []string `xml:"author_teim,omitempty"`
	BiogHist_teim          []string `xml:"bioghist_teim,omitempty"`
	ChronList_teim         []string `xml:"chronlist_teim,omitempty"`
	Collection_sim         []string `xml:"collection_sim,omitempty"`
	Collection_ssm         []string `xml:"collection_ssm,omitempty"`
	Collection_teim        []string `xml:"collection_teim,omitempty"`
	CorpName_ssm           []string `xml:"corpname_ssm,omitempty"`
	CorpName_teim          []string `xml:"corpname_teim,omitempty"`
	Creator_sim            []string `xml:"creator_sim,omitempty"`
	Creator_ssm            []string `xml:"creator_ssm,omitempty"`
	Creator_teim           []string `xml:"creator_teim,omitempty"`
	CustodHist_teim        []string `xml:"custodhist_teim,omitempty"`
	DAO_sim                string   `xml:"dao_sim,omitempty"`
	DateRange_sim          []string `xml:"date_range_sim,omitempty"`
	EAD_ssi                string   `xml:"ead_ssi"`
	FamName_ssm            []string `xml:"famname_ssm,omitempty"`
	FamName_teim           []string `xml:"famname_teim,omitempty"`
	Format_ii              string   `xml:"format_ii,omitempty"`
	Format_sim             string   `xml:"format_sim,omitempty"`
	Format_ssm             string   `xml:"format_ssm,omitempty"`
	Function_ssm           []string `xml:"function_ssm,omitempty"`
	Function_teim          []string `xml:"function_teim,omitempty"`
	GenreForm_ssm          []string `xml:"genreform_ssm,omitempty"`
	GenreForm_teim         []string `xml:"genreform_teim,omitempty"`
	GeogName_ssm           []string `xml:"geogname_ssm,omitempty"`
	GeogName_teim          []string `xml:"geogname_teim,omitempty"`
	Heading_ssm            []string `xml:"heading_ssm,omitempty"`
	ID                     string   `xml:"id"`
	Language_sim           string   `xml:"language_sim,omitempty"`
	Language_ssm           string   `xml:"language_ssm,omitempty"`
	MaterialType_sim       []string `xml:"material_type_sim,omitempty"`
	MaterialType_ssm       []string `xml:"material_type_ssm,omitempty"`
	Name_sim               []string `xml:"name_sim,omitempty"`
	Name_ssm               []string `xml:"name_ssm,omitempty"`
	Name_teim              []string `xml:"name_teim,omitempty"`
	Note_ssm               []string `xml:"note_ssm,omitempty"`
	Note_teim              []string `xml:"note_teim,omitempty"`
	Occupation_ssm         []string `xml:"occupation_ssm,omitempty"`
	Occupation_teim        []string `xml:"occupation_teim,omitempty"`
	PersName_ssm           []string `xml:"persname_ssm,omitempty"`
	PersName_teim          []string `xml:"persname_teim,omitempty"`
	PhysTech_teim          []string `xml:"phystech_teim,omitempty"`
	Place_sim              []string `xml:"place_sim,omitempty"`
	Repository_sim         string   `xml:"repository_sim"`
	Repository_ssi         string   `xml:"repository_ssi"`
	Repository_ssm         string   `xml:"repository_ssm"`
	ScopeContent_teim      []string `xml:"scopecontent_teim,omitempty"`
	Subject_sim            []string `xml:"subject_sim,omitempty"`
	Subject_ssm            []string `xml:"subject_ssm,omitempty"`
	Subject_teim           []string `xml:"subject_teim,omitempty"`
	Title_ssm              []string `xml:"title_ssm,omitempty"`
	Title_teim             []string `xml:"title_teim,omitempty"`
	UnitDateBulk_teim      []string `xml:"unitdate_bulk_teim,omitempty"`
	UnitDateEnd_si         []string `xml:"unitdate_end_si,omitempty"`
	UnitDateEnd_sim        []string `xml:"unitdate_end_sim,omitempty"`
	UnitDateEnd_ssm        []string `xml:"unitdate_end_ssm,omitempty"`
	UnitDateInclusive_teim []string `xml:"unitdate_inclusive_teim,omitempty"`
	UnitDateNormal_sim     []string `xml:"unitdate_normal_sim,omitempty"`
	UnitDateNormal_ssm     []string `xml:"unitdate_normal_ssm,omitempty"`
	UnitDateNormal_teim    []string `xml:"unitdate_normal_teim,omitempty"`
	UnitDate_ssm           []string `xml:"unitdate_ssm,omitempty"`
	UnitDateStart_si       []string `xml:"unitdate_start_si,omitempty"`
	UnitDateStart_sim      []string `xml:"unitdate_start_sim,omitempty"`
	UnitDateStart_ssm      []string `xml:"unitdate_start_ssm,omitempty"`
	UnitDate_teim          []string `xml:"unitdate_teim,omitempty"`
	UnitID_ssm             []string `xml:"unitid_ssm,omitempty"`
	UnitID_teim            []string `xml:"unitid_teim,omitempty"`
	UnitTitle_ssm          []string `xml:"unittitle_ssm,omitempty"`
	UnitTitle_teim         []string `xml:"unittitle_teim,omitempty"`
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

	docElement.ChronList_teim = append(docElement.ChronList_teim, collectionDoc.Parts.ChronList.Values...)

	docElement.Collection_sim = append(docElement.Collection_sim, collectionDoc.Parts.UnitTitle.Values...)
	docElement.Collection_ssm = append(docElement.Collection_ssm, collectionDoc.Parts.UnitTitle.Values...)
	docElement.Collection_teim = append(docElement.Collection_teim, collectionDoc.Parts.UnitTitle.Values...)

	docElement.CorpName_ssm = append(docElement.CorpName_ssm, collectionDoc.Parts.CorpNameNotInDSC.Values...)
	docElement.CorpName_teim = append(docElement.CorpName_teim, collectionDoc.Parts.CorpNameNotInDSC.Values...)

	// See 2nd `Creator_ssm` append below.
	docElement.Creator_sim = append(docElement.Creator_sim, collectionDoc.Parts.CreatorComplex.Values...)
	docElement.Creator_ssm = append(docElement.Creator_ssm, collectionDoc.Parts.Creator.Values...)
	// TODO: is this duplication done in v1 indexer a bug that needs to be added
	// to https://jira.nyu.edu/browse/DLFA-211?
	docElement.Creator_ssm = append(docElement.Creator_ssm, collectionDoc.Parts.CreatorComplex.Values...)
	docElement.Creator_teim = append(docElement.Creator_teim, collectionDoc.Parts.Creator.Values...)

	docElement.CustodHist_teim = append(docElement.CustodHist_teim, collectionDoc.Parts.CustodHist.Values...)

	docElement.DAO_sim = collectionDoc.Parts.DAO.Values[0]

	docElement.DateRange_sim = append(docElement.DateRange_sim, collectionDoc.Parts.DateRange.Values...)

	docElement.EAD_ssi = collectionDoc.Parts.EADID.Values[0]

	docElement.FamName_ssm = append(docElement.FamName_ssm, collectionDoc.Parts.FamName.Values...)
	docElement.FamName_teim = append(docElement.FamName_teim, collectionDoc.Parts.FamName.Values...)

	docElement.Format_ii = strconv.Itoa(collectionDoc.Parts.FormatForSort)
	docElement.Format_sim = collectionDoc.Parts.FormatForDisplay
	docElement.Format_ssm = collectionDoc.Parts.FormatForDisplay

	docElement.Function_ssm = append(docElement.Function_ssm, collectionDoc.Parts.Function.Values...)
	docElement.Function_teim = append(docElement.Function_teim, collectionDoc.Parts.Function.Values...)

	docElement.GenreForm_ssm = append(docElement.GenreForm_ssm, collectionDoc.Parts.GenreForm.Values...)
	docElement.GenreForm_teim = append(docElement.GenreForm_teim, collectionDoc.Parts.GenreForm.Values...)

	docElement.GeogName_ssm = append(docElement.GeogName_ssm, collectionDoc.Parts.GeogName.Values...)
	docElement.GeogName_teim = append(docElement.GeogName_teim, collectionDoc.Parts.GeogName.Values...)

	docElement.Heading_ssm = append(docElement.Heading_ssm, collectionDoc.Parts.Heading.Values...)

	docElement.ID = collectionDoc.Parts.EADID.Values[0]

	if len(collectionDoc.Parts.Language.Values) > 0 {
		docElement.Language_sim = collectionDoc.Parts.Language.Values[0]
		docElement.Language_ssm = collectionDoc.Parts.Language.Values[0]
	}

	docElement.MaterialType_sim = append(docElement.MaterialType_sim, collectionDoc.Parts.MaterialType.Values...)
	docElement.MaterialType_ssm = append(docElement.MaterialType_ssm, collectionDoc.Parts.MaterialType.Values...)

	docElement.Name_sim = append(docElement.Name_sim, collectionDoc.Parts.NameNotInDSC.Values...)
	docElement.Name_ssm = append(docElement.Name_ssm, collectionDoc.Parts.NameNotInDSC.Values...)
	docElement.Name_teim = append(docElement.Name_teim, collectionDoc.Parts.Name.Values...)

	docElement.Note_ssm = append(docElement.Note_ssm, collectionDoc.Parts.NoteNotInDSC.Values...)
	docElement.Note_teim = append(docElement.Note_teim, collectionDoc.Parts.NoteNotInDSC.Values...)

	docElement.Occupation_ssm = append(docElement.Occupation_ssm, collectionDoc.Parts.OccupationNotInDSC.Values...)
	docElement.Occupation_teim = append(docElement.Occupation_teim, collectionDoc.Parts.OccupationNotInDSC.Values...)

	docElement.PersName_ssm = append(docElement.PersName_ssm, collectionDoc.Parts.PersName.Values...)
	docElement.PersName_teim = append(docElement.PersName_teim, collectionDoc.Parts.PersName.Values...)

	docElement.PhysTech_teim = append(docElement.PhysTech_teim, collectionDoc.Parts.Phystech.Values...)

	docElement.Place_sim = append(docElement.Place_sim, collectionDoc.Parts.Place.Values...)

	docElement.Repository_sim = collectionDoc.Parts.RepositoryCode.Values[0]
	docElement.Repository_ssi = collectionDoc.Parts.RepositoryCode.Values[0]
	docElement.Repository_ssm = collectionDoc.Parts.RepositoryCode.Values[0]

	docElement.ScopeContent_teim = append(docElement.ScopeContent_teim, collectionDoc.Parts.ScopeContent.Values...)

	// See 2nd `Subject_teim` append below.
	docElement.Subject_sim = append(docElement.Subject_sim, collectionDoc.Parts.SubjectForFacets.Values...)
	docElement.Subject_ssm = append(docElement.Subject_ssm, collectionDoc.Parts.SubjectNotInDSC.Values...)
	docElement.Subject_teim = append(docElement.Subject_teim, collectionDoc.Parts.SubjectNotInDSC.Values...)
	// TODO: is this duplication done in v1 indexer a bug that needs to be added
	// to https://jira.nyu.edu/browse/DLFA-211?
	docElement.Subject_teim = append(docElement.Subject_teim, collectionDoc.Parts.SubjectForFacets.Values...)

	docElement.Title_ssm = append(docElement.Title_ssm, collectionDoc.Parts.TitleNotInDSC.Values...)
	docElement.Title_teim = append(docElement.Title_teim, collectionDoc.Parts.TitleNotInDSC.Values...)

	docElement.UnitDateBulk_teim = append(docElement.UnitDateBulk_teim, collectionDoc.Parts.UnitDateBulk.Values...)

	docElement.UnitDateEnd_si = append(docElement.UnitDateEnd_si, collectionDoc.Parts.UnitDateEnd.Values...)
	docElement.UnitDateEnd_sim = append(docElement.UnitDateEnd_sim, collectionDoc.Parts.UnitDateEnd.Values...)
	docElement.UnitDateEnd_ssm = append(docElement.UnitDateEnd_ssm, collectionDoc.Parts.UnitDateEnd.Values...)

	docElement.UnitDateInclusive_teim = append(docElement.UnitDateInclusive_teim, collectionDoc.Parts.UnitDateInclusive.Values...)

	docElement.UnitDateNormal_sim = append(docElement.UnitDateNormal_sim, collectionDoc.Parts.UnitDateNormal.Values...)
	docElement.UnitDateNormal_ssm = append(docElement.UnitDateNormal_ssm, collectionDoc.Parts.UnitDateNormal.Values...)
	docElement.UnitDateNormal_teim = append(docElement.UnitDateNormal_teim, collectionDoc.Parts.UnitDateNormal.Values...)

	docElement.UnitDate_ssm = append(docElement.UnitDate_ssm, collectionDoc.Parts.UnitDateDisplay.Values...)
	docElement.UnitDate_teim = append(docElement.UnitDate_teim, collectionDoc.Parts.UnitDateNoTypeAttribute.Values...)

	docElement.UnitDateStart_si = append(docElement.UnitDateStart_si, collectionDoc.Parts.UnitDateStart.Values...)
	docElement.UnitDateStart_sim = append(docElement.UnitDateStart_sim, collectionDoc.Parts.UnitDateStart.Values...)
	docElement.UnitDateStart_ssm = append(docElement.UnitDateStart_ssm, collectionDoc.Parts.UnitDateStart.Values...)

	docElement.UnitID_ssm = append(docElement.UnitID_ssm, collectionDoc.Parts.UnitID.Values...)
	docElement.UnitID_teim = append(docElement.UnitID_teim, collectionDoc.Parts.UnitID.Values...)

	docElement.UnitTitle_ssm = append(docElement.UnitTitle_ssm, collectionDoc.Parts.UnitTitleHTML.Values...)
	docElement.UnitTitle_teim = append(docElement.UnitTitle_teim, collectionDoc.Parts.UnitTitle.Values...)
}

func (solrAddMessage SolrAddMessage) String() string {
	return "test"
}
