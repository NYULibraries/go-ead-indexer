package component

import (
	"github.com/lestrrat-go/libxml2/types"
	"go-ead-indexer/pkg/ead/util"
)

type Component struct {
	ID             string         `json:"id"`
	Parts          ComponentParts `json:"parts"`
	SolrAddMessage SolrAddMessage `json:"solr_add_message"`
}

// For now, no struct tags for the `Component*` fields.  Keep it flat.
type ComponentParts struct {
	ComponentComplexParts
	ComponentXPathParts
	Containers     []Container   `json:"containers"`
	RepositoryCode ComponentPart `json:"repository_code"`
}

type ComponentComplexParts struct {
	CreatorComplex   ComponentPart `json:"creator_complex"`
	DAO              ComponentPart `json:"dao"`
	Format           ComponentPart `json:"format"`
	Location         ComponentPart `json:"location"`
	Name             ComponentPart `json:"name"`
	Place            ComponentPart `json:"place"`
	SubjectForFacets ComponentPart `json:"subject_for_facets"`
}

type ComponentXPathParts struct {
	Address                       ComponentPart `json:"address"`
	Appraisal                     ComponentPart `json:"appraisal"`
	BiogHist                      ComponentPart `json:"biog_hist"`
	ChronList                     ComponentPart `json:"chron_list"`
	Collection                    ComponentPart `json:"collection"`
	CollectionUnitID              ComponentPart `json:"collection_unit_id"`
	CorpName                      ComponentPart `json:"corpname"`
	CorpNameNotInRepository       ComponentPart `json:"corp_name_not_in_repository"`
	Creator                       ComponentPart `json:"creator"`
	CreatorCorpName               ComponentPart `json:"creator_corp_name"`
	CreatorFamName                ComponentPart `json:"creator_fam_name"`
	CreatorPersName               ComponentPart `json:"creator_pers_name"`
	DAODescriptionParagraph       ComponentPart `json:"dao"`
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
	NameElementAll                ComponentPart `json:"name_element_all"`
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

type Container struct {
	ID        string `json:"id"`
	Parent    string `json:"parent"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	XMLString string `json:"xmlstring"`
}

const CElementName = "c"

// See `ead.new()` comment on why we have to pass in `repositoryCode` as an argument.
func MakeComponents(repositoryCode string, node types.Node) (*[]Component, error) {
	xpathResult, err := node.Find("//" + CElementName)
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
		// Remove child <c> nodes from to prevent duplication from overlapping
		// node trees.
		cNodeChildCNodesRemoved, err := removeChildCNodes(cNode)
		if err != nil {
			return &components, err
		}

		newComponent, err := MakeComponent(repositoryCode, cNodeChildCNodesRemoved)
		if err != nil {
			return &components, err
		}

		components = append(components, newComponent)
	}

	return &components, nil
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

	err := component.setParts(node)
	if err != nil {
		return component, err
	}

	component.ID = component.Parts.EADID.Values[0] +
		component.Parts.Ref.Values[0]

	return component, nil
}

func (component *Component) setParts(node types.Node) error {
	err := component.setXPathSimpleParts(node)
	if err != nil {
		return err
	}

	err = component.setContainersPart(node)
	if err != nil {
		return err
	}

	err = component.setComplexParts()
	if err != nil {
		return err
	}

	return nil
}

// TODO: `removeChildCNodes()` adds `xmlns:xlink="http://www.w3.org/1999/xlink"`
// to the <c>.  Should we leave it, or strip it?  It's added in the `resultNode`
// defensive copy; it doesn't happen happen when `node` is mutated directly.
func removeChildCNodes(node types.Node) (types.Node, error) {
	return util.RemoveChildNodesMatchingName(node, CElementName)
}
