package component

import (
	"github.com/lestrrat-go/libxml2/types"
	"go-ead-indexer/pkg/ead/util"
)

func (component *Component) setXPathSimpleParts(node types.Node) error {
	var err error

	parts := &component.Parts

	parts.Address.Source = "//address/p"
	parts.Address.Values, parts.Address.XMLStrings, err = util.GetValuesForXPathQuery(parts.Address.Source, node)
	if err != nil {
		return err
	}

	parts.Appraisal.Source = "//appraisal/p"
	parts.Appraisal.Values, parts.Appraisal.XMLStrings, err = util.GetValuesForXPathQuery(parts.Appraisal.Source, node)
	if err != nil {
		return err
	}

	parts.BiogHist.Source = "//bioghist/p"
	parts.BiogHist.Values, parts.BiogHist.XMLStrings, err = util.GetValuesForXPathQuery(parts.BiogHist.Source, node)
	if err != nil {
		return err
	}

	parts.ChronList.Source = "//chronlist/chronitem//text()"
	parts.ChronList.Values, parts.ChronList.XMLStrings, err = util.GetValuesForXPathQuery(parts.ChronList.Source, node)
	if err != nil {
		return err
	}

	parts.Collection.Source = "//archdesc/did/unittitle"
	parts.Collection.Values, parts.Collection.XMLStrings, err = util.GetValuesForXPathQuery(parts.Collection.Source, node)
	if err != nil {
		return err
	}

	parts.CollectionUnitID.Source = "//archdesc/did/unitid"
	parts.CollectionUnitID.Values, parts.CollectionUnitID.XMLStrings, err = util.GetValuesForXPathQuery(parts.CollectionUnitID.Source, node)
	if err != nil {
		return err
	}

	parts.CorpName.Source = "//corpname"
	parts.CorpName.Values, parts.CorpName.XMLStrings, err = util.GetValuesForXPathQuery(parts.CorpName.Source, node)
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

	parts.DAO.Source = "//dao/daodesc/p"
	parts.DAO.Values, parts.DAO.XMLStrings, err = util.GetValuesForXPathQuery(parts.DAO.Source, node)
	if err != nil {
		return err
	}

	parts.DIDUnitID.Source = "//did/unitid"
	parts.DIDUnitID.Values, parts.DIDUnitID.XMLStrings, err = util.GetValuesForXPathQuery(parts.DIDUnitID.Source, node)
	if err != nil {
		return err
	}

	parts.DIDUnitTitle.Source = "//did/unittitle"
	parts.DIDUnitTitle.Values, parts.DIDUnitTitle.XMLStrings, err = util.GetValuesForXPathQuery(parts.DIDUnitTitle.Source, node)
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

	parts.Function.Source = "//function"
	parts.Function.Values, parts.Function.XMLStrings, err = util.GetValuesForXPathQuery(parts.Function.Source, node)
	if err != nil {
		return err
	}

	parts.GenreForm.Source = "//genreform"
	parts.GenreForm.Values, parts.GenreForm.XMLStrings, err = util.GetValuesForXPathQuery(parts.GenreForm.Source, node)
	if err != nil {
		return err
	}

	parts.GeogName.Source = "//geogname"
	parts.GeogName.Values, parts.GeogName.XMLStrings, err = util.GetValuesForXPathQuery(parts.GeogName.Source, node)
	if err != nil {
		return err
	}

	parts.Heading.Source = "//archdesc[@level='collection']/did/unittitle"
	parts.Heading.Values, parts.Heading.XMLStrings, err = util.GetValuesForXPathQuery(parts.Heading.Source, node)
	if err != nil {
		return err
	}

	parts.Language.Source = "//did/langmaterial/language/@langcode"
	parts.Language.Values, parts.Language.XMLStrings, err = util.GetValuesForXPathQuery(parts.Language.Source, node)
	if err != nil {
		return err
	}

	parts.Level.Source = "///c/@level"
	parts.Level.Values, parts.Level.XMLStrings, err = util.GetValuesForXPathQuery(parts.Level.Source, node)
	if err != nil {
		return err
	}

	parts.NameElementAll.Source = "//name"
	parts.NameElementAll.Values, parts.NameElementAll.XMLStrings, err = util.GetValuesForXPathQuery(parts.NameElementAll.Source, node)
	if err != nil {
		return err
	}

	parts.Note.Source = "//note"
	parts.Note.Values, parts.Note.XMLStrings, err = util.GetValuesForXPathQuery(parts.Note.Source, node)
	if err != nil {
		return err
	}

	parts.Occupation.Source = "//occupation"
	parts.Occupation.Values, parts.Occupation.XMLStrings, err = util.GetValuesForXPathQuery(parts.Occupation.Source, node)
	if err != nil {
		return err
	}

	parts.PersName.Source = "//persname"
	parts.PersName.Values, parts.PersName.XMLStrings, err = util.GetValuesForXPathQuery(parts.PersName.Source, node)
	if err != nil {
		return err
	}

	parts.PhysTech.Source = "//phystech/p"
	parts.PhysTech.Values, parts.PhysTech.XMLStrings, err = util.GetValuesForXPathQuery(parts.PhysTech.Source, node)
	if err != nil {
		return err
	}

	parts.Ref.Source = "///c/@id"
	parts.Ref.Values, parts.Ref.XMLStrings, err = util.GetValuesForXPathQuery(parts.Ref.Source, node)
	if err != nil {
		return err
	}

	parts.ScopeContent.Source = "//scopecontent/p"
	parts.ScopeContent.Values, parts.ScopeContent.XMLStrings, err = util.GetValuesForXPathQuery(parts.ScopeContent.Source, node)
	if err != nil {
		return err
	}

	parts.Subject.Source = "//subject"
	parts.Subject.Values, parts.Subject.XMLStrings, err = util.GetValuesForXPathQuery(parts.Subject.Source, node)
	if err != nil {
		return err
	}

	parts.SubjectOrFunctionOrOccupation.Source = "//*[local-name()='subject' or local-name()='function' or local-name() = 'occupation']"
	parts.SubjectOrFunctionOrOccupation.Values, parts.SubjectOrFunctionOrOccupation.XMLStrings, err = util.GetValuesForXPathQuery(parts.SubjectOrFunctionOrOccupation.Source, node)
	if err != nil {
		return err
	}

	parts.Title.Source = "//title"
	parts.Title.Values, parts.Title.XMLStrings, err = util.GetValuesForXPathQuery(parts.Title.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateNotType.Source = "//did/unitdate[not(@type)]"
	parts.UnitDateNotType.Values, parts.UnitDateNotType.XMLStrings, err = util.GetValuesForXPathQuery(parts.UnitDateNotType.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateBulk.Source = "//archdesc[@level='collection']/did/unitdate[@type='bulk']"
	parts.UnitDateBulk.Values, parts.UnitDateBulk.XMLStrings, err = util.GetValuesForXPathQuery(parts.UnitDateBulk.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateNormal.Source = "//did/unitdate/@normal"
	parts.UnitDateNormal.Values, parts.UnitDateNormal.XMLStrings, err = util.GetValuesForXPathQuery(parts.UnitDateNormal.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateInclusive.Source = "//archdesc[@level='collection']/did/unitdate[@type='inclusive']"
	parts.UnitDateInclusive.Values, parts.UnitDateInclusive.XMLStrings, err = util.GetValuesForXPathQuery(parts.UnitDateInclusive.Source, node)
	if err != nil {
		return err
	}

	return nil
}
