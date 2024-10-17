package ead

import (
	"github.com/lestrrat-go/libxml2/types"
	"go-ead-indexer/pkg/util"
)

type Collection struct {
	SolrAddMessage string
	Parts          CollectionParts
}

type CollectionParts struct {
	XPath CollectionXPath
}

type CollectionXPath struct {
	Composite CollectionXPathCompositeParts
	Simple    CollectionXPathSimpleParts
}

type CollectionXPathCompositeParts struct {
	Creator []string
	Name    []string
}

type CollectionXPathSimpleParts struct {
	Abstract                CollectionXPathPart
	AcqInfo                 CollectionXPathPart
	Appraisal               CollectionXPathPart
	Author                  CollectionXPathPart
	BiogHist                CollectionXPathPart
	ChronList               CollectionXPathPart
	Collection              CollectionXPathPart
	CorpNameNotInRepository CollectionXPathPart
	CorpNameNotInDSC        CollectionXPathPart
	Creator                 CollectionXPathPart
	CreatorCorpName         CollectionXPathPart
	CreatorFamName          CollectionXPathPart
	CreatorPersName         CollectionXPathPart
	CustodHist              CollectionXPathPart
	EADID                   CollectionXPathPart
	FamName                 CollectionXPathPart
	FamNameNotInDSC         CollectionXPathPart
	FunctionNotInDSC        CollectionXPathPart
	GenreForm               CollectionXPathPart
	GenreFormNotInDSC       CollectionXPathPart
	GeogNameNotInDSC        CollectionXPathPart
	Geogname                CollectionXPathPart
	Heading                 CollectionXPathPart
	LangCode                CollectionXPathPart
	NameNotInDSC            CollectionXPathPart
	NoteNotInDSC            CollectionXPathPart
	OccupationNotInDSC      CollectionXPathPart
	PersName                CollectionXPathPart
	PersNameNotInDSC        CollectionXPathPart
	Phystech                CollectionXPathPart
	ScopeContent            CollectionXPathPart
	SubjectForFacets        CollectionXPathPart
	SubjectNotInDSC         CollectionXPathPart
	TitleNotInDSC           CollectionXPathPart
	UnitDateBulk            CollectionXPathPart
	UnitDateInclusive       CollectionXPathPart
	UnitDateNormal          CollectionXPathPart
	UnitDateNotType         CollectionXPathPart
	UnitID                  CollectionXPathPart
	UnitTitle               CollectionXPathPart
	Unitdate_normal         CollectionXPathPart
	Unitdate_start          CollectionXPathPart
}

type CollectionXPathPart struct {
	Query  string
	Values []string
}

func MakeCollection(node types.Node) (Collection, error) {
	newCollection := Collection{
		SolrAddMessage: "",
		Parts:          CollectionParts{},
	}

	err := newCollection.populateParts(node)
	if err != nil {
		return newCollection, err
	}

	return newCollection, nil
}

func (collection *Collection) populateParts(node types.Node) error {
	err := collection.populateXPathSimpleParts(node)
	if err != nil {
		return err
	}

	collection.populateXPathCompositeParts()

	return nil
}

func (collection *Collection) populateXPathCompositeParts() {
	xpc := &collection.Parts.XPath.Composite
	xps := &collection.Parts.XPath.Simple

	//  Creator
	creator := []string{}
	creator = append(creator, xps.CreatorCorpName.Values...)
	creator = append(creator, xps.CreatorFamName.Values...)
	creator = append(creator, xps.CreatorPersName.Values...)
	xpc.Creator = creator

	//  Name
	name := []string{}
	name = append(name, xps.FamName.Values...)
	name = append(name, xps.PersName.Values...)
	name = append(name, xps.CorpNameNotInRepository.Values...)
	name = replaceMARCSubfieldDemarcatorsInSlice(name)
	name = util.CompactStringSlicePreserveOrder(name)
	xpc.Name = name
}

func (collection *Collection) populateXPathSimpleParts(node types.Node) error {
	var err error

	xp := &collection.Parts.XPath.Simple

	xp.Abstract.Query = "//archdesc[@level='collection']/did/abstract"
	xp.Abstract.Values, err = getValuesForXPathQuery(xp.Abstract.Query, node)
	if err != nil {
		return err
	}

	xp.AcqInfo.Query = "//archdesc[@level='collection']/acqinfo/p"
	xp.AcqInfo.Values, err = getValuesForXPathQuery(xp.AcqInfo.Query, node)
	if err != nil {
		return err
	}

	xp.Appraisal.Query = "//archdesc[@level='collection']/appraisal/p"
	xp.Appraisal.Values, err = getValuesForXPathQuery(xp.Appraisal.Query, node)
	if err != nil {
		return err
	}

	xp.Author.Query = "//filedesc/titlestmt/author"
	xp.Author.Values, err = getValuesForXPathQuery(xp.Author.Query, node)
	if err != nil {
		return err
	}

	xp.BiogHist.Query = "//archdesc[@level='collection']/bioghist/p"
	xp.BiogHist.Values, err = getValuesForXPathQuery(xp.BiogHist.Query, node)
	if err != nil {
		return err
	}

	xp.ChronList.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//chronlist/chronitem//text()"
	xp.ChronList.Values, err = getValuesForXPathQuery(xp.ChronList.Query, node)
	if err != nil {
		return err
	}

	xp.CorpNameNotInRepository.Query = "//*[local-name()!='repository']/corpname"
	xp.CorpNameNotInRepository.Values, err = getValuesForXPathQuery(xp.CorpNameNotInRepository.Query, node)
	if err != nil {
		return err
	}

	xp.CorpNameNotInDSC.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//corpname"
	xp.CorpNameNotInDSC.Values, err = getValuesForXPathQuery(xp.CorpNameNotInDSC.Query, node)
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
	xp.Creator.Query = "//archdesc[@level='collection']/did/origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/*[name() = 'corpname' or name() = 'famname' or name() = 'persname']"
	xp.Creator.Values, err = getValuesForXPathQuery(xp.Creator.Query, node)
	if err != nil {
		return err
	}

	xp.CreatorCorpName.Query = "//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/corpname"
	xp.CreatorCorpName.Values, err = getValuesForXPathQuery(xp.CreatorCorpName.Query, node)
	if err != nil {
		return err
	}

	xp.CreatorFamName.Query = "//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/famname"
	xp.CreatorFamName.Values, err = getValuesForXPathQuery(xp.CreatorFamName.Query, node)
	if err != nil {
		return err
	}

	xp.CreatorPersName.Query = "//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/persname"
	xp.CreatorPersName.Values, err = getValuesForXPathQuery(xp.CreatorPersName.Query, node)
	if err != nil {
		return err
	}

	xp.CustodHist.Query = "//archdesc[@level='collection']/custodhist/p"
	xp.CustodHist.Values, err = getValuesForXPathQuery(xp.CustodHist.Query, node)
	if err != nil {
		return err
	}

	xp.EADID.Query = "//eadid"
	xp.EADID.Values, err = getValuesForXPathQuery(xp.EADID.Query, node)
	if err != nil {
		return err
	}

	xp.FamName.Query = "//famname"
	xp.FamName.Values, err = getValuesForXPathQuery(xp.FamName.Query, node)
	if err != nil {
		return err
	}

	xp.FamNameNotInDSC.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//famname"
	xp.FamNameNotInDSC.Values, err = getValuesForXPathQuery(xp.FamNameNotInDSC.Query, node)
	if err != nil {
		return err
	}

	xp.FunctionNotInDSC.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//function"
	xp.FunctionNotInDSC.Values, err = getValuesForXPathQuery(xp.FunctionNotInDSC.Query, node)
	if err != nil {
		return err
	}

	xp.GenreForm.Query = "//genreform"
	xp.GenreForm.Values, err = getValuesForXPathQuery(xp.GenreForm.Query, node)
	if err != nil {
		return err
	}

	xp.GenreFormNotInDSC.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//genreform"
	xp.GenreFormNotInDSC.Values, err = getValuesForXPathQuery(xp.GenreFormNotInDSC.Query, node)
	if err != nil {
		return err
	}

	xp.Geogname.Query = "//geogname"
	xp.Geogname.Values, err = getValuesForXPathQuery(xp.Geogname.Query, node)
	if err != nil {
		return err
	}

	xp.GeogNameNotInDSC.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//geogname"
	xp.GeogNameNotInDSC.Values, err = getValuesForXPathQuery(xp.GeogNameNotInDSC.Query, node)
	if err != nil {
		return err
	}

	xp.LangCode.Query = "//archdesc[@level='collection']/did/langmaterial/language/@langcode"
	xp.LangCode.Values, err = getValuesForXPathQuery(xp.LangCode.Query, node)
	if err != nil {
		return err
	}

	xp.NameNotInDSC.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//name"
	xp.NameNotInDSC.Values, err = getValuesForXPathQuery(xp.NameNotInDSC.Query, node)
	if err != nil {
		return err
	}

	xp.NoteNotInDSC.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//note"
	xp.NoteNotInDSC.Values, err = getValuesForXPathQuery(xp.NoteNotInDSC.Query, node)
	if err != nil {
		return err
	}

	xp.OccupationNotInDSC.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//occupation"
	xp.OccupationNotInDSC.Values, err = getValuesForXPathQuery(xp.OccupationNotInDSC.Query, node)
	if err != nil {
		return err
	}

	xp.PersName.Query = "//persname"
	xp.PersName.Values, err = getValuesForXPathQuery(xp.PersName.Query, node)
	if err != nil {
		return err
	}

	xp.PersNameNotInDSC.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//persname"
	xp.PersNameNotInDSC.Values, err = getValuesForXPathQuery(xp.PersNameNotInDSC.Query, node)
	if err != nil {
		return err
	}

	xp.Phystech.Query = "//archdesc[@level='collection']/phystech/p"
	xp.Phystech.Values, err = getValuesForXPathQuery(xp.Phystech.Query, node)
	if err != nil {
		return err
	}

	xp.ScopeContent.Query = "//archdesc[@level='collection']/scopecontent/p"
	xp.ScopeContent.Values, err = getValuesForXPathQuery(xp.ScopeContent.Query, node)
	if err != nil {
		return err
	}

	xp.SubjectForFacets.Query = "//*[local-name()='subject' or local-name()='function' or local-name() = 'occupation']"
	xp.SubjectForFacets.Values, err = getValuesForXPathQuery(xp.SubjectForFacets.Query, node)
	if err != nil {
		return err
	}

	xp.SubjectNotInDSC.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//subject"
	xp.SubjectNotInDSC.Values, err = getValuesForXPathQuery(xp.SubjectNotInDSC.Query, node)
	if err != nil {
		return err
	}

	xp.TitleNotInDSC.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//title"
	xp.TitleNotInDSC.Values, err = getValuesForXPathQuery(xp.TitleNotInDSC.Query, node)
	if err != nil {
		return err
	}

	xp.UnitDateBulk.Query = "//archdesc[@level='collection']/did/unitdate[@type='bulk']"
	xp.UnitDateBulk.Values, err = getValuesForXPathQuery(xp.UnitDateBulk.Query, node)
	if err != nil {
		return err
	}

	xp.UnitDateInclusive.Query = "//archdesc[@level='collection']/did/unitdate[@type='inclusive']"
	xp.UnitDateInclusive.Values, err = getValuesForXPathQuery(xp.UnitDateInclusive.Query, node)
	if err != nil {
		return err
	}

	xp.UnitDateNormal.Query = "//archdesc[@level='collection']/did/unitdate/@normal"
	xp.UnitDateNormal.Values, err = getValuesForXPathQuery(xp.UnitDateNormal.Query, node)
	if err != nil {
		return err
	}

	xp.UnitDateNotType.Query = "//archdesc[@level='collection']/did/unitdate[not(@type)]"
	xp.UnitDateNotType.Values, err = getValuesForXPathQuery(xp.UnitDateNotType.Query, node)
	if err != nil {
		return err
	}

	xp.UnitID.Query = "//archdesc[@level='collection']/did/unitid"
	xp.UnitID.Values, err = getValuesForXPathQuery(xp.UnitID.Query, node)
	if err != nil {
		return err
	}

	xp.UnitTitle.Query = "//archdesc[@level='collection']/did/unittitle"
	xp.UnitTitle.Values, err = getValuesForXPathQuery(xp.UnitTitle.Query, node)
	if err != nil {
		return err
	}

	return nil
}
