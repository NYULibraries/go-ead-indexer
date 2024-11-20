package component

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const ARCHIVAL_OBJECT_FORMAT = "Archival Object"
const ARCHIVAL_SERIES_FORMAT = "Archival Series"

var archivalSeriesRegExp = regexp.MustCompile(`/\Aseries|subseries/`)

// This algorithm is based on the one used here:
// https://github.com/NYULibraries/dlts-finding-aids-ead-go-packages/blob/7baee7dfde24a01422ec8e6470fdc8a76d84b3fb/ead/modify/modify.go#L153-L180
func (component *Component) makeRootContainerSliceAndParentChildContainerMap() (
	[]Container, map[string]Container, error) {
	rootContainers := []Container{}
	parentChildContainerMap := map[string]Container{}

	mappingErrors := []string{}
	for _, container := range component.Parts.Containers {
		parentID := container.Parent
		if parentID != "" {
			// Has a parent, so must be a child container.  Map it to its parent.

			// There should be no sibling relationships.  If there is already a
			// child container mapped to `parentID`, that's an error condition.
			if _, ok := parentChildContainerMap[parentID]; ok {
				mappingErrors = append(mappingErrors,
					fmt.Sprintf("A child <container> element has already"+
						` been mapped to a parent <container> with @id="%s""`, parentID))
			}

			parentChildContainerMap[parentID] = container
		} else {
			// No parent, so must be a root container.
			rootContainers = append(rootContainers, container)
		}
	}

	var err error
	if len(mappingErrors) > 0 {
		err = errors.New(strings.Join(mappingErrors, "; "))
	}

	return rootContainers, parentChildContainerMap, err
}

// TODO: Do we need to have anything in `CollectionDoc.Part.Source` for these?
func (component *Component) setComplexParts() error {
	component.setFormat()
	err := component.setLocation()
	if err != nil {
		return err
	}

	return nil
}

func (component *Component) setFormat() {
	parts := &component.Parts

	var formatLevel string
	level := parts.Level.Values[0]
	if archivalSeriesRegExp.Match([]byte(level)) {
		formatLevel = ARCHIVAL_SERIES_FORMAT
	} else {
		formatLevel = ARCHIVAL_OBJECT_FORMAT
	}

	parts.Format.Values = []string{formatLevel}
}

// TODO: DLFA-243
// Change this from <container> hierarchy based on occurrence order of elements
// to using `parent` attribute linking to the `id` attribute of direct parent
// <container>.
func (component *Component) setLocation() error {
	parts := &component.Parts

	// TODO: DLFA-243
	// Delete this:
	locationValues, err := component.getLocationValuesInOccurrenceOrder()
	// ...and uncomment this:
	//locationValues, err := component.getLocationValues()
	if err != nil {
		return err
	} else {
		parts.Location.Values = locationValues

		return nil
	}
}

func (component *Component) getLocationValues() ([]string, error) {
	locationValues := []string{}
	rootContainersSlice, parentChildContainerMap, err :=
		component.makeRootContainerSliceAndParentChildContainerMap()
	if err != nil {
		return locationValues, err
	}

	for _, rootContainer := range rootContainersSlice {
		locationValue := fmt.Sprintf("%s: %s", rootContainer.Type, rootContainer.Value)
		currentParentID := rootContainer.ID
		for {
			childContainer, ok := parentChildContainerMap[currentParentID]
			if ok {
				locationValue += fmt.Sprintf(", %s: %s", childContainer.Type, childContainer.Value)
				currentParentID = childContainer.ID
			} else {
				break
			}
		}

		locationValues = append(locationValues, locationValue)
	}

	return locationValues, nil
}

// TODO: DLFA-243
// Remove this after passing the DLFA-201 acceptance test.  See function header
// for `Component.setLocation()` above.
func (component *Component) getLocationValuesInOccurrenceOrder() ([]string, error) {
	parts := &component.Parts

	locationValues := []string{}

	currentLocationValue := ""
	for _, container := range parts.Containers {
		// This is a root <container>.  Commit the in-progress location value,
		// if any, and start up a new one.
		if container.Parent == "" {
			if currentLocationValue != "" {
				locationValues = append(locationValues, currentLocationValue)
				currentLocationValue = ""
			}

			currentLocationValue += fmt.Sprintf("%s: %s", container.Type, container.Value)
		} else {
			currentLocationValue += fmt.Sprintf(", %s: %s", container.Type, container.Value)
		}
	}

	// Commit the in-progress location value, if any.
	if currentLocationValue != "" {
		locationValues = append(locationValues, currentLocationValue)
		currentLocationValue = ""
	}

	// This function can't really return an error, but when we switch over
	// to the permanent `getLocationValues()` we'll need to trap errors in
	// mapping parent to child containers, so we match its signature to make
	// transition easier and faster.
	return locationValues, nil
}
