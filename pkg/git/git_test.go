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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
)

var thisPath string
var gitSourceRepoPathAbsolute string
var gitRepoTestGitRepoPathAbsolute string
var gitRepoTestGitRepoPathRelative string
var gitRepoTestGitRepoDotGitDirectory string
var gitRepoTestGitRepoHiddenGitDirectory string

// this code is based on that in the debug package, written by David Arjanik
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
	gitSourceRepoPathAbsolute = filepath.Join(thisPath, "testdata", "fixtures", "git-repo")

	// This could be done as a const at top level, but assigning it here to keep
	// all this path stuff in one place.
	gitRepoTestGitRepoPathAbsolute = filepath.Join(thisPath, "testdata", "fixtures", "test-git-repo")
	gitRepoTestGitRepoPathRelative = filepath.Join(".", "testdata", "fixtures", "test-git-repo")
	gitRepoTestGitRepoDotGitDirectory = filepath.Join(gitRepoTestGitRepoPathAbsolute, "dot-git")
	gitRepoTestGitRepoHiddenGitDirectory = filepath.Join(gitRepoTestGitRepoPathAbsolute, ".git")
}

func TestCheckout(t *testing.T) {
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	scenarios := []struct {
		Hash         string
		FileRelPaths []string
	}{
		{"a5ca6cca30fc08cfc13e4f1492dbfbbf3ec7cf63", []string{"fales/mss_001.xml"}},
		{"33ac5f1415ac8fe611944bad4925528b62e845c8", []string{"archives/mc_1.xml", "fales/mss_005.xml", "tamwag/aia_002.xml"}},
		{"382c67e2ac64323e328506c85f97e229070a46cc", []string{"archives/cap_1.xml", "fales/mss_004.xml", "tamwag/aia_001.xml"}},
		{"2f531fc31b82cb128428c83e11d1e3f79b0da394", []string{"fales/mss_002.xml", "fales/mss_003.xml"}},
		{"7e65f35361c9a2d7fc48bece8f04856b358620bf", []string{"fales/mss_001.xml"}},
	}

	for _, scenario := range scenarios {
		err := Checkout(gitRepoTestGitRepoPathAbsolute, scenario.Hash)
		if err != nil {
			t.Errorf("unexpected error: %v for commit hash %s", err, scenario.Hash)
			continue
		}

		for _, fileRelPath := range scenario.FileRelPaths {
			filePath := filepath.Join(gitRepoTestGitRepoPathAbsolute, fileRelPath)
			_, err := os.Stat(filePath)
			if err != nil {
				t.Errorf("expected file '%s' to exist for commit hash '%s'", filePath, scenario.Hash)
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
	err := Checkout(gitRepoTestGitRepoPathAbsolute, badHash)
	if err == nil {
		t.Errorf("expected error but no error generated")
		return
	}

	exp := fmt.Sprintf("problem checking out hash '%s', error: 'reference not found'", badHash)
	if err.Error() != exp {
		t.Errorf("expected error message '%s', got '%s'", exp, err.Error())
	}
}

func TestCheckout_BadPath(t *testing.T) {
	err := Checkout("this-is-not-a-real-path", "33ac5f1415ac8fe611944bad4925528b62e845c8")
	if err == nil {
		t.Errorf("expected error but no error generated")
		return
	}
	exp := "repository does not exist"
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
		eadID, err := EADPathToEADID(scenario.Path)
		if err != nil {
			t.Errorf("unexpected error: %v for path %s", err, scenario.Path)
			continue
		}
		if eadID != scenario.ExpectedEADID {
			t.Errorf("expected EADID '%s', got '%s' for path '%s'", scenario.ExpectedEADID, eadID, scenario.Path)
		}
	}

	_, err := EADPathToEADID("this-is-not-a-real-path")
	if err == nil {
		t.Errorf("expected error but no error generated")
		return
	}
	exp := "unable to extract EADID from path 'this-is-not-a-real-path'"
	if err.Error() != exp {
		t.Errorf("expected error message '%s', got '%s'", exp, err.Error())
	}
}

func TestListEADFilesForCommit(t *testing.T) {
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	scenarios := []struct {
		Hash  string
		Steps []IndexerStep
	}{
		{"a5ca6cca30fc08cfc13e4f1492dbfbbf3ec7cf63", []IndexerStep{{Add, "fales/mss_001.xml"}}},
		{"33ac5f1415ac8fe611944bad4925528b62e845c8", []IndexerStep{{Add, "archives/mc_1.xml"}, {Delete, "fales/mss_002.xml"}, {Add, "fales/mss_005.xml"}, {Add, "tamwag/aia_002.xml"}}},
		{"382c67e2ac64323e328506c85f97e229070a46cc", []IndexerStep{{Add, "archives/cap_1.xml"}, {Add, "fales/mss_004.xml"}, {Add, "tamwag/aia_001.xml"}}},
		{"2f531fc31b82cb128428c83e11d1e3f79b0da394", []IndexerStep{{Add, "fales/mss_002.xml"}, {Add, "fales/mss_003.xml"}}},
		{"7e65f35361c9a2d7fc48bece8f04856b358620bf", []IndexerStep{{Add, "fales/mss_001.xml"}}},
	}

	for _, scenario := range scenarios {
		commitSteps, err := ListEADFilesForCommit(gitRepoTestGitRepoPathAbsolute, scenario.Hash)
		if err != nil {
			t.Errorf("unexpected error: %v for commit %s", err, scenario.Hash)
			continue
		}
		if len(commitSteps) != len(scenario.Steps) {
			t.Errorf("expected %d steps, got %d for commit hash %s", len(scenario.Steps), len(commitSteps), scenario.Hash)
			continue
		}
		for i, scenarioStep := range scenario.Steps {
			if scenarioStep.Operation != commitSteps[i].Operation {
				t.Errorf("expected operation '%s', got '%s'", scenarioStep.Operation, commitSteps[i].Operation)
			}
			if scenarioStep.FilePath != commitSteps[i].FilePath {
				t.Errorf("expected path '%s', got '%s'", scenarioStep.FilePath, commitSteps[i].FilePath)
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
			`deleteEnabledHiddenGitDirectory() failed with error "%s", remove %s manually`,
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
