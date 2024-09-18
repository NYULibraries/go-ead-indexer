package language

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrEmptyLanguageCode  = errors.New("language code cannot be empty")
	ErrInvalidCharacters  = errors.New("invalid characters provided")
	ErrInvalidLength      = errors.New("language code must be either 2 or 3 characters long")
	ErrInternalWhitespace = errors.New("language code contains invalid whitespace characters")
	ErrLanguageNotFound   = errors.New("language code not found. Please refer to ISO-639-2 language code table")
)

func isAlphabetic(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func GetLanguageForLanguageCode(languageCode string) (string, error) {

	languageCode = strings.ToLower(strings.TrimSpace(languageCode))

	if languageCode == "" {
		return "", ErrEmptyLanguageCode
	}

	if strings.Contains(languageCode, " ") || strings.Contains(languageCode, "\t") {
		return "", ErrInternalWhitespace
	}

	if !isAlphabetic(languageCode) {
		return "", ErrInvalidCharacters
	}

	if len(languageCode) < 2 || len(languageCode) > 3 {
		return "", ErrInvalidLength
	}

	if language, found := ISO639_2_DB[languageCode]; found {
		return language.Name, nil
	}

	return "", ErrLanguageNotFound
}
