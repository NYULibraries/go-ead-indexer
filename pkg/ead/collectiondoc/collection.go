package collectiondoc

import (
	"github.com/lestrrat-go/libxml2/types"
	"go-ead-indexer/pkg/ead/util"
)

type CollectionDoc struct {
	SolrAddMessage string             `json:"solr_add_message"`
	Parts          CollectionDocParts `json:"parts"`
}

// For now, no struct tags for the `CollectionDoc*` fields.  Keep it flat.
type CollectionDocParts struct {
	CollectionDocComplexParts
	CollectionDocHardcodedParts
	CollectionDocXPathParts
	RepositoryCode CollectionDocPart `json:"repository_code"`
}

type CollectionDocComplexParts struct {
	CreatorComplex CollectionDocPart `json:"creator_complex"`
	DateRange      CollectionDocPart `json:"date_range"`
	MaterialType   CollectionDocPart `json:"material_type"`
	Name           CollectionDocPart `json:"name"`
	Place          CollectionDocPart `json:"place"`
	OnlineAccess   CollectionDocPart `json:"online_access"`
	Subject        CollectionDocPart `json:"subject"`
	UnitDateEnd    CollectionDocPart `json:"unit_date_end"`
	UnitDateStart  CollectionDocPart `json:"unit_date_start"`
	UnitTitleHTML  CollectionDocPart `json:"unit_title_html"`
}

type CollectionDocHardcodedParts struct {
	FormatForDisplay string `json:"format_for_display"`
	FormatForSort    int    `json:"format_for_sort"`
}

type CollectionDocXPathParts struct {
	Abstract                CollectionDocPart `json:"abstract"`
	AcqInfo                 CollectionDocPart `json:"acq_info"`
	Appraisal               CollectionDocPart `json:"appraisal"`
	Author                  CollectionDocPart `json:"author"`
	BiogHist                CollectionDocPart `json:"biog_hist"`
	ChronList               CollectionDocPart `json:"chron_list"`
	Collection              CollectionDocPart `json:"collection"`
	CorpNameNotInRepository CollectionDocPart `json:"corp_name_not_in_repository"`
	CorpNameNotInDSC        CollectionDocPart `json:"corp_name_not_in_dsc"`
	Creator                 CollectionDocPart `json:"creator"`
	CreatorCorpName         CollectionDocPart `json:"creator_corp_name"`
	CreatorFamName          CollectionDocPart `json:"creator_fam_name"`
	CreatorPersName         CollectionDocPart `json:"creator_pers_name"`
	CustodHist              CollectionDocPart `json:"custod_hist"`
	DAO                     CollectionDocPart `json:"dao"`
	EADID                   CollectionDocPart `json:"eadid"`
	FamName                 CollectionDocPart `json:"fam_name"`
	FamNameNotInDSC         CollectionDocPart `json:"fam_name_not_in_dsc"`
	FunctionNotInDSC        CollectionDocPart `json:"function_not_in_dsc"`
	GenreForm               CollectionDocPart `json:"genre_form"`
	GenreFormNotInDSC       CollectionDocPart `json:"genre_form_not_in_dsc"`
	GeogNameNotInDSC        CollectionDocPart `json:"geog_name_not_in_dsc"`
	GeogName                CollectionDocPart `json:"geog_name"`
	Heading                 CollectionDocPart `json:"heading"`
	LangCode                CollectionDocPart `json:"lang_code"`
	Language                CollectionDocPart `json:"language"`
	NameNotInDSC            CollectionDocPart `json:"name_not_in_dsc"`
	NoteNotInDSC            CollectionDocPart `json:"note_not_in_dsc"`
	OccupationNotInDSC      CollectionDocPart `json:"occupation_not_in_dsc"`
	PersName                CollectionDocPart `json:"pers_name"`
	PersNameNotInDSC        CollectionDocPart `json:"pers_name_not_in_dsc"`
	Phystech                CollectionDocPart `json:"phystech"`
	ScopeContent            CollectionDocPart `json:"scope_content"`
	SubjectForFacets        CollectionDocPart `json:"subject_for_facets"`
	SubjectNotInDSC         CollectionDocPart `json:"subject_not_in_dsc"`
	TitleNotInDSC           CollectionDocPart `json:"title_not_in_dsc"`
	UnitDateDisplay         CollectionDocPart `json:"unit_date_display"`
	UnitDateBulk            CollectionDocPart `json:"unit_date_bulk"`
	UnitDateInclusive       CollectionDocPart `json:"unit_date_inclusive"`
	UnitDateNormal          CollectionDocPart `json:"unit_date_normal"`
	UnitDateNoTypeAttribute CollectionDocPart `json:"unit_date_no_type_attribute"`
	UnitID                  CollectionDocPart `json:"unit_id"`
	UnitTitle               CollectionDocPart `json:"unit_title"`
}

type CollectionDocPart struct {
	Source     string   `json:"source"`
	Values     []string `json:"values"`
	XMLStrings []string `json:"xml_strings"`
}

// See `ead.new()` comment on why we have to pass in `repositoryCode` as an argument.
func MakeCollectionDoc(repositoryCode string, node types.Node) (CollectionDoc, error) {
	newCollectionDoc := CollectionDoc{
		SolrAddMessage: "",
		Parts: CollectionDocParts{
			RepositoryCode: CollectionDocPart{
				Values: []string{repositoryCode},
			},
		},
	}

	err := newCollectionDoc.setParts(node)
	if err != nil {
		return newCollectionDoc, err
	}

	return newCollectionDoc, nil
}

func (collectionDoc *CollectionDoc) setParts(node types.Node) error {
	err := collectionDoc.setXPathSimpleParts(node)
	if err != nil {
		return err
	}

	collectionDoc.setComplexParts()
	collectionDoc.setHardcodedParts()

	return nil
}

// TODO: Do we need to have anything in `CollectionDoc.Part.Source` for these?
func (collectionDoc *CollectionDoc) setComplexParts() []error {
	errs := []error{}

	collectionDoc.setCreator()
	collectionDoc.setDateRange()
	languageErrors := collectionDoc.setLanguage()
	if len(languageErrors) > 0 {
		errs = append(errs, languageErrors...)
	}
	collectionDoc.setMaterialType()
	collectionDoc.setName()
	collectionDoc.setOnlineAccess()
	collectionDoc.setPlace()
	collectionDoc.setSubject()
	collectionDoc.setUnitDateEnd()
	collectionDoc.setUnitDateStart()
	unitTitleHTMLError := collectionDoc.setUnitTitleHTML()
	if unitTitleHTMLError != nil {
		errs = append(errs, unitTitleHTMLError)
	}

	return errs
}

func (collectionDoc *CollectionDoc) setHardcodedParts() {
	collectionDoc.Parts.FormatForDisplay = "Archival Collection"
	collectionDoc.Parts.FormatForSort = 0
}

func (collectionDoc *CollectionDoc) setXPathSimpleParts(node types.Node) error {
	var err error

	parts := &collectionDoc.Parts

	parts.Abstract.Source = "//archdesc[@level='collection']/did/abstract"
	parts.Abstract.Values, parts.Abstract.XMLStrings, err = util.GetValuesForXPathQuery(parts.Abstract.Source, node)
	if err != nil {
		return err
	}

	parts.AcqInfo.Source = "//archdesc[@level='collection']/acqinfo/p"
	parts.AcqInfo.Values, parts.AcqInfo.XMLStrings, err = util.GetValuesForXPathQuery(parts.AcqInfo.Source, node)
	if err != nil {
		return err
	}

	parts.Appraisal.Source = "//archdesc[@level='collection']/appraisal/p"
	parts.Appraisal.Values, parts.Appraisal.XMLStrings, err = util.GetValuesForXPathQuery(parts.Appraisal.Source, node)
	if err != nil {
		return err
	}

	parts.Author.Source = "//filedesc/titlestmt/author"
	parts.Author.Values, parts.Author.XMLStrings, err = util.GetValuesForXPathQuery(parts.Author.Source, node)
	if err != nil {
		return err
	}

	parts.BiogHist.Source = "//archdesc[@level='collection']/bioghist/p"
	parts.BiogHist.Values, parts.BiogHist.XMLStrings, err = util.GetValuesForXPathQuery(parts.BiogHist.Source, node)
	if err != nil {
		return err
	}

	parts.ChronList.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//chronlist/chronitem//text()"
	parts.ChronList.Values, parts.ChronList.XMLStrings, err = util.GetValuesForXPathQuery(parts.ChronList.Source, node)
	if err != nil {
		return err
	}

	parts.CorpNameNotInRepository.Source = "//*[local-name()!='repository']/corpname"
	parts.CorpNameNotInRepository.Values, parts.CorpNameNotInRepository.XMLStrings, err = util.GetValuesForXPathQuery(parts.CorpNameNotInRepository.Source, node)
	if err != nil {
		return err
	}

	parts.CorpNameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//corpname"
	parts.CorpNameNotInDSC.Values, parts.CorpNameNotInDSC.XMLStrings, err = util.GetValuesForXPathQuery(parts.CorpNameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	// We need to be able to find elements with `label="Creator"` and `label="creator"`.
	// For details, see email thread starting with email sent by Joe on Mon, Aug 28, 2023, 12:56PM
	// with subject:
	// "FADESIGN: ead-publisher taken offline, full site rebuild in progress, missing creator facet"
	// ...and Jira ticket: https: //jira.nyu.edu/browse/FADESIGN-843.
	//
	// Note that XPath 2.0 functions `matches` and `lower-case` don't work for
	// here.  `matches(@label,'creator','i')` fails with compile errors:
	//
	//           xmlXPathCompOpEval: function matches not found
	//           XPath error : Unregistered function
	//
	// ...`lower-case(@label)='creator'`, the same.  Presumably this is because
	// the libxml2 package we are using doesn't support XPath 2.0.
	//
	// The `translate` solution we use below for the `Creator*` fields seems
	// to be the common method for who don't have XPath 2.0 options:
	// "Case insensitive xpaths"
	// https://groups.google.com/g/selenium-users/c/Lcvbjisk4qE
	// "case-insensitive matching in XPath?"
	// https://stackoverflow.com/questions/2893551/case-insensitive-matching-in-xpath
	parts.Creator.Source = "//archdesc[@level='collection']/did/origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/*[name() = 'corpname' or name() = 'famname' or name() = 'persname']"
	parts.Creator.Values, parts.Creator.XMLStrings, err = util.GetValuesForXPathQuery(parts.Creator.Source, node)
	if err != nil {
		return err
	}

	parts.CreatorCorpName.Source = "//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/corpname"
	parts.CreatorCorpName.Values, parts.CreatorCorpName.XMLStrings, err = util.GetValuesForXPathQuery(parts.CreatorCorpName.Source, node)
	if err != nil {
		return err
	}

	parts.CreatorFamName.Source = "//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/famname"
	parts.CreatorFamName.Values, parts.CreatorFamName.XMLStrings, err = util.GetValuesForXPathQuery(parts.CreatorFamName.Source, node)
	if err != nil {
		return err
	}

	parts.CreatorPersName.Source = "//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/persname"
	parts.CreatorPersName.Values, parts.CreatorPersName.XMLStrings, err = util.GetValuesForXPathQuery(parts.CreatorPersName.Source, node)
	if err != nil {
		return err
	}

	parts.CustodHist.Source = "//archdesc[@level='collection']/custodhist/p"
	parts.CustodHist.Values, parts.CustodHist.XMLStrings, err = util.GetValuesForXPathQuery(parts.CustodHist.Source, node)
	if err != nil {
		return err
	}

	parts.DAO.Source = "//dao"
	parts.DAO.Values, parts.DAO.XMLStrings, err = util.GetValuesForXPathQuery(parts.DAO.Source, node)
	if err != nil {
		return err
	}

	parts.EADID.Source = "//eadid"
	parts.EADID.Values, parts.EADID.XMLStrings, err = util.GetValuesForXPathQuery(parts.EADID.Source, node)
	if err != nil {
		return err
	}

	parts.FamName.Source = "//famname"
	parts.FamName.Values, parts.FamName.XMLStrings, err = util.GetValuesForXPathQuery(parts.FamName.Source, node)
	if err != nil {
		return err
	}

	parts.FamNameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//famname"
	parts.FamNameNotInDSC.Values, parts.FamNameNotInDSC.XMLStrings, err = util.GetValuesForXPathQuery(parts.FamNameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.FunctionNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//function"
	parts.FunctionNotInDSC.Values, parts.FunctionNotInDSC.XMLStrings, err = util.GetValuesForXPathQuery(parts.FunctionNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.GenreForm.Source = "//genreform"
	parts.GenreForm.Values, parts.GenreForm.XMLStrings, err = util.GetValuesForXPathQuery(parts.GenreForm.Source, node)
	if err != nil {
		return err
	}

	parts.GenreFormNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//genreform"
	parts.GenreFormNotInDSC.Values, parts.GenreFormNotInDSC.XMLStrings, err = util.GetValuesForXPathQuery(parts.GenreFormNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.GeogName.Source = "//geogname"
	parts.GeogName.Values, parts.GeogName.XMLStrings, err = util.GetValuesForXPathQuery(parts.GeogName.Source, node)
	if err != nil {
		return err
	}

	parts.GeogNameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//geogname"
	parts.GeogNameNotInDSC.Values, parts.GeogNameNotInDSC.XMLStrings, err = util.GetValuesForXPathQuery(parts.GeogNameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.LangCode.Source = "//archdesc[@level='collection']/did/langmaterial/language/@langcode"
	parts.LangCode.Values, parts.LangCode.XMLStrings, err = util.GetValuesForXPathQuery(parts.LangCode.Source, node)
	if err != nil {
		return err
	}

	parts.NameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//name"
	parts.NameNotInDSC.Values, parts.NameNotInDSC.XMLStrings, err = util.GetValuesForXPathQuery(parts.NameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.NoteNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//note"
	parts.NoteNotInDSC.Values, parts.NoteNotInDSC.XMLStrings, err = util.GetValuesForXPathQuery(parts.NoteNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.OccupationNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//occupation"
	parts.OccupationNotInDSC.Values, parts.OccupationNotInDSC.XMLStrings, err = util.GetValuesForXPathQuery(parts.OccupationNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.PersName.Source = "//persname"
	parts.PersName.Values, parts.PersName.XMLStrings, err = util.GetValuesForXPathQuery(parts.PersName.Source, node)
	if err != nil {
		return err
	}

	parts.PersNameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//persname"
	parts.PersNameNotInDSC.Values, parts.PersNameNotInDSC.XMLStrings, err = util.GetValuesForXPathQuery(parts.PersNameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.Phystech.Source = "//archdesc[@level='collection']/phystech/p"
	parts.Phystech.Values, parts.Phystech.XMLStrings, err = util.GetValuesForXPathQuery(parts.Phystech.Source, node)
	if err != nil {
		return err
	}

	parts.ScopeContent.Source = "//archdesc[@level='collection']/scopecontent/p"
	parts.ScopeContent.Values, parts.ScopeContent.XMLStrings, err = util.GetValuesForXPathQuery(parts.ScopeContent.Source, node)
	if err != nil {
		return err
	}

	parts.SubjectForFacets.Source = "//*[local-name()='subject' or local-name()='function' or local-name() = 'occupation']"
	parts.SubjectForFacets.Values, parts.SubjectForFacets.XMLStrings, err = util.GetValuesForXPathQuery(parts.SubjectForFacets.Source, node)
	if err != nil {
		return err
	}

	parts.SubjectNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//subject"
	parts.SubjectNotInDSC.Values, parts.SubjectNotInDSC.XMLStrings, err = util.GetValuesForXPathQuery(parts.SubjectNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.TitleNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//title"
	parts.TitleNotInDSC.Values, parts.TitleNotInDSC.XMLStrings, err = util.GetValuesForXPathQuery(parts.TitleNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateBulk.Source = "//archdesc[@level='collection']/did/unitdate[@type='bulk']"
	parts.UnitDateBulk.Values, parts.UnitDateBulk.XMLStrings, err = util.GetValuesForXPathQuery(parts.UnitDateBulk.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateInclusive.Source = "//archdesc[@level='collection']/did/unitdate[@type='inclusive']"
	parts.UnitDateInclusive.Values, parts.UnitDateInclusive.XMLStrings, err = util.GetValuesForXPathQuery(parts.UnitDateInclusive.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateNormal.Source = "//archdesc[@level='collection']/did/unitdate/@normal"
	parts.UnitDateNormal.Values, parts.UnitDateNormal.XMLStrings, err = util.GetValuesForXPathQuery(parts.UnitDateNormal.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateNoTypeAttribute.Source = "//archdesc[@level='collection']/did/unitdate[not(@type)]"
	parts.UnitDateNoTypeAttribute.Values, parts.UnitDateNoTypeAttribute.XMLStrings, err = util.GetValuesForXPathQuery(parts.UnitDateNoTypeAttribute.Source, node)
	if err != nil {
		return err
	}

	parts.UnitID.Source = "//archdesc[@level='collection']/did/unitid"
	parts.UnitID.Values, parts.UnitID.XMLStrings, err = util.GetValuesForXPathQuery(parts.UnitID.Source, node)
	if err != nil {
		return err
	}

	parts.UnitTitle.Source = "//archdesc[@level='collection']/did/unittitle"
	parts.UnitTitle.Values, parts.UnitTitle.XMLStrings, err = util.GetValuesForXPathQuery(parts.UnitTitle.Source, node)
	if err != nil {
		return err
	}

	// Proxy for UnitTitle
	parts.Collection.Source = parts.UnitTitle.Source
	parts.Collection.Values, parts.Collection.XMLStrings, err = util.GetValuesForXPathQuery(parts.Collection.Source, node)
	if err != nil {
		return err
	}

	// Proxy for UnitTitle
	parts.Heading.Source = parts.UnitTitle.Source
	parts.Heading.Values, parts.Heading.XMLStrings, err = util.GetValuesForXPathQuery(parts.Heading.Source, node)
	if err != nil {
		return err
	}

	return nil
}
