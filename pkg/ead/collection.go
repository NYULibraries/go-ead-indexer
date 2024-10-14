package ead

import (
	"github.com/lestrrat-go/libxml2/types"
)

type Collection struct {
	SolrAddMessage string
	XPathParts     CollectionXPathParts
}

type CollectionXPathParts struct {
	Abstract           CollectionXPathPart
	AcqInfo            CollectionXPathPart
	Appraisal          CollectionXPathPart
	Author             CollectionXPathPart
	BiogHist           CollectionXPathPart
	ChronList          CollectionXPathPart
	Collection         CollectionXPathPart
	CorpNameNotInDSC   CollectionXPathPart
	Creator            CollectionXPathPart
	CustodHist         CollectionXPathPart
	EADID              CollectionXPathPart
	FamNameNotInDSC    CollectionXPathPart
	FunctionNotInDSC   CollectionXPathPart
	GenreForm          CollectionXPathPart
	GenreFormNotInDSC  CollectionXPathPart
	GeogNameNotInDSC   CollectionXPathPart
	Geogname           CollectionXPathPart
	Heading            CollectionXPathPart
	LangCode           CollectionXPathPart
	NameNotInDSC       CollectionXPathPart
	NoteNotInDSC       CollectionXPathPart
	OccupationNotInDSC CollectionXPathPart
	PersnameNotInDSC   CollectionXPathPart
	Phystech           CollectionXPathPart
	ScopeContent       CollectionXPathPart
	SubjectForFacets   CollectionXPathPart
	SubjectNotInDSC    CollectionXPathPart
	TitleNotInDSC      CollectionXPathPart
	UnitDateBulk       CollectionXPathPart
	UnitDateInclusive  CollectionXPathPart
	UnitDateNormal     CollectionXPathPart
	UnitDateNotType    CollectionXPathPart
	UnitID             CollectionXPathPart
	UnitTitle          CollectionXPathPart
	Unitdate_normal    CollectionXPathPart
	Unitdate_start     CollectionXPathPart
}

type CollectionXPathPart struct {
	Query  string
	Values []string
}

func (collection *Collection) populateXPathParts(node types.Node) error {
	var err error

	xp := &collection.XPathParts

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

	xp.CorpNameNotInDSC.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//corpname"
	xp.CorpNameNotInDSC.Values, err = getValuesForXPathQuery(xp.CorpNameNotInDSC.Query, node)
	if err != nil {
		return err
	}

	xp.Creator.Query = "//archdesc[@level='collection']/did/origination[@label='creator']/*[name() = 'corpname' or name() = 'famname' or name() = 'persname']"
	xp.Creator.Values, err = getValuesForXPathQuery(xp.Creator.Query, node)
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

	xp.PersnameNotInDSC.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//persname"
	xp.PersnameNotInDSC.Values, err = getValuesForXPathQuery(xp.PersnameNotInDSC.Query, node)
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

func getValuesForXPathQuery(query string, node types.Node) ([]string, error) {
	var values []string

	xpathResult, err := node.Find(query)
	if err != nil {
		return nil, err
	}

	for _, node = range xpathResult.NodeList() {
		values = append(values, node.NodeValue())
	}

	return values, nil
}
