package util

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
)

func TestCompactStringSlicePreserveOrder(t *testing.T) {
	testCases := []struct {
		in  []string
		out []string
	}{
		{
			[]string{},
			[]string{},
		},
		{
			[]string{"a"},
			[]string{"a"},
		},
		{
			[]string{"a", "a", "a", "b", "b", "c", "c"},
			[]string{"a", "b", "c"},
		},

		{
			[]string{"a", "b", "c", "a", "d", "b", "d", "c", "c", "a", "e"},
			[]string{"a", "b", "c", "d", "e"},
		},
	}

	for _, testCase := range testCases {
		actual := CompactStringSlicePreserveOrder(testCase.in)
		if slices.Compare(actual, testCase.out) != 0 {
			t.Errorf("Expected compacted slice %v for input slice %v, got %v",
				testCase.out, testCase.in, actual)
		}
	}
}

func TestGetCallerFunctionNameError(t *testing.T) {
	expectedError := "runtime.Caller(999) failed"

	_, _, err := GetCallerFunctionName(999)
	if err == nil {
		t.Errorf("Expected GetCallerFunctionName(2) to return error"+
			` "%s", but no error was returned.`, expectedError)
	} else {
		actualError := err.Error()
		if actualError != expectedError {
			t.Errorf("Expected GetCallerFunctionName(2) to return error"+
				` "%s", but instead got error "%s"`, expectedError, actualError)
		}
	}
}

func TestGetCallerFunctionNameSkip1(t *testing.T) {
	const expectedPackage = "github.com/nyulibraries/go-ead-indexer/pkg/util"
	const expectedFunction = "TestGetCallerFunctionNameSkip1"

	actualPackage, actualFunction, err := GetCallerFunctionName(1)
	if err != nil {
		t.Errorf("Expected GetCallerFunction(1) to not return an error"+
			` but an error was returned: "%s""`, err.Error())
	}

	if actualPackage != expectedPackage {
		t.Errorf(`Expected package name "%s" for skip 1, but got "%s"`,
			expectedPackage, actualPackage)
	}
	if actualFunction != expectedFunction {
		t.Errorf(`Expected function name "%s" for skip 1, but got "%s"`,
			expectedFunction, actualFunction)
	}
}

func TestGetCallerFunctionNameSkip3(t *testing.T) {
	testGetCallerFunctionNameSkip3_nestedFunc1(t)
}

func testGetCallerFunctionNameSkip3_nestedFunc1(t *testing.T) {
	testGetCallerFunctionNameSkip3_nestedFunc2(t)
}

func testGetCallerFunctionNameSkip3_nestedFunc2(t *testing.T) {
	const expectedPackage = "github.com/nyulibraries/go-ead-indexer/pkg/util"
	const expectedFunction = "TestGetCallerFunctionNameSkip3"

	actualPackage, actualFunction, err := GetCallerFunctionName(3)
	if err != nil {
		t.Errorf("Expected GetCallerFunction(3) to not return an error"+
			` but an error was returned: "%s""`, err.Error())
	}

	if actualPackage != expectedPackage {
		t.Errorf(`Expected package name "%s" for skip 3, but got "%s"`,
			expectedPackage, actualPackage)
	}
	if actualFunction != expectedFunction {
		t.Errorf(`Expected function name "%s" for skip 3, but got "%s"`,
			expectedFunction, actualFunction)
	}
}

func TestGetRepositoryCode(t *testing.T) {
	testCases := []struct {
		input              string
		expected           string
		expectedErrMessage string
	}{
		{"", "", "EAD file path must have at least two non-empty components, the last of which is a .xml file: ''"},
		{string(os.PathSeparator), "", fmt.Sprintf("EAD file path must have at least two non-empty components, the last of which is a .xml file: '%s'", string(os.PathSeparator))},
		{strings.Join([]string{"a", "b"}, string(os.PathSeparator)), "", fmt.Sprintf("EAD file path must have at least two non-empty components, the last of which is a .xml file: '%s'", strings.Join([]string{"a", "b"}, string(os.PathSeparator)))},
		{strings.Join([]string{"a", "b.xml"}, string(os.PathSeparator)), "a", ""},
		{string(os.PathSeparator) + strings.Join([]string{"a", "b", "c", "d.xml"}, string(os.PathSeparator)), "c", ""},
	}

	for _, testCase := range testCases {
		actual, err := GetRepositoryCode(testCase.input)
		if err != nil {
			if err.Error() != testCase.expectedErrMessage {
				t.Errorf("Expected error message '%s' for input '%s', but got '%s'",
					testCase.expectedErrMessage, testCase.input, err.Error())
			}
		} else {
			if actual != testCase.expected {
				t.Errorf("Expected repository code '%s' for input '%s', but got '%s'",
					testCase.expected, testCase.input, actual)
			}
		}
	}
}
