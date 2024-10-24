package ead

func (collectionDoc *CollectionDoc) setCreator() {
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
		getDateRange(collectionDoc.Parts.UnitDateNormal.Values)
}

func (collectionDoc *CollectionDoc) setMaterialType() {
	collectionDoc.Parts.MaterialType.Values =
		convertToFacetSlice(collectionDoc.Parts.GenreForm.Values)
}

func (collectionDoc *CollectionDoc) setLanguage() []error {
	language, errs := getLanguage(collectionDoc.Parts.LangCode.Values)
	if len(errs) > 0 {
		return errs
	}

	collectionDoc.Parts.Language.Values = language

	return nil
}

func (collectionDoc *CollectionDoc) setName() {
	parts := &collectionDoc.Parts

	nameValues := []string{}
	nameValues = append(nameValues, parts.FamName.Values...)
	nameValues = append(nameValues, parts.PersName.Values...)
	nameValues = append(nameValues, parts.CorpNameNotInRepository.Values...)

	nameValues = convertToFacetSlice(nameValues)

	parts.Name.Values = nameValues
}

func (collectionDoc *CollectionDoc) setOnlineAccess() {
	if len(collectionDoc.Parts.DAO.Values) > 0 {
		collectionDoc.Parts.OnlineAccess.Values = []string{"Online Access"}
	}
}

func (collectionDoc *CollectionDoc) setPlace() {
	collectionDoc.Parts.Place.Values =
		convertToFacetSlice(collectionDoc.Parts.GeogName.Values)
}

	parts.Place.Values = placeValues
}

func (collectionDoc *CollectionDoc) setUnitDateDisplay() {
	parts := &collectionDoc.Parts

	parts.UnitDateDisplay.Values[0] = getUnitDateDisplay(parts.UnitDateNoTypeAttribute.Values,
		parts.UnitDateInclusive.Values, parts.UnitDateBulk.Values)
}

func (collectionDoc *CollectionDoc) setUnitDateEnd() {
	parts := &collectionDoc.Parts

	unitDateEndValues := []string{}
	for _, unitDateNormal := range parts.UnitDateNormal.Values {
		unitDateEnd := getDateParts(unitDateNormal).End
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
		unitDateStart := getDateParts(unitDateNormal).Start
		if unitDateStart != "" {
			unitDateStartValues = append(unitDateStartValues, unitDateStart)
		}
	}

	parts.UnitDateStart.Values = unitDateStartValues
}
