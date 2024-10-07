package testutils

import (
	_ "embed"
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
	eadFixturesDirPath = filepath.Join(testutilsPath, "testdata", "fixtures", "ead-files")
	goldenFilesDirPath = filepath.Join(testutilsPath, "testdata", "golden")
}

func GetEADFixtureValue(eadID string) (string, error) {
	return GetTestdataFileContents(eadFixturePath(eadID))
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

func GetGoldenFileIDs(eadID string) []string {
	goldenFileIDs := []string{}

	err := filepath.WalkDir(filepath.Join(goldenFilesDirPath, eadID),
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

func GetTestEADIDs() []string {
	testEADIDs := []string{}

	err := filepath.WalkDir(eadFixturesDirPath, func(path string, dirEntry fs.DirEntry, err error) error {
		if !dirEntry.IsDir() && filepath.Ext(path) == ".xml" {
			testEADIDs = append(testEADIDs, strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return testEADIDs
}

func GoldenFilePath(eadID string, fileID string) string {
	return filepath.Join(goldenFilesDirPath, eadID, fileID+".xml")
}

func UpdateGoldenFile(eadID string, fileID string, data string) error {
	return os.WriteFile(GoldenFilePath(eadID, fileID), []byte(data), 0644)
}

func eadFixturePath(eadID string) string {
	return filepath.Join(eadFixturesDirPath, eadID+".xml")
}
