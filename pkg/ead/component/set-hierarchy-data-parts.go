package component

import (
	"github.com/lestrrat-go/libxml2/types"
	"go-ead-indexer/pkg/ead/eadutil"
	"slices"
	"strings"
)

func (component *Component) setHierarchyDataParts(node types.Node) error {
	err := component.setComponentChildren(node)
	if err != nil {
		return err
	}

	// It is not possible for a <c> node to have no parent.  If we can't get the
	// parent node for the `node` arg, there's no point doing any of the processing
	// in this method.
	_, err = node.ParentNode()
	if err != nil {
		return err
	}

	// Shouldn't be possible for this to return an error given the early error
	// above, but check anyway.
	err = component.setParentForSort(node)
	if err != nil {
		return err
	}

	component.setParentForDisplay(node)
	// Depends on `Component.setParentForDisplay()`
	component.setComponentLevel()

	return nil
}

func (component *Component) setComponentChildren(node types.Node) error {
	childNodes, err := node.ChildNodes()
	if err != nil {
		return err
	}

	component.Parts.ComponentChildren = slices.ContainsFunc(childNodes, func(node types.Node) bool {
		return node.NodeName() == CElementName
	})

	return nil
}

// Depends on `Component.setParentForDisplay()`
func (component *Component) setComponentLevel() {
	component.Parts.ComponentLevel = len(component.Parts.ParentForDisplay.Values) + 1
}

func (component *Component) setParentForDisplay(node types.Node) {
	parentIDList := []string{}

	currentNode := node
	for {
		parentNode, err := currentNode.ParentNode()
		if err != nil {
			break
		}

		parentNodeIDAttributeNode, err := parentNode.(types.Element).GetAttribute("id")
		if err == nil {
			parentIDList = append(parentIDList, parentNodeIDAttributeNode.Value())
		} else {
			break
		}

		currentNode = parentNode
	}

	slices.Reverse(parentIDList)

	component.Parts.ParentForDisplay.Values = parentIDList
}

func (component *Component) setParentForSort(node types.Node) error {
	parentNode, err := node.ParentNode()
	if err != nil {
		return err
	}

	parentNodeIDAttributeNode, err := parentNode.(types.Element).GetAttribute("id")
	if err == nil {
		component.Parts.ParentForSort = parentNodeIDAttributeNode.Value()
	} else {
		// Parent is <dsc>, which has no `id` attribute.
	}

	return nil
}

func (component *Component) setSeriesForSort() error {
	parts := &component.Parts

	ancestorUnitTitleHTMLValues := []string{}
	for _, ancestorUnitTitle := range parts.AncestorUnitTitleList {
		ancestorUnitTitleHTMLValue, err := eadutil.MakeTitleHTML(ancestorUnitTitle)
		if err != nil {
			return err
		}

		ancestorUnitTitleHTMLValues = append(ancestorUnitTitleHTMLValues,
			ancestorUnitTitleHTMLValue)
	}

	parts.SeriesForSort = strings.Join(ancestorUnitTitleHTMLValues, " >> ")

	return nil
}
