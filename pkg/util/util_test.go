package util

import (
	"slices"
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
