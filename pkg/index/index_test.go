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

func TestIndexEADFile_EADFileDoesNotExist(t *testing.T) {

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

func TestIndexEADFile_Success(t *testing.T) {

	repositoryCode := "fales"
	eadid := "mss_460"
	testEAD := filepath.Join(repositoryCode, eadid)

	var eadPath = eadtestutils.EadFixturePath(testEAD)

	sc := testutils.GetSolrClientMock()

	// load the Solr POST body expectations
	err := sc.InitMockForIndexing(testEAD)
	if err != nil {
		t.Errorf("Error setting Solr client: %s", err)
		t.FailNow()
	}

	// set expectations
	sc.ExpectedCallOrder.Delete = 1
	sc.ExpectedCallOrder.Commit = sc.NumberOfFilesToIndex + 2
	sc.ExpectedDeleteArgument = eadid

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	err = IndexEADFile(eadPath)
	if err != nil {
		t.Errorf("Error indexing EAD file: %s", err)
	}

	err = sc.CheckAssertions()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}
}

func TestIndexEADFile_RollbackOnBadDelete(t *testing.T) {

	repositoryCode := "fales"
	eadid := "mss_460"
	testEAD := filepath.Join(repositoryCode, eadid)

	var eadPath = eadtestutils.EadFixturePath(testEAD)

	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForIndexing(testEAD)
	if err != nil {
		t.Errorf("Error setting Solr client: %s", err)
		t.FailNow()
	}

	// set up expected call orders
	sc.ExpectedCallOrder.Delete = 1   // delete is always called first
	sc.ExpectedCallOrder.Rollback = 2 // rollback = delete + rollback = 2
	sc.ExpectedDeleteArgument = eadid

	// setup error events
	var errorEvents []testutils.ErrorEvent

	errorEvents = append(errorEvents, testutils.ErrorEvent{CallerName: "Delete", ErrorMessage: "error during Delete", CallCount: 1})
	sc.ErrorEvents = errorEvents

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	sc.ActualError = IndexEADFile(eadPath)

	// check that all expectations were met
	err = sc.CheckAssertions()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}
}

func TestIndexEADFile_RollbackOnBadCollectionIndex(t *testing.T) {

	repositoryCode := "fales"
	eadid := "mss_460"
	testEAD := filepath.Join(repositoryCode, eadid)

	var eadPath = eadtestutils.EadFixturePath(testEAD)

	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForIndexing(testEAD)
	if err != nil {
		t.Errorf("Error setting Solr client: %s", err)
		t.FailNow()
	}

	// set up expected call orders
	sc.ExpectedCallOrder.Delete = 1   // delete is always called first
	sc.ExpectedCallOrder.Rollback = 3 // rollback = delete + Add(CollectionDoc) + rollback = 3
	sc.ExpectedDeleteArgument = eadid

	// setup error events
	var errorEvents []testutils.ErrorEvent

	errorEvents = append(errorEvents, testutils.ErrorEvent{CallerName: "Add", ErrorMessage: "error during Add", CallCount: 2})
	sc.ErrorEvents = errorEvents

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	sc.ActualError = IndexEADFile(eadPath)

	// check that all expectations were met
	err = sc.CheckAssertions()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}
}

func TestIndexEADFile_RollbackOnBadComponentIndex(t *testing.T) {

	// specify which calls to Add() will return an error
	errorCallCounts := []int{11, 226, 333, 444, 555, 666, 777, 888, 999, 1000, 1208}

	repositoryCode := "nyhs"
	eadid := "ms347_foundling_hospital"
	testEAD := filepath.Join(repositoryCode, eadid)

	var eadPath = eadtestutils.EadFixturePath(testEAD)

	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForIndexing(testEAD)
	if err != nil {
		t.Errorf("Error initializing Solr Client Mock: %s", err)
		t.FailNow()
	}

	// set expectations
	// (note: Commit() is not called because there were errors during component-level indexing)
	sc.ExpectedCallOrder.Delete = 1                             // delete is always called first
	sc.ExpectedCallOrder.Rollback = sc.NumberOfFilesToIndex + 2 // rollback = delete + number of files + rollback = number of files + 2
	sc.ExpectedDeleteArgument = eadid

	// setup error events
	var errorEvents []testutils.ErrorEvent

	for _, errorCallCount := range errorCallCounts {
		errorEvents = append(errorEvents, testutils.ErrorEvent{CallerName: "Add", ErrorMessage: fmt.Sprintf("error during Add: %d", errorCallCount), CallCount: errorCallCount})
	}

	sc.ErrorEvents = testutils.SortErrorEventsByCallCount(errorEvents)

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	sc.ActualError = IndexEADFile(eadPath)

	err = sc.CheckAssertions()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}
}

func TestIndexEADFile_RollbackOnBadCommit(t *testing.T) {

	repositoryCode := "fales"
	eadid := "mss_460"

	testEAD := filepath.Join(repositoryCode, eadid)
	var eadPath = eadtestutils.EadFixturePath(testEAD)

	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForIndexing(testEAD)
	if err != nil {
		t.Errorf("Error setting Solr client: %s", err)
		t.FailNow()
	}

	// set expectations
	sc.ExpectedCallOrder.Delete = 1                             // delete is always called first
	sc.ExpectedCallOrder.Commit = sc.NumberOfFilesToIndex + 2   // commit   = delete + number of files + commit = number of files + 2
	sc.ExpectedCallOrder.Rollback = sc.NumberOfFilesToIndex + 3 // rollback = delete + number of files + commit + rollback = number of files + 3
	sc.ExpectedDeleteArgument = eadid

	// setup error events
	var errorEvents []testutils.ErrorEvent
	errorEvents = append(errorEvents, testutils.ErrorEvent{CallerName: "Commit", ErrorMessage: "error during Commit", CallCount: sc.ExpectedCallOrder.Commit})
	sc.ErrorEvents = testutils.SortErrorEventsByCallCount(errorEvents)

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	sc.ActualError = IndexEADFile(eadPath)

	// check that all expectations were met
	err = sc.CheckAssertions()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}
}

func TestDeleteEADFileDataFromIndex_Success(t *testing.T) {

	eadid := "mss_460"

	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForDelete()
	if err != nil {
		t.Errorf("Error initializing the Solr client for delete testing: %s", err)
		t.FailNow()
	}

	// set expectations
	sc.ExpectedCallOrder.Delete = 1 // delete is always called first
	sc.ExpectedDeleteArgument = eadid

	// Set the Solr client
	SetSolrClient(sc)

	// Delete the data for the EADID
	sc.ActualError = DeleteEADFileDataFromIndex(eadid)

	// check that all expectations were met
	err = sc.CheckAssertions()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}
}

func TestDeleteEADFileDataFromIndex_RollbackOnBadDelete(t *testing.T) {

	eadid := "mss_460"

	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForDelete()
	if err != nil {
		t.Errorf("Error initializing the Solr client for delete testing: %s", err)
		t.FailNow()
	}

	// set expectations
	// (note: Commit() is not called because there were errors during component-level indexing)
	sc.ExpectedCallOrder.Delete = 1   // delete is always called first
	sc.ExpectedCallOrder.Rollback = 2 // rollback = delete + rollback = 2
	sc.ExpectedDeleteArgument = eadid

	// setup error events
	var errorEvents []testutils.ErrorEvent

	errorEvents = append(errorEvents, testutils.ErrorEvent{CallerName: "Delete", ErrorMessage: "error during Delete", CallCount: 1})
	sc.ErrorEvents = errorEvents

	// Set the Solr client
	SetSolrClient(sc)

	// Delete the data for the EADID
	sc.ActualError = DeleteEADFileDataFromIndex(eadid)

	// check that all expectations were met
	err = sc.CheckAssertions()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}
}
