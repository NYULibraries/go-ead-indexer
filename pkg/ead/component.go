package ead

import "github.com/lestrrat-go/libxml2/types"

type Component struct {
	ID             string
	Parts          ComponentParts
	SolrAddMessage string
}

type ComponentParts struct {
	XPath ComponentXPath
}

type ComponentXPath struct {
	Simple ComponentXPathSimpleParts
}

type ComponentXPathSimpleParts struct {
	Address                       ComponentXPathPart
	Appraisal                     ComponentXPathPart
	BiogHist                      ComponentXPathPart
	ChronList                     ComponentXPathPart
	Collection                    ComponentXPathPart
	CollectionUnitID              ComponentXPathPart
	Corpname                      ComponentXPathPart
	Creator                       ComponentXPathPart
	DAO                           ComponentXPathPart
	DIDUnitID                     ComponentXPathPart
	DIDUnitTitle                  ComponentXPathPart
	EADID                         ComponentXPathPart
	FamName                       ComponentXPathPart
	Function                      ComponentXPathPart
	GenreForm                     ComponentXPathPart
	GeogName                      ComponentXPathPart
	Heading                       ComponentXPathPart
	Language                      ComponentXPathPart
	Level                         ComponentXPathPart
	Name                          ComponentXPathPart
	Note                          ComponentXPathPart
	Occupation                    ComponentXPathPart
	PersName                      ComponentXPathPart
	PhysTech                      ComponentXPathPart
	Ref                           ComponentXPathPart
	ScopeContent                  ComponentXPathPart
	Subject                       ComponentXPathPart
	SubjectOrFunctionOrOccupation ComponentXPathPart
	Title                         ComponentXPathPart
	UnitDateNotType               ComponentXPathPart
	UnitDateBulk                  ComponentXPathPart
	UnitDateNormal                ComponentXPathPart
	UnitdateInclusive             ComponentXPathPart
}

type ComponentXPathPart struct {
	Query  string
	Values []string
}

func MakeComponents(node types.Node) (*[]Component, error) {
	xpathResult, err := node.Find("//c")
	if err != nil {
		return nil, err
	}

	// Note: can't do `&xpathResult.NodeList()`
	// See https://groups.google.com/g/golang-nuts/c/reaIlFdibWU
	cNodes := xpathResult.NodeList()

	// Early exit
	if len(cNodes) == 0 {
		return nil, nil
	}

	components := []Component{}
	for _, cNode := range cNodes {
		newComponent, err := MakeComponent(cNode)
		if err != nil {
			return &components, err
		}

		components = append(components, newComponent)
	}

	return &components, nil
}

func (component *Component) populateXPathSimpleParts(node types.Node) error {
	var err error

	xp := &component.Parts.XPath.Simple

	xp.Address.Query = "//address/p"
	xp.Address.Values, err = getValuesForXPathQuery(xp.Address.Query, node)
	if err != nil {
		return err
	}

	xp.Appraisal.Query = "//appraisal/p"
	xp.Appraisal.Values, err = getValuesForXPathQuery(xp.Appraisal.Query, node)
	if err != nil {
		return err
	}

	xp.BiogHist.Query = "//bioghist/p"
	xp.BiogHist.Values, err = getValuesForXPathQuery(xp.BiogHist.Query, node)
	if err != nil {
		return err
	}

	xp.ChronList.Query = "//chronlist/chronitem//text()"
	xp.ChronList.Values, err = getValuesForXPathQuery(xp.ChronList.Query, node)
	if err != nil {
		return err
	}

	xp.Collection.Query = "//archdesc/did/unittitle"
	xp.Collection.Values, err = getValuesForXPathQuery(xp.Collection.Query, node)
	if err != nil {
		return err
	}

	xp.CollectionUnitID.Query = "//archdesc/did/unitid"
	xp.CollectionUnitID.Values, err = getValuesForXPathQuery(xp.CollectionUnitID.Query, node)
	if err != nil {
		return err
	}

	xp.Corpname.Query = "//corpname"
	xp.Corpname.Values, err = getValuesForXPathQuery(xp.Corpname.Query, node)
	if err != nil {
		return err
	}

	xp.Creator.Query = "//archdesc[@level='collection']/did/origination[@label='creator']/*[name() = 'corpname' or name() = 'famname' or name() = 'persname']"
	xp.Creator.Values, err = getValuesForXPathQuery(xp.Creator.Query, node)
	if err != nil {
		return err
	}

	xp.DAO.Query = "//dao/daodesc/p"
	xp.DAO.Values, err = getValuesForXPathQuery(xp.DAO.Query, node)
	if err != nil {
		return err
	}

	xp.DIDUnitID.Query = "//did/unitid"
	xp.DIDUnitID.Values, err = getValuesForXPathQuery(xp.DIDUnitID.Query, node)
	if err != nil {
		return err
	}

	xp.DIDUnitTitle.Query = "//did/unittitle"
	xp.DIDUnitTitle.Values, err = getValuesForXPathQuery(xp.DIDUnitTitle.Query, node)
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

	xp.Function.Query = "//function"
	xp.Function.Values, err = getValuesForXPathQuery(xp.Function.Query, node)
	if err != nil {
		return err
	}

	xp.GenreForm.Query = "//genreform"
	xp.GenreForm.Values, err = getValuesForXPathQuery(xp.GenreForm.Query, node)
	if err != nil {
		return err
	}

	xp.GeogName.Query = "//geogname"
	xp.GeogName.Values, err = getValuesForXPathQuery(xp.GeogName.Query, node)
	if err != nil {
		return err
	}

	xp.Heading.Query = "//archdesc[@level='collection']/did/unittitle"
	xp.Heading.Values, err = getValuesForXPathQuery(xp.Heading.Query, node)
	if err != nil {
		return err
	}

	xp.Language.Query = "//did/langmaterial/language/@langcode"
	xp.Language.Values, err = getValuesForXPathQuery(xp.Language.Query, node)
	if err != nil {
		return err
	}

	xp.Level.Query = "///c/@level"
	xp.Level.Values, err = getValuesForXPathQuery(xp.Level.Query, node)
	if err != nil {
		return err
	}

	xp.Name.Query = "//name"
	xp.Name.Values, err = getValuesForXPathQuery(xp.Name.Query, node)
	if err != nil {
		return err
	}

	xp.Note.Query = "//note"
	xp.Note.Values, err = getValuesForXPathQuery(xp.Note.Query, node)
	if err != nil {
		return err
	}

	xp.Occupation.Query = "//occupation"
	xp.Occupation.Values, err = getValuesForXPathQuery(xp.Occupation.Query, node)
	if err != nil {
		return err
	}

	xp.PersName.Query = "//persname"
	xp.PersName.Values, err = getValuesForXPathQuery(xp.PersName.Query, node)
	if err != nil {
		return err
	}

	xp.PhysTech.Query = "//phystech/p"
	xp.PhysTech.Values, err = getValuesForXPathQuery(xp.PhysTech.Query, node)
	if err != nil {
		return err
	}

	xp.Ref.Query = "///c/@id"
	xp.Ref.Values, err = getValuesForXPathQuery(xp.Ref.Query, node)
	if err != nil {
		return err
	}

	xp.ScopeContent.Query = "//scopecontent/p"
	xp.ScopeContent.Values, err = getValuesForXPathQuery(xp.ScopeContent.Query, node)
	if err != nil {
		return err
	}

	xp.Subject.Query = "//subject"
	xp.Subject.Values, err = getValuesForXPathQuery(xp.Subject.Query, node)
	if err != nil {
		return err
	}

	xp.SubjectOrFunctionOrOccupation.Query = "//*[local-name()='subject' or local-name()='function' or local-name() = 'occupation']"
	xp.SubjectOrFunctionOrOccupation.Values, err = getValuesForXPathQuery(xp.SubjectOrFunctionOrOccupation.Query, node)
	if err != nil {
		return err
	}

	xp.Title.Query = "//title"
	xp.Title.Values, err = getValuesForXPathQuery(xp.Title.Query, node)
	if err != nil {
		return err
	}

	xp.UnitDateNotType.Query = "//did/unitdate[not(@type)]"
	xp.UnitDateNotType.Values, err = getValuesForXPathQuery(xp.UnitDateNotType.Query, node)
	if err != nil {
		return err
	}

	xp.UnitDateBulk.Query = "//archdesc[@level='collection']/did/unitdate[@type='bulk']"
	xp.UnitDateBulk.Values, err = getValuesForXPathQuery(xp.UnitDateBulk.Query, node)
	if err != nil {
		return err
	}

	xp.UnitDateNormal.Query = "//did/unitdate/@normal"
	xp.UnitDateNormal.Values, err = getValuesForXPathQuery(xp.UnitDateNormal.Query, node)
	if err != nil {
		return err
	}

	xp.UnitdateInclusive.Query = "//archdesc[@level='collection']/did/unitdate[@type='inclusive']"
	xp.UnitdateInclusive.Values, err = getValuesForXPathQuery(xp.UnitdateInclusive.Query, node)
	if err != nil {
		return err
	}

	return nil
}

func MakeComponent(node types.Node) (Component, error) {
	component := Component{}

	err := component.populateXPathSimpleParts(node)
	if err != nil {
		return component, err
	}

	component.ID = component.Parts.XPath.Simple.EADID.Values[0] +
		component.Parts.XPath.Simple.Ref.Values[0]

	return component, nil
}
