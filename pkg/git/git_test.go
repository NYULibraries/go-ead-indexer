package git

/*
  PLEASE NOTE: if you regenerate the "simple-repo" via the testsupport/gen-repo.bash script
               THE COMMIT HASHES WILL CHANGE! Therefore, you will need to update the hash
               values in the "scenarios" slice below.

a5ca6cca30fc08cfc13e4f1492dbfbbf3ec7cf63 2025-02-13 14:16:29 -0500 | Updating file fales/mss_001.xml (HEAD -> main) [jgpawletko]
33ac5f1415ac8fe611944bad4925528b62e845c8 2025-02-13 14:16:29 -0500 | Updating file archives/mc_1.xml, Deleting file fales/mss_002.xml EADID='mss_002', Updating file fales/mss_005.xml, Updating file tamwag/aia_002.xml [jgpawletko]
382c67e2ac64323e328506c85f97e229070a46cc 2025-02-13 14:16:29 -0500 | Updating file archives/cap_1.xml, Updating file fales/mss_004.xml, Updating file tamwag/aia_001.xml [jgpawletko]
2f531fc31b82cb128428c83e11d1e3f79b0da394 2025-02-13 14:16:29 -0500 | Updating file fales/mss_002.xml, Updating file fales/mss_003.xml [jgpawletko]
7e65f35361c9a2d7fc48bece8f04856b358620bf 2025-02-13 14:16:29 -0500 | Updating file fales/mss_001.xml [jgpawletko]

*/

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
)

var thisPath string
var gitRepoPathAbsolute string
var gitRepoPathRelative string
var gitRepoDotGitDirectory string
var gitRepoEnabledHiddenGitDirectory string

// this code was copied from the debug package, written by David Arjanik
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
	thisPath = filepath.Dir(filename)
	// Get testdata directory paths
	gitRepoPathAbsolute = filepath.Join(thisPath, "testdata", "fixtures", "git-repo")
	// This could be done as a const at top level, but assigning it here to keep
	// all this path stuff in one place.
	gitRepoPathRelative = filepath.Join(".", "testdata", "fixtures", "git-repo")
	gitRepoDotGitDirectory = filepath.Join(gitRepoPathAbsolute, "dot-git")
	gitRepoEnabledHiddenGitDirectory = filepath.Join(gitRepoPathAbsolute, ".git")
}

func TestListEADFilesForCommit(t *testing.T) {
	// cleanup any leftovers from interrupted tests
	deleteEnabledHiddenGitDirectory(t)

	createEnabledHiddenGitDirectory(t)
	defer deleteEnabledHiddenGitDirectory(t)

	scenarios := []struct {
		Hash       string
		Operations map[string]IndexerOperation
	}{
		{"a5ca6cca30fc08cfc13e4f1492dbfbbf3ec7cf63", map[string]IndexerOperation{"fales/mss_001.xml": Add}},
		{"33ac5f1415ac8fe611944bad4925528b62e845c8", map[string]IndexerOperation{"archives/mc_1.xml": Add, "fales/mss_002.xml": Delete, "fales/mss_005.xml": Add, "tamwag/aia_002.xml": Add}},
		{"382c67e2ac64323e328506c85f97e229070a46cc", map[string]IndexerOperation{"archives/cap_1.xml": Add, "fales/mss_004.xml": Add, "tamwag/aia_001.xml": Add}},
		{"2f531fc31b82cb128428c83e11d1e3f79b0da394", map[string]IndexerOperation{"fales/mss_002.xml": Add, "fales/mss_003.xml": Add}},
		{"7e65f35361c9a2d7fc48bece8f04856b358620bf", map[string]IndexerOperation{"fales/mss_001.xml": Add}},
	}

	for _, scenario := range scenarios {
		operations, err := ListEADFilesForCommit(gitRepoPathAbsolute, scenario.Hash)
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

func TestListEADFilesForCommitBadHash(t *testing.T) {
	// cleanup any leftovers from interrupted tests
	deleteEnabledHiddenGitDirectory(t)

	createEnabledHiddenGitDirectory(t)
	defer deleteEnabledHiddenGitDirectory(t)

	scenarios := []struct {
		Hash           string
		ExpectedErrMsg string
	}{
		{"e2e97a13e88e7a13a7c85f2c96293c7c2714a801", "problem getting commit object for commit hash e2e97a13e88e7a13a7c85f2c96293c7c2714a801: object not found"},
	}

	for _, scenario := range scenarios {
		_, err := ListEADFilesForCommit(gitRepoPathRelative, scenario.Hash)
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

func TestListEADFilesBadRepoPath(t *testing.T) {

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
		From         diff.File
		To           diff.File
		ExpectedOp   IndexerOperation
		ExpectedPath string
	}{
		{nil, nil, "unknown", ""},
		{nil, diffFileMock{ThisPath: "fales/mss_001.xml"}, "add", "fales/mss_001.xml"},
		{diffFileMock{ThisPath: "fales/mss_001.xml"}, nil, "delete", "fales/mss_001.xml"},
		{diffFileMock{ThisPath: "fales/mss_001.xml"}, diffFileMock{ThisPath: "fales/mss_001.xml"}, "add", "fales/mss_001.xml"},
		{diffFileMock{ThisPath: "fales/mss_001.xml"}, diffFileMock{ThisPath: "fales/mss_002.xml"}, "unknown", ""},
	}

	for _, scenario := range scenarios {
		result, indexerOp := classifyFileChange(scenario.From, scenario.To)
		if indexerOp != scenario.ExpectedOp {
			t.Errorf("expected operation '%s', got '%s'", scenario.ExpectedOp, indexerOp)
		}
		if result != scenario.ExpectedPath {
			t.Errorf("expected path '%s', got '%s'", scenario.ExpectedPath, result)
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

// ------------------------------------------------------------------------------
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

// ------------------------------------------------------------------------------
// Mock type for testing classifyFileChange function
// ------------------------------------------------------------------------------
type diffFileMock struct {
	ThisPath string
}

func (f diffFileMock) Hash() plumbing.Hash {
	return plumbing.NewHash("a5ca6cca30fc08cfc13e4f1492dbfbbf3ec7cf63")
}

func (f diffFileMock) Mode() filemode.FileMode {
	return filemode.Regular
}

func (f diffFileMock) Path() string {
	return f.ThisPath
}
