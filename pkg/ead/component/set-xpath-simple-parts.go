package component

import (
	"github.com/lestrrat-go/libxml2/types"
	"go-ead-indexer/pkg/ead/eadutil"
)

// Even though `node` is technically a <c> element node, XPath expressions that
// start with "//" when passed to `node.Find()` seem to query in the context of
// the entire EAD document.  This is the case even though `node.Find()` makes
// a next `Context` based on `node` (which is <c>).  There likely is a way
// to force querying relative to the <c> node using the LibXML2 API to make a
// new context based on copy, but the cheaper way for us would seem to be to
// just add "." to the beginning of XPath expressions which want running in the
// context of the <c> node.
// See https://stackoverflow.com/questions/33416524/libxml2-xpath-relative-to-sub-node
// for demonstration and explanation of the options.
func (component *Component) setXPathSimpleParts(node types.Node) error {
	var err error

	parts := &component.Parts

	parts.Address.Source = ".//address/p"
	parts.Address.Values, parts.Address.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.Address.Source, node)
	if err != nil {
		return err
	}

	parts.Appraisal.Source = ".//appraisal/p"
	parts.Appraisal.Values, parts.Appraisal.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.Appraisal.Source, node)
	if err != nil {
		return err
	}

	parts.BiogHist.Source = ".//bioghist/p"
	parts.BiogHist.Values, parts.BiogHist.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.BiogHist.Source, node)
	if err != nil {
		return err
	}

	parts.ChronList.Source = ".//chronlist/chronitem//text()"
	parts.ChronList.Values, parts.ChronList.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.ChronList.Source, node)
	if err != nil {
		return err
	}

	parts.CorpName.Source = ".//corpname"
	parts.CorpName.Values, parts.CorpName.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.CorpName.Source, node)
	if err != nil {
		return err
	}

	parts.CorpNameNotInRepository.Source = ".//*[local-name()!='repository']/corpname"
	parts.CorpNameNotInRepository.Values, parts.CorpNameNotInRepository.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.CorpNameNotInRepository.Source, node)
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
	parts.Creator.Source = ".//did/origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/*[name() = 'corpname' or name() = 'famname' or name() = 'persname']"
	parts.Creator.Values, parts.Creator.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.Creator.Source, node)
	if err != nil {
		return err
	}

	parts.CreatorCorpName.Source = ".//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/corpname"
	parts.CreatorCorpName.Values, parts.CreatorCorpName.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.CreatorCorpName.Source, node)
	if err != nil {
		return err
	}

	parts.CreatorFamName.Source = ".//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/famname"
	parts.CreatorFamName.Values, parts.CreatorFamName.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.CreatorFamName.Source, node)
	if err != nil {
		return err
	}

	parts.CreatorPersName.Source = ".//origination[translate(@label, 'ABCDEFGHIJKLMNOPQRSTUVWXYZ','abcdefghijklmnopqrstuvwxyz')='creator']/persname"
	parts.CreatorPersName.Values, parts.CreatorPersName.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.CreatorPersName.Source, node)
	if err != nil {
		return err
	}

	parts.DAODescriptionParagraph.Source = ".//dao/daodesc/p"
	parts.DAODescriptionParagraph.Values, parts.DAODescriptionParagraph.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.DAODescriptionParagraph.Source, node)
	if err != nil {
		return err
	}

	parts.DIDUnitID.Source = ".//did/unitid"
	parts.DIDUnitID.Values, parts.DIDUnitID.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.DIDUnitID.Source, node)
	if err != nil {
		return err
	}

	parts.DIDUnitTitle.Source = ".//did/unittitle"
	parts.DIDUnitTitle.Values, parts.DIDUnitTitle.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.DIDUnitTitle.Source, node)
	if err != nil {
		return err
	}

	parts.EADID.Source = "//eadid"
	parts.EADID.Values, parts.EADID.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.EADID.Source, node)
	if err != nil {
		return err
	}

	parts.FamName.Source = ".//famname"
	parts.FamName.Values, parts.FamName.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.FamName.Source, node)
	if err != nil {
		return err
	}

	parts.Function.Source = ".//function"
	parts.Function.Values, parts.Function.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.Function.Source, node)
	if err != nil {
		return err
	}

	parts.GenreForm.Source = ".//genreform"
	parts.GenreForm.Values, parts.GenreForm.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.GenreForm.Source, node)
	if err != nil {
		return err
	}

	parts.GeogName.Source = ".//geogname"
	parts.GeogName.Values, parts.GeogName.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.GeogName.Source, node)
	if err != nil {
		return err
	}

	parts.LangCode.Source = ".//did/langmaterial/language/@langcode"
	parts.LangCode.Values, parts.LangCode.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.LangCode.Source, node)
	if err != nil {
		return err
	}

	parts.Level.Source = "./@level"
	parts.Level.Values, parts.Level.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.Level.Source, node)
	if err != nil {
		return err
	}

	parts.NameElementAll.Source = ".//name"
	parts.NameElementAll.Values, parts.NameElementAll.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.NameElementAll.Source, node)
	if err != nil {
		return err
	}

	parts.Note.Source = ".//note"
	parts.Note.Values, parts.Note.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.Note.Source, node)
	if err != nil {
		return err
	}

	parts.Occupation.Source = ".//occupation"
	parts.Occupation.Values, parts.Occupation.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.Occupation.Source, node)
	if err != nil {
		return err
	}

	parts.PersName.Source = ".//persname"
	parts.PersName.Values, parts.PersName.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.PersName.Source, node)
	if err != nil {
		return err
	}

	parts.PhysTech.Source = ".//phystech/p"
	parts.PhysTech.Values, parts.PhysTech.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.PhysTech.Source, node)
	if err != nil {
		return err
	}

	parts.Ref.Source = "./@id"
	parts.Ref.Values, parts.Ref.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.Ref.Source, node)
	if err != nil {
		return err
	}

	parts.ScopeContent.Source = ".//scopecontent/p"
	parts.ScopeContent.Values, parts.ScopeContent.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.ScopeContent.Source, node)
	if err != nil {
		return err
	}

	parts.Subject.Source = ".//subject"
	parts.Subject.Values, parts.Subject.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.Subject.Source, node)
	if err != nil {
		return err
	}

	parts.SubjectOrFunctionOrOccupation.Source = ".//*[local-name()='subject' or local-name()='function' or local-name() = 'occupation']"
	parts.SubjectOrFunctionOrOccupation.Values, parts.SubjectOrFunctionOrOccupation.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.SubjectOrFunctionOrOccupation.Source, node)
	if err != nil {
		return err
	}

	parts.Title.Source = ".//title"
	parts.Title.Values, parts.Title.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.Title.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateBulk.Source = ".//did/unitdate[@type='bulk']"
	parts.UnitDateBulk.Values, parts.UnitDateBulk.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.UnitDateBulk.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateInclusive.Source = ".//did/unitdate[@type='inclusive']"
	parts.UnitDateInclusive.Values, parts.UnitDateInclusive.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.UnitDateInclusive.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateNormal.Source = ".//did/unitdate/@normal"
	parts.UnitDateNormal.Values, parts.UnitDateNormal.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.UnitDateNormal.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateNoTypeAttribute.Source = ".//did/unitdate[not(@type)]"
	parts.UnitDateNoTypeAttribute.Values, parts.UnitDateNoTypeAttribute.XMLStrings, err = eadutil.GetNodeValuesAndXMLStrings(parts.UnitDateNoTypeAttribute.Source, node)
	if err != nil {
		return err
	}

	return nil
}
