package collectiondoc

import (
	"github.com/nyulibraries/go-ead-indexer/pkg/ead/eadutil"
	"github.com/nyulibraries/go-ead-indexer/pkg/util"
	"strings"
)

// TODO: Do we need to have anything in `CollectionDoc.Part.Source` for these?
// TODO: figure out whether to keep return value of `setComplexParts()` as
// `[]error` or to change it to a single error, which is what the caller
// `setParts()` returns to its caller.  If a single error, should it be an
// early exit true single error, or use `errors.Join()` to wrap the slice of
// all accumulated errors?
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
		if util.IsNonEmptyString(chronListValue) {
			chronListComplexValues = append(chronListComplexValues, strings.TrimSpace(chronListValue))
		}
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
		eadutil.GetDateRange(collectionDoc.Parts.UnitDateNormal.Values)
}

func (collectionDoc *CollectionDoc) setLanguage() []error {
	language, errs := eadutil.GetLanguage(collectionDoc.Parts.LangCode.Values)
	if len(errs) > 0 {
		return errs
	}

	collectionDoc.Parts.Language.Values = language

	return nil
}

func (collectionDoc *CollectionDoc) setMaterialType() {
	collectionDoc.Parts.MaterialType.Values =
		eadutil.ConvertToFacetSlice(collectionDoc.Parts.GenreForm.Values)
}

func (collectionDoc *CollectionDoc) setName() {
	parts := &collectionDoc.Parts

	nameValues := []string{}
	nameValues = append(nameValues, parts.CorpNameNotInRepository.Values...)
	nameValues = append(nameValues, parts.FamName.Values...)
	nameValues = append(nameValues, parts.PersName.Values...)

	nameValues = eadutil.ConvertToFacetSlice(nameValues)

	parts.Name.Values = nameValues
}

func (collectionDoc *CollectionDoc) setOnlineAccess() {
	if len(collectionDoc.Parts.DAO.Values) > 0 {
		collectionDoc.Parts.OnlineAccess.Values = []string{"Online Access"}
	}
}

func (collectionDoc *CollectionDoc) setPlace() {
	collectionDoc.Parts.Place.Values =
		eadutil.ConvertToFacetSlice(collectionDoc.Parts.GeogName.Values)
}

func (collectionDoc *CollectionDoc) setSubjectForFacets() {
	collectionDoc.Parts.SubjectForFacets.Values =
		eadutil.ConvertToFacetSlice(collectionDoc.Parts.Subject.Values)
}

func (collectionDoc *CollectionDoc) setUnitDateDisplay() {
	parts := &collectionDoc.Parts

	parts.UnitDateDisplay.Values = []string{
		eadutil.GetUnitDateDisplay(parts.UnitDateNoTypeAttribute.Values,
			parts.UnitDateInclusive.Values, parts.UnitDateBulk.Values),
	}
}

func (collectionDoc *CollectionDoc) setUnitDateEnd() {
	parts := &collectionDoc.Parts

	unitDateEndValues := []string{}
	for _, unitDateNormal := range parts.UnitDateNormal.Values {
		unitDateEnd := eadutil.GetDateParts(unitDateNormal).End
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
		unitDateStart := eadutil.GetDateParts(unitDateNormal).Start
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
