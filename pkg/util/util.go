package util

import (
	"go-ead-indexer/pkg/util/diff"
	"os"
)

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
