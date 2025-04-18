package component

import "github.com/lestrrat-go/libxml2/types"

func (component *Component) setContainersPart(node types.Node) error {
	containers := []Container{}

	xpathResult, err := node.Find(".//container")
	if err != nil {
		return err
	}
	defer xpathResult.Free()

	containerNodes := xpathResult.NodeList()
	for _, containerNode := range containerNodes {
		container, err := makeContainer(containerNode)
		if err != nil {
			return err
		}
		containers = append(containers, container)
	}

	component.Parts.Containers = containers

	return nil
}

func makeContainer(containerNode types.Node) (Container, error) {
	container := Container{
		Value:     containerNode.NodeValue(),
		XMLString: containerNode.String(),
	}

	idAttributeNode, err := containerNode.(types.Element).GetAttribute("id")
	if err == nil {
		container.ID = idAttributeNode.Value()
	} else {
		return container, err
	}

	typeAttributeNode, err := containerNode.(types.Element).GetAttribute("type")
	if err == nil {
		container.Type = typeAttributeNode.Value()
	} else if err.Error() == "attribute not found" {
		// TODO: DLFA-238
		// This is a unique DLFA_238 case, in that we have identified a DLFA-211
		// bug, but the specifics of what the bug is depends on whether stakeholders
		// consider the `type` attribute to be required in the <container> element
		// or not required, and if the former, whether that invalidates the entire
		// EAD file.  For details, see:
		// 	  * https://jira.nyu.edu/browse/DLFA-277
		//        * "Debug indexer error "attribute not found" for tamwag/aia_003.xml"
		//    * https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=11923734&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-11923734
		//        * "[bug] Missing <container> attribute `type` leads to ": " at beginning of <location_ssm> and <location_ssi> values"
		// For now we are matching v1 indexer's behavior of not invalidating the
		// <container> and making `location_ssm` and `location_si` Solr fields
		// as best we can.  The v1 indexer looks to be doing this by accident,
		// and adds ": " to the beginning of the field values.  We will strip out
		// this erroneous blank `type` intro string, abd instead just use an empty
		// string to represent a missing `type`.
		// The DLFA-238 steps will be one of these options:
		//     * Keep this code as-is and add a golden file test for tamwag/aia_003.
		//     * Use some other handling non-fatal handling for missing `type`.
		//         * Use a label/intro like "[Type unspecified]: "
		//         * Keep code as-is but but `log.Warn()` (would need to enable `warn` logging level)
		//     * Remove this code and let the EAD file fail completely and return a more comprehensible wrapped error.
		//         * Note that we will probably do an error reporting overhaul,
		//        	 so might want to hold off on crafting the "perfect" error.
		//     * Something else
		container.Type = ""
	} else {
		return container, err
	}

	parentAttributeNode, err := containerNode.(types.Element).GetAttribute("parent")
	if err == nil {
		container.Parent = parentAttributeNode.Value()
	} else {
		// Do nothing.  This is a root node.
	}

	return container, nil
}
