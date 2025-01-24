package git

/*
* b13375421fdfd7b417b5f3571bfeffea2d030547 2025-01-23 19:36:18 -0500 | Updating file archives/mc_1.xml, Deleting file fales/mss_002.xml EADID=mss_002, Updating file fales/mss_005.xml, Updating file tamwag/aia_002.xml (HEAD -> main) [jgpawletko]
* 58ba9870dcfd02a1ee95f30dcad9380e0bbf5f80 2025-01-23 19:36:18 -0500 | Updating file archives/cap_1.xml, Updating file fales/mss_004.xml, Updating file tamwag/aia_001.xml [jgpawletko]
* 8cdbf6b645cd89db709ea8a5196fab9cba194826 2025-01-23 19:36:18 -0500 | Updating file fales/mss_002.xml, Updating file fales/mss_003.xml [jgpawletko]
* 153ee1db908614837afa8edd29aec69060f6574b 2025-01-23 19:36:17 -0500 | Updating file fales/mss_001.xml [jgpawletko]
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

	operations, err := ListEADFilesForCommit("testdata/simple-repo", "b13375421fdfd7b417b5f3571bfeffea2d030547")
	if err != nil {
		t.Fatal(err)
	}
	if len(operations) != 4 {
		t.Errorf("expected 4 operations, got %d", len(operations))
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
