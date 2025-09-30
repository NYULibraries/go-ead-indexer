package testutils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"runtime"
	"strings"
	"testing"

	eadtestutils "github.com/nyulibraries/go-ead-indexer/pkg/ead/testutils"
)

type FunctionName string

const IGNORE_CALL_ORDER = -1
const Add = FunctionName("Add")
const Commit = FunctionName("Commit")
const Delete = FunctionName("Delete")
const Rollback = FunctionName("Rollback")

// ------------------------------------------------------------------------------
// type definitions
// ------------------------------------------------------------------------------
type CallOrder struct {
	Commit   int
	Delete   int
	Rollback int
}

type Event struct {
	Args      []string
	CallCount int
	Err       error
	FuncName  FunctionName
}

type ErrorEvent struct {
	CallCount    int
	FuncName     string
	ErrorMessage string
}

type SolrClientMock struct {
	GoldenFileHashes       map[string]string // Hashes of the golden files
	NumberOfFilesToIndex   int
	CallCount              int
	ActualCallOrder        CallOrder
	ExpectedCallOrder      CallOrder
	ActualDeleteArgument   string
	ExpectedDeleteArgument string
	ActualEvents           []Event
	ExpectedEvents         []Event
	ErrorEvents            []ErrorEvent
	ActualError            error
	expectedCallCount      int
	sut                    string
	urlOrigin              string
}

// ------------------------------------------------------------------------------
// git repo fixture constants shared by cmd/index and pkg/index tests
// ------------------------------------------------------------------------------

/*
	# Commit history from test fixture
	6696e0513a6dcb38e14a1da46ac7ba44611c6f90 Updating README.md (HEAD -> master)
	598ce06b5bf534e9dec0db5fd64bee88020c6571 Updating nyuad/ad_mc_019.xml, Deleting file tamwag/tam_143.xml EADID='tam_143', Updating edip/mos_2024.xml, Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Updating akkasah/ad_mc_030.xml
	50fc07058d893854b2eab1ce6285aa98d6596a16 Updating tamwag/tam_143.xml, Updating cbh/arc_212_plymouth_beecher.xml
	244e53e7827640496ead934516ccb68d5d25cb96 Deleting file fales/mss_420.xml EADID='mss_420', Updating fales/mss_420.xml
	dc63b18f64864f2bdcaffee758e4c590dac8f5ab Updating fales/mss_420.xml
	cb2d1300d7c5572bed7a6f2ec5aa67f023fe087c Deleting file fales/mss_460.xml EADID='mss_460'
	52ac657cc70005670c2ba97c23fba68ce8f1f9de Updating fales/mss_460.xml
	6c82536efc4149599c6d341e34dcc1255131c365 Deleting file akkasah/ad_mc_030.xml EADID='ad_mc_030', Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Deleting file edip/mos_2024.xml EADID='mos_2024', Deleting file fales/mss_420.xml EADID='mss_420', Deleting file fales/mss_460.xml EADID='mss_460', Deleting file nyhs/ms256_harmon_hendricks_goldstone.xml EADID='ms256_harmon_hendricks_goldstone', Deleting file nyhs/ms347_foundling_hospital.xml EADID='ms347_foundling_hospital', Deleting file nyuad/ad_mc_019.xml EADID='ad_mc_019', Deleting file tamwag/tam_143.xml EADID='tam_143'
	6c814c9836fc2abfa89d49f548fcd9cb11eae78a Updating akkasah/ad_mc_030.xml, Updating cbh/arc_212_plymouth_beecher.xml, Updating edip/mos_2024.xml, Updating fales/mss_420.xml, Updating fales/mss_460.xml, Updating nyhs/ms256_harmon_hendricks_goldstone.xml, Updating nyhs/ms347_foundling_hospital.xml, Updating nyuad/ad_mc_019.xml, Updating tamwag/tam_143.xml
*/

// hashes from the git-repo fixture (in order of commits)
/*
	# Commit history from test fixture
	6696e0513a6dcb38e14a1da46ac7ba44611c6f90 Updating README.md (HEAD -> master)
	598ce06b5bf534e9dec0db5fd64bee88020c6571 Updating nyuad/ad_mc_019.xml, Deleting file tamwag/tam_143.xml EADID='tam_143', Updating edip/mos_2024.xml, Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Updating akkasah/ad_mc_030.xml
	50fc07058d893854b2eab1ce6285aa98d6596a16 Updating tamwag/tam_143.xml, Updating cbh/arc_212_plymouth_beecher.xml
	244e53e7827640496ead934516ccb68d5d25cb96 Deleting file fales/mss_420.xml EADID='mss_420', Updating fales/mss_420.xml
	dc63b18f64864f2bdcaffee758e4c590dac8f5ab Updating fales/mss_420.xml
	cb2d1300d7c5572bed7a6f2ec5aa67f023fe087c Deleting file fales/mss_460.xml EADID='mss_460'
	52ac657cc70005670c2ba97c23fba68ce8f1f9de Updating fales/mss_460.xml
	6c82536efc4149599c6d341e34dcc1255131c365 Deleting file akkasah/ad_mc_030.xml EADID='ad_mc_030', Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Deleting file edip/mos_2024.xml EADID='mos_2024', Deleting file fales/mss_420.xml EADID='mss_420', Deleting file fales/mss_460.xml EADID='mss_460', Deleting file nyhs/ms256_harmon_hendricks_goldstone.xml EADID='ms256_harmon_hendricks_goldstone', Deleting file nyhs/ms347_foundling_hospital.xml EADID='ms347_foundling_hospital', Deleting file nyuad/ad_mc_019.xml EADID='ad_mc_019', Deleting file tamwag/tam_143.xml EADID='tam_143'
	6c814c9836fc2abfa89d49f548fcd9cb11eae78a Updating akkasah/ad_mc_030.xml, Updating cbh/arc_212_plymouth_beecher.xml, Updating edip/mos_2024.xml, Updating fales/mss_420.xml, Updating fales/mss_460.xml, Updating nyhs/ms256_harmon_hendricks_goldstone.xml, Updating nyhs/ms347_foundling_hospital.xml, Updating nyuad/ad_mc_019.xml, Updating tamwag/tam_143.xml
*/

// hashes from the git-repo fixture (in order of commits)
const AddAllHash = "6c814c9836fc2abfa89d49f548fcd9cb11eae78a"
const DeleteAllHash = "6c82536efc4149599c6d341e34dcc1255131c365"
const AddOneHash = "52ac657cc70005670c2ba97c23fba68ce8f1f9de"
const DeleteOneHash = "cb2d1300d7c5572bed7a6f2ec5aa67f023fe087c"
const DeleteModifyAddHash = "244e53e7827640496ead934516ccb68d5d25cb96"
const AddTwoHash = "50fc07058d893854b2eab1ce6285aa98d6596a16"
const AddThreeDeleteTwoHash = "598ce06b5bf534e9dec0db5fd64bee88020c6571"
const NoEADFilesInCommitHash = "6696e0513a6dcb38e14a1da46ac7ba44611c6f90"

// ------------------------------------------------------------------------------
// public functions
// ------------------------------------------------------------------------------

func AssertCallCount(t *testing.T, expectedCallCount, actualCallCount int) {
	if actualCallCount != expectedCallCount {
		t.Errorf("error: actual CallCount '%d' does not match expected CallCount '%d'", actualCallCount, expectedCallCount)
	}
}

func AssertError(t *testing.T, fname string, err error) {
	if err == nil {
		t.Errorf("error: expected '%s' to return an error, but nothing was returned", fname)
	}
}

func AssertErrorMessageContainsString(t *testing.T, fname string, err error, str string) {
	emsg := err.Error()
	if !strings.Contains(emsg, str) {
		t.Errorf("error: expected function '%s' to return an error with message containing '%s', but got: '%s'", fname, str, emsg)
	}
}

func EventsToString(events []Event) string {
	var str string
	for _, e := range events {
		if e.FuncName != Delete {
			continue
		}
		args := ""
		if e.FuncName != Add {
			args = strings.Join(e.Args, ",")
		}
		err := ""
		if e.Err != nil {
			err = e.Err.Error()
		}
		str += fmt.Sprintf("%5d  %-8s  %-35s  %s\n", e.CallCount, e.FuncName, args, err)
	}
	return str
}

func GetSolrClientMock() *SolrClientMock {
	sc := &SolrClientMock{
		GoldenFileHashes: make(map[string]string),
	}
	sc.Reset()
	return sc
}

func SortErrorEventsByCallCount(events []ErrorEvent) []ErrorEvent {
	// sort the error events by CallCount
	// using bubble sort
	for range events {
		for j := range len(events) - 1 {
			if events[j].CallCount > events[j+1].CallCount {
				events[j], events[j+1] = events[j+1], events[j]
			}
		}
	}
	return events
}

// ------------------------------------------------------------------------------
// SolrClientMock methods
// ------------------------------------------------------------------------------
func (sc *SolrClientMock) Add(xmlPostBody string) error {
	sc.CallCount++

	err := sc.updateHash(xmlPostBody)
	if err != nil {
		return err
	}

	err = sc.checkForErrorEvent()
	sc.updateEvents(Add, []string{xmlPostBody}, err)
	return err
}

func (sc *SolrClientMock) CheckAssertions() error {
	errs := []error{}

	// if an error was NOT expected, and an error was found, return error
	// if an error was expected, and no error was found, return error
	if len(sc.ErrorEvents) == 0 {
		// no errors expected
		if sc.ActualError != nil {
			return fmt.Errorf("error: expected operation NOT to return an error, but an error was returned: %v", sc.ActualError)
		}
	} else {
		// errors expected
		if sc.ActualError == nil {
			return fmt.Errorf("error: expected operation to return an error, but nothing was returned")
		}
	}

	// If there were files to be indexed, assert that all were indexed
	if sc.ExpectedCallOrder.Commit != IGNORE_CALL_ORDER && sc.NumberOfFilesToIndex > 0 {
		if !sc.IsComplete() {
			errs = append(errs, fmt.Errorf("not all files were added to the Solr index. Remaining values: \n%s", sc.GoldenFileHashesToString()))
		}
	}

	// Delete() calls
	if sc.ExpectedCallOrder.Delete != IGNORE_CALL_ORDER {
		if sc.ActualCallOrder.Delete != sc.ExpectedCallOrder.Delete {
			errs = append(errs, fmt.Errorf("Delete() was not called in the correct sequence. Expected: %d Actual: %d", sc.ExpectedCallOrder.Delete, sc.ActualCallOrder.Delete))
		}

		if sc.ActualDeleteArgument != sc.ExpectedDeleteArgument {
			errs = append(errs, fmt.Errorf("Delete() was not called with the correct argument. Expected: %s, got: %s", sc.ExpectedDeleteArgument, sc.ActualDeleteArgument))
		}
	}

	// Commit() calls
	if sc.ExpectedCallOrder.Commit != IGNORE_CALL_ORDER && sc.ActualCallOrder.Commit != sc.ExpectedCallOrder.Commit {
		errs = append(errs, fmt.Errorf("Commit() was not called in the correct sequence. Expected: %d Actual: %d", sc.ExpectedCallOrder.Commit, sc.ActualCallOrder.Commit))
	}

	// Rollback() calls
	if sc.ExpectedCallOrder.Rollback != IGNORE_CALL_ORDER && sc.ActualCallOrder.Rollback != sc.ExpectedCallOrder.Rollback {
		errs = append(errs, fmt.Errorf("Rollback() was not called in the correct sequence. Expected: %d Actual: %d", sc.ExpectedCallOrder.Rollback, sc.ActualCallOrder.Rollback))
	}

	// if there were expected errors during the operation...
	if len(sc.ErrorEvents) > 0 {
		// check that the expected errors were returned
		for i, errString := range strings.Split(sc.ActualError.Error(), "\n") {
			if errString != sc.ErrorEvents[i].ErrorMessage {
				errs = append(errs, fmt.Errorf("error: expected IndexEADFile to return an error with message '%s', but got: '%s'", sc.ErrorEvents[i].ErrorMessage, errString))
			}
		}
	}

	// Check for any failed assertions
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	// all assertions passed
	return nil
}

func (sc *SolrClientMock) CheckAssertionsViaEvents() error {
	errs := []error{}

	// fmt.Println("---------------------------------------------------------------")
	// fmt.Printf("EXPECTED:\n%v\n", EventsToString(sc.ExpectedEvents))
	// fmt.Printf("ACTUAL  :\n%v\n", EventsToString(sc.ActualEvents))

	if len(sc.ExpectedEvents) != len(sc.ActualEvents) {
		return fmt.Errorf("error: %s : mismatched events array length: expected %d events, but got %d", sc.sut, len(sc.ExpectedEvents), len(sc.ActualEvents))
	}

	for i, expectedEvent := range sc.ExpectedEvents {
		actualEvent := sc.ActualEvents[i]
		eventErrs := assertEventMatch(expectedEvent, actualEvent, sc.sut)
		if len(eventErrs) > 0 {
			errs = append(errs, eventErrs...)
		}
	}

	// Check for any failed assertions
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	// all assertions passed
	return nil
}

func (sc *SolrClientMock) Commit() error {

	sc.CallCount++

	sc.ActualCallOrder.Commit = sc.CallCount
	err := sc.checkForErrorEvent()
	sc.updateEvents(Commit, []string{}, err)
	return err
}

func (sc *SolrClientMock) Delete(eadid string) error {
	sc.CallCount++

	sc.ActualCallOrder.Delete = sc.CallCount
	sc.ActualDeleteArgument = eadid

	err := sc.checkForErrorEvent()
	sc.updateEvents(Delete, []string{eadid}, err)
	return err
}

func (sc *SolrClientMock) GetPostRequest(string) (*http.Request, error) {
	return nil, nil
}

func (sc *SolrClientMock) GetSolrURLOrigin() string {
	return sc.urlOrigin
}

func (sc *SolrClientMock) GoldenFileHashesToString() string {
	var str string
	for k, v := range sc.GoldenFileHashes {
		str += fmt.Sprintf("%s  %s\n", k, v)
	}
	return str
}

// testEAD = repositoryCode+filesystem separator+eadID (e.g. "fales/mss_460")
func (sc *SolrClientMock) InitMockForIndexing(testEAD string) error {
	// reset the solr client mock
	sc.Reset()

	// update the golden file hashes
	err := sc.updateGoldenFileHashes(testEAD)
	if err != nil {
		return err
	}

	// record the number of files to index
	sc.NumberOfFilesToIndex = len(sc.GoldenFileHashes)
	return nil
}

// function signature mirrors InitMockForIndexing()
// in case we want to add more sophisticated logic
// in the future
func (sc *SolrClientMock) InitMockForDelete(sut string) error {
	sc.Reset()
	sc.sut = sut
	return nil
}

func (sc *SolrClientMock) IsComplete() bool {
	return len(sc.GoldenFileHashes) == 0
}

func (sc *SolrClientMock) Reset() {
	// reset the solr client mock
	clear(sc.GoldenFileHashes)

	// reset the call count
	sc.CallCount = 0
	sc.expectedCallCount = 0

	// reset the call order values
	sc.ActualCallOrder.Commit = IGNORE_CALL_ORDER
	sc.ActualCallOrder.Delete = IGNORE_CALL_ORDER
	sc.ActualCallOrder.Rollback = IGNORE_CALL_ORDER
	sc.ExpectedCallOrder.Commit = IGNORE_CALL_ORDER
	sc.ExpectedCallOrder.Delete = IGNORE_CALL_ORDER
	sc.ExpectedCallOrder.Rollback = IGNORE_CALL_ORDER

	// reset the delete arguments
	sc.ActualDeleteArgument = ""
	sc.ExpectedDeleteArgument = ""

	sc.ErrorEvents = []ErrorEvent{}
	sc.ActualEvents = []Event{}
	sc.ExpectedEvents = []Event{}
	sc.ActualError = nil
	sc.sut = ""

	sc.urlOrigin = "http://www.example.com"
}

func (sc *SolrClientMock) Rollback() error {
	sc.CallCount++
	sc.ActualCallOrder.Rollback = sc.CallCount

	err := sc.checkForErrorEvent()
	sc.updateEvents(Rollback, []string{}, err)
	return err
}

func (sc *SolrClientMock) SetSolrURLOrigin(url string) {
	sc.urlOrigin = url
}

func (sc *SolrClientMock) UpdateMockForIndexEADFile(testEAD, eadid string) error {

	// snapshot the length of the golden file hashes before updating
	initialGoldenFileHashesLength := len(sc.GoldenFileHashes)
	err := sc.updateGoldenFileHashes(testEAD)
	if err != nil {
		return err
	}

	// update the expected events
	sc.addDeleteEvent(eadid)
	for i := initialGoldenFileHashesLength; i < len(sc.GoldenFileHashes); i++ {
		sc.addAddEvent()
	}

	sc.addCommitEvent()
	return nil
}

func (sc *SolrClientMock) UpdateMockForDeleteEADFileDataFromIndex(eadid string) error {
	// update the expected events
	sc.addDeleteEvent(eadid)
	sc.addCommitEvent()
	return nil
}

// ------------------------------------------------------------------------------
// private functions
// ------------------------------------------------------------------------------
func assertEventMatch(expectedEvent Event, actualEvent Event, sut string) []error {
	/*
		check that the function name is the same
		check that the call count is the same
		check that all expected args are present, except for Add
		check that all expected error substrings are present
	*/

	errs := []error{}

	if expectedEvent.FuncName != actualEvent.FuncName {
		errs = append(errs, fmt.Errorf("error: %s : expected function '%s', but got '%s'", sut, expectedEvent.FuncName, actualEvent.FuncName))
	}
	if expectedEvent.CallCount != actualEvent.CallCount {
		errs = append(errs, fmt.Errorf("error: %s : expected call count '%d', but got '%d'", sut, expectedEvent.CallCount, actualEvent.CallCount))
	}
	if len(expectedEvent.Args) != len(actualEvent.Args) {
		errs = append(errs, fmt.Errorf("error: %s : expected %d args, but got %d", sut, len(expectedEvent.Args), len(actualEvent.Args)))
	}
	// run argument comparisons for non-Add functions
	if expectedEvent.FuncName != Add {
		for i, expectedArg := range expectedEvent.Args {
			if expectedArg != actualEvent.Args[i] {
				errs = append(errs, fmt.Errorf("error: %s : expected arg '%s', but got '%s'", sut, expectedArg, actualEvent.Args[i]))
			}
		}
	}

	if expectedEvent.Err == nil && actualEvent.Err != nil {
		errs = append(errs, fmt.Errorf("error: %s : expected no error, but got '%v'", sut, actualEvent.Err))
	}

	if expectedEvent.Err != nil && actualEvent.Err == nil {
		errs = append(errs, fmt.Errorf("error: %s : expected error '%v', but got none", sut, expectedEvent.Err))
	}

	if expectedEvent.Err != nil && actualEvent.Err != nil && !strings.Contains(actualEvent.Err.Error(), expectedEvent.Err.Error()) {
		errs = append(errs, fmt.Errorf("error: %s : expected error '%v', but got '%v'", sut, expectedEvent.Err, actualEvent.Err))
	}

	return errs
}

func formattedHashSum(h hash.Hash) string {
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ------------------------------------------------------------------------------
// private SolrClientMock methods
// ------------------------------------------------------------------------------
func (sc *SolrClientMock) addAddEvent() {
	sc.expectedCallCount++
	sc.ExpectedEvents = append(sc.ExpectedEvents, Event{
		Args:      []string{"XMLPostBody"},
		CallCount: sc.expectedCallCount,
		Err:       nil,
		FuncName:  "Add",
	})
}

func (sc *SolrClientMock) addCommitEvent() {
	sc.expectedCallCount++
	sc.ExpectedEvents = append(sc.ExpectedEvents, Event{
		CallCount: sc.expectedCallCount,
		Err:       nil,
		FuncName:  "Commit",
	})
}

func (sc *SolrClientMock) addDeleteEvent(eadid string) {
	sc.expectedCallCount++
	sc.ExpectedEvents = append(sc.ExpectedEvents, Event{
		Args:      []string{eadid},
		CallCount: sc.expectedCallCount,
		Err:       nil,
		FuncName:  Delete,
	})
}

func (sc *SolrClientMock) checkForErrorEvent() error {
	// scan the error events to see if there is a match between the caller
	// and CallerName and the CallCount
	// if so, return the error message
	// iterate through range of ErrorEvents
	// if the caller name and call count match, return the error message
	// if no match, return nil
	fullyQualifiedCallerName := ""
	callerName := ""
	pc, _, _, ok := runtime.Caller(1) // 1 means caller of this function
	if ok {
		fullyQualifiedCallerName = runtime.FuncForPC(pc).Name()
		parts := strings.Split(fullyQualifiedCallerName, ".")
		if len(parts) > 0 {
			callerName = parts[len(parts)-1]
		}
	}
	if callerName == "" {
		return fmt.Errorf("unable to determine caller name from %s", fullyQualifiedCallerName)
	}

	// iterate through the error events
	// looking for a matching event
	for _, event := range sc.ErrorEvents {
		if event.FuncName == callerName && event.CallCount == sc.CallCount {
			return fmt.Errorf(event.ErrorMessage)
		}
	}

	return nil
}

func (sc *SolrClientMock) updateEvents(funcName FunctionName, args []string, err error) {
	event := Event{
		Args:      args,
		CallCount: sc.CallCount,
		Err:       err,
		FuncName:  funcName,
	}
	sc.ActualEvents = append(sc.ActualEvents, event)
}

func (sc *SolrClientMock) updateGoldenFileHashes(testEAD string) error {
	// load the golden file IDs
	// assumes all files in the directory are golden files
	// and that all files will be consumed by the test
	goldenFileIDs := eadtestutils.GetGoldenFileIDs(testEAD)

	// load the golden file hashes map
	h := md5.New()
	for _, goldenFileID := range goldenFileIDs {
		goldenFileContents, err := eadtestutils.GetGoldenFileValue(testEAD, goldenFileID)
		if err != nil {
			return err
		}

		h.Reset()
		h.Write([]byte(goldenFileContents))
		sum := formattedHashSum(h)
		goldenFilePath := eadtestutils.GoldenFilePath(testEAD, goldenFileID)

		if sc.GoldenFileHashes[sum] != "" {
			return fmt.Errorf("duplicate hash '%s' found in golden file hashes for file: %s, file already in hash: %s", sum, goldenFilePath, sc.GoldenFileHashes[sum])
		}
		// no collision, add the hash to the golden file hash map
		sc.GoldenFileHashes[sum] = goldenFilePath
	}

	return nil
}

func (sc *SolrClientMock) updateHash(xmlPostBody string) error {
	h := md5.New()
	io.WriteString(h, xmlPostBody)

	hash := formattedHashSum(h)
	if _, ok := sc.GoldenFileHashes[hash]; !ok {
		return fmt.Errorf("hash '%s' not found in golden file hashes", hash)
	}
	// remove the hash from the golden file hash map
	delete(sc.GoldenFileHashes, hash)
	return nil
}

func CallCountToIdx(callCount int) int {
	return callCount - 1
}
