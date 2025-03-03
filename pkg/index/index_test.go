package index

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	eadtestutils "github.com/nyulibraries/go-ead-indexer/pkg/ead/testutils"
	"github.com/nyulibraries/go-ead-indexer/pkg/index/testutils"
)

func TestEADFileDoesNotExist(t *testing.T) {
	sc := testutils.GetSolrClientMock()
	SetSolrClient(sc)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Errorf("ERROR: `runtime.Caller(0)` failed")
		t.FailNow()
	}

	dir := filepath.Dir(filename)
	eadPath := filepath.Join(dir, "this-file-does-not-exist.xml")

	err := IndexEADFile(eadPath)
	if err == nil {
		t.Errorf("error: expected IndexEADFile to return an error, but nothing was returned: %v", err)
	}
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("error: expected IndexEADFile to return an error with message containing 'no such file or directory', but got: %v", err)
	}

	if sc.CallCount != 0 {
		t.Errorf("error: expected IndexEADFile to not call any SolrClient methods, but it did: %v", sc.CallCount)
	}

}

func TestAdd(t *testing.T) {

	//	var testFixturePath string
	repositoryCode := "fales"
	eadid := "mss_460"

	testEAD := filepath.Join(repositoryCode, eadid)
	var eadPath = eadtestutils.EadFixturePath(testEAD)

	sc := testutils.GetSolrClientMock()
	err := sc.InitMock(testEAD)
	if err != nil {
		t.Errorf("Error setting Solr client: %s", err)
		t.FailNow()
	}

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	err = IndexEADFile(eadPath)
	if err != nil {
		t.Errorf("Error indexing EAD file: %s", err)
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

	// specify which calls to Add() will return an error
	errorCallCounts := []int{11, 226, 333, 444, 555, 666, 777, 888, 999, 1000, 1208}

	repositoryCode := "nyhs"
	eadid := "ms347_foundling_hospital"
	testEAD := filepath.Join(repositoryCode, eadid)

	var eadPath = eadtestutils.EadFixturePath(testEAD)

	sc := testutils.GetSolrClientMock()
	err := sc.InitMock(testEAD)
	if err != nil {
		t.Errorf("Error setting Solr client: %s", err)
		t.FailNow()
	}

	// setup error events
	var errorEvents []testutils.ErrorEvent

	for _, errorCallCount := range errorCallCounts {
		errorEvents = append(errorEvents, testutils.ErrorEvent{CallerName: "Add", ErrorMessage: fmt.Sprintf("error during Add: %d", errorCallCount), CallCount: errorCallCount})
	}

	sc.ErrorEvents = testutils.SortErrorEvents(errorEvents)

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	err = IndexEADFile(eadPath)
	if err == nil {
		t.Errorf("error: expected IndexEADFile to return an error, but nothing was returned: %v", err)
		t.FailNow()
	}

	// check that the expected error message was returned
	for i, errString := range strings.Split(err.Error(), "\n") {
		if errString != errorEvents[i].ErrorMessage {
			t.Errorf("error: expected IndexEADFile to return an error with message '%s', but got: '%s'", errorEvents[i].ErrorMessage, errString)
		}
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
		t.Errorf("Rollback was not called at the expected time. Expected: %d, got: %d", sc.NumberOfFilesToIndex+2, sc.RollbackCallOrder)
	}
}
