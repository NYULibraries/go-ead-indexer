package language

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

const MIN_LANGUAGE_CODE_LENGTH = 2
const MAX_LANGUAGE_CODE_LENGTH = 3

var (
	ErrEmptyLanguageCode = errors.New("language code cannot be empty")
	ErrInvalidCharacters = errors.New("language code contains invalid characters")
	ErrInvalidLength     = errors.New(
		fmt.Sprintf("language code must be between %d-%d characters long, inclusive",
			MIN_LANGUAGE_CODE_LENGTH, MAX_LANGUAGE_CODE_LENGTH))
	ErrInternalWhitespace = errors.New("language code contains internal whitespace")
	// We can't refer the user to the official ISO-639-2 language code list because
	// the list changes over time and is not versioned.  We are in fact currently
	// using an out of date ISO-639-2 language code map in order to exactly match
	// the v1 indexer translations.  For details, see this comment in DLFA-224:
	// https://jira.nyu.edu/browse/DLFA-224?focusedCommentId=9870333&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-9870333
	ErrLanguageNotFound = errors.New("language code not found.  Please refer to the source code of this package for accepted language codes.")
)

func GetLanguageForLanguageCode(languageCode string) (string, error) {
	languageCode = strings.ToLower(strings.TrimSpace(languageCode))

	if languageCode == "" {
		return "", ErrEmptyLanguageCode
	}

	if containsWhitespace(languageCode) {
		return "", ErrInternalWhitespace
	}

	if !isAlphabetic(languageCode) {
		return "", ErrInvalidCharacters
	}

	if len(languageCode) < MIN_LANGUAGE_CODE_LENGTH ||
		len(languageCode) > MAX_LANGUAGE_CODE_LENGTH {
		return "", ErrInvalidLength
	}

	if language, found := ISO639_2_DB[languageCode]; found {
		return language.Name, nil
	}

	return "", ErrLanguageNotFound
}

func containsWhitespace(s string) bool {
	for _, r := range s {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}

func isAlphabetic(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
