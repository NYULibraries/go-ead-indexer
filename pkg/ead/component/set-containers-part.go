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
