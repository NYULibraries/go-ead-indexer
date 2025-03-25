package index

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	eadtestutils "github.com/nyulibraries/go-ead-indexer/pkg/ead/testutils"
	"github.com/nyulibraries/go-ead-indexer/pkg/index/testutils"
)

var thisPath string
var gitSourceRepoPathAbsolute string
var gitRepoTestGitRepoPathAbsolute string
var gitRepoTestGitRepoPathRelative string
var gitRepoTestGitRepoDotGitDirectory string
var gitRepoTestGitRepoHiddenGitDirectory string

// this code is based on that in the debug package, written by David Arjanik
// We need to get the absolute path to this package in order to enable the function
// for golden file and fixture file retrieval to be called from other packages
// which would not be able to resolve the hardcoded relative paths used here.
func init() {
	// The `filename` string is the absolute path to this source file, which should
	// be located at the root of the package directory.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("ERROR: `runtime.Caller(0)` failed")
	}

	// Get the path to the parent directory of this file.  Again, this is assuming
	// that this `init()` function is defined in a package top level file -- or
	// more precisely, that this file is in the same directory at the `testdata/`
	// directory that is referenced in the relative paths used in the functions
	// defined in this file.
	thisPath = filepath.Dir(filename)

	// Get testdata directory paths
	gitSourceRepoPathAbsolute = filepath.Join(thisPath, "testdata", "fixtures", "git-repo")

	// This could be done as a const at top level, but assigning it here to keep
	// all this path stuff in one place.
	gitRepoTestGitRepoPathAbsolute = filepath.Join(thisPath, "testdata", "fixtures", "test-git-repo")
	gitRepoTestGitRepoPathRelative = filepath.Join(".", "testdata", "fixtures", "test-git-repo")
	gitRepoTestGitRepoDotGitDirectory = filepath.Join(gitRepoTestGitRepoPathAbsolute, "dot-git")
	gitRepoTestGitRepoHiddenGitDirectory = filepath.Join(gitRepoTestGitRepoPathAbsolute, ".git")
}

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
	solrClientExpectedEvents := []testutils.Event{
		{FuncName: "Delete", Args: []string{eadid}, CallCount: 1, Err: fmt.Errorf("error during Delete")},
		{FuncName: "Rollback", CallCount: 2},
	}
	sc.ExpectedEvents = solrClientExpectedEvents

	// setup error events
	solrClientErrorEvents := []testutils.ErrorEvent{
		{FuncName: "Delete", ErrorMessage: "error during Delete", CallCount: 1},
	}
	sc.ErrorEvents = solrClientErrorEvents

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

	// set up the Solr client mock
	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForDelete(sut)
	if err != nil {
		t.Errorf("Error initializing the Solr client for delete testing: %s", err)
		t.FailNow()
	}

	// set up expected events
	solrClientExpectedEvents := []testutils.Event{
		{FuncName: "Delete", Args: []string{eadid}, CallCount: 1, Err: fmt.Errorf("error during Delete")},
		{FuncName: "Rollback", CallCount: 2, Err: fmt.Errorf("error during Rollback")},
	}
	sc.ExpectedEvents = solrClientExpectedEvents

	// setup error events
	solrClientErrorEvents := []testutils.ErrorEvent{
		{FuncName: "Delete", ErrorMessage: "error during Delete", CallCount: 1},
		{FuncName: "Rollback", ErrorMessage: "error during Rollback", CallCount: 2},
	}
	sc.ErrorEvents = solrClientErrorEvents

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

func TestIndexEADFile_RollbackOnBadDelete(t *testing.T) {

	repositoryCode := "fales"
	eadid := "mss_460"
	testEAD := filepath.Join(repositoryCode, eadid)
	var eadPath = eadtestutils.EadFixturePath(testEAD)

	// set up the Solr client mock
	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForIndexing(testEAD)
	if err != nil {
		t.Errorf("Error initializing Solr Client Mock: %s", err)
		t.FailNow()
	}

	// setup expectations
	solrClientExpectedEvents := []testutils.Event{
		{FuncName: "Delete", Args: []string{eadid}, CallCount: 1, Err: fmt.Errorf("error during Delete")},
		{FuncName: "Rollback", CallCount: 2},
	}
	sc.ExpectedEvents = solrClientExpectedEvents

	// setup error events
	solrClientErrorEvents := []testutils.ErrorEvent{
		{FuncName: "Delete", ErrorMessage: "error during Delete", CallCount: 1},
	}
	sc.ErrorEvents = solrClientErrorEvents

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	IndexEADFile(eadPath)

	// check that all expectations were met
	err = sc.CheckAssertionsViaEvents()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}
}

func TestIndexEADFile_RollbackOnBadCollectionIndex(t *testing.T) {

	repositoryCode := "fales"
	eadid := "mss_460"
	testEAD := filepath.Join(repositoryCode, eadid)
	var eadPath = eadtestutils.EadFixturePath(testEAD)

	// set up the Solr client mock
	sc := testutils.GetSolrClientMock()
	err := sc.InitMockForIndexing(testEAD)
	if err != nil {
		t.Errorf("Error initializing Solr Client Mock: %s", err)
		t.FailNow()
	}

	// set up expected events
	solrClientExpectedEvents := []testutils.Event{
		{FuncName: "Delete", Args: []string{eadid}, CallCount: 1},
		{FuncName: "Add", CallCount: 2, Args: []string{"XMLPostBody"}, Err: fmt.Errorf("error during Add")},
		{FuncName: "Rollback", CallCount: 3},
	}
	sc.ExpectedEvents = solrClientExpectedEvents

	// setup error events
	solrClientErrorEvents := []testutils.ErrorEvent{
		{FuncName: "Add", ErrorMessage: "error during Add", CallCount: 2},
	}
	sc.ErrorEvents = solrClientErrorEvents

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	IndexEADFile(eadPath)

	// check that all expectations were met
	err = sc.CheckAssertionsViaEvents()
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
	runningCallCountInitialValue := 1
	runningCallCount := runningCallCountInitialValue
	solrClientExpectedEvents := []testutils.Event{
		{FuncName: "Delete", Args: []string{eadid}, CallCount: runningCallCount},
	}
	// generate Add events
	for range sc.NumberOfFilesToIndex {
		runningCallCount++
		solrClientExpectedEvents = append(solrClientExpectedEvents, testutils.Event{FuncName: "Add", CallCount: runningCallCount, Args: []string{"XMLPostBody"}})
	}
	// add Rollback event
	runningCallCount++
	solrClientExpectedEvents = append(solrClientExpectedEvents, testutils.Event{FuncName: "Rollback", CallCount: runningCallCount})
	sc.ExpectedEvents = solrClientExpectedEvents

	// setup error events
	var solrClientErrorEvents []testutils.ErrorEvent
	for _, errorCallCount := range errorCallCounts {
		emsg := fmt.Sprintf("error during Add: %d", errorCallCount)
		solrClientErrorEvents = append(solrClientErrorEvents, testutils.ErrorEvent{FuncName: "Add", ErrorMessage: emsg, CallCount: errorCallCount})
		sc.ExpectedEvents[errorCallCount-runningCallCountInitialValue].Err = fmt.Errorf("%s", emsg)
	}
	sc.ErrorEvents = testutils.SortErrorEventsByCallCount(solrClientErrorEvents)

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	IndexEADFile(eadPath)

	// check that all expectations were met
	err = sc.CheckAssertionsViaEvents()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}

	// even though there were errors during component-level indexing, everything should have been indexed
	if !sc.IsComplete() {
		t.Errorf("not all files were added to the Solr index. Remaining values: %v", sc.GoldenFileHashes)
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
	solrClientErrorEvents := []testutils.ErrorEvent{
		{FuncName: "Commit", ErrorMessage: "error during Commit", CallCount: sc.ExpectedCallOrder.Commit},
	}
	sc.ErrorEvents = solrClientErrorEvents

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

func TestIndexEADFile_Success_Events(t *testing.T) {

	repositoryCode := "fales"
	eadid := "mss_460"
	testEAD := filepath.Join(repositoryCode, eadid)
	var eadPath = eadtestutils.EadFixturePath(testEAD)

	sc := testutils.GetSolrClientMock()
	sc.Reset()
	err := sc.UpdateMockForIndexEADFile(testEAD, eadid)
	if err != nil {
		t.Errorf("Error updating the SolrClientMock: %s", err)
		t.FailNow()
	}

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	err = IndexEADFile(eadPath)
	if err != nil {
		t.Errorf("Error indexing EAD file: %s", err)
	}

	err = sc.CheckAssertionsViaEvents()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}

	if !sc.IsComplete() {
		t.Errorf("not all files were added to the Solr index. Remaining values: %v", sc.GoldenFileHashes)
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

func TestIndexGitCommit_AddAll(t *testing.T) {
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	/*
			set up git repo
			get absolute repo path
			set up expectations
			run transaction
			check expectations

			# Commit history replicated in repo (NOTE: commit hashes WILL differ)
		    # 5546ffda27581c4933aeb4102f6a0107c3e522ff 2025-03-24 19:53:30 -0400 | Updating akkasah/ad_mc_030.xml, Updating cbh/arc_212_plymouth_beecher.xml, Updating edip/mos_2024.xml, Updating fales/mss_420.xml, Updating fales/mss_460.xml, Updating nyhs/ms256_harmon_hendricks_goldstone.xml, Updating nyhs/ms347_foundling_hospital.xml, Updating nyuad/ad_mc_019.xml, Updating tamwag/tam_143.xml [jgpawletko]
	*/
	sc := testutils.GetSolrClientMock()
	sc.Reset()

	testEADs := [][]string{
		{"akkasah", "ad_mc_030"},
		{"cbh", "arc_212_plymouth_beecher"},
		{"edip", "mos_2024"},
		{"fales", "mss_420"},
		{"fales", "mss_460"},
		{"nyhs", "ms256_harmon_hendricks_goldstone"},
		{"nyhs", "ms347_foundling_hospital"},
		{"nyuad", "ad_mc_019"},
		{"tamwag", "tam_143"},
	}

	for _, testEAD := range testEADs {
		repositoryCode := testEAD[0]
		eadid := testEAD[1]
		testEAD := filepath.Join(repositoryCode, eadid)
		err := sc.UpdateMockForIndexEADFile(testEAD, eadid)
		if err != nil {
			t.Errorf("Error updating the SolrClientMock: %s", err)
			t.FailNow()
		}
	}

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	err := IndexGitCommit(gitRepoTestGitRepoPathAbsolute, "5546ffda27581c4933aeb4102f6a0107c3e522ff")
	if err != nil {
		t.Errorf("Error indexing EAD file: %s", err)
	}

	err = sc.CheckAssertionsViaEvents()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}

	if !sc.IsComplete() {
		t.Errorf("not all files were added to the Solr index. Remaining values: \n%v", sc.GoldenFileHashesToString())
	}
}

func TestIndexGitCommit_AddOne(t *testing.T) {
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	/*
		set up git repo
		get absolute repo path
		set up expectations
		run transaction
		check expectations

		# Commit history replicated in repo (NOTE: commit hashes WILL differ)
		# fdd7ce5e54b88894460b52dd0dd27055ffb3bbdd 2025-03-24 19:53:30 -0400 | Updating fales/mss_460.xml [jgpawletko]
	*/
	sc := testutils.GetSolrClientMock()
	sc.Reset()

	repositoryCode := "fales"
	eadid := "mss_460"
	testEAD := filepath.Join(repositoryCode, eadid)
	err := sc.UpdateMockForIndexEADFile(testEAD, eadid)
	if err != nil {
		t.Errorf("Error updating the SolrClientMock: %s", err)
		t.FailNow()
	}

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	err = IndexGitCommit(gitRepoTestGitRepoPathAbsolute, "fdd7ce5e54b88894460b52dd0dd27055ffb3bbdd")
	if err != nil {
		t.Errorf("Error indexing EAD file: %s", err)
	}

	err = sc.CheckAssertionsViaEvents()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}

	if !sc.IsComplete() {
		t.Errorf("not all files were added to the Solr index. Remaining values: \n%v", sc.GoldenFileHashesToString())
	}
}

func TestIndexGitCommit_DeleteOne(t *testing.T) {
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	/*
		set up git repo
		get absolute repo path
		set up expectations
		run transaction
		check expectations

		# Commit history replicated in repo (NOTE: commit hashes WILL differ)
		# fdd7ce5e54b88894460b52dd0dd27055ffb3bbdd 2025-03-24 19:53:30 -0400 | Updating fales/mss_460.xml [jgpawletko]
	*/
	sc := testutils.GetSolrClientMock()
	sc.Reset()

	eadid := "mss_460"
	err := sc.UpdateMockForDeleteEADFileDataFromIndex(eadid)
	if err != nil {
		t.Errorf("Error updating the SolrClientMock: %s", err)
		t.FailNow()
	}

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	err = IndexGitCommit(gitRepoTestGitRepoPathAbsolute, "2fee15ffc217a86d19756a6c816f59ca86e23893")
	if err != nil {
		t.Errorf("Error indexing EAD file: %s", err)
	}

	err = sc.CheckAssertionsViaEvents()
	if err != nil {
		t.Errorf("Assertions failed: %s", err)
	}

	if !sc.IsComplete() {
		t.Errorf("not all files were added to the Solr index. Remaining values: \n%v", sc.GoldenFileHashesToString())
	}
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
// 	var solrClientErrorEvents []testutils.ErrorEvent

// 	solrClientErrorEvents = append(solrClientErrorEvents, testutils.ErrorEvent{FuncName: "Delete", ErrorMessage: "error during Delete", CallCount: 1})
// 	solrClientErrorEvents = append(solrClientErrorEvents, testutils.ErrorEvent{FuncName: "Rollback", ErrorMessage: "error during Rollback", CallCount: 2})
// 	sc.ErrorEvents = solrClientErrorEvents

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
func createTestGitRepo(t *testing.T) {
	gitSourceRepoPathAbsoluteFS := os.DirFS(gitSourceRepoPathAbsolute)
	err := os.CopyFS(gitRepoTestGitRepoPathAbsolute, gitSourceRepoPathAbsoluteFS)
	if err != nil {
		t.Errorf(
			`Unexpected error returned by os.CopyFS(gitSourceRepoPathAbsoluteFS, gitRepoTestGitRepoPathAbsolute): %s`,
			err.Error())
		t.FailNow()
	}

	err = os.Rename(gitRepoTestGitRepoDotGitDirectory, gitRepoTestGitRepoHiddenGitDirectory)
	if err != nil {
		t.Errorf(
			`Unexpected error returned by os.Rename(gitRepoTestGitRepoDotGitDirectory, gitRepoTestGitRepoHiddenGitDirectory): %s`,
			err.Error())
		t.FailNow()
	}
}

func deleteTestGitRepo(t *testing.T) {
	err := os.RemoveAll(gitRepoTestGitRepoPathAbsolute)
	if err != nil {
		t.Errorf(
			`deleteEnabledHiddenGitDirectory() failed with error "%s", remove %s manually`,
			err.Error(), gitRepoTestGitRepoPathAbsolute)
		t.FailNow()
	}
}
