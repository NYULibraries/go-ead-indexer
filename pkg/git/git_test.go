//go:generate testsupport/update-test-fixtures.bash

package git

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/format/diff"

	"github.com/nyulibraries/go-ead-indexer/pkg/ead/eadutil"
)

var thisPath string
var gitSourceRepoPathAbsolute string
var gitRepoTestGitRepoPathAbsolute string
var gitRepoTestGitRepoPathRelative string
var gitRepoTestGitRepoDotGitDirectory string
var gitRepoTestGitRepoHiddenGitDirectory string

// this code is based on that in the debug package, written by David Arjanik
// We need to get the absolute path to this package in order to enable the
// function for golden file and fixture file retrieval to be called from other
// packages which would not be able to resolve the hardcoded relative paths
// used here.
func init() {
	// The `filename` string is the absolute path to this source file,
	// which should be located at the root of the package directory.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("ERROR: `runtime.Caller(0)` failed")
	}

	// Get the path to the parent directory of this file.  Again, this is assuming
	// that this `init()` function is defined in a package top level file -- or
	// more precisely, that this file is in the same directory at the `testdata/`
	// directory that is referenced in the relative paths used in the functions
	// defined in this file.
	thisPath = filepath.Dir(filename)

	// Get testdata directory paths
	gitSourceRepoPathAbsolute = filepath.Join(thisPath, "testdata", "fixtures", "git-repo")

	// This could be done as a const at top level, but assigning it here to
	// keep all this path stuff in one place.
	gitRepoTestGitRepoPathAbsolute = filepath.Join(thisPath, "testdata", "fixtures", "test-git-repo")
	gitRepoTestGitRepoPathRelative = filepath.Join(".", "testdata", "fixtures", "test-git-repo")
	gitRepoTestGitRepoDotGitDirectory = filepath.Join(gitRepoTestGitRepoPathAbsolute, "dot-git")
	gitRepoTestGitRepoHiddenGitDirectory = filepath.Join(gitRepoTestGitRepoPathAbsolute, ".git")
}

func TestCheckoutMergeReset(t *testing.T) {
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	scenarios := []struct {
		Hash         string
		FileRelPaths []string
	}{
		{Commit4Hash,
			[]string{".circleci/config.yml",
				"README.md",
				"archives/cap_1.xml",
				"archives/mc_1.xml",
				"fales/mss_001.xml",
				"fales/mss_003.xml",
				"fales/mss_004.xml",
				"fales/mss_005.xml",
				"tamwag/aia_001.xml",
				"tamwag/aia_002.xml"},
		},
		{Commit3Hash,
			[]string{".circleci/config.yml",
				"README.md",
				"archives/cap_1.xml",
				"fales/mss_001.xml",
				"fales/mss_002.xml",
				"fales/mss_003.xml",
				"fales/mss_004.xml",
				"tamwag/aia_001.xml"},
		},
		{Commit2Hash,
			[]string{".circleci/config.yml",
				"README.md",
				"fales/mss_001.xml",
				"fales/mss_002.xml",
				"fales/mss_003.xml"},
		},
		{Commit1Hash,
			[]string{".circleci/config.yml",
				"README.md",
				"fales/mss_001.xml"},
		},
	}

	for _, scenario := range scenarios {
		err := CheckoutMergeReset(gitRepoTestGitRepoPathAbsolute, scenario.Hash)
		if err != nil {
			t.Errorf("unexpected error: %v for commit hash %s", err,
				scenario.Hash)
			continue
		}

		actualFileRelPaths := make([]string, 0)
		err = filepath.Walk(gitRepoTestGitRepoPathAbsolute,
			func(path string, info os.FileInfo, err error) error {
				if info.IsDir() &&
					path == gitRepoTestGitRepoHiddenGitDirectory {
					return filepath.SkipDir
				}
				if !info.IsDir() {
					relativePath, _ := strings.CutPrefix(path,
						gitRepoTestGitRepoPathAbsolute+
							string(os.PathSeparator))
					actualFileRelPaths = append(actualFileRelPaths,
						relativePath)
				}
				return nil
			})

		if err != nil {
			t.Errorf("error walking directory: %v", err)
			t.FailNow()
		}

		if len(actualFileRelPaths) != len(scenario.FileRelPaths) {
			t.Errorf("expected %d files, got %d for commit hash %s",
				len(scenario.FileRelPaths),
				len(actualFileRelPaths),
				scenario.Hash)
			t.Errorf("expected files: %v", scenario.FileRelPaths)
			t.Errorf("actual files: %v", actualFileRelPaths)
			// use continue here to avoid a panic in the actual vs. expected
			// file comparison below
			continue
		}

		// compare the actual file paths to the expected file paths
		for i, fileRelPath := range scenario.FileRelPaths {
			if actualFileRelPaths[i] != fileRelPath {
				t.Errorf("expected file '%s', got '%s' for commit hash '%s'",
					fileRelPath, actualFileRelPaths[i], scenario.Hash)
			}
		}
	}
}

func TestCheckout_BadHash(t *testing.T) {
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	badHash := "this is not a valid hash"
	err := CheckoutMergeReset(gitRepoTestGitRepoPathAbsolute, badHash)
	if err == nil {
		t.Errorf("expected error but no error generated")
		return
	}

	exp := fmt.Sprintf(errNotAValidCommitHashStringTemplate, badHash)
	if err.Error() != exp {
		t.Errorf("expected error message '%s', got '%s'", exp, err.Error())
	}
}

func TestCheckout_BadPath(t *testing.T) {
	err := CheckoutMergeReset("this-is-not-a-real-path", "33ac5f1415ac8fe611944bad4925528b62e845c8")
	if err == nil {
		t.Errorf("expected error but no error generated")
		return
	}
	exp := "repository does not exist"
	if err.Error() != exp {
		t.Errorf("expected error message '%s', got '%s'", exp, err.Error())
	}
}

func TestCheckout_CommitHashNotFound(t *testing.T) {
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	commitHashNotFound := "1111111111111111111111111111111111111111"
	err := CheckoutMergeReset(gitRepoTestGitRepoPathAbsolute, commitHashNotFound)
	if err == nil {
		t.Errorf("expected error but no error generated")
		return
	}

	exp := "problem checking out hash '1111111111111111111111111111111111111111', error: 'object not found'"
	if err.Error() != exp {
		t.Errorf("expected error message '%s', got '%s'", exp, err.Error())
	}
}

func TestEADPathToEADID(t *testing.T) {
	scenarios := []struct {
		Path          string
		ExpectedEADID string
	}{
		{"fales/mss_001.xml", "mss_001"},
		{"archives/mc_1.xml", "mc_1"},
		{"fales/mss_002.xml", "mss_002"},
		{"fales/mss_005.xml", "mss_005"},
		{"tamwag/aia_002.xml", "aia_002"},
	}

	for _, scenario := range scenarios {
		eadID, err := eadutil.EADPathToEADID(scenario.Path)
		if err != nil {
			t.Errorf("unexpected error: %v for path %s", err, scenario.Path)
			continue
		}
		if eadID != scenario.ExpectedEADID {
			t.Errorf("expected EADID '%s', got '%s' for path '%s'", scenario.ExpectedEADID, eadID, scenario.Path)
		}
	}

	_, err := eadutil.EADPathToEADID("this-is-not-a-real-path")
	if err == nil {
		t.Errorf("expected error but no error generated")
		return
	}
	expectedError := "invalid EADID: this-is-not-a-real-path"
	if err.Error() != expectedError {
		t.Errorf("expected error message '%s', got '%s'",
			expectedError, err.Error())
	}
}

func TestListEADFilesForCommit(t *testing.T) {
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	scenarios := []struct {
		Hash       string
		Operations map[string]IndexerOperation
	}{
		{Commit10Hash, map[string]IndexerOperation{
			"archives/cap_001.xml": Add,
		}},
		{Commit9Hash, map[string]IndexerOperation{
			"archives/cap_001.xml": Delete,
		}},
		{Commit8Hash, map[string]IndexerOperation{
			"archives/mc_1.xml":    Delete,
			"archives/mc_001.xml":  Add,
			"archives/cap_1.xml":   Delete,
			"archives/cap_001.xml": Add,
		}},
		{Commit7Hash, map[string]IndexerOperation{"archives/mc_1.xml": Add}},
		{Commit6Hash, map[string]IndexerOperation{}},
		{Commit5Hash, map[string]IndexerOperation{"fales/mss_001.xml": Add}},
		{Commit4Hash, map[string]IndexerOperation{"archives/mc_1.xml": Add, "fales/mss_002.xml": Delete, "fales/mss_005.xml": Add, "tamwag/aia_002.xml": Add}},
		{Commit3Hash, map[string]IndexerOperation{"archives/cap_1.xml": Add, "fales/mss_004.xml": Add, "tamwag/aia_001.xml": Add}},
		{Commit2Hash, map[string]IndexerOperation{"fales/mss_002.xml": Add, "fales/mss_003.xml": Add}},
		{Commit1Hash, map[string]IndexerOperation{"fales/mss_001.xml": Add}},
	}

	for _, scenario := range scenarios {
		operations, err := ListEADFilesForCommit(gitRepoTestGitRepoPathAbsolute, scenario.Hash)
		if err != nil {
			t.Errorf("unexpected error: %v for commit hash %s", err, scenario.Hash)
			continue
		}
		if len(operations) != len(scenario.Operations) {
			t.Errorf("expected %d operations, got %d for commit hash %s", len(scenario.Operations), len(operations), scenario.Hash)
			continue
		}
		for file, expectedOp := range scenario.Operations {
			op, ok := operations[file]
			if !ok {
				t.Errorf("missing operation for file '%s' for commit hash '%s'", file, scenario.Hash)
			}
			if op != expectedOp {
				t.Errorf("expected operation '%s' for file '%s', got '%s' for commit hash '%s'", expectedOp, file, op, scenario.Hash)
			}
		}
	}
}

func TestListEADFilesForCommit_BadHash(t *testing.T) {
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	scenarios := []struct {
		Hash           string
		ExpectedErrMsg string
	}{
		{"e2e97a13e88e7a13a7c85f2c96293c7c2714a801", "problem getting commit object for commit hash e2e97a13e88e7a13a7c85f2c96293c7c2714a801: object not found"},
	}

	for _, scenario := range scenarios {
		_, err := ListEADFilesForCommit(gitRepoTestGitRepoPathRelative, scenario.Hash)
		if err == nil {
			t.Errorf("expected error but no error generated for commit hash %s", scenario.Hash)
			continue
		}
		if err.Error() != scenario.ExpectedErrMsg {
			t.Errorf("expected error message '%s' but got error message '%s' for hash '%s'", scenario.ExpectedErrMsg, err.Error(), scenario.Hash)
			continue
		}
	}
}

func TestListEADFilesForCommit_BadRepoPath(t *testing.T) {

	scenarios := []struct {
		ExpectedErrMsg string
	}{
		{"repository does not exist"},
	}

	for _, scenario := range scenarios {
		_, err := ListEADFilesForCommit("this-is-not-a-real-path", "a5ca6cca30fc08cfc13e4f1492dbfbbf3ec7cf63")
		if err == nil {
			t.Errorf("expected error but no error generated")
			continue
		}
		if err.Error() != scenario.ExpectedErrMsg {
			t.Errorf("expected error message '%s' but got error message '%s'", scenario.ExpectedErrMsg, err.Error())
			continue
		}
	}
}

func Test_classifyFileChange(t *testing.T) {

	scenarios := []struct {
		From       diff.File
		To         diff.File
		ExpectedOp IndexerOperation
	}{
		{nil, nil, Unknown},
		{nil, diffFileMock{ThisPath: "fales/mss_001.xml"}, Add},
		{diffFileMock{ThisPath: "fales/mss_001.xml"}, nil, Delete},
		{diffFileMock{ThisPath: "fales/mss_001.xml"}, diffFileMock{ThisPath: "fales/mss_001.xml"}, Add},
		{diffFileMock{ThisPath: "fales/mss_001.xml"}, diffFileMock{ThisPath: "fales/mss_002.xml"}, Rename},
	}

	for _, scenario := range scenarios {
		indexerOp := classifyFileChange(scenario.From, scenario.To)
		if indexerOp != scenario.ExpectedOp {
			t.Errorf("expected operation '%s', got '%s'", scenario.ExpectedOp, indexerOp)
		}
	}

}

func Test_getPath(t *testing.T) {

	scenarios := []struct {
		File         diff.File
		ExpectedPath string
	}{
		{nil, ""},
		{diffFileMock{ThisPath: "fales/mss_001.xml"}, "fales/mss_001.xml"},
	}

	for _, scenario := range scenarios {
		result := getPath(scenario.File)
		if result != scenario.ExpectedPath {
			t.Errorf("expected path '%s', got '%s'", scenario.ExpectedPath, result)
		}
	}

}

func createTestGitRepo(t *testing.T) {
	gitSourceRepoPathAbsoluteFS := os.DirFS(gitSourceRepoPathAbsolute)
	err := os.CopyFS(gitRepoTestGitRepoPathAbsolute, gitSourceRepoPathAbsoluteFS)
	if err != nil {
		t.Errorf(
			`Unexpected error returned by os.CopyFS(gitSourceRepoPathAbsoluteFS, gitRepoTestGitRepoPathAbsolute): %s`,
			err.Error())
		t.FailNow()
	}

	err = os.Rename(gitRepoTestGitRepoDotGitDirectory, gitRepoTestGitRepoHiddenGitDirectory)
	if err != nil {
		t.Errorf(
			`Unexpected error returned by os.Rename(gitRepoTestGitRepoDotGitDirectory, gitRepoTestGitRepoHiddenGitDirectory): %s`,
			err.Error())
		t.FailNow()
	}
}

func deleteTestGitRepo(t *testing.T) {
	err := os.RemoveAll(gitRepoTestGitRepoPathAbsolute)
	if err != nil {
		t.Errorf(
			`deleteTestGitRepo() failed with error "%s", remove %s manually`,
			err.Error(), gitRepoTestGitRepoPathAbsolute)
		t.FailNow()
	}
}

// ------------------------------------------------------------------------------
// Mock type for testing classifyFileChange function
// ------------------------------------------------------------------------------
type diffFileMock struct {
	ThisPath string
}

func (f diffFileMock) Hash() plumbing.Hash {
	// dummy value
	return plumbing.NewHash("a5ca6cca30fc08cfc13e4f1492dbfbbf3ec7cf63")
}

func (f diffFileMock) Mode() filemode.FileMode {
	// dummy value
	return filemode.Regular
}

func (f diffFileMock) Path() string {
	return f.ThisPath
}
