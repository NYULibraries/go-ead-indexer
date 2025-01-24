package git

/*
  PLEASE NOTE: if you regenerate the "simple-repo" via the testsupport/gen-repo.bash script
               THE COMMIT HASHES WILL CHANGE! Therefore, you will need to update the hash
               values in the "scenarios" slice below.

* ac3dde7f32f91ccca7dabb247e22ca131429f31d 2025-01-24 08:48:46 -0500 | Updating file fales/mss_001.xml (HEAD -> main) [jgpawletko]
* 0c56161d1bc28581c69ae93729ec2039117c3f00 2025-01-24 08:48:46 -0500 | Updating file archives/mc_1.xml, Deleting file fales/mss_002.xml EADID=mss_002, Updating file fales/mss_005.xml, Updating file tamwag/aia_002.xml [jgpawletko]
* bee0fd5241b8952fd1bca35cb0fc4314fb52652b 2025-01-24 08:48:46 -0500 | Updating file archives/cap_1.xml, Updating file fales/mss_004.xml, Updating file tamwag/aia_001.xml [jgpawletko]
* f860b602eeb315d5d142228a9e1fe72a818b4c4c 2025-01-24 08:48:46 -0500 | Updating file fales/mss_002.xml, Updating file fales/mss_003.xml [jgpawletko]
* 4565d2689e24beafee83e99652436f3de7eae738 2025-01-24 08:48:46 -0500 | Updating file fales/mss_001.xml [jgpawletko]

*/

import (
	"os"
	"testing"

	"github.com/c4milo/unpackit"
)

func TestListEADFilesForCommit(t *testing.T) {
	err := extractRepo("testdata/simple-repo.tar.gz", "testdata")
	if err != nil {
		t.Fatal(err)
	}
	defer teardownRepo("testdata/simple-repo")

	scenarios := []struct {
		Hash       string
		Operations map[string]IndexerOperation
	}{
		{"ac3dde7f32f91ccca7dabb247e22ca131429f31d", map[string]IndexerOperation{"fales/mss_001.xml": Add}},
		{"0c56161d1bc28581c69ae93729ec2039117c3f00", map[string]IndexerOperation{"archives/mc_1.xml": Add, "fales/mss_002.xml": Delete, "fales/mss_005.xml": Add, "tamwag/aia_002.xml": Add}},
		{"bee0fd5241b8952fd1bca35cb0fc4314fb52652b", map[string]IndexerOperation{"archives/cap_1.xml": Add, "fales/mss_004.xml": Add, "tamwag/aia_001.xml": Add}},
		{"f860b602eeb315d5d142228a9e1fe72a818b4c4c", map[string]IndexerOperation{"fales/mss_002.xml": Add, "fales/mss_003.xml": Add}},
		{"4565d2689e24beafee83e99652436f3de7eae738", map[string]IndexerOperation{"fales/mss_001.xml": Add}},
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
				t.Errorf("expected operation '%d' for file '%s', got '%d' for commit hash '%s'", expectedOp, file, op, scenario.Hash)
			}
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
