package component

import (
	"errors"
	"fmt"
	"github.com/lestrrat-go/libxml2/types"
	"go-ead-indexer/pkg/ead/eadutil"
	"regexp"
)

type Component struct {
	ID             string         `json:"id"`
	IDAttribute    string         `json:"IDAttribute"`
	Parts          ComponentParts `json:"parts"`
	SolrAddMessage SolrAddMessage `json:"solr_add_message"`
}

// For now, no struct tags for the `Component*` fields.  Keep it flat.
type ComponentParts struct {
	Collection       string `json:"collection"`
	CollectionUnitID string `json:"collection_unit_id"`
	ComponentComplexParts
	ComponentHierarchyParts
	ComponentXPathParts
	Containers     []Container   `json:"containers"`
	RepositoryCode ComponentPart `json:"repository_code"`
}

type ComponentComplexParts struct {
	ChronListText    ComponentPart `json:"chron_list_text"`
	CreatorComplex   ComponentPart `json:"creator_complex"`
	DAO              ComponentPart `json:"dao"`
	DateRange        ComponentPart `json:"date_range"`
	Format           ComponentPart `json:"format"`
	Heading          ComponentPart `json:"heading"`
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
	ComponentChildren     bool          `json:"component_children"`
	ComponentLevel        int           `json:"component_level"`
	ParentForDisplay      ComponentPart `json:"parent_for_display"`
	ParentForSort         ComponentPart `json:"parent_for_sort"`
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
	UnitDateNoTypeAttribute       ComponentPart `json:"unit_date_not_type"`
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

var unitDateOpenTagRegExp = regexp.MustCompile("^<unitdate[^>]*>")

// See `ead.new()` comment on why we have to pass in `repositoryCode` as an argument.
func MakeComponents(repositoryCode string, collection string, collectionUnitID string,
	node types.Node) (*[]Component, error) {
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

	ancestorUnitTitleListMap, err := makeAncestorUnitTitleListMap(node)
	if err != nil {
		return nil, err
	}

	components := []Component{}
	for _, cNode := range cNodes {
		newComponent, err := MakeComponent(repositoryCode, collection,
			collectionUnitID, cNode)
		if err != nil {
			return &components, err
		}

		if ancestorUnitTitleList, ok :=
			ancestorUnitTitleListMap[newComponent.IDAttribute]; ok {
			newComponent.Parts.AncestorUnitTitleList = ancestorUnitTitleList
		}

		components = append(components, newComponent)
	}

	return &components, nil
}

// See `ead.new()` comment on why we have to pass in `repositoryCode` as an argument.
func MakeComponent(repositoryCode string, collection string, collectionUnitID string,
	node types.Node) (Component, error) {
	component := Component{
		Parts: ComponentParts{
			Collection:       collection,
			CollectionUnitID: collectionUnitID,
			RepositoryCode: ComponentPart{
				Values: []string{repositoryCode},
			},
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
	// These deleted <c> nodes are later passed into this method when by the
	// outer loop in which this method is called.
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

func getAncestorUnitTitle(node types.Node) (string, error) {
	var ancestorUnitTitle string

	// Try `<unittitle>` first, then try <unitdate>, and if neither worked, return
	// no title available.
	xpathResult, err := node.Find("did/unittitle")
	if err != nil {
		return ancestorUnitTitle, err
	}

	unitTitleNodes := xpathResult.NodeList()
	if len(unitTitleNodes) > 0 {
		unitTitleXMLString := unitTitleNodes[0].String()
		unitTitleContents := eadutil.StripOpenAndCloseTags(unitTitleXMLString)
		// TODO: DLFA-243
		// Replace this with `util.IsNonEmptyString(unitDateContents)`
		if unitTitleContents != "" {
			ancestorUnitTitle = unitTitleContents
		}
	}

	// <unittitle> didn't work.  Try <unitdate>.
	if ancestorUnitTitle == "" {
		xpathResult, err := node.Find("did/unitdate")
		if err != nil {
			return ancestorUnitTitle, err
		}

		unitDateNodes := xpathResult.NodeList()
		if len(unitDateNodes) > 0 {
			unitDateXMLString := unitDateNodes[0].String()
			unitDateContents := eadutil.StripOpenAndCloseTags(unitDateXMLString)
			// TODO: DLFA-243
			// Replace this with `util.IsNonEmptyString(unitDateContents)`
			if unitDateContents != "" {
				ancestorUnitTitle = unitDateContents
			}
		}
	}

	// Can't create a title.
	if ancestorUnitTitle == "" {
		return "[No title available]", nil
	}

	ancestorUnitTitleHTML, err := eadutil.MakeTitleHTML(ancestorUnitTitle)
	if err != nil {
		return ancestorUnitTitle, err
	}

	return ancestorUnitTitleHTML, nil
}

func makeAncestorUnitTitleListMap(node types.Node) (map[string][]string, error) {
	ancestorUnitTitleListMap := map[string][]string{}

	dscNode, err := eadutil.GetFirstNodeForXPathQuery("//dsc", node)
	if err != nil {
		return ancestorUnitTitleListMap, err
	}

	if dscNode == nil {
		return ancestorUnitTitleListMap,
			errors.New("makeAncestorUnitTitleListMap() error: no <dsc> element found")
	}

	childCNodes, err := eadutil.GetNodeListForXPathQuery("child::c", dscNode)
	if err != nil {
		return ancestorUnitTitleListMap, err
	}

	if len(childCNodes) == 0 {
		return ancestorUnitTitleListMap, nil
	}

	for _, cNode := range childCNodes {
		ancestorUnitTitleList := []string{}
		err = makeAncestorUnitTitleListMap_add(ancestorUnitTitleListMap,
			ancestorUnitTitleList, cNode)
		if err != nil {
			return ancestorUnitTitleListMap, err
		}
	}

	return ancestorUnitTitleListMap, nil
}

func makeAncestorUnitTitleListMap_add(ancestorUnitTitleListMap map[string][]string,
	ancestorUnitTitleList []string, node types.Node) error {

	ancestorUnitTitle, err := getAncestorUnitTitle(node)
	if err != nil {
		return err
	}
	ancestorUnitTitleList = append(ancestorUnitTitleList, ancestorUnitTitle)

	idNode, err := node.(types.Element).GetAttribute("id")
	if err != nil {
		return errors.New(
			fmt.Sprintf("makeAncestorUnitTitleListMap_add() error: can't get `id` attribute of node: %s",
				node.String()))
	}

	ancestorUnitTitleListMap[idNode.NodeValue()] = ancestorUnitTitleList

	childCNodes, err := eadutil.GetNodeListForXPathQuery("child::c", node)
	if err != nil {
		return err
	}

	if len(childCNodes) == 0 {
		return nil
	}

	for _, childCNode := range childCNodes {
		err = makeAncestorUnitTitleListMap_add(ancestorUnitTitleListMap,
			ancestorUnitTitleList, childCNode)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeChildCNodes(node types.Node) error {
	return eadutil.RemoveChildNodesMatchingName(node, CElementName)
}
