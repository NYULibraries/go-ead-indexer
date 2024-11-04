package collectiondoc

import (
	"github.com/lestrrat-go/libxml2/types"
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
