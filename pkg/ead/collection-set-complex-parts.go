package ead

import (
	"go-ead-indexer/pkg/util"
)

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

func (collectionDoc *CollectionDoc) setName() {
	parts := &collectionDoc.Parts

	nameValues := []string{}
	nameValues = append(nameValues, parts.FamName.Values...)
	nameValues = append(nameValues, parts.PersName.Values...)
	nameValues = append(nameValues, parts.CorpNameNotInRepository.Values...)
	nameValues = replaceMARCSubfieldDemarcatorsInSlice(nameValues)
	nameValues = util.CompactStringSlicePreserveOrder(nameValues)
	parts.Name.Values = nameValues
}

func (collectionDoc *CollectionDoc) setOnlineAccess() {
	if len(collectionDoc.Parts.DAO.Values) > 0 {
		collectionDoc.Parts.OnlineAccess.Values = []string{"Online Access"}
	}
}

func (collectionDoc *CollectionDoc) setUnitDateDisplay() {
	parts := &collectionDoc.Parts

	parts.UnitDateDisplay.Values[0] = getUnitDateDisplay(parts.UnitDateNoTypeAttribute.Values,
		parts.UnitDateInclusive.Values, parts.UnitDateBulk.Values)
}
