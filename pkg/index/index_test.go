package index

import (
	"math/rand"
	"path/filepath"
	"testing"

	"github.com/nyulibraries/go-ead-indexer/pkg/index/testutils"
)

func TestAdd(t *testing.T) {

	var testFixturePath string
	repositoryCode := "fales"
	eadid := "mss_460"

	pwd, err := filepath.Abs(".")
	if err != nil {
		t.Errorf("Error getting absolute path for pwd: %s", err)
		t.FailNow()
	}

	// pwd should be /root/path/to/go-ead-indexer/pkg/index/
	// need to get to: /root/path/to/go-ead-indexer/pkg/ead/testdata/fixtures/
	testFixturePath = filepath.Join(pwd, "..", "ead", "testdata")

	var eadPath = filepath.Join(testFixturePath, "fixtures", "ead-files", repositoryCode, eadid+".xml")
	var xmlDir = filepath.Join(testFixturePath, "golden", repositoryCode, eadid)

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

func TestRollbackOnBadAdd(t *testing.T) {

	var testFixturePath string
	repositoryCode := "nyhs"
	eadid := "ms347_foundling_hospital"

	pwd, err := filepath.Abs(".")
	if err != nil {
		t.Errorf("Error getting absolute path for pwd: %s", err)
		t.FailNow()
	}

	// pwd should be /root/path/to/go-ead-indexer/pkg/index/
	// need to get to: /root/path/to/go-ead-indexer/pkg/ead/testdata/fixtures/
	testFixturePath = filepath.Join(pwd, "..", "ead", "testdata")

	var eadPath = filepath.Join(testFixturePath, "fixtures", "ead-files", repositoryCode, eadid+".xml")
	var xmlDir = filepath.Join(testFixturePath, "golden", repositoryCode, eadid)

	sc := testutils.GetSolrClientMock()
	err = sc.InitMock(xmlDir)
	if err != nil {
		t.Errorf("Error setting Solr client: %s", err)
		t.FailNow()
	}

	// get a random number to simulate an error during Add
	randomNumber := rand.Intn(sc.NumberOfFilesToIndex) + 1

	// setup error events
	var errorEvents []testutils.ErrorEvent
	//	errorEvents = append(errorEvents, testutils.ErrorEvent{CallerName: "Delete", ErrorMessage: "error during initial Delete", CallCount: 1})
	errorEvents = append(errorEvents, testutils.ErrorEvent{CallerName: "Add", ErrorMessage: "error during Add", CallCount: randomNumber})
	errorEvents = append(errorEvents, testutils.ErrorEvent{CallerName: "Add", ErrorMessage: "error during Add", CallCount: randomNumber + 20})
	sc.ErrorEvents = errorEvents

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	errs := IndexEADFile(eadPath)
	if len(errs) == 0 {
		t.Errorf("error: expected IndexEADFile to return an error, but nothing was returned: %v", errs)
		t.FailNow()
	}

	// check that the expected error message was returned
	if errs[0].Error() != "error during Add" {
		t.Errorf("error: expected IndexEADFile to return an error with message 'error during initial Delete', but got: %s", errs[0].Error())
	}
	if errs[1].Error() != "error during Add" {
		t.Errorf("error: expected IndexEADFile to return an error with message 'error during initial Delete', but got: %s", errs[0].Error())
	}

	// check that delete was called first
	if sc.DeleteCallOrder != 1 {
		t.Errorf("Delete was not called first. Call order: %d", sc.DeleteCallOrder)
	}
	if sc.DeleteArgument != eadid {
		t.Errorf("Delete was not called with the correct argument. expected: %s, got: %s", eadid, sc.DeleteArgument)
	}

	// check that rollback was called in the expected sequence
	// the mock increments the call count before storing the value
	// so: delete + all Add() operations + rollback = Number of files to index + 2
	if sc.RollbackCallOrder != sc.NumberOfFilesToIndex+2 {
		t.Errorf("Rollback was not called at the expected time. Expected: 2, got: %d", sc.RollbackCallOrder)
	}
}
