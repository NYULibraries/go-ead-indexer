package git

/*
  PLEASE NOTE: if you regenerate the "simple-repo" via the testsupport/gen-repo.bash script
               THE COMMIT HASHES WILL CHANGE! Therefore, you will need to update the hash
               values in the "scenarios" slice below.

 * 95f2f904ad261e7d31632021fa10768d2b4096c9 2025-01-24 17:10:44 -0500 | Updating file fales/mss_001.xml (HEAD -> main) [jgpawletko]
 * aa58b2314e11ae5af61129ebfe1ceb07b49c2d33 2025-01-24 17:10:44 -0500 | Updating file archives/mc_1.xml, Deleting file fales/mss_002.xml EADID='mss_002', Updating file fales/mss_005.xml, Updating file tamwag/aia_002.xml [jgpawletko]
 * 3dc6fabe0fcd990e95cdd3f88cff821196fccdbd 2025-01-24 17:10:44 -0500 | Updating file archives/cap_1.xml, Updating file fales/mss_004.xml, Updating file tamwag/aia_001.xml [jgpawletko]
 * 7fe6de7c56d30149889f8d24eaf2fa66ed9f2e2d 2025-01-24 17:10:44 -0500 | Updating file fales/mss_002.xml, Updating file fales/mss_003.xml [jgpawletko]
 * 155309f674b5acffd7473c1648f3647a2a3d242b 2025-01-24 17:10:44 -0500 | Updating file fales/mss_001.xml [jgpawletko]

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
		{"95f2f904ad261e7d31632021fa10768d2b4096c9", map[string]IndexerOperation{"fales/mss_001.xml": Add}},
		{"aa58b2314e11ae5af61129ebfe1ceb07b49c2d33", map[string]IndexerOperation{"archives/mc_1.xml": Add, "fales/mss_002.xml": Delete, "fales/mss_005.xml": Add, "tamwag/aia_002.xml": Add}},
		{"3dc6fabe0fcd990e95cdd3f88cff821196fccdbd", map[string]IndexerOperation{"archives/cap_1.xml": Add, "fales/mss_004.xml": Add, "tamwag/aia_001.xml": Add}},
		{"7fe6de7c56d30149889f8d24eaf2fa66ed9f2e2d", map[string]IndexerOperation{"fales/mss_002.xml": Add, "fales/mss_003.xml": Add}},
		{"155309f674b5acffd7473c1648f3647a2a3d242b", map[string]IndexerOperation{"fales/mss_001.xml": Add}},
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
