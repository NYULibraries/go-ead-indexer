package collectiondoc

import (
	"github.com/lestrrat-go/libxml2/types"
	"go-ead-indexer/pkg/ead/eadutil"
)

func (collectionDoc *CollectionDoc) setXPathSimpleParts(node types.Node) error {
	var err error

	parts := &collectionDoc.Parts

	parts.Abstract.Source = "//archdesc[@level='collection']/did/abstract"
	parts.Abstract.Values, parts.Abstract.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.Abstract.Source, node)
	if err != nil {
		return err
	}

	parts.AcqInfo.Source = "//archdesc[@level='collection']/acqinfo/p"
	parts.AcqInfo.Values, parts.AcqInfo.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.AcqInfo.Source, node)
	if err != nil {
		return err
	}

	parts.Appraisal.Source = "//archdesc[@level='collection']/appraisal/p"
	parts.Appraisal.Values, parts.Appraisal.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.Appraisal.Source, node)
	if err != nil {
		return err
	}

	parts.Author.Source = "//filedesc/titlestmt/author"
	parts.Author.Values, parts.Author.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.Author.Source, node)
	if err != nil {
		return err
	}

	parts.BiogHist.Source = "//archdesc[@level='collection']/bioghist/p"
	parts.BiogHist.Values, parts.BiogHist.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.BiogHist.Source, node)
	if err != nil {
		return err
	}

	parts.ChronList.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//chronlist/chronitem//text()"
	parts.ChronList.Values, parts.ChronList.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.ChronList.Source, node)
	if err != nil {
		return err
	}

	parts.CorpNameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//corpname"
	parts.CorpNameNotInDSC.Values, parts.CorpNameNotInDSC.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.CorpNameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.CorpNameNotInRepository.Source = "//*[local-name()!='repository']/corpname"
	parts.CorpNameNotInRepository.Values, parts.CorpNameNotInRepository.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.CorpNameNotInRepository.Source, node)
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
	parts.Creator.Values, parts.Creator.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.Creator.Source, node)
	if err != nil {
		return err
	}

	parts.CreatorCorpName.Source = "//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/corpname"
	parts.CreatorCorpName.Values, parts.CreatorCorpName.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.CreatorCorpName.Source, node)
	if err != nil {
		return err
	}

	parts.CreatorFamName.Source = "//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/famname"
	parts.CreatorFamName.Values, parts.CreatorFamName.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.CreatorFamName.Source, node)
	if err != nil {
		return err
	}

	parts.CreatorPersName.Source = "//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/persname"
	parts.CreatorPersName.Values, parts.CreatorPersName.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.CreatorPersName.Source, node)
	if err != nil {
		return err
	}

	parts.CustodHist.Source = "//archdesc[@level='collection']/custodhist/p"
	parts.CustodHist.Values, parts.CustodHist.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.CustodHist.Source, node)
	if err != nil {
		return err
	}

	parts.DAO.Source = "//dao"
	parts.DAO.Values, parts.DAO.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.DAO.Source, node)
	if err != nil {
		return err
	}

	parts.EADID.Source = "//eadid"
	parts.EADID.Values, parts.EADID.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.EADID.Source, node)
	if err != nil {
		return err
	}

	parts.FamName.Source = "//famname"
	parts.FamName.Values, parts.FamName.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.FamName.Source, node)
	if err != nil {
		return err
	}

	parts.FamNameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//famname"
	parts.FamNameNotInDSC.Values, parts.FamNameNotInDSC.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.FamNameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.Function.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//function"
	parts.Function.Values, parts.Function.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.Function.Source, node)
	if err != nil {
		return err
	}

	parts.GenreForm.Source = "//genreform"
	parts.GenreForm.Values, parts.GenreForm.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.GenreForm.Source, node)
	if err != nil {
		return err
	}

	parts.GenreFormNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//genreform"
	parts.GenreFormNotInDSC.Values, parts.GenreFormNotInDSC.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.GenreFormNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.GeogName.Source = "//geogname"
	parts.GeogName.Values, parts.GeogName.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.GeogName.Source, node)
	if err != nil {
		return err
	}

	parts.GeogNameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//geogname"
	parts.GeogNameNotInDSC.Values, parts.GeogNameNotInDSC.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.GeogNameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.LangCode.Source = "//archdesc[@level='collection']/did/langmaterial/language/@langcode"
	parts.LangCode.Values, parts.LangCode.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.LangCode.Source, node)
	if err != nil {
		return err
	}

	parts.NameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//name"
	parts.NameNotInDSC.Values, parts.NameNotInDSC.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.NameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.NoteNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//note"
	parts.NoteNotInDSC.Values, parts.NoteNotInDSC.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.NoteNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.OccupationNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//occupation"
	parts.OccupationNotInDSC.Values, parts.OccupationNotInDSC.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.OccupationNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.PersName.Source = "//persname"
	parts.PersName.Values, parts.PersName.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.PersName.Source, node)
	if err != nil {
		return err
	}

	parts.PersNameNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//persname"
	parts.PersNameNotInDSC.Values, parts.PersNameNotInDSC.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.PersNameNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.Phystech.Source = "//archdesc[@level='collection']/phystech/p"
	parts.Phystech.Values, parts.Phystech.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.Phystech.Source, node)
	if err != nil {
		return err
	}

	parts.ScopeContent.Source = "//archdesc[@level='collection']/scopecontent/p"
	parts.ScopeContent.Values, parts.ScopeContent.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.ScopeContent.Source, node)
	if err != nil {
		return err
	}

	parts.Subject.Source = "//*[local-name()='subject' or local-name()='function' or local-name() = 'occupation']"
	parts.Subject.Values, parts.Subject.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.Subject.Source, node)
	if err != nil {
		return err
	}

	parts.SubjectNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//subject"
	parts.SubjectNotInDSC.Values, parts.SubjectNotInDSC.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.SubjectNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.TitleNotInDSC.Source = "//archdesc[@level='collection']/*[name() != 'dsc']//title"
	parts.TitleNotInDSC.Values, parts.TitleNotInDSC.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.TitleNotInDSC.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateBulk.Source = "//archdesc[@level='collection']/did/unitdate[@type='bulk']"
	parts.UnitDateBulk.Values, parts.UnitDateBulk.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.UnitDateBulk.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateInclusive.Source = "//archdesc[@level='collection']/did/unitdate[@type='inclusive']"
	parts.UnitDateInclusive.Values, parts.UnitDateInclusive.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.UnitDateInclusive.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateNormal.Source = "//archdesc[@level='collection']/did/unitdate/@normal"
	parts.UnitDateNormal.Values, parts.UnitDateNormal.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.UnitDateNormal.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateNoTypeAttribute.Source = "//archdesc[@level='collection']/did/unitdate[not(@type)]"
	parts.UnitDateNoTypeAttribute.Values, parts.UnitDateNoTypeAttribute.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.UnitDateNoTypeAttribute.Source, node)
	if err != nil {
		return err
	}

	parts.UnitID.Source = "//archdesc[@level='collection']/did/unitid[not(@type='aspace_uri')]"
	parts.UnitID.Values, parts.UnitID.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.UnitID.Source, node)
	if err != nil {
		return err
	}

	parts.UnitTitle.Source = "//archdesc[@level='collection']/did/unittitle"
	parts.UnitTitle.Values, parts.UnitTitle.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.UnitTitle.Source, node)
	if err != nil {
		return err
	}

	// PROXY PARTS
	// These parts depend on other simple XPath parts, so we do them at the end.

	// Proxy for UnitTitle
	parts.Collection.Source = parts.UnitTitle.Source
	parts.Collection.Values, parts.Collection.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.Collection.Source, node)
	if err != nil {
		return err
	}

	// Proxy for UnitTitle
	parts.Heading.Source = parts.UnitTitle.Source
	parts.Heading.Values, parts.Heading.XMLStrings, err = eadutil.GetValuesForXPathQuery(parts.Heading.Source, node)
	if err != nil {
		return err
	}

	return nil
}
