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

/*
	# Commit history from test fixture (NOTE: double check commit hashes)
	00fd44a8e69285cf3789be3e7bc0e4e88d5f6dd8 2025-03-26 13:18:07 -0400 | Updating nyuad/ad_mc_019.xml, Deleting file tamwag/tam_143.xml EADID='tam_143', Updating edip/mos_2024.xml, Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Updating akkasah/ad_mc_030.xml (HEAD -> main) [jgpawletko]
	5fec61740cb7e4f05bbfa77548b42be2003e278b 2025-03-26 13:18:07 -0400 | Updating tamwag/tam_143.xml, Updating cbh/arc_212_plymouth_beecher.xml [jgpawletko]
	d8144b3136ef4a9abf0613a1302606644f90bd6c 2025-03-26 13:18:07 -0400 | Deleting file fales/mss_420.xml EADID='mss_420', Updating fales/mss_420.xml [jgpawletko]
	0afcf14e99bbd6e158f486090877fbd50370494c 2025-03-26 13:18:07 -0400 | Updating fales/mss_420.xml [jgpawletko]
	aee0af16b6d92444326eea4847893844f3ca59ae 2025-03-26 13:18:07 -0400 | Deleting file fales/mss_460.xml EADID='mss_460' [jgpawletko]
	e24960b1dc934c628d4475cb4537f7e21f54032c 2025-03-26 13:18:07 -0400 | Updating fales/mss_460.xml [jgpawletko]
	0fcdd54abaeb3b2f15b50f8eb5ef903ba2231896 2025-03-26 13:18:07 -0400 | Deleting file akkasah/ad_mc_030.xml EADID='ad_mc_030', Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Deleting file edip/mos_2024.xml EADID='mos_2024', Deleting file fales/mss_420.xml EADID='mss_420', Deleting file fales/mss_460.xml EADID='mss_460', Deleting file nyhs/ms256_harmon_hendricks_goldstone.xml EADID='ms256_harmon_hendricks_goldstone', Deleting file nyhs/ms347_foundling_hospital.xml EADID='ms347_foundling_hospital', Deleting file nyuad/ad_mc_019.xml EADID='ad_mc_019', Deleting file tamwag/tam_143.xml EADID='tam_143' [jgpawletko]
	7fdb03f4ab09f0eddf9b3c0e77ba50f5d036b2e9 2025-03-26 13:18:07 -0400 | Updating akkasah/ad_mc_030.xml, Updating cbh/arc_212_plymouth_beecher.xml, Updating edip/mos_2024.xml, Updating fales/mss_420.xml, Updating fales/mss_460.xml, Updating nyhs/ms256_harmon_hendricks_goldstone.xml, Updating nyhs/ms347_foundling_hospital.xml, Updating nyuad/ad_mc_019.xml, Updating tamwag/tam_143.xml [jgpawletko]
*/

// hashes from the git-repo fixture (in order of commits)
var addAllHash = "7fdb03f4ab09f0eddf9b3c0e77ba50f5d036b2e9"
var deleteAllHash = "0fcdd54abaeb3b2f15b50f8eb5ef903ba2231896"
var addOneHash = "e24960b1dc934c628d4475cb4537f7e21f54032c"
var deleteOneHash = "aee0af16b6d92444326eea4847893844f3ca59ae"
var deleteModifyAddHash = "d8144b3136ef4a9abf0613a1302606644f90bd6c"
var addTwoHash = "5fec61740cb7e4f05bbfa77548b42be2003e278b"
var addThreeDeleteTwoHash = "00fd44a8e69285cf3789be3e7bc0e4e88d5f6dd8"

// test git repo paths
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

func TestIndexEADFile_Success(t *testing.T) {

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

func TestIndexGitCommit_AddAll(t *testing.T) {
	/*
	   # Commit history replicated in repo (NOTE: commit hashes WILL differ)
	   # 5546ffda27581c4933aeb4102f6a0107c3e522ff 2025-03-24 19:53:30 -0400 | Updating akkasah/ad_mc_030.xml, Updating cbh/arc_212_plymouth_beecher.xml, Updating edip/mos_2024.xml, Updating fales/mss_420.xml, Updating fales/mss_460.xml, Updating nyhs/ms256_harmon_hendricks_goldstone.xml, Updating nyhs/ms347_foundling_hospital.xml, Updating nyuad/ad_mc_019.xml, Updating tamwag/tam_143.xml [jgpawletko]
	*/
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

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
	err := IndexGitCommit(gitRepoTestGitRepoPathAbsolute, addAllHash)
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
	err = IndexGitCommit(gitRepoTestGitRepoPathAbsolute, addOneHash)
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
func TestIndexGitCommit_AddThreeDeleteTwo(t *testing.T) {
	/*
	   # Commit history replicated in repo (NOTE: commit hashes WILL differ)
	   # b2456cf44f6ff4cefeb621ef2f4cde76218327d5 2025-03-25 20:36:43 -0400 | Updating akkasah/ad_mc_030.xml, Deleting file nyuad/ad_mc_019.xml EADID='ad_mc_019', Updating cbh/arc_212_plymouth_beecher.xml, Deleting file tamwag/tam_143.xml EADID='tam_143', Updating edip/mos_2024.xml (HEAD -> main) [jgpawletko]
	*/
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	sc := testutils.GetSolrClientMock()
	sc.Reset()

	// NOTE: the commits will always be returned in alphabetical order by relative path
	ops := [][]string{
		{"akkasah", "ad_mc_030", "Add"},
		{"cbh", "arc_212_plymouth_beecher", "Delete"},
		{"edip", "mos_2024", "Add"},
		{"nyuad", "ad_mc_019", "Add"},
		{"tamwag", "tam_143", "Delete"},
	}

	for _, op := range ops {
		repositoryCode := op[0]
		eadid := op[1]
		testEAD := filepath.Join(repositoryCode, eadid)
		if op[2] == "Add" {
			err := sc.UpdateMockForIndexEADFile(testEAD, eadid)
			if err != nil {
				t.Errorf("Error updating the SolrClientMock: %s", err)
				t.FailNow()
			}
		}
		if op[2] == "Delete" {
			err := sc.UpdateMockForDeleteEADFileDataFromIndex(eadid)
			if err != nil {
				t.Errorf("Error updating the SolrClientMock: %s", err)
				t.FailNow()
			}
		}
	}

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	err := IndexGitCommit(gitRepoTestGitRepoPathAbsolute, addThreeDeleteTwoHash)
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

func TestIndexGitCommit_AddTwo(t *testing.T) {
	/*
	   # Commit history replicated in repo (NOTE: commit hashes WILL differ)
	   # 5fec61740cb7e4f05bbfa77548b42be2003e278b 2025-03-26 13:18:07 -0400 | Updating tamwag/tam_143.xml, Updating cbh/arc_212_plymouth_beecher.xml [jgpawletko]
	*/
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	sc := testutils.GetSolrClientMock()
	sc.Reset()

	testEADs := [][]string{
		{"cbh", "arc_212_plymouth_beecher"},
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
	err := IndexGitCommit(gitRepoTestGitRepoPathAbsolute, addTwoHash)
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

func TestIndexGitCommit_DeleteAll(t *testing.T) {
	/*
	   # Commit history replicated in repo (NOTE: commit hashes WILL differ)
	   # e4fe6008decb5f26382fae903de40a4f3470d509 2025-03-24 19:53:30 -0400 | Deleting file akkasah/ad_mc_030.xml EADID='ad_mc_030', Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Deleting file edip/mos_2024.xml EADID='mos_2024', Deleting file fales/mss_420.xml EADID='mss_420', Deleting file fales/mss_460.xml EADID='mss_460', Deleting file nyhs/ms256_harmon_hendricks_goldstone.xml EADID='ms256_harmon_hendricks_goldstone', Deleting file nyhs/ms347_foundling_hospital.xml EADID='ms347_foundling_hospital', Deleting file nyuad/ad_mc_019.xml EADID='ad_mc_019', Deleting file tamwag/tam_143.xml EADID='tam_143' [jgpawletko]
	*/
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

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
		eadid := testEAD[1]
		err := sc.UpdateMockForDeleteEADFileDataFromIndex(eadid)
		if err != nil {
			t.Errorf("Error updating the SolrClientMock: %s", err)
			t.FailNow()
		}
	}

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	err := IndexGitCommit(gitRepoTestGitRepoPathAbsolute, deleteAllHash)
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

func TestIndexGitCommit_DeleteModifyAdd(t *testing.T) {
	/*
	   # Commit history replicated in repo (NOTE: commit hashes WILL differ)
	   c834255a61231deb3d090ef3a8578b43cccaa5fb 2025-03-26 12:34:08 -0400 | Deleting file fales/mss_420.xml EADID='mss_420', Updating fales/mss_420.xml [jgpawletko]
	*/
	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	sc := testutils.GetSolrClientMock()
	sc.Reset()

	// even though we are deleting the file during the commit-staging process
	// (in the testsupport/gen-repo.bash script), the script adds a modified
	// version of the file back to the staging area before the commit.
	// Therefore, the state of the Git staging area at the time of the commit only
	// contains the file in the "modified" state and the commit boils down to single
	// "add" operation
	ops := [][]string{
		{"fales", "mss_420", "Add"},
	}

	for _, op := range ops {
		repositoryCode := op[0]
		eadid := op[1]
		testEAD := filepath.Join(repositoryCode, eadid)
		if op[2] == "Add" {
			err := sc.UpdateMockForIndexEADFile(testEAD, eadid)
			if err != nil {
				t.Errorf("Error updating the SolrClientMock: %s", err)
				t.FailNow()
			}
		}
		if op[2] == "Delete" {
			err := sc.UpdateMockForDeleteEADFileDataFromIndex(eadid)
			if err != nil {
				t.Errorf("Error updating the SolrClientMock: %s", err)
				t.FailNow()
			}
		}
	}

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	err := IndexGitCommit(gitRepoTestGitRepoPathAbsolute, deleteModifyAddHash)
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
	err = IndexGitCommit(gitRepoTestGitRepoPathAbsolute, deleteOneHash)
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

func TestIndexGitCommit_FailFast(t *testing.T) {
	/*
	   # Commit history replicated in repo (NOTE: commit hashes WILL differ)
	   # b2456cf44f6ff4cefeb621ef2f4cde76218327d5 2025-03-25 20:36:43 -0400 | Updating akkasah/ad_mc_030.xml, Deleting file nyuad/ad_mc_019.xml EADID='ad_mc_019', Updating cbh/arc_212_plymouth_beecher.xml, Deleting file tamwag/tam_143.xml EADID='tam_143', Updating edip/mos_2024.xml (HEAD -> main) [jgpawletko]
	*/

	errorEventCallCount := 300    // this is in the middle of the akkasah/ad_mc_030.xml file component indexing
	errorRollbackCallCount := 630 // delete + collection + components + rollback

	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)

	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	sc := testutils.GetSolrClientMock()
	sc.Reset()

	// NOTE: the commits will always be returned in alphabetical order by relative path
	ops := [][]string{
		{"akkasah", "ad_mc_030", "Add"},
		{"cbh", "arc_212_plymouth_beecher", "Delete"},
		{"edip", "mos_2024", "Add"},
	}

	for _, op := range ops {
		repositoryCode := op[0]
		eadid := op[1]
		testEAD := filepath.Join(repositoryCode, eadid)
		if op[2] == "Add" {
			err := sc.UpdateMockForIndexEADFile(testEAD, eadid)
			if err != nil {
				t.Errorf("Error updating the SolrClientMock: %s", err)
				t.FailNow()
			}
		}
		if op[2] == "Delete" {
			err := sc.UpdateMockForDeleteEADFileDataFromIndex(eadid)
			if err != nil {
				t.Errorf("Error updating the SolrClientMock: %s", err)
				t.FailNow()
			}
		}
	}

	solrClientErrorEvents := []testutils.ErrorEvent{
		{FuncName: "Add", ErrorMessage: "error during Add", CallCount: errorEventCallCount},
	}
	sc.ErrorEvents = solrClientErrorEvents

	// Set the Solr client
	SetSolrClient(sc)

	// Index the EAD file
	err := IndexGitCommit(gitRepoTestGitRepoPathAbsolute, addThreeDeleteTwoHash)
	if err == nil {
		t.Errorf("Expected error from IndexGitCommit() but no error was returned.")
	}

	if err.Error() != "error during Add" {
		t.Errorf("Expected error message 'error during Add' but got '%s'", err.Error())
	}

	// check that rollback was called
	if sc.ActualEvents[errorRollbackCallCount].FuncName != "Rollback" {
		t.Errorf("Expected Rollback() to be called at call count %d but it was no", errorRollbackCallCount)
	}

	// indexing should NOT have completed
	if sc.IsComplete() {
		t.Errorf("All files were added to the Solr index when indexing should have halted.")
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
