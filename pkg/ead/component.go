package ead

import "github.com/lestrrat-go/libxml2/types"

type Component struct {
	ID             string         `json:"id"`
	Parts          ComponentParts `json:"parts"`
	SolrAddMessage string         `json:"solr_add_message"`
}

// For now, no struct tags for the `Component*` fields.  Keep it flat.
type ComponentParts struct {
	ComponentXPathDirectQueryParts
	RepositoryCode ComponentPart `json:"repository_code"`
}

type ComponentXPathDirectQueryParts struct {
	Address                       ComponentPart `json:"address"`
	Appraisal                     ComponentPart `json:"appraisal"`
	BiogHist                      ComponentPart `json:"biog_hist"`
	ChronList                     ComponentPart `json:"chron_list"`
	Collection                    ComponentPart `json:"collection"`
	CollectionUnitID              ComponentPart `json:"collection_unit_id"`
	Corpname                      ComponentPart `json:"corpname"`
	Creator                       ComponentPart `json:"creator"`
	DAO                           ComponentPart `json:"dao"`
	DIDUnitID                     ComponentPart `json:"did_unit_id"`
	DIDUnitTitle                  ComponentPart `json:"did_unit_title"`
	EADID                         ComponentPart `json:"eadid"`
	FamName                       ComponentPart `json:"fam_name"`
	Function                      ComponentPart `json:"function"`
	GenreForm                     ComponentPart `json:"genre_form"`
	GeogName                      ComponentPart `json:"geog_name"`
	Heading                       ComponentPart `json:"heading"`
	Language                      ComponentPart `json:"language"`
	Level                         ComponentPart `json:"level"`
	Name                          ComponentPart `json:"name"`
	Note                          ComponentPart `json:"note"`
	Occupation                    ComponentPart `json:"occupation"`
	PersName                      ComponentPart `json:"pers_name"`
	PhysTech                      ComponentPart `json:"phys_tech"`
	Ref                           ComponentPart `json:"ref"`
	ScopeContent                  ComponentPart `json:"scope_content"`
	Subject                       ComponentPart `json:"subject"`
	SubjectOrFunctionOrOccupation ComponentPart `json:"subject_or_function_or_occupation"`
	Title                         ComponentPart `json:"title"`
	UnitDateNotType               ComponentPart `json:"unit_date_not_type"`
	UnitDateBulk                  ComponentPart `json:"unit_date_bulk"`
	UnitDateNormal                ComponentPart `json:"unit_date_normal"`
	UnitDateInclusive             ComponentPart `json:"unit_date_inclusive"`
}

type ComponentPart struct {
	Source     string   `json:"source"`
	Values     []string `json:"values"`
	XMLStrings []string `json:"xml_strings"`
}

// See `ead.new()` comment on why we have to pass in `repositoryCode` as an argument.
func MakeComponents(repositoryCode string, node types.Node) (*[]Component, error) {
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
		newComponent, err := MakeComponent(repositoryCode, cNode)
		if err != nil {
			return &components, err
		}

		components = append(components, newComponent)
	}

	return &components, nil
}

func (component *Component) setXPathDirectQueryParts(node types.Node) error {
	var err error

	parts := &component.Parts

	parts.Address.Source = "//address/p"
	parts.Address.Values, parts.Address.XMLStrings, err = getValuesForXPathQuery(parts.Address.Source, node)
	if err != nil {
		return err
	}

	parts.Appraisal.Source = "//appraisal/p"
	parts.Appraisal.Values, parts.Appraisal.XMLStrings, err = getValuesForXPathQuery(parts.Appraisal.Source, node)
	if err != nil {
		return err
	}

	parts.BiogHist.Source = "//bioghist/p"
	parts.BiogHist.Values, parts.BiogHist.XMLStrings, err = getValuesForXPathQuery(parts.BiogHist.Source, node)
	if err != nil {
		return err
	}

	parts.ChronList.Source = "//chronlist/chronitem//text()"
	parts.ChronList.Values, parts.ChronList.XMLStrings, err = getValuesForXPathQuery(parts.ChronList.Source, node)
	if err != nil {
		return err
	}

	parts.Collection.Source = "//archdesc/did/unittitle"
	parts.Collection.Values, parts.Collection.XMLStrings, err = getValuesForXPathQuery(parts.Collection.Source, node)
	if err != nil {
		return err
	}

	parts.CollectionUnitID.Source = "//archdesc/did/unitid"
	parts.CollectionUnitID.Values, parts.CollectionUnitID.XMLStrings, err = getValuesForXPathQuery(parts.CollectionUnitID.Source, node)
	if err != nil {
		return err
	}

	parts.Corpname.Source = "//corpname"
	parts.Corpname.Values, parts.Corpname.XMLStrings, err = getValuesForXPathQuery(parts.Corpname.Source, node)
	if err != nil {
		return err
	}

	parts.Creator.Source = "//archdesc[@level='collection']/did/origination[@label='creator']/*[name() = 'corpname' or name() = 'famname' or name() = 'persname']"
	parts.Creator.Values, parts.Creator.XMLStrings, err = getValuesForXPathQuery(parts.Creator.Source, node)
	if err != nil {
		return err
	}

	parts.DAO.Source = "//dao/daodesc/p"
	parts.DAO.Values, parts.DAO.XMLStrings, err = getValuesForXPathQuery(parts.DAO.Source, node)
	if err != nil {
		return err
	}

	parts.DIDUnitID.Source = "//did/unitid"
	parts.DIDUnitID.Values, parts.DIDUnitID.XMLStrings, err = getValuesForXPathQuery(parts.DIDUnitID.Source, node)
	if err != nil {
		return err
	}

	parts.DIDUnitTitle.Source = "//did/unittitle"
	parts.DIDUnitTitle.Values, parts.DIDUnitTitle.XMLStrings, err = getValuesForXPathQuery(parts.DIDUnitTitle.Source, node)
	if err != nil {
		return err
	}

	parts.EADID.Source = "//eadid"
	parts.EADID.Values, parts.EADID.XMLStrings, err = getValuesForXPathQuery(parts.EADID.Source, node)
	if err != nil {
		return err
	}

	parts.FamName.Source = "//famname"
	parts.FamName.Values, parts.FamName.XMLStrings, err = getValuesForXPathQuery(parts.FamName.Source, node)
	if err != nil {
		return err
	}

	parts.Function.Source = "//function"
	parts.Function.Values, parts.Function.XMLStrings, err = getValuesForXPathQuery(parts.Function.Source, node)
	if err != nil {
		return err
	}

	parts.GenreForm.Source = "//genreform"
	parts.GenreForm.Values, parts.GenreForm.XMLStrings, err = getValuesForXPathQuery(parts.GenreForm.Source, node)
	if err != nil {
		return err
	}

	parts.GeogName.Source = "//geogname"
	parts.GeogName.Values, parts.GeogName.XMLStrings, err = getValuesForXPathQuery(parts.GeogName.Source, node)
	if err != nil {
		return err
	}

	parts.Heading.Source = "//archdesc[@level='collection']/did/unittitle"
	parts.Heading.Values, parts.Heading.XMLStrings, err = getValuesForXPathQuery(parts.Heading.Source, node)
	if err != nil {
		return err
	}

	parts.Language.Source = "//did/langmaterial/language/@langcode"
	parts.Language.Values, parts.Language.XMLStrings, err = getValuesForXPathQuery(parts.Language.Source, node)
	if err != nil {
		return err
	}

	parts.Level.Source = "///c/@level"
	parts.Level.Values, parts.Level.XMLStrings, err = getValuesForXPathQuery(parts.Level.Source, node)
	if err != nil {
		return err
	}

	parts.Name.Source = "//name"
	parts.Name.Values, parts.Name.XMLStrings, err = getValuesForXPathQuery(parts.Name.Source, node)
	if err != nil {
		return err
	}

	parts.Note.Source = "//note"
	parts.Note.Values, parts.Note.XMLStrings, err = getValuesForXPathQuery(parts.Note.Source, node)
	if err != nil {
		return err
	}

	parts.Occupation.Source = "//occupation"
	parts.Occupation.Values, parts.Occupation.XMLStrings, err = getValuesForXPathQuery(parts.Occupation.Source, node)
	if err != nil {
		return err
	}

	parts.PersName.Source = "//persname"
	parts.PersName.Values, parts.PersName.XMLStrings, err = getValuesForXPathQuery(parts.PersName.Source, node)
	if err != nil {
		return err
	}

	parts.PhysTech.Source = "//phystech/p"
	parts.PhysTech.Values, parts.PhysTech.XMLStrings, err = getValuesForXPathQuery(parts.PhysTech.Source, node)
	if err != nil {
		return err
	}

	parts.Ref.Source = "///c/@id"
	parts.Ref.Values, parts.Ref.XMLStrings, err = getValuesForXPathQuery(parts.Ref.Source, node)
	if err != nil {
		return err
	}

	parts.ScopeContent.Source = "//scopecontent/p"
	parts.ScopeContent.Values, parts.ScopeContent.XMLStrings, err = getValuesForXPathQuery(parts.ScopeContent.Source, node)
	if err != nil {
		return err
	}

	parts.Subject.Source = "//subject"
	parts.Subject.Values, parts.Subject.XMLStrings, err = getValuesForXPathQuery(parts.Subject.Source, node)
	if err != nil {
		return err
	}

	parts.SubjectOrFunctionOrOccupation.Source = "//*[local-name()='subject' or local-name()='function' or local-name() = 'occupation']"
	parts.SubjectOrFunctionOrOccupation.Values, parts.SubjectOrFunctionOrOccupation.XMLStrings, err = getValuesForXPathQuery(parts.SubjectOrFunctionOrOccupation.Source, node)
	if err != nil {
		return err
	}

	parts.Title.Source = "//title"
	parts.Title.Values, parts.Title.XMLStrings, err = getValuesForXPathQuery(parts.Title.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateNotType.Source = "//did/unitdate[not(@type)]"
	parts.UnitDateNotType.Values, parts.UnitDateNotType.XMLStrings, err = getValuesForXPathQuery(parts.UnitDateNotType.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateBulk.Source = "//archdesc[@level='collection']/did/unitdate[@type='bulk']"
	parts.UnitDateBulk.Values, parts.UnitDateBulk.XMLStrings, err = getValuesForXPathQuery(parts.UnitDateBulk.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateNormal.Source = "//did/unitdate/@normal"
	parts.UnitDateNormal.Values, parts.UnitDateNormal.XMLStrings, err = getValuesForXPathQuery(parts.UnitDateNormal.Source, node)
	if err != nil {
		return err
	}

	parts.UnitDateInclusive.Source = "//archdesc[@level='collection']/did/unitdate[@type='inclusive']"
	parts.UnitDateInclusive.Values, parts.UnitDateInclusive.XMLStrings, err = getValuesForXPathQuery(parts.UnitDateInclusive.Source, node)
	if err != nil {
		return err
	}

	return nil
}

// See `ead.new()` comment on why we have to pass in `repositoryCode` as an argument.
func MakeComponent(repositoryCode string, node types.Node) (Component, error) {
	component := Component{
		Parts: ComponentParts{
			RepositoryCode: ComponentPart{
				Values: []string{repositoryCode},
			},
		},
	}

	err := component.setXPathDirectQueryParts(node)
	if err != nil {
		return component, err
	}

	component.ID = component.Parts.EADID.Values[0] +
		component.Parts.Ref.Values[0]

	return component, nil
}
