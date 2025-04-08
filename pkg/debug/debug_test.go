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

const testCommitHash = "b9cd08d511316cb311e0d9461aa29647c510087b"
const forEADFileGoldenFile = "for-ead-file.json"
const forGitCommitGoldenFile = "for-git-commit.json"

var debugPath string
var gitRepoSourcePath string
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
	gitRepoSourcePath = filepath.Join(debugPath, "testdata", "fixtures", "git-repo-source")
	gitRepoPathAbsolute = filepath.Join(debugPath, "testdata", "fixtures", "test-git-repo")
	// This could be done as a const at top level, but assigning it here to keep
	// all this path stuff in one place.
	gitRepoPathRelative = filepath.Join(".", "testdata", "fixtures", "test-git-repo")
	gitRepoDotGitDirectory = filepath.Join(gitRepoPathAbsolute, "dot-git")
	gitRepoEnabledHiddenGitDirectory = filepath.Join(gitRepoPathAbsolute, ".git")
	goldenFileDirPath = filepath.Join(debugPath, "testdata", "golden")
}

func TestDumpSolrIndexerHTTPRequestsForEADFile(t *testing.T) {
	resetRepo(t)
	defer func() { deleteRepo(t) }()

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
	t.Run("Absolute git repo path", func(t *testing.T) {
		testDumpSolrIndexerHTTPRequestsForGitCommit(gitRepoPathAbsolute, t)
	})
	t.Run("Relative git repo path", func(t *testing.T) {
		testDumpSolrIndexerHTTPRequestsForGitCommit(gitRepoPathRelative, t)
	})
}

func testDumpSolrIndexerHTTPRequestsForGitCommit(gitRepoPath string, t *testing.T) {
	resetRepo(t)
	defer func() { deleteRepo(t) }()

	goldenFilePath := path.Join(goldenFileDirPath, forGitCommitGoldenFile)

	actual, err := DumpSolrIndexerHTTPRequestsForGitCommit(gitRepoPath, testCommitHash)
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

func createRepo(t *testing.T) {
	gitRepoSourcePathFS := os.DirFS(gitRepoSourcePath)
	err := os.CopyFS(gitRepoPathAbsolute, gitRepoSourcePathFS)
	if err != nil {
		t.Errorf(
			`Unexpected error returned by os.CopyFS(gitRepoPathAbsolute, gitRepoSourcePathFS): %s`,
			err.Error())
		t.FailNow()
	}

	err = os.Rename(gitRepoDotGitDirectory, gitRepoEnabledHiddenGitDirectory)
	if err != nil {
		t.Errorf(
			`Unexpected error returned by os.Rename(gitRepoDotGitDirectory, gitRepoEnabledHiddenGitDirectory): %s`,
			err.Error())
		t.FailNow()
	}
}

func deleteRepo(t *testing.T) {
	err := os.RemoveAll(gitRepoPathAbsolute)
	if err != nil {
		t.Errorf(
			`deleteRepo() failed with error "%s", remove %s manually`,
			err.Error(), gitRepoPathAbsolute)
		t.FailNow()
	}
}

func resetRepo(t *testing.T) {
	deleteRepo(t)
	createRepo(t)
}
