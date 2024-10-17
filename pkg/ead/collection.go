package ead

import (
	"github.com/lestrrat-go/libxml2/types"
	"go-ead-indexer/pkg/util"
)

type CollectionDoc struct {
	SolrAddMessage string
	Parts          CollectionDocParts
}

type CollectionDocParts struct {
	CollectionDocXPathCompositeParts
	CollectionDocXPathSimpleParts
}

type CollectionDocXPathCompositeParts struct {
	CreatorComposite CollectionDocXPathPart
	Name             CollectionDocXPathPart
}

type CollectionDocXPathSimpleParts struct {
	Abstract                CollectionDocXPathPart
	AcqInfo                 CollectionDocXPathPart
	Appraisal               CollectionDocXPathPart
	Author                  CollectionDocXPathPart
	BiogHist                CollectionDocXPathPart
	ChronList               CollectionDocXPathPart
	Collection              CollectionDocXPathPart
	CorpNameNotInRepository CollectionDocXPathPart
	CorpNameNotInDSC        CollectionDocXPathPart
	Creator                 CollectionDocXPathPart
	CreatorCorpName         CollectionDocXPathPart
	CreatorFamName          CollectionDocXPathPart
	CreatorPersName         CollectionDocXPathPart
	CustodHist              CollectionDocXPathPart
	EADID                   CollectionDocXPathPart
	FamName                 CollectionDocXPathPart
	FamNameNotInDSC         CollectionDocXPathPart
	FunctionNotInDSC        CollectionDocXPathPart
	GenreForm               CollectionDocXPathPart
	GenreFormNotInDSC       CollectionDocXPathPart
	GeogNameNotInDSC        CollectionDocXPathPart
	Geogname                CollectionDocXPathPart
	Heading                 CollectionDocXPathPart
	LangCode                CollectionDocXPathPart
	NameNotInDSC            CollectionDocXPathPart
	NoteNotInDSC            CollectionDocXPathPart
	OccupationNotInDSC      CollectionDocXPathPart
	PersName                CollectionDocXPathPart
	PersNameNotInDSC        CollectionDocXPathPart
	Phystech                CollectionDocXPathPart
	ScopeContent            CollectionDocXPathPart
	SubjectForFacets        CollectionDocXPathPart
	SubjectNotInDSC         CollectionDocXPathPart
	TitleNotInDSC           CollectionDocXPathPart
	UnitDateBulk            CollectionDocXPathPart
	UnitDateInclusive       CollectionDocXPathPart
	UnitDateNormal          CollectionDocXPathPart
	UnitDateNotType         CollectionDocXPathPart
	UnitID                  CollectionDocXPathPart
	UnitTitle               CollectionDocXPathPart
	Unitdate_normal         CollectionDocXPathPart
	Unitdate_start          CollectionDocXPathPart
}

type CollectionDocXPathPart struct {
	Source string
	Values []string
}

func MakeCollectionDoc(node types.Node) (CollectionDoc, error) {
	newCollectionDoc := CollectionDoc{
		SolrAddMessage: "",
		Parts:          CollectionDocParts{},
	}

	err := newCollectionDoc.populateParts(node)
	if err != nil {
		return newCollectionDoc, err
	}

	return newCollectionDoc, nil
}

func (collectionDoc *CollectionDoc) populateParts(node types.Node) error {
	err := collectionDoc.populateXPathSimpleParts(node)
	if err != nil {
		return err
	}

	collectionDoc.populateXPathCompositeParts()

	return nil
}

func (collectionDoc *CollectionDoc) populateXPathCompositeParts() {
	// TODO: Do we need to have anything in `CollectionDocXPathPart.Source`?

	parts := &collectionDoc.Parts

	//  CreatorComposite
	creatorCompositeValues := []string{}
	creatorCompositeValues = append(creatorCompositeValues, parts.CreatorCorpName.Values...)
	creatorCompositeValues = append(creatorCompositeValues, parts.CreatorFamName.Values...)
	creatorCompositeValues = append(creatorCompositeValues, parts.CreatorPersName.Values...)
	parts.CreatorComposite.Values = creatorCompositeValues

	//  Name
	nameValues := []string{}
	nameValues = append(nameValues, parts.FamName.Values...)
	nameValues = append(nameValues, parts.PersName.Values...)
	nameValues = append(nameValues, parts.CorpNameNotInRepository.Values...)
	nameValues = replaceMARCSubfieldDemarcatorsInSlice(nameValues)
	nameValues = util.CompactStringSlicePreserveOrder(nameValues)
	parts.Name.Values = nameValues
}

func (collectionDoc *CollectionDoc) populateXPathSimpleParts(node types.Node) error {
	var err error

	parts := &collectionDoc.Parts

	parts.Abstract.Source = "//archdesc[@level='collection']/did/abstract"
	parts.Abstract.Values, err = getValuesForXPathQuery(parts.Abstract.Source, node)
	if err != nil {
		return err
	}

	parts.AcqInfo.Source = "//archdesc[@level='collection']/acqinfo/p"
	parts.AcqInfo.Values, err = getValuesForXPathQuery(parts.AcqInfo.Source, node)
	if err != nil {
		return err
	}

	parts.Appraisal.Source = "//archdesc[@level='collection']/appraisal/p"
	parts.Appraisal.Values, err = getValuesForXPathQuery(parts.Appraisal.Source, node)
	if err != nil {
		return err
	}

	parts.Author.Source = "//filedesc/titlestmt/author"
	parts.Author.Values, err = getValuesForXPathQuery(parts.Author.Source, node)
	if err != nil {
		return err
	}

	parts.BiogHist.Source = "//archdesc[@level='collection']/bioghist/p"
	parts.BiogHist.Values, err = getValuesForXPathQuery(parts.BiogHist.Source, node)
	if err != nil {
		return err
	}

	parts.ChronList.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//chronlist/chronitem//text()"
	parts.ChronList.Values, err = getValuesForXPathQuery(parts.ChronList.Source, node)
	if err != nil {
		return err
	}

	parts.CorpNameNotInRepository.Source = "//*[local-name()!='repository']/corpname"
	parts.CorpNameNotInRepository.Values, err = getValuesForXPathQuery(parts.CorpNameNotInRepository.Source, node)
	if err != nil {
		return err
	}

	parts.CorpNameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//corpname"
	parts.CorpNameNotInDSC.Values, err = getValuesForXPathQuery(parts.CorpNameNotInDSC.Source, node)
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
	parts.Creator.Values, err = getValuesForXPathQuery(parts.Creator.Source, node)
	if err != nil {
		return err
	}

	parts.CreatorCorpName.Source = "//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/corpname"
	parts.CreatorCorpName.Values, err = getValuesForXPathQuery(parts.CreatorCorpName.Source, node)
	if err != nil {
		return err
	}

	parts.CreatorFamName.Source = "//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/famname"
	parts.CreatorFamName.Values, err = getValuesForXPathQuery(parts.CreatorFamName.Source, node)
	if err != nil {
		return err
	}

	parts.CreatorPersName.Source = "//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/persname"
	parts.CreatorPersName.Values, err = getValuesForXPathQuery(parts.CreatorPersName.Source, node)
	if err != nil {
		return err
	}

	parts.CustodHist.Source = "//archdesc[@level='collection']/custodhist/p"
	parts.CustodHist.Values, err = getValuesForXPathQuery(parts.CustodHist.Source, node)
	if err != nil {
		return err
	}

	parts.EADID.Source = "//eadid"
	parts.EADID.Values, err = getValuesForXPathQuery(parts.EADID.Source, node)
	if err != nil {
		return err
	}

	parts.FamName.Source = "//famname"
	parts.FamName.Values, err = getValuesForXPathQuery(parts.FamName.Source, node)
	if err != nil {
		return err
	}

	parts.FamNameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//famname"
	parts.FamNameNotInDSC.Values, err = getValuesForXPathQuery(parts.FamNameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.FunctionNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//function"
	parts.FunctionNotInDSC.Values, err = getValuesForXPathQuery(parts.FunctionNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.GenreForm.Source = "//genreform"
	parts.GenreForm.Values, err = getValuesForXPathQuery(parts.GenreForm.Source, node)
	if err != nil {
		return err
	}

	parts.GenreFormNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//genreform"
	parts.GenreFormNotInDSC.Values, err = getValuesForXPathQuery(parts.GenreFormNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.Geogname.Source = "//geogname"
	parts.Geogname.Values, err = getValuesForXPathQuery(parts.Geogname.Source, node)
	if err != nil {
		return err
	}

	parts.GeogNameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//geogname"
	parts.GeogNameNotInDSC.Values, err = getValuesForXPathQuery(parts.GeogNameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.LangCode.Source = "//archdesc[@level='collection']/did/langmaterial/language/@langcode"
	parts.LangCode.Values, err = getValuesForXPathQuery(parts.LangCode.Source, node)
	if err != nil {
		return err
	}

	parts.NameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//name"
	parts.NameNotInDSC.Values, err = getValuesForXPathQuery(parts.NameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.NoteNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//note"
	parts.NoteNotInDSC.Values, err = getValuesForXPathQuery(parts.NoteNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.OccupationNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//occupation"
	parts.OccupationNotInDSC.Values, err = getValuesForXPathQuery(parts.OccupationNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.PersName.Source = "//persname"
	parts.PersName.Values, err = getValuesForXPathQuery(parts.PersName.Source, node)
	if err != nil {
		return err
	}

	parts.PersNameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//persname"
	parts.PersNameNotInDSC.Values, err = getValuesForXPathQuery(parts.PersNameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.Phystech.Source = "//archdesc[@level='collection']/phystech/p"
	parts.Phystech.Values, err = getValuesForXPathQuery(parts.Phystech.Source, node)
	if err != nil {
		return err
	}

	parts.ScopeContent.Source = "//archdesc[@level='collection']/scopecontent/p"
	parts.ScopeContent.Values, err = getValuesForXPathQuery(parts.ScopeContent.Source, node)
	if err != nil {
		return err
	}

	parts.SubjectForFacets.Source = "//*[local-name()='subject' or local-name()='function' or local-name() = 'occupation']"
	parts.SubjectForFacets.Values, err = getValuesForXPathQuery(parts.SubjectForFacets.Source, node)
	if err != nil {
		return err
	}

	parts.SubjectNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//subject"
	parts.SubjectNotInDSC.Values, err = getValuesForXPathQuery(parts.SubjectNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.TitleNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//title"
	parts.TitleNotInDSC.Values, err = getValuesForXPathQuery(parts.TitleNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateBulk.Source = "//archdesc[@level='collection']/did/unitdate[@type='bulk']"
	parts.UnitDateBulk.Values, err = getValuesForXPathQuery(parts.UnitDateBulk.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateInclusive.Source = "//archdesc[@level='collection']/did/unitdate[@type='inclusive']"
	parts.UnitDateInclusive.Values, err = getValuesForXPathQuery(parts.UnitDateInclusive.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateNormal.Source = "//archdesc[@level='collection']/did/unitdate/@normal"
	parts.UnitDateNormal.Values, err = getValuesForXPathQuery(parts.UnitDateNormal.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateNotType.Source = "//archdesc[@level='collection']/did/unitdate[not(@type)]"
	parts.UnitDateNotType.Values, err = getValuesForXPathQuery(parts.UnitDateNotType.Source, node)
	if err != nil {
		return err
	}

	parts.UnitID.Source = "//archdesc[@level='collection']/did/unitid"
	parts.UnitID.Values, err = getValuesForXPathQuery(parts.UnitID.Source, node)
	if err != nil {
		return err
	}

	parts.UnitTitle.Source = "//archdesc[@level='collection']/did/unittitle"
	parts.UnitTitle.Values, err = getValuesForXPathQuery(parts.UnitTitle.Source, node)
	if err != nil {
		return err
	}

	return nil
}