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

var archivalSeriesRegExp = regexp.MustCompile(`\Aseries|subseries`)

// TODO: Do we need to have anything in `CollectionDoc.Part.Source` for these?
// TODO: figure out whether to keep return value of `setComplexParts()` as a
// single error or to change it to an error slice.  Note that the caller
// `setParts()` returns a single error to its caller.  If this method should
// return a single error, should it be an early exit single error as it is now,
// or should we use `errors.Join()` to wrap the slice of all accumulated errors?
func (component *Component) setComplexParts() error {
	component.setChronListComplex()
	component.setCreatorComplex()
	component.setDAO()

	// TODO: DLFA-238
	// Remove this override which adds  left- and right- padding for matching v1
	// indexer bug behavior described here:
	// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10849506&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10849506
	component.setDAODescriptionParagraph()

	component.setDateRange()

	// TODO: DLFA-238
	// Remove this override which adds  left- and right- padding for matching v1
	// indexer bug behavior described here:
	// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10849506&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10849506
	component.setDIDUnitTitle()

	component.setFormat()
	component.setHeading()
	component.setLanguage()
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

func (component *Component) setChronListComplex() {
	parts := &component.Parts

	chronListComplexValues := []string{}
	for _, chronListValue := range parts.ChronList.Values {
		if util.IsNonEmptyString(chronListValue) {
			chronListComplexValues = append(chronListComplexValues, strings.TrimSpace(chronListValue))
		}
	}

	parts.ChronListComplex.Values = chronListComplexValues
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

func (component *Component) setDAO() {
	parts := &component.Parts

	if len(parts.DAODescriptionParagraph.Values) > 0 {
		parts.DAO.Values = []string{"Online Access"}
	} else {
		// No value
	}
}

// TODO: DLFA-238
// Remove this override which adds  left- and right- padding for matching v1
// indexer bug behavior described here:
// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10849506&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10849506
func (component *Component) setDAODescriptionParagraph() {
	parts := &component.Parts

	paddedDAODescriptionParagraph := []string{}
	numDAODescriptionParagraph := len(parts.DAODescriptionParagraph.Values)
	for i := 0; i < numDAODescriptionParagraph; i++ {
		daoDescriptionParagraphXMLString := eadutil.StripOpenAndCloseTags(parts.DAODescriptionParagraph.XMLStrings[i])

		daoDescriptionParagraphValue := eadutil.PadDAODescriptionParagraphIfNeeded(
			daoDescriptionParagraphXMLString,
			eadutil.StripOpenAndCloseTags(parts.DAODescriptionParagraph.Values[i]))

		paddedDAODescriptionParagraph = append(paddedDAODescriptionParagraph, daoDescriptionParagraphValue)
	}

	parts.DAODescriptionParagraph.Values = paddedDAODescriptionParagraph
}

func (component *Component) setDateRange() {
	component.Parts.DateRange.Values =
		eadutil.GetDateRange(component.Parts.UnitDateNormal.Values)
}

// TODO: DLFA-238
// Remove this override which adds  left- and right- padding for matching v1
// indexer bug behavior described here:
// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10849506&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10849506
func (component *Component) setDIDUnitTitle() {
	parts := &component.Parts

	paddedDIDUnitTitleValues := []string{}
	numDIDUnitTitleValues := len(parts.DIDUnitTitle.Values)
	for i := 0; i < numDIDUnitTitleValues; i++ {
		unitTitleXMLString := eadutil.StripOpenAndCloseTags(parts.DIDUnitTitle.XMLStrings[i])
		unitTitleValue := eadutil.PadUnitTitleIfNeeded(unitTitleXMLString,
			eadutil.StripOpenAndCloseTags(parts.DIDUnitTitle.Values[i]))

		paddedDIDUnitTitleValues = append(paddedDIDUnitTitleValues, unitTitleValue)
	}

	parts.DIDUnitTitle.Values = paddedDIDUnitTitleValues
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

func (component *Component) setLanguage() []error {
	language, errs := eadutil.GetLanguage(component.Parts.LangCode.Values)
	if len(errs) > 0 {
		return errs
	}

	component.Parts.Language.Values = language

	return nil
}

// TODO: DLFA-238
// Change this from <container> hierarchy based on occurrence order of elements
// to using `parent` attribute linking to the `id` attribute of direct parent
// <container>.
func (component *Component) setLocation() error {
	parts := &component.Parts

	// TODO: DLFA-238
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

// TODO: DLFA-238
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
		eadutil.ConvertToFacetSlice(component.Parts.SubjectOrFunctionOrOccupation.Values)
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

		// TODO: DLFA-238
		// Remove this left- and right- padding for matching v1 indexer bug
		// behavior described here:
		// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10849506&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10849506
		unitTitleHTMLValue = eadutil.PadUnitTitleIfNeeded(unitTitleContents, unitTitleHTMLValue)

		// TODO: DLFA-238
		// This is tough one to call.  Not sure if this absolutely needs to be
		// removed or not.  Go's `xml.EscapeText()` automatically escapes single
		// and double-quotes because they could potentially break tag attributes,
		// but these <unittitle> values are never going to be in attributes.
		// If there's no harm in letting them stay escaped, however, probably
		// better to not undo it.
		unitTitleHTMLValue = strings.ReplaceAll(unitTitleHTMLValue, "&#39;", "'")
		unitTitleHTMLValue = strings.ReplaceAll(unitTitleHTMLValue, "&#34;", `"`)

		unitTitleHTMLValues = append(unitTitleHTMLValues, unitTitleHTMLValue)
	}

	parts.UnitTitleHTML.Values = unitTitleHTMLValues

	return nil
}
