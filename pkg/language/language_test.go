package language

import (
	"regexp"
	"strings"
	"testing"
)

func assertError(t *testing.T, err error, expectedErr error, testCase string) {
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("Expected error for test case \"%s\": \"%s\", got: \"%s\"",
			testCase, expectedErr, err)
	}
}

func assertLanguage(t *testing.T, language string, expectedLanguage string, err error, testCase string) {
	if err != nil {
		t.Errorf("Unexpected error for test case \"%s\": \"%v\" ", testCase, err)
	} else if language != expectedLanguage {
		t.Errorf("Expected language for test case \"%s\": \"%s\", got: \"%s\"",
			testCase, expectedLanguage, language)
	}
}

// Convert newlines, carriage returns, and tabs to escapes "\n", "\r", and "\t" for
// safe and more readable printing in test results and logs.
func escapeForReadability(code string) string {
	carriageReturnRegexp := regexp.MustCompile(`\r`)
	newlineRegexp := regexp.MustCompile(`\n`)
	tabRegexp := regexp.MustCompile(`\t`)

	code = carriageReturnRegexp.ReplaceAllString(code, "\\r")
	code = newlineRegexp.ReplaceAllString(code, "\\n")
	code = tabRegexp.ReplaceAllString(code, "\\t")

	return code
}

func TestGetLanguageForLanguageCode_StaticCases(t *testing.T) {
	var staticErrorTests = []struct {
		name          string
		languageCode  string
		expectedError error
	}{
		{"invalid length", "abcd", ErrInvalidLength},
		{"empty string", "", ErrEmptyLanguageCode},
		{"invalid language code with leading space", " qA", ErrLanguageNotFound},
		{"invalid language code with trailing space", "Ra ", ErrLanguageNotFound},
		{"non-existing language code", "zpy", ErrLanguageNotFound},
		{"language code contains invalid characters", "en1", ErrInvalidCharacters},
		{"language code contains invalid characters", "e#}", ErrInvalidCharacters},
	}

	for _, test := range staticErrorTests {
		t.Run(test.name, func(t *testing.T) {
			_, err := GetLanguageForLanguageCode(test.languageCode)
			assertError(t, err, test.expectedError, test.name+" for "+test.languageCode)
		})
	}

}
func TestGetLanguageForLanguageCode_Errors(t *testing.T) {
	var tests = []struct {
		name          string
		expectedError error
		inputModifier func(string) string
	}{
		{"internal whitespace", ErrInternalWhitespace, func(code string) string { return code[:1] + " " + code[1:] }},
		{"carriage return character in between", ErrInternalWhitespace, func(code string) string { return code[:1] + "\r" + code[1:] }},
		{"tab character in between", ErrInternalWhitespace, func(code string) string { return code[:1] + "\t" + code[1:] }},
		{"new line in between", ErrInternalWhitespace, func(code string) string { return code[:1] + "\n" + code[1:] }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, languageCode := range GetTestLanguageCodes() {
				modifiedCode := test.inputModifier(languageCode)
				_, err := GetLanguageForLanguageCode(modifiedCode)
				assertError(t, err, test.expectedError, test.name+" for '"+escapeForReadability(languageCode)+"'")
			}
		})
	}
}

func TestGetLanguageForLanguageCode(t *testing.T) {
	var tests = []struct {
		name          string
		inputModifier func(string) string
	}{
		{"lowercase", func(code string) string { return strings.ToLower(code) }},
		{"uppercase", func(code string) string { return strings.ToUpper(code) }},
		{"mixedcase", func(code string) string {
			if len(code) > 1 {
				return strings.ToLower(code[:1]) + strings.ToUpper(code[1:])
			}
			return strings.ToLower(code)
		}},
		{"lowercase with whitespace", func(code string) string { return " " + code + " " }},
		{"valid code with new lines", func(code string) string { return code + "\n" }},
		{"valid code with carriage return", func(code string) string { return code + "\r" }},
		{"valid code with new lines", func(code string) string { return code + "\n\r" }},
		{"valid code with tab", func(code string) string { return code + "\t" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for code, expectedLanguage := range ExpectedTestLanguages {
				modifiedCode := test.inputModifier(code)
				result, err := GetLanguageForLanguageCode(modifiedCode)
				assertLanguage(t, result, expectedLanguage, err, test.name+" for '"+escapeForReadability(modifiedCode)+"'")
			}
		})
	}
}
