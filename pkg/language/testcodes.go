package language

import (
	"maps"
	"slices"
)

var ExpectedTestLanguages = map[string]string{
	"en":  "English",
	"eng": "English",

	"deu": "German",
	"ger": "German",

	"es":  "Spanish; Castilian",
	"spa": "Spanish; Castilian",

	"fr":  "French",
	"fra": "French",
	"fre": "French",

	"ar":  "Arabic",
	"ara": "Arabic",

	"uk":  "Ukrainian",
	"ukr": "Ukrainian",

	"ru":  "Russian",
	"rus": "Russian",

	"por": "Portuguese",
	"pt":  "Portuguese",
}

func GetTestLanguageCodes() []string {
	return slices.Sorted(maps.Keys(ExpectedTestLanguages))
}
