package component

import (
	"slices"

	"github.com/lestrrat-go/libxml2/types"
	"github.com/nyulibraries/go-ead-indexer/pkg/ead/eadutil"
)

type componentsOnlineAccessMap map[string]map[string]int

func incrementOnlineAccessMapEntryCount(onlineAccessMap map[string]map[string]int,
	componentID string, onlineAccessDirect string) {
	_, ok := onlineAccessMap[componentID]
	if ok {
		_, ok := onlineAccessMap[componentID][onlineAccessDirect]
		if ok {
			onlineAccessMap[componentID][onlineAccessDirect]++
		} else {
			onlineAccessMap[componentID][onlineAccessDirect] = 1
		}
	} else {
		onlineAccessMap[componentID] = map[string]int{}
		onlineAccessMap[componentID][onlineAccessDirect] = 1
	}
}

// This map calculates the rollup numbers that are needed for the Digital
// Content facet in the discovery portal UI.  The rollup numbers are calculated
// by tallying the direct totals and then adding that tally to the indirect
// totals of every single ancestor node.
func makeOnlineAccessCountsMap(components []Component) componentsOnlineAccessMap {
	onlineAccessMap := componentsOnlineAccessMap{}
	for _, component := range components {
		for _, onlineAccessDirect := range component.Parts.OnlineAccessDirect.Values {
			incrementOnlineAccessMapEntryCount(onlineAccessMap, component.IDAttribute,
				onlineAccessDirect)

			for _, parentID := range component.Parts.ComponentHierarchyParts.ParentForDisplay.Values {
				incrementOnlineAccessMapEntryCount(onlineAccessMap, parentID,
					onlineAccessDirect)
			}
		}
	}

	return onlineAccessMap
}

// Online access list for the node itself, with nothing rolled up from child
// nodes (which should have already been removed from the node before it was
// passed in).  This list is used when calculating the rolled up lists/counts
// for each node in a final pass after all node lists have been compiled.
func (component *Component) setOnlineAccessDirectPart(node types.Node) error {
	parts := &component.Parts

	xpathResult, err := node.Find(".//dao")
	if err != nil {
		return err
	}
	defer xpathResult.Free()

	daoNodes := xpathResult.NodeList()
	if len(daoNodes) == 0 {
		return nil
	}

	onlineAccessValues := []string{}

	for _, resultNode := range daoNodes {
		roleAttribute, err := resultNode.(types.Element).GetAttribute("xlink:role")
		if err != nil {
			if err.Error() == "attribute not found" {
				continue
			} else {
				return err
			}
		}

		role := roleAttribute.Value()
		roleValue, ok := eadutil.OnlineAccessRolesToFaceValues[role]
		if ok {
			onlineAccessValues = append(onlineAccessValues, roleValue)
		}
	}
	// Sort and deduplicate values
	slices.Sort(onlineAccessValues)
	parts.OnlineAccessDirect.Values = slices.Compact(onlineAccessValues)

	return nil
}
