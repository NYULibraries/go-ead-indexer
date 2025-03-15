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

const IGNORE_CALL_ORDER = -1

type CallOrder struct {
	Commit   int
	Delete   int
	Rollback int
}

type Event struct {
	Args         []string
	CallCount    int
	ErrorMessage string
	FunctionName string
}

type ErrorEvent struct {
	CallCount    int
	CallerName   string
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
	ErrorEvents            []ErrorEvent
	ActualError            error
	urlOrigin              string
}

func GetSolrClientMock() *SolrClientMock {
	sc := &SolrClientMock{
		GoldenFileHashes: make(map[string]string),
	}
	sc.Reset()
	return sc
}

func (sc *SolrClientMock) Add(xmlPostBody string) error {
	sc.CallCount++
	err := sc.updateHash(xmlPostBody)
	if err != nil {
		return err
	}

	return sc.checkForErrorEvent()
}

func (sc *SolrClientMock) Commit() error {
	sc.CallCount++
	sc.ActualCallOrder.Commit = sc.CallCount
	return sc.checkForErrorEvent()
}

func (sc *SolrClientMock) Delete(eadid string) error {
	sc.CallCount++
	sc.ActualCallOrder.Delete = sc.CallCount
	sc.ActualDeleteArgument = eadid

	return sc.checkForErrorEvent()
}

func (sc *SolrClientMock) GetPostRequest(string) (*http.Request, error) {
	return nil, nil
}

func (sc *SolrClientMock) GetSolrURLOrigin() string {
	return sc.urlOrigin
}

// testEAD = repositoryCode+filesystem separator+eadID (e.g. "fales/mss_460")
func (sc *SolrClientMock) InitMockForIndexing(testEAD string) error {

	sc.Reset()

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

	// record the number of files to index
	sc.NumberOfFilesToIndex = len(sc.GoldenFileHashes)
	return nil
}

// function signature mirrors InitMockForIndexing()
// in case we want to add more sophisticated logic
// in the future
func (sc *SolrClientMock) InitMockForDelete() error {
	sc.Reset()
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

	sc.urlOrigin = "http://www.example.com"
}

func (sc *SolrClientMock) Rollback() error {
	sc.CallCount++
	sc.ActualCallOrder.Rollback = sc.CallCount
	return sc.checkForErrorEvent()
}

func (sc *SolrClientMock) SetSolrURLOrigin(url string) {
	sc.urlOrigin = url
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

func formattedHashSum(h hash.Hash) string {
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (sc *SolrClientMock) checkForErrorEvent() error {
	// scan the error events to see if there is a match between the caller and CallerName
	// and the CallCount
	// if so, return the error message
	// iterate through range of ErrorEvents
	// if the caller name and call count match, return the error message
	// if no match, return nil
	callerName := ""
	pc, _, _, ok := runtime.Caller(1) // 1 means caller of the caller
	if ok {
		callerName = runtime.FuncForPC(pc).Name()
	}

	// iterate through the error events
	// looking for a matching event
	if callerName != "" {
		for _, event := range sc.ErrorEvents {
			if ("github.com/nyulibraries/go-ead-indexer/pkg/index/testutils.(*SolrClientMock)."+event.CallerName) == callerName && event.CallCount == sc.CallCount {
				return fmt.Errorf(event.ErrorMessage)
			}
		}
	}
	return nil
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

func (sc *SolrClientMock) CheckAssertions() error {
	errs := []error{}

	// if an error was expected, and no error was found, return error
	// if an error was NOT expected, and an error was found, return error
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
			errs = append(errs, fmt.Errorf("not all files were added to the Solr index. Remaining values: %v", sc.GoldenFileHashes))
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

func AssertCallCount(t *testing.T, expectedCallCount, actualCallCount int) {
	if actualCallCount != expectedCallCount {
		t.Errorf("error: actual CallCount '%d' does not match expected CallCount '%d'", actualCallCount, expectedCallCount)
	}
}
