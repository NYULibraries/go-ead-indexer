package index

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	eadtestutils "github.com/nyulibraries/go-ead-indexer/pkg/ead/testutils"
	"github.com/nyulibraries/go-ead-indexer/pkg/index/testutils"
)

// func TestDeleteEADFileDataFromIndex_RollbackOnBadDelete(t *testing.T) {

// 	eadid := "mss_460"

// 	sc := testutils.GetSolrClientMock()
// 	err := sc.InitMockForDelete()
// 	if err != nil {
// 		t.Errorf("Error initializing the Solr client for delete testing: %s", err)
// 		t.FailNow()
// 	}

// 	// set expectations
// 	// (note: Commit() is not called because there were errors during component-level indexing)
// 	sc.ExpectedCallOrder.Delete = 1   // delete is always called first
// 	sc.ExpectedCallOrder.Rollback = 2 // rollback = delete + rollback = 2
// 	sc.ExpectedDeleteArgument = eadid

// 	// setup error events
// 	var errorEvents []testutils.ErrorEvent

// 	errorEvents = append(errorEvents, testutils.ErrorEvent{FuncName: "Delete", ErrorMessage: "error during Delete", CallCount: 1})
// 	sc.ErrorEvents = errorEvents

// 	// Set the Solr client
// 	SetSolrClient(sc)

// 	// Delete the data for the EADID
// 	sc.ActualError = DeleteEADFileDataFromIndex(eadid)

//		// check that all expectations were met
//		err = sc.CheckAssertions()
//		if err != nil {
//			t.Errorf("Assertions failed: %s", err)
//		}
//	}
func TestDeleteEADFileDataFromIndex_RollbackOnBadDelete(t *testing.T) {

	sut := "DeleteEADFileDataFromIndex"
	eadid := "mss_460"

	// set up the Solr client mock
	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForDelete(sut)
	if err != nil {
		t.Errorf("Error initializing the Solr client for delete testing: %s", err)
		t.FailNow()
	}

	// set up expected events
	expectedEvents := []testutils.Event{
		{FuncName: "Delete", Args: []string{eadid}, CallCount: 1, Err: fmt.Errorf("error during Delete")},
		{FuncName: "Rollback", CallCount: 2},
	}
	sc.ExpectedEvents = expectedEvents

	// setup error events
	var errorEvents []testutils.ErrorEvent
	errorEvents = append(errorEvents, testutils.ErrorEvent{FuncName: "Delete", ErrorMessage: "error during Delete", CallCount: 1})
	sc.ErrorEvents = errorEvents

	// Set the Solr client
	SetSolrClient(sc)

	// Delete the data for the EADID
	DeleteEADFileDataFromIndex(eadid)

	// check that all expectations were met
	err = sc.CheckAssertionsViaEvents()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}
}
func TestDeleteEADFileDataFromIndex_BadEADID(t *testing.T) {

	eadid := "waffles!@#%"
	sut := "DeleteEADFileDataFromIndex"
	expectedErrStringFragment := fmt.Sprintf("invalid EADID: %s", eadid)
	expectedCallCount := 0

	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForDelete(sut)
	if err != nil {
		t.Errorf("Error initializing the Solr client for delete testing: %s", err)
		t.FailNow()
	}

	// Set the Solr client
	SetSolrClient(sc)

	// Delete the data for the EADID
	err = DeleteEADFileDataFromIndex(eadid)

	testutils.AssertError(t, sut, err)
	testutils.AssertErrorMessageContainsString(t, sut, err, expectedErrStringFragment)
	testutils.AssertCallCount(t, expectedCallCount, sc.CallCount)
}

func TestDeleteEADFileDataFromIndex_SolrClientNotSet(t *testing.T) {

	sut := "DeleteEADFileDataFromIndex"
	expectedErrStringFragment := "you must call `SetSolrClient()` before calling any indexing functions"
	eadid := "mss_460"

	SetSolrClient(nil)

	// Delete the data for the EADID
	err := DeleteEADFileDataFromIndex(eadid)

	testutils.AssertError(t, sut, err)
	testutils.AssertErrorMessageContainsString(t, sut, err, expectedErrStringFragment)
}

func TestDeleteEADFileDataFromIndex_SolrClientMissingOriginURL(t *testing.T) {

	sut := "DeleteEADFileDataFromIndex"
	expectedErrStringFragment := "the SolrClient URL origin is not set"
	eadid := "mss_460"

	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForDelete(sut)
	if err != nil {
		t.Errorf("Error initializing the Solr client for delete testing: %s", err)
		t.FailNow()
	}

	sc.SetSolrURLOrigin("")
	SetSolrClient(sc)

	// trigger the error
	err = DeleteEADFileDataFromIndex(eadid)

	testutils.AssertError(t, sut, err)
	testutils.AssertErrorMessageContainsString(t, sut, err, expectedErrStringFragment)
}

func TestDeleteEADFileDataFromIndex_ErrorOnRollback(t *testing.T) {

	sut := "DeleteEADFileDataFromIndex"
	eadid := "mss_460"

	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForDelete(sut)
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

	errorEvents = append(errorEvents, testutils.ErrorEvent{FuncName: "Delete", ErrorMessage: "error during Delete", CallCount: 1})
	errorEvents = append(errorEvents, testutils.ErrorEvent{FuncName: "Rollback", ErrorMessage: "error during Rollback", CallCount: 2})
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

func TestDeleteEADFileDataFromIndex_ErrorOnRollbackEVENTS(t *testing.T) {

	sut := "DeleteEADFileDataFromIndex"
	eadid := "mss_460"

	// set up the Solr client mock
	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForDelete(sut)
	if err != nil {
		t.Errorf("Error initializing the Solr client for delete testing: %s", err)
		t.FailNow()
	}

	// set up expected events
	expectedEvents := []testutils.Event{
		{FuncName: "Delete", Args: []string{eadid}, CallCount: 1, Err: fmt.Errorf("error during Delete")},
		{FuncName: "Rollback", CallCount: 2, Err: fmt.Errorf("error during Rollback")},
	}
	sc.ExpectedEvents = expectedEvents

	// setup error events
	var errorEvents []testutils.ErrorEvent

	errorEvents = append(errorEvents, testutils.ErrorEvent{FuncName: "Delete", ErrorMessage: "error during Delete", CallCount: 1})
	errorEvents = append(errorEvents, testutils.ErrorEvent{FuncName: "Rollback", ErrorMessage: "error during Rollback", CallCount: 2})
	sc.ErrorEvents = errorEvents

	// Set the Solr client
	SetSolrClient(sc)

	// Delete the data for the EADID
	DeleteEADFileDataFromIndex(eadid)

	// check that all expectations were met
	err = sc.CheckAssertionsViaEvents()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}
}

func TestIndexEADFile_EADFileDoesNotExist(t *testing.T) {

	sut := "IndexEADFile"
	expectedErrStringFragment := "no such file or directory"
	expectedCallCount := 0

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

	testutils.AssertError(t, sut, err)
	testutils.AssertErrorMessageContainsString(t, sut, err, expectedErrStringFragment)
	testutils.AssertCallCount(t, expectedCallCount, sc.CallCount)
}

func TestIndexEADFile_EADFilePathIsAbsolute(t *testing.T) {

	sut := "IndexEADFile"
	expectedErrStringFragment := "EAD file path must be absolute:"
	expectedCallCount := 0

	sc := testutils.GetSolrClientMock()
	SetSolrClient(sc)

	eadPath := filepath.Join(".", "this-file-does-not-exist.xml")

	err := IndexEADFile(eadPath)
	testutils.AssertError(t, sut, err)
	testutils.AssertErrorMessageContainsString(t, sut, err, expectedErrStringFragment)
	testutils.AssertCallCount(t, expectedCallCount, sc.CallCount)
}

func TestIndexEADFile_SolrClientNotSet(t *testing.T) {

	sut := "IndexEADFile"
	expectedErrStringFragment := "you must call `SetSolrClient()` before calling any indexing functions"

	repositoryCode := "fales"
	eadid := "mss_460"
	testEAD := filepath.Join(repositoryCode, eadid)

	var eadPath = eadtestutils.EadFixturePath(testEAD)

	SetSolrClient(nil)
	err := IndexEADFile(eadPath)

	testutils.AssertError(t, sut, err)
	testutils.AssertErrorMessageContainsString(t, sut, err, expectedErrStringFragment)
}
func TestIndexEADFile_ErrorExtractingRepositoryCode(t *testing.T) {

	sut := "IndexEADFile"
	expectedErrStringFragment := "EAD file path must have at least two non-empty components"
	expectedCallCount := 0

	sc := testutils.GetSolrClientMock()
	SetSolrClient(sc)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Errorf("ERROR: `runtime.Caller(0)` failed")
		t.FailNow()
	}

	dir := filepath.Dir(filename)
	eadPath := filepath.Join(dir, "testdata", "fixtures", "edip", "test-file-for-bad-repo-extraction.empty")

	err := IndexEADFile(eadPath)
	testutils.AssertError(t, sut, err)
	testutils.AssertErrorMessageContainsString(t, sut, err, expectedErrStringFragment)
	testutils.AssertCallCount(t, expectedCallCount, sc.CallCount)
}

func TestIndexEADFile_ErrorDuringEADParsing(t *testing.T) {

	sut := "IndexEADFile"
	expectedErrStringFragment := `"THIS!IS#AND$INVALID)EADID" is not a valid EAD ID`
	expectedCallCount := 0

	sc := testutils.GetSolrClientMock()
	SetSolrClient(sc)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Errorf("ERROR: `runtime.Caller(0)` failed")
		t.FailNow()
	}

	dir := filepath.Dir(filename)
	eadPath := filepath.Join(dir, "testdata", "fixtures", "edip", "this-is-an-invalid-eadid.xml")

	err := IndexEADFile(eadPath)
	testutils.AssertError(t, sut, err)
	testutils.AssertErrorMessageContainsString(t, sut, err, expectedErrStringFragment)
	testutils.AssertCallCount(t, expectedCallCount, sc.CallCount)
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
		t.Errorf("Error initializing the SolrClientMock: %s", err)
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
		t.Errorf("Error initializing Solr Client Mock: %s", err)
		t.FailNow()
	}

	// set up expected call orders
	sc.ExpectedCallOrder.Delete = 1   // delete is always called first
	sc.ExpectedCallOrder.Rollback = 2 // rollback = delete + rollback = 2
	sc.ExpectedDeleteArgument = eadid

	// setup error events
	var errorEvents []testutils.ErrorEvent

	errorEvents = append(errorEvents, testutils.ErrorEvent{FuncName: "Delete", ErrorMessage: "error during Delete", CallCount: 1})
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
		t.Errorf("Error initializing Solr Client Mock: %s", err)
		t.FailNow()
	}

	// set up expected call orders
	sc.ExpectedCallOrder.Delete = 1   // delete is always called first
	sc.ExpectedCallOrder.Rollback = 3 // rollback = delete + Add(CollectionDoc) + rollback = 3
	sc.ExpectedDeleteArgument = eadid

	// setup error events
	var errorEvents []testutils.ErrorEvent

	errorEvents = append(errorEvents, testutils.ErrorEvent{FuncName: "Add", ErrorMessage: "error during Add", CallCount: 2})
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
		errorEvents = append(errorEvents, testutils.ErrorEvent{FuncName: "Add", ErrorMessage: fmt.Sprintf("error during Add: %d", errorCallCount), CallCount: errorCallCount})
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
		t.Errorf("Error initializing Solr Client Mock: %s", err)
		t.FailNow()
	}

	// set expectations
	sc.ExpectedCallOrder.Delete = 1                             // delete is always called first
	sc.ExpectedCallOrder.Commit = sc.NumberOfFilesToIndex + 2   // commit   = delete + number of files + commit = number of files + 2
	sc.ExpectedCallOrder.Rollback = sc.NumberOfFilesToIndex + 3 // rollback = delete + number of files + commit + rollback = number of files + 3
	sc.ExpectedDeleteArgument = eadid

	// setup error events
	var errorEvents []testutils.ErrorEvent
	errorEvents = append(errorEvents, testutils.ErrorEvent{FuncName: "Commit", ErrorMessage: "error during Commit", CallCount: sc.ExpectedCallOrder.Commit})
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

	sut := "DeleteEADFileDataFromIndex"
	eadid := "mss_460"

	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForDelete(sut)
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

func TestIndexGitCommit_SolrClientNotSet(t *testing.T) {

	sut := "IndexGitCommit"
	expectedErrStringFragment := "you must call `SetSolrClient()` before calling any indexing functions"
	repoPath := "/foo/bar"
	commit := "a5ca6cca30fc08cfc13e4f1492dbfbbf3ec7cf63"

	SetSolrClient(nil)

	// trigger the error
	err := IndexGitCommit(repoPath, commit)

	testutils.AssertError(t, sut, err)
	testutils.AssertErrorMessageContainsString(t, sut, err, expectedErrStringFragment)
}

func TestIndexGitCommit_SolrClientMissingOriginURL(t *testing.T) {

	sut := "IndexGitCommit"
	expectedErrStringFragment := "the SolrClient URL origin is not set"
	repoPath := "/foo/bar"
	commit := "a5ca6cca30fc08cfc13e4f1492dbfbbf3ec7cf63"

	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForDelete(sut)
	if err != nil {
		t.Errorf("Error initializing the Solr client for delete testing: %s", err)
		t.FailNow()
	}

	sc.SetSolrURLOrigin("")
	SetSolrClient(sc)

	// trigger the error
	err = IndexGitCommit(repoPath, commit)

	testutils.AssertError(t, sut, err)
	testutils.AssertErrorMessageContainsString(t, sut, err, expectedErrStringFragment)
}

// func TestIndexGitCommit_ErrorOnRollback(t *testing.T) {

// 	repoPath := "/foo/bar"
// 	commit := "a5ca6cca30fc08cfc13e4f1492dbfbbf3ec7cf63"

// 	sc := testutils.GetSolrClientMock()
// 	err := sc.InitMockForDelete()
// 	if err != nil {
// 		t.Errorf("Error initializing the Solr client for delete testing: %s", err)
// 		t.FailNow()
// 	}

// 	// set expectations
// 	// (note: Commit() is not called because there were errors during component-level indexing)
// 	sc.ExpectedCallOrder.Delete = 1   // delete is always called first
// 	sc.ExpectedCallOrder.Rollback = 2 // rollback = delete + rollback = 2
// 	sc.ExpectedDeleteArgument = eadid

// 	// setup error events
// 	var errorEvents []testutils.ErrorEvent

// 	errorEvents = append(errorEvents, testutils.ErrorEvent{FuncName: "Delete", ErrorMessage: "error during Delete", CallCount: 1})
// 	errorEvents = append(errorEvents, testutils.ErrorEvent{FuncName: "Rollback", ErrorMessage: "error during Rollback", CallCount: 2})
// 	sc.ErrorEvents = errorEvents

// 	// Set the Solr client
// 	SetSolrClient(sc)

// 	// Delete the data for the EADID
// 	sc.ActualError = DeleteEADFileDataFromIndex(eadid)

// 	// check that all expectations were met
// 	err = sc.CheckAssertions()
// 	if err != nil {
// 		t.Errorf("Assertions failed: %s", err)
// 	}
// }
