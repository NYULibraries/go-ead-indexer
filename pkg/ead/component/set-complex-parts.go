package component

import (
	"errors"
	"fmt"
	"go-ead-indexer/pkg/ead/eadutil"
	"go-ead-indexer/pkg/util"
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
	component.setChronListText()
	component.setCreatorComplex()
	component.setDAO()
	component.setFormat()
	component.setHeading()
	err := component.setLocation()
	if err != nil {
		return err
	}
	component.setMaterialType()
	component.setName()
	component.setPlace()
	component.setSubjectForFacets()
	err = component.setUnitTitleHTML()
	if err != nil {
		return err
	}
	component.setUnitDateDisplay()
	component.setUnitDateEnd()
	component.setUnitDateStart()

	return nil
}

func (component *Component) setCreatorComplex() {
	parts := &component.Parts

	// CreatorComplex
	creatorComplexValues := []string{}
	creatorComplexValues = append(creatorComplexValues, parts.CreatorCorpName.Values...)
	creatorComplexValues = append(creatorComplexValues, parts.CreatorFamName.Values...)
	creatorComplexValues = append(creatorComplexValues, parts.CreatorPersName.Values...)
	parts.CreatorComplex.Values = creatorComplexValues
}

func (component *Component) setChronListText() {
	parts := &component.Parts

	chronListTextValues := []string{}
	for _, chronListValue := range parts.ChronList.Values {
		if util.IsNonEmptyString(chronListValue) {
			chronListTextValues = append(chronListTextValues, chronListValue)
		}
	}

	parts.ChronListText.Values = chronListTextValues
}

func (component *Component) setDAO() {
	parts := &component.Parts

	if len(parts.DAODescriptionParagraph.Values) > 0 {
		parts.DAO.Values = []string{"Online Access"}
	} else {
		// No value
	}
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

func (component *Component) setHeading() {
	component.Parts.Heading.Values = append(component.Parts.Heading.Values,
		component.Parts.DIDUnitTitle.Values...)
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

func (component *Component) setMaterialType() {
	component.Parts.MaterialType.Values =
		eadutil.ConvertToFacetSlice(component.Parts.GenreForm.Values)
}

func (component *Component) setName() {
	parts := &component.Parts

	nameValues := []string{}
	nameValues = append(nameValues, parts.CorpNameNotInRepository.Values...)
	nameValues = append(nameValues, parts.FamName.Values...)
	nameValues = append(nameValues, parts.PersName.Values...)

	nameValues = eadutil.ConvertToFacetSlice(nameValues)

	parts.Name.Values = nameValues
}

func (component *Component) setPlace() {
	component.Parts.Place.Values =
		eadutil.ConvertToFacetSlice(component.Parts.GeogName.Values)
}

func (component *Component) setSubjectForFacets() {
	component.Parts.SubjectForFacets.Values =
		eadutil.ConvertToFacetSlice(component.Parts.Subject.Values)
}

func (component *Component) setUnitDateDisplay() {
	parts := &component.Parts

	parts.UnitDateDisplay.Values = []string{
		eadutil.GetUnitDateDisplay(parts.UnitDateNoTypeAttribute.Values,
			parts.UnitDateInclusive.Values, parts.UnitDateBulk.Values),
	}
}

func (component *Component) setUnitDateEnd() {
	parts := &component.Parts

	unitDateEndValues := []string{}
	for _, unitDateNormal := range parts.UnitDateNormal.Values {
		unitDateEnd := eadutil.GetDateParts(unitDateNormal).End
		if unitDateEnd != "" {
			unitDateEndValues = append(unitDateEndValues, unitDateEnd)
		}
	}

	parts.UnitDateEnd.Values = unitDateEndValues
}

func (component *Component) setUnitDateStart() {
	parts := &component.Parts

	unitDateStartValues := []string{}
	for _, unitDateNormal := range parts.UnitDateNormal.Values {
		unitDateStart := eadutil.GetDateParts(unitDateNormal).Start
		if unitDateStart != "" {
			unitDateStartValues = append(unitDateStartValues, unitDateStart)
		}
	}

	parts.UnitDateStart.Values = unitDateStartValues
}

func (component *Component) setUnitTitleHTML() error {
	parts := &component.Parts

	unitTitleHTMLValues := []string{}
	for _, unitTitle := range parts.DIDUnitTitle.XMLStrings {
		// `eadutil.MakeTitleHTML()` will in most if not all cases strip out the
		// open and close tags, but better safe than sorry.
		unitTitleContents := eadutil.StripOpenAndCloseTags(unitTitle)

		unitTitleHTMLValue, err := eadutil.MakeTitleHTML(unitTitleContents)
		if err != nil {
			return err
		}

		unitTitleHTMLValues = append(unitTitleHTMLValues, unitTitleHTMLValue)
	}

	parts.UnitTitleHTML.Values = unitTitleHTMLValues

	return nil
}
