package language

import (
	"maps"
	"slices"
)

var ExpectedTestLanguages = map[string]string{
	"en":  "English",
	"eng": "English",
	"ger": "German",
	"deu": "German",
	"spa": "Spanish; Castilian",
	"es":  "Spanish; Castilian",
	"fre": "French",
	"fra": "French",
	"fr":  "French",
	"ara": "Arabic",
	"ar":  "Arabic",
	"ukr": "Ukrainian",
	"uk":  "Ukrainian",
	"rus": "Russian",
	"ru":  "Russian",
	"por": "Portuguese",
	"pt":  "Portuguese",
}

func GetTestLanguageCodes() []string {
	return slices.Sorted(maps.Keys(ExpectedTestLanguages))
}
