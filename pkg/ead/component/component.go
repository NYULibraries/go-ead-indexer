package component

import (
	"errors"
	"fmt"
	"github.com/lestrrat-go/libxml2/types"
	"go-ead-indexer/pkg/ead/eadutil"
)

type Component struct {
	ID             string         `json:"id"`
	IDAttribute    string         `json:"IDAttribute"`
	Parts          ComponentParts `json:"parts"`
	SolrAddMessage SolrAddMessage `json:"solr_add_message"`
}

// For now, no struct tags for the `Component*` fields.  Keep it flat.
type ComponentParts struct {
	ComponentCollectionDocParts
	ComponentComplexParts
	ComponentHierarchyParts
	ComponentXPathParts
	Containers []Container `json:"containers"`

	Sort int `json:"sort"`
}

type ComponentCollectionDocParts struct {
	// TODO: DLFA-238
	// We collect this Author information but don't include it in `SolrAddMessage`.
	// See: "Solr field `author` in Component Solr doc is never populated"
	// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=8577864&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-8577864
	// After passing DLFA-201 add `author` to `SolrAddMessage`.
	Author           []string `json:"author"`
	Collection       string   `json:"collection"`
	CollectionUnitID string   `json:"collection_unit_id"`
	RepositoryCode   string   `json:"repository_code"`
}

type ComponentComplexParts struct {
	ChronListComplex ComponentPart `json:"chron_list_complex"`
	CreatorComplex   ComponentPart `json:"creator_complex"`
	DAO              ComponentPart `json:"dao"`
	DateRange        ComponentPart `json:"date_range"`
	Format           ComponentPart `json:"format"`
	Heading          ComponentPart `json:"heading"`
	Language         ComponentPart `json:"language"`
	Location         ComponentPart `json:"location"`
	MaterialType     ComponentPart `json:"material_type"`
	Name             ComponentPart `json:"name"`
	Place            ComponentPart `json:"place"`
	SubjectForFacets ComponentPart `json:"subject_for_facets"`
	UnitDateDisplay  ComponentPart `json:"unit_date_display"`
	UnitDateEnd      ComponentPart `json:"unit_date_end"`
	UnitDateStart    ComponentPart `json:"unit_date_start"`
	UnitTitleHTML    ComponentPart `json:"unit_title_html"`
}

type ComponentHierarchyParts struct {
	AncestorUnitTitleList []string      `json:"ancestor_unit_title_list"`
	ComponentChildren     string        `json:"component_children"`
	ComponentLevel        string        `json:"component_level"`
	ParentForDisplay      ComponentPart `json:"parent_for_display"`
	ParentForSort         string        `json:"parent_for_sort"`
	SeriesForSort         string        `json:"series_for_sort"`
}

type ComponentXPathParts struct {
	Address                       ComponentPart `json:"address"`
	Appraisal                     ComponentPart `json:"appraisal"`
	BiogHist                      ComponentPart `json:"biog_hist"`
	ChronList                     ComponentPart `json:"chron_list"`
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
	LangCode                      ComponentPart `json:"language"`
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
	UnitDateBulk                  ComponentPart `json:"unit_date_bulk"`
	UnitDateInclusive             ComponentPart `json:"unit_date_inclusive"`
	UnitDateNoTypeAttribute       ComponentPart `json:"unit_date_not_type"`
	UnitDateNormal                ComponentPart `json:"unit_date_normal"`
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

func MakeComponents(collectionDocParts ComponentCollectionDocParts, node types.Node) (*[]Component, error) {
	xpathResult, err := node.Find("//" + CElementName)
	if err != nil {
		return nil, err
	}
	defer xpathResult.Free()

	// Note: can't do `&xpathResult.NodeList()`
	// See https://groups.google.com/g/golang-nuts/c/reaIlFdibWU
	cNodes := xpathResult.NodeList()

	// Early exit
	if len(cNodes) == 0 {
		return nil, nil
	}

	ancestorUnitTitleListMap, err := makeAncestorUnitTitleListMap(node)
	if err != nil {
		return nil, err
	}

	components := []Component{}
	sort := 0
	for _, cNode := range cNodes {
		sort += 1
		newComponent, err := MakeComponent(collectionDocParts, sort, cNode)
		if err != nil {
			return &components, err
		}

		if ancestorUnitTitleList, ok :=
			ancestorUnitTitleListMap[newComponent.IDAttribute]; ok {
			newComponent.Parts.AncestorUnitTitleList = ancestorUnitTitleList
			err = newComponent.setSeriesForSort()
			if err != nil {
				return &components, err
			}
		}

		// This depends on `newComponent.Parts.AncestorUnitTitleList`
		newComponent.setSolrAddMessage()

		components = append(components, newComponent)
	}

	return &components, nil
}

func MakeComponent(collectionDocParts ComponentCollectionDocParts, sort int,
	node types.Node) (Component, error) {
	component := Component{
		Parts: ComponentParts{
			ComponentCollectionDocParts: collectionDocParts,
			Sort:                        sort,
		},
	}

	idNode, err := node.(types.Element).GetAttribute("id")
	if err != nil {
		return component, errors.New(
			fmt.Sprintf("Can't get `id` attribute of <c> element: %s", node.String()))
	}
	component.IDAttribute = idNode.NodeValue()

	err = component.setParts(node)
	if err != nil {
		return component, err
	}

	component.ID = component.Parts.EADID.Values[0] + component.IDAttribute

	return component, nil
}

func (component *Component) setParts(node types.Node) error {
	err := component.setHierarchyDataParts(node)
	if err != nil {
		return err
	}

	// Create a copy of the node with child <c> nodes for subsequent processing
	// to prevent duplication from overlapping node trees.  Defensive copying is
	// necessary because some processes require access to the parent node data,
	// which `removeChildCNodes()` deletes from the child <c> nodes deleted.
	// These deleted <c> nodes are later passed into this method by the outer
	// loop in which this method is called.
	// As a rule, we would want to take care of all processing which requires
	// parent node data before getting here, in which case we could keep using
	// the original node, but there's no harm in being careful, and having the
	// original node on hand for before/after comparison could come in handy.
	nodeCopyWithChildCNodesRemoved, err := node.Copy()
	if err != nil {
		return err
	}

	err = removeChildCNodes(nodeCopyWithChildCNodesRemoved)
	if err != nil {
		return err
	}

	err = component.setXPathSimpleParts(nodeCopyWithChildCNodesRemoved)
	if err != nil {
		return err
	}

	err = component.setContainersPart(nodeCopyWithChildCNodesRemoved)
	if err != nil {
		return err
	}

	err = component.setComplexParts()
	if err != nil {
		return err
	}

	return nil
}

func removeChildCNodes(node types.Node) error {
	return eadutil.RemoveChildNodesMatchingName(node, CElementName)
}
