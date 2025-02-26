package index

import (
	"path/filepath"
	"testing"

	"github.com/nyulibraries/go-ead-indexer/pkg/index/testutils"
)

func TestAdd(t *testing.T) {

	var testFixturePath string
	partner := "fales"
	eadid := "mss_460"

	pwd, err := filepath.Abs(".")
	if err != nil {
		t.Errorf("Error getting absolute path for pwd: %s", err)
		t.FailNow()
	}

	// pwd should be /root/path/to/go-ead-indexer/pkg/index/
	// need to get to: /root/path/to/go-ead-indexer/pkg/ead/testdata/fixtures/
	testFixturePath = filepath.Join(pwd, "..", "ead", "testdata")

	var eadPath = filepath.Join(testFixturePath, "fixtures", "ead-files", partner, eadid+".xml")
	var xmlDir = filepath.Join(testFixturePath, "golden", partner, eadid)

	sc := testutils.GetSolrClientMock()
	err = sc.InitMock(xmlDir)
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

	// check that delete was called first
	if sc.DeleteCallOrder != 1 {
		t.Errorf("Delete was not called first. Call order: %d", sc.DeleteCallOrder)
	}
	if sc.DeleteArgument != eadid {
		t.Errorf("Delete was not called with the correct argument. expected: %s, got: %s", eadid, sc.DeleteArgument)
	}

	// check that commit was called in the expected sequence
	// the mock increments the call count before storing the value
	// so: delete + number of files + commit = number of files + 2
	if sc.CommitCallOrder != sc.NumberOfFilesToIndex+2 {
		t.Errorf("Commit was not called at the expected time. Expected: %d, got: %d", sc.NumberOfFilesToIndex+1, sc.CommitCallOrder)
	}
}

func TestRollback(t *testing.T) {

	var testFixturePath string
	partner := "nyhs"
	eadid := "ms347_foundling_hospital"

	pwd, err := filepath.Abs(".")
	if err != nil {
		t.Errorf("Error getting absolute path for pwd: %s", err)
		t.FailNow()
	}

	// pwd should be /root/path/to/go-ead-indexer/pkg/index/
	// need to get to: /root/path/to/go-ead-indexer/pkg/ead/testdata/fixtures/
	testFixturePath = filepath.Join(pwd, "..", "ead", "testdata")

	var eadPath = filepath.Join(testFixturePath, "fixtures", "ead-files", partner, eadid+".xml")
	var xmlDir = filepath.Join(testFixturePath, "golden", partner, eadid)

	sc := testutils.GetSolrClientMock()
	err = sc.InitMock(xmlDir)
	if err != nil {
		t.Errorf("Error setting Solr client: %s", err)
		t.FailNow()
	}

	// setup error events
	var errorEvents []testutils.ErrorEvent
	errorEvents = append(errorEvents, testutils.ErrorEvent{CallerName: "Delete", ErrorMessage: "error during initial Delete", CallCount: 1})
	sc.ErrorEvents = errorEvents

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

	// check that delete was called first
	if sc.DeleteCallOrder != 1 {
		t.Errorf("Delete was not called first. Call order: %d", sc.DeleteCallOrder)
	}
	if sc.DeleteArgument != eadid {
		t.Errorf("Delete was not called with the correct argument. expected: %s, got: %s", eadid, sc.DeleteArgument)
	}

	// check that commit was called in the expected sequence
	// the mock increments the call count before storing the value
	// so: delete + number of files + commit = number of files + 2
	if sc.CommitCallOrder != sc.NumberOfFilesToIndex+2 {
		t.Errorf("Commit was not called at the expected time. Expected: %d, got: %d", sc.NumberOfFilesToIndex+1, sc.CommitCallOrder)
	}
}
