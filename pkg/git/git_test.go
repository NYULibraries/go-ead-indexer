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
	"testing"

	"github.com/c4milo/unpackit"
)

func TestListEADFilesForCommit(t *testing.T) {
	// cleanup any leftovers from interrupted tests
	err := teardownRepo("testdata/simple-repo")
	if err != nil {
		t.Fatal(err)
	}

	err = extractRepo("testdata/simple-repo.tar.gz", "testdata")
	if err != nil {
		t.Fatal(err)
	}
	defer teardownRepo("testdata/simple-repo")

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
		operations, err := ListEADFilesForCommit("testdata/simple-repo", scenario.Hash)
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
	err := teardownRepo("testdata/simple-repo")
	if err != nil {
		t.Fatal(err)
	}

	err = extractRepo("testdata/simple-repo.tar.gz", "testdata")
	if err != nil {
		t.Fatal(err)
	}
	defer teardownRepo("testdata/simple-repo")

	scenarios := []struct {
		Hash           string
		ExpectedErrMsg string
	}{
		{"e2e97a13e88e7a13a7c85f2c96293c7c2714a801", "problem getting commit object for commit hash e2e97a13e88e7a13a7c85f2c96293c7c2714a801: object not found"},
	}

	for _, scenario := range scenarios {
		_, err := ListEADFilesForCommit("testdata/simple-repo", scenario.Hash)
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

func extractRepo(tarball, targetDir string) error {
	file, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer file.Close()

	err = unpackit.Unpack(file, targetDir)
	if err != nil {
		return err
	}

	return nil
}

func teardownRepo(targetDir string) error {
	err := os.RemoveAll(targetDir)
	if err != nil {
		return err
	}
	return nil
}
