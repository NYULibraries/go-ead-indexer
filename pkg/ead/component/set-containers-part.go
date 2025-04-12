package component

import (
	"fmt"
	"log"

	"github.com/lestrrat-go/libxml2/types"
)

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
			log.Printf("Warning: Skipping problematic container: %v", err)
			return fmt.Errorf("error creating container: %v", err)
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
		log.Printf("Error getting container type - Container ID: %s, XML: %s, Error: %v", container.ID, container.XMLString, err)
		return container, fmt.Errorf("missing required 'type' attribute for container ID %s: %v", container.ID, err)
	}

	parentAttributeNode, err := containerNode.(types.Element).GetAttribute("parent")
	if err == nil {
		container.Parent = parentAttributeNode.Value()
	} else {
		// Do nothing.  This is a root node.
	}

	return container, nil
}
