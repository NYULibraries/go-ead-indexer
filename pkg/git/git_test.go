package git

/*
* ba7d9b00023a8d6ed962e46465d800265a6d06b9 2025-01-03 17:05:48 -0500 | Updating file archives/mc_1.xml, Deleting file fales/mss_002.xml EADID='mss_002', Updating file fales/mss_005.xml, Updating file tamwag/aia_002.xml (HEAD -> main) [jgpawletko]
* ca2dd426b9f22b52e101e71ce8db83c80508df06 2025-01-03 17:00:32 -0500 | Updating file archives/cap_1.xml, Updating file fales/mss_004.xml, Updating file tamwag/aia_001.xml [jgpawletko]
* ae25d50165da2befdfc21624ba52241ad36070de 2025-01-03 16:55:07 -0500 | Updating file fales/mss_002.xml, Updating file fales/mss_003.xml [jgpawletko]
* d9fa76ef7c89994d8d3ed458e5c06b2c5bb9f414 2025-01-03 16:54:16 -0500 | Updating file fales/mss_001.xml [jgpawletko]
 */

// TODO: update simple repo with non-null files
//       switch to using unpackit package https://github.com/c4milo/unpackit/blob/master/unpackit_test.go

import (
	"os"
	"testing"

	"github.com/walle/targz"
)

func TestListEADFilesForCommit(t *testing.T) {
	err := extractRepo("testdata/simple-repo.tar.gz", "testdata/simple-repo")
	if err != nil {
		t.Fatal(err)
	}

	operations, err := ListEADFilesForCommit("testdata/simple-repo", "ba7d9b00023a8d6ed962e46465d800265a6d06b9")
	if err != nil {
		t.Fatal(err)
	}
	if len(operations) != 4 {
		t.Fatalf("expected 4 operations, got %d", len(operations))
	}

	err = teardownRepo("testdata/simple-repo")
	if err != nil {
		t.Fatal(err)
	}
}

func extractRepo(tarball, targetDir string) error {
	err := targz.Extract(tarball, targetDir)
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
