package index

import (
	"go-ead-indexer/pkg/index/testutils"
	"path/filepath"
	"testing"
)

func TestAdd(t *testing.T) {

	var testFixturePath string

	pwd, err := filepath.Abs(".")
	if err != nil {
		t.Errorf("Error getting absolute path for pwd: %s", err)
		t.FailNow()
	}

	// pwd should be /root/path/to/go-ead-indexer/pkg/index/
	// need to get to: /root/path/to/go-ead-indexer/pkg/ead/testdata/fixtures/
	testFixturePath = filepath.Join(pwd, "..", "ead", "testdata")

	var eadPath = filepath.Join(testFixturePath, "fixtures", "ead-files", "fales", "mss_460.xml")
	var xmlDir = filepath.Join(testFixturePath, "golden", "fales", "mss_460")

	sc := testutils.GetSolrClientMock()
	err = sc.SetupMock(xmlDir)
	if err != nil {
		t.Errorf("Error setting Solr client: %s", err)
		t.FailNow()
	}

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	errs := IndexEADFile(eadPath)
	if len(errs) > 0 {
		t.Errorf("Error indexing EAD file: %s", errs)
	}

	// Check if the operation is complete from the Solr client perspective
	if !sc.IsComplete() {
		t.Errorf("Not all files were added to the Solr index. Remaining values: %v", sc.GoldenFileHashes)
	}
}
