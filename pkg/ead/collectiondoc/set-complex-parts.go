package collectiondoc

import (
	"go-ead-indexer/pkg/ead/util"
	"strings"
)

// TODO: Do we need to have anything in `CollectionDoc.Part.Source` for these?
func (collectionDoc *CollectionDoc) setComplexParts() []error {
	errs := []error{}

	collectionDoc.setChronListComplex()
	collectionDoc.setCreatorComplex()
	collectionDoc.setDateRange()
	languageErrors := collectionDoc.setLanguage()
	if len(languageErrors) > 0 {
		errs = append(errs, languageErrors...)
	}
	collectionDoc.setMaterialType()
	collectionDoc.setName()
	collectionDoc.setOnlineAccess()
	collectionDoc.setPlace()
	collectionDoc.setSubjectForFacets()
	collectionDoc.setUnitDateEnd()
	collectionDoc.setUnitDateStart()
	collectionDoc.setUnitDateDisplay()
	unitTitleHTMLError := collectionDoc.setUnitTitleHTML()
	if unitTitleHTMLError != nil {
		errs = append(errs, unitTitleHTMLError)
	}

	return errs
}

func (collectionDoc *CollectionDoc) setChronListComplex() {
	parts := &collectionDoc.Parts

	chronListComplexValues := []string{}
	for _, chronListValue := range parts.ChronList.Values {
		chronListComplexValues = append(chronListComplexValues, strings.TrimSpace(chronListValue))
	}

	parts.ChronListComplex.Values = chronListComplexValues
}

func (collectionDoc *CollectionDoc) setCreatorComplex() {
	parts := &collectionDoc.Parts

	// CreatorComplex
	creatorComplexValues := []string{}
	creatorComplexValues = append(creatorComplexValues, parts.CreatorCorpName.Values...)
	creatorComplexValues = append(creatorComplexValues, parts.CreatorFamName.Values...)
	creatorComplexValues = append(creatorComplexValues, parts.CreatorPersName.Values...)
	parts.CreatorComplex.Values = creatorComplexValues
}

func (collectionDoc *CollectionDoc) setDateRange() {
	collectionDoc.Parts.DateRange.Values =
		util.GetDateRange(collectionDoc.Parts.UnitDateNormal.Values)
}

func (collectionDoc *CollectionDoc) setMaterialType() {
	collectionDoc.Parts.MaterialType.Values =
		util.ConvertToFacetSlice(collectionDoc.Parts.GenreForm.Values)
}

func (collectionDoc *CollectionDoc) setLanguage() []error {
	language, errs := util.GetLanguage(collectionDoc.Parts.LangCode.Values)
	if len(errs) > 0 {
		return errs
	}

	collectionDoc.Parts.Language.Values = language

	return nil
}

func (collectionDoc *CollectionDoc) setName() {
	parts := &collectionDoc.Parts

	nameValues := []string{}
	nameValues = append(nameValues, parts.CorpNameNotInRepository.Values...)
	nameValues = append(nameValues, parts.FamName.Values...)
	nameValues = append(nameValues, parts.PersName.Values...)

	nameValues = util.ConvertToFacetSlice(nameValues)

	parts.Name.Values = nameValues
}

func (collectionDoc *CollectionDoc) setOnlineAccess() {
	if len(collectionDoc.Parts.DAO.Values) > 0 {
		collectionDoc.Parts.OnlineAccess.Values = []string{"Online Access"}
	}
}

func (collectionDoc *CollectionDoc) setPlace() {
	collectionDoc.Parts.Place.Values =
		util.ConvertToFacetSlice(collectionDoc.Parts.GeogName.Values)
}

func (collectionDoc *CollectionDoc) setSubjectForFacets() {
	collectionDoc.Parts.SubjectForFacets.Values =
		util.ConvertToFacetSlice(collectionDoc.Parts.Subject.Values)
}

func (collectionDoc *CollectionDoc) setUnitDateDisplay() {
	parts := &collectionDoc.Parts

	parts.UnitDateDisplay.Values = []string{
		util.GetUnitDateDisplay(parts.UnitDateNoTypeAttribute.Values,
			parts.UnitDateInclusive.Values, parts.UnitDateBulk.Values),
	}
}

func (collectionDoc *CollectionDoc) setUnitDateEnd() {
	parts := &collectionDoc.Parts

	unitDateEndValues := []string{}
	for _, unitDateNormal := range parts.UnitDateNormal.Values {
		unitDateEnd := util.GetDateParts(unitDateNormal).End
		if unitDateEnd != "" {
			unitDateEndValues = append(unitDateEndValues, unitDateEnd)
		}
	}

	parts.UnitDateEnd.Values = unitDateEndValues
}

func (collectionDoc *CollectionDoc) setUnitDateStart() {
	parts := &collectionDoc.Parts

	unitDateStartValues := []string{}
	for _, unitDateNormal := range parts.UnitDateNormal.Values {
		unitDateStart := util.GetDateParts(unitDateNormal).Start
		if unitDateStart != "" {
			unitDateStartValues = append(unitDateStartValues, unitDateStart)
		}
	}

	parts.UnitDateStart.Values = unitDateStartValues
}

func (collectionDoc *CollectionDoc) setUnitTitleHTML() error {
	parts := &collectionDoc.Parts

	unitTitleHTMLValues := []string{}
	for _, unitTitle := range parts.UnitTitle.XMLStrings {
		unitTitleContents := strings.TrimSuffix(
			strings.TrimPrefix(unitTitle, "<unittitle>"),
			"</unittitle>")
		converted, err := util.ConvertEADToHTML(unitTitleContents)
		if err != nil {
			return err
		}

		unitTitleHTMLValue, err := util.StripTags(converted)
		if err != nil {
			return err
		}

		unitTitleHTMLValues = append(unitTitleHTMLValues, unitTitleHTMLValue)
	}

	parts.UnitTitleHTML.Values = unitTitleHTMLValues

	return nil
}
