package debug

import (
	"flag"
	eadtestutils "github.com/nyulibraries/go-ead-indexer/pkg/ead/testutils"
	"github.com/nyulibraries/go-ead-indexer/pkg/util"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

const commitHash = "cb2edf315b46e969b65f67806b6d9ce612d1e275"
const forEADFileGoldenFile = "for-ead-file.json"
const forGitCommitGoldenFile = "for-git-commit.json"

var debugPath string
var gitRepoPathAbsolute string
var gitRepoPathRelative string
var gitRepoDotGitDirectory string
var gitRepoEnabledHiddenGitDirectory string
var goldenFileDirPath string

var updateGoldenFiles = flag.Bool("update-golden-files", false, "update the golden files")

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
	debugPath = filepath.Dir(filename)
	// Get testdata directory paths
	gitRepoPathAbsolute = filepath.Join(debugPath, "testdata", "fixtures", "git-repo")
	// This could be done as a const at top level, but assigning it here to keep
	// all this path stuff in one place.
	gitRepoPathRelative = filepath.Join(".", "testdata", "fixtures", "git-repo")
	gitRepoDotGitDirectory = filepath.Join(gitRepoPathAbsolute, "dot-git")
	gitRepoEnabledHiddenGitDirectory = filepath.Join(gitRepoPathAbsolute, ".git")
	goldenFileDirPath = filepath.Join(debugPath, "testdata", "golden")
}

func TestDumpSolrIndexerHTTPRequestsForEADFile(t *testing.T) {
	eadFilePath := eadtestutils.EadFixturePath("edip/mos_2024")

	goldenFilePath := path.Join(goldenFileDirPath, forEADFileGoldenFile)

	actual, err := DumpSolrIndexerHTTPRequestsForEADFile(eadFilePath)
	if err != nil {
		t.Errorf(`Unexpected error returned by DumpSolrIndexerHTTPRequestsForEADFile(): %s`,
			err.Error())
	}

	if *updateGoldenFiles {
		err := os.WriteFile(goldenFilePath, []byte(actual), 0644)
		if err != nil {
			t.Fatalf("Error updating golden file: %s", err)
		}
	}

	goldenBytes, err := os.ReadFile(goldenFilePath)
	if err != nil {
		t.Errorf(`Error retrieving golden file data: %s`, err.Error())
	}
	expected := string(goldenBytes)

	if actual != expected {
		diff := util.DiffStrings(forEADFileGoldenFile,
			expected,
			"actual",
			actual,
		)
		t.Errorf(`Golden and actual dumped HTTP requests JSON do not match: %s`,
			diff)
	}
}

func TestDumpSolrIndexerHTTPRequestsForGitCommit(t *testing.T) {
	deleteEnabledHiddenGitDirectory(t)

	createEnabledHiddenGitDirectory(t)

	t.Run("Absolute git repo path", func(t *testing.T) {
		testDumpSolrIndexerHTTPRequestsForGitCommit(gitRepoPathAbsolute, t)
	})
	t.Run("Relative git repo path", func(t *testing.T) {
		testDumpSolrIndexerHTTPRequestsForGitCommit(gitRepoPathRelative, t)
	})

	deleteEnabledHiddenGitDirectory(t)
}

func testDumpSolrIndexerHTTPRequestsForGitCommit(gitRepoPath string, t *testing.T) {
	goldenFilePath := path.Join(goldenFileDirPath, forGitCommitGoldenFile)

	actual, err := DumpSolrIndexerHTTPRequestsForGitCommit(gitRepoPath, commitHash)
	if err != nil {
		t.Errorf(`Unexpected error returned by DumpSolrIndexerHTTPRequestsForGitCommit(): %s`,
			err.Error())
	}

	if *updateGoldenFiles {
		err := os.WriteFile(goldenFilePath, []byte(actual), 0644)
		if err != nil {
			t.Fatalf("Error updating golden file: %s", err)
		}
	}

	goldenBytes, err := os.ReadFile(goldenFilePath)
	if err != nil {
		t.Errorf(`Error retrieving golden file data: %s`, err.Error())
	}
	expected := string(goldenBytes)

	if actual != expected {
		diff := util.DiffStrings(forEADFileGoldenFile,
			expected,
			"actual",
			actual,
		)
		t.Errorf(`Golden and actual dumped HTTP requests JSON do not match: %s`,
			diff)
	}
}

func createEnabledHiddenGitDirectory(t *testing.T) {
	gitRepoDotGitDirectoryFS := os.DirFS(gitRepoDotGitDirectory)
	err := os.CopyFS(gitRepoEnabledHiddenGitDirectory, gitRepoDotGitDirectoryFS)
	if err != nil {
		t.Errorf(
			`Unexpected error returned by os.CopyFS(gitRepoEnabledHiddenGitDirectory, gitRepoDotGitDirectoryFS): %s`,
			err.Error())
		t.FailNow()
	}
}

func deleteEnabledHiddenGitDirectory(t *testing.T) {
	err := os.RemoveAll(gitRepoEnabledHiddenGitDirectory)
	if err != nil {
		t.Errorf(
			`deleteEnabledHiddenGitDirectory() failed with error "%s", remove %s manually`,
			err.Error(), gitRepoEnabledHiddenGitDirectory)
		t.FailNow()
	}
}
