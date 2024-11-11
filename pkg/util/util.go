package util

import (
	"go-ead-indexer/pkg/util/diff"
	"os"
)

// Replicate https://github.com/NYULibraries/ead_indexer/blob/a367ab8cc791376f0d8a287cbcd5b6ee43d5c04f/lib/ead_indexer/behaviors.rb#L137
// TODO: If it turns out that we don't need to preserve order in any of our uniq'ed
// slices, remove this.
func CompactStringSlicePreserveOrder(stringSlice []string) []string {
	compactedSlice := []string{}
	seen := map[string]bool{}
	for _, element := range stringSlice {
		if seen[element] {
			continue
		}
		compactedSlice = append(compactedSlice, element)
		seen[element] = true
	}

	return compactedSlice
}

func DiffFiles(path1 string, path2 string) (string, error) {
	bytes1, err := os.ReadFile(path1)
	if err != nil {
		return "", err
	}

	bytes2, err := os.ReadFile(path2)
	if err != nil {
		return "", err
	}

	diffBytes := diff.Diff(path1, bytes1, path2, bytes2)

	return string(diffBytes), nil
}

func DiffStrings(label1 string, string1 string, label2, string2 string) string {
	diffString := string(diff.Diff(label1, []byte(string1), label2, []byte(string2)))

	return diffString
}
