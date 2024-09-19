package language

import (
	"testing"
)

func assertError(t *testing.T, err error, expectedErr error, testCase string) {
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("Expected error: %s, got: %s for test case: %s", expectedErr, err, testCase)
	}
}

func assertLanguage(t *testing.T, language string, expectedLanguage string, err error, testCase string) {
	if err != nil {
		t.Errorf("Unexpected error: %v for test case: %s", err, testCase)
	} else if language != expectedLanguage {
		t.Errorf("Expected language: %s, got: %s for test case: %s", expectedLanguage, language, testCase)
	}
}

func TestGetLanguageForLanguageCode_Errors(t *testing.T) {
	var tests = []struct {
		name          string
		expectedError error
		languageCode  string
	}{
		{"invalid length", ErrInvalidLength, "abcd"},
		{"empty string", ErrEmptyLanguageCode, ""},
		{"internal whitespace", ErrInternalWhitespace, "a a"},
		{"invalid language code with leading space", ErrLanguageNotFound, " aR"},
		{"invalid language code with trailing space", ErrLanguageNotFound, "Ra "},
		{"non-existing language code", ErrLanguageNotFound, "zpy"},
		{"language code contains invalid characters", ErrInvalidCharacters, "en1"},
		{"language code contains invalid characters", ErrInvalidCharacters, "e#}"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := GetLanguageForLanguageCode(test.languageCode)
			assertError(t, err, test.expectedError, test.name)
		})
	}
}

func TestGetLanguageForLanguageCode(t *testing.T) {
	var tests = []struct {
		name             string
		expectedLanguage string
		languageCode     string
	}{
		{"lowercase", "aar", "Afar"},
		{"uppercase", "ABK", "Abkhazian"},
		{"mixedcase", "aFr", "Afrikaans"},
		{"lowercase with whitespace", " alb ", "Albanian"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := GetLanguageForLanguageCode(test.languageCode)
			assertLanguage(t, result, test.expectedLanguage, err, test.languageCode)
		})
	}
}
