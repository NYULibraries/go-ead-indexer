package testutils

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type SolrClientMock struct {
	fileDir              string            // Directory containing the "golden files" for the test
	GoldenFileHashes     map[string]string // Hashes of the golden files
	NumberOfFilesToIndex int
	CallCount            int
	CommitCallOrder      int
	DeleteCallOrder      int
	RollbackCallOrder    int
	DeleteArgument       string
	ErrorEvents          []ErrorEvent
}

type ErrorEvent struct {
	CallerName   string
	ErrorMessage string
	CallCount    int
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
	sc.CommitCallOrder = sc.CallCount
	return sc.checkForErrorEvent()
}

func (sc *SolrClientMock) Delete(eadid string) error {
	sc.CallCount++
	sc.DeleteCallOrder = sc.CallCount
	sc.DeleteArgument = eadid

	return sc.checkForErrorEvent()
}

func (sc *SolrClientMock) GetPostRequest(string) (*http.Request, error) {
	return nil, nil
}

func (sc *SolrClientMock) GetSolrURLOrigin() string {
	return "http://www.example.com"
}

func (sc *SolrClientMock) Reset() {
	// reset the solr client mock
	sc.fileDir = ""
	clear(sc.GoldenFileHashes)

	// reset the call count
	sc.CallCount = 0

	// reset the call order values
	sc.CommitCallOrder = -1
	sc.DeleteCallOrder = -1
	sc.RollbackCallOrder = -1

	// reset the delete argument
	sc.DeleteArgument = ""
	sc.ErrorEvents = []ErrorEvent{}
}

func (sc *SolrClientMock) Rollback() error {
	sc.CallCount++
	sc.RollbackCallOrder = sc.CallCount
	return nil
}

func (sc *SolrClientMock) IsComplete() bool {
	return len(sc.GoldenFileHashes) == 0
}

//func (sc *SolrClientMock) SetError(function, message string) {

func (sc *SolrClientMock) InitMock(goldenFileDir string) error {

	sc.Reset()

	// assumes all files in the directory are golden files
	// and that all files will be consumed by the test
	sc.fileDir = goldenFileDir

	files, err := os.ReadDir(goldenFileDir)
	if err != nil {
		return err
	}

	// load the golden file hashes map
	h := md5.New()
	for _, file := range files {
		filePath := filepath.Join(goldenFileDir, file.Name())

		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(h, f); err != nil {
			return err
		}

		sum := formattedHashSum(h)

		// look for collisions
		if sc.GoldenFileHashes[sum] != "" {
			return fmt.Errorf("duplicate hash '%s' found in golden file hashes for file: %s, file already in hash: %s", sum, filePath, sc.GoldenFileHashes[sum])
		}
		// no collision, add the hash to the golden file hash map
		sc.GoldenFileHashes[sum] = filePath
		h.Reset()
	}

	// record the number of files to index
	sc.NumberOfFilesToIndex = len(sc.GoldenFileHashes)
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

// Create a new random number generator using the source.
var randomSeed = time.Now().UnixNano()
var randomSource = rand.NewSource(randomSeed)
var randomGenerator = rand.New(randomSource)

func GetRandomNumber(max int) int {
	return randomGenerator.Intn(max)
}

func SortErrorEvents(events []ErrorEvent) []ErrorEvent {
	// sort the error events by CallCount
	// using bubble sort
	for i := 0; i < len(events); i++ {
		for j := 0; j < len(events)-1; j++ {
			if events[j].CallCount > events[j+1].CallCount {
				events[j], events[j+1] = events[j+1], events[j]
			}
		}
	}
	return events
}
