package language

var LanguageTestCodes = []string{
	"en", "eng", "ger", "deu", "spa", "es", "fre", "fra", "fr", "ara", "ar",
	"ukr", "uk", "rus", "ru", "por", "pt",
}

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
