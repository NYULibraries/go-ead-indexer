package testutils

import (
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var eadFixturesDirPath string
var goldenFilesDirPath string
var testutilsPath string

// We need to get the absolute path to this package in order to enable the function
// for golden file and fixture file retrieval to be called from other packages
// which would not be able to resolve the hardcoded relative paths used here.
func init() {
	// The `filename` string is the absolute path to this source file, which should
	// be located at the root of the package directory.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("ERROR: `runtime.Caller(0)` failed")
	}

	// Get the path to the parent directory of this file.  Again, this is assuming
	// that this `init()` function is defined in a package top level file -- or
	// more precisely, that this file is in the same directory at the `testdata/`
	// directory that is referenced in the relative paths used in the functions
	// defined in this file.
	testutilsPath = filepath.Dir(filename)
	// Get testdata directory paths
	eadFixturesDirPath = filepath.Join(testutilsPath, "..", "testdata", "fixtures", "ead-files")
	goldenFilesDirPath = filepath.Join(testutilsPath, "..", "testdata", "golden")
}

func EadFixturePath(testEAD string) string {
	return filepath.Join(eadFixturesDirPath, testEAD+".xml")
}

func GetEADFixtureValue(testEAD string) (string, error) {
	return GetTestdataFileContents(EadFixturePath(testEAD))
}

// A "testEAD" is a repository code + "/" + EAD ID.  Example: "edip/mos_2024"
func GetGoldenFileIDs(testEAD string) []string {
	goldenFileIDs := []string{}

	err := filepath.WalkDir(filepath.Join(goldenFilesDirPath, testEAD),
		func(path string, dirEntry fs.DirEntry, err error) error {
			if !dirEntry.IsDir() && filepath.Ext(path) == ".xml" {
				goldenFileIDs = append(goldenFileIDs, strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
			}
			return nil
		})
	if err != nil {
		panic(err)
	}

	return goldenFileIDs
}

func GetGoldenFileValue(eadID string, fileID string) (string, error) {
	return GetTestdataFileContents(GoldenFilePath(eadID, fileID))
}

func GetTestdataFileContents(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return filename, err
	}

	return string(bytes), nil
}

func GetTestEADs() []string {
	testEADs := []string{}

	err := filepath.WalkDir(eadFixturesDirPath, func(path string, dirEntry fs.DirEntry, err error) error {
		if !dirEntry.IsDir() && filepath.Ext(path) == ".xml" {
			repositoryCode := filepath.Base(filepath.Dir(path))
			eadID := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			testEADs = append(testEADs, fmt.Sprintf("%s/%s", repositoryCode, eadID))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return testEADs
}

func GoldenFilePath(testEAD string, fileID string) string {
	return filepath.Join(goldenFilesDirPath, testEAD, fileID+".xml")
}

func ParseEADID(testEAD string) string {
	return filepath.Base(testEAD)
}

func ParseRepositoryCode(testEAD string) string {
	return filepath.Dir(testEAD)
}

func UpdateGoldenFile(testEAD string, fileID string, data string) error {
	return os.WriteFile(GoldenFilePath(testEAD, fileID), []byte(data), 0644)
}
