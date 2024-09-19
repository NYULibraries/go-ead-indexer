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
	_, err := GetLanguageForLanguageCode("abcd")
	assertError(t, err, ErrInvalidLength, "invalid length")

	_, err = GetLanguageForLanguageCode("")
	assertError(t, err, ErrEmptyLanguageCode, "empty string")

	_, err = GetLanguageForLanguageCode("a a")
	assertError(t, err, ErrInternalWhitespace, "internal whitespace")

	_, err = GetLanguageForLanguageCode(" aR")
	assertError(t, err, ErrLanguageNotFound, "invalid language code with leading space")

	_, err = GetLanguageForLanguageCode("Ra ")
	assertError(t, err, ErrLanguageNotFound, "invalid language code with trailing space")

	_, err = GetLanguageForLanguageCode("zpy")
	assertError(t, err, ErrLanguageNotFound, "non-existing language code")

	_, err = GetLanguageForLanguageCode("en1")
	assertError(t, err, ErrInvalidCharacters, "language code contains invalid characters")

	_, err = GetLanguageForLanguageCode("!es#")
	assertError(t, err, ErrInvalidCharacters, "language code contains invalid characters")
}
