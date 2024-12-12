package solr

import (
	"flag"
	"fmt"
	"go-ead-indexer/pkg/ead/testutils"
	"go-ead-indexer/pkg/util"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
)

const testEAD = "appdev/mos_2021"

var tmpFilesDirPath = filepath.Join("testdata", "tmp", "actual")

const fakeSolrHostAndPort = "fake-solr-host.library.nyu.edu:8080"

var idRegExp = regexp.MustCompile(`<field name="id">([a-z0-9_-]+)</field>`)

var postRequestHTTPHeadersString = fmt.Sprintf(
	`POST /solr/findingaids/update?wt=json HTTP/1.1
Content-Type: text/xml
Connection: close
Host: %s
Content-Length: `, fakeSolrHostAndPort)

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func TestAdd(t *testing.T) {
	err := clean()
	if err != nil {
		t.Errorf("clean() failed with error: %s", err)
	}

	fakeSolrServer := makeSolrFake(t)
	defer fakeSolrServer.Close()

	err = SetSolrURLOrigin(fakeSolrServer.URL)
	if err != nil {
		t.Fatalf(`Setup of Solr fake failed with error: %s`, err)
	}

	goldenFileIDs := testutils.GetGoldenFileIDs(testEAD)
	for _, goldenFileID := range goldenFileIDs {
		testAdd(goldenFileID, t)
	}
}

func testAdd(goldenFileID string, t *testing.T) {
	postBody, err := testutils.GetGoldenFileValue(testEAD, goldenFileID)
	if err != nil {
		t.Fatalf("testutils.GetGoldenFileValue(testEAD, goldenFileID) failed with error: %s", err)
	}

	err = Add(postBody)
	if err != nil {
		t.Fatalf("Expected no error for %s, got: %s", goldenFileID, err)
	}

	actualRequest, err := getActualFileContents(testEAD, goldenFileID)
	if err != nil {
		t.Fatalf("getActualFileContents(testEAD, goldenFileID) failed with error: %s", err)
	}

	expectedRequest := getExpectedPOSTRequestString(postBody)
	diff := util.DiffStrings("expected", expectedRequest,
		"actual", actualRequest)
	if diff != "" {
		t.Errorf(`%s fail: actual request does not match expected: %s`,
			goldenFileID, diff)
	}
}

func TestSetSolrURL(t *testing.T) {
	testCases := []struct {
		origin        string
		expectedError string
	}{
		{
			origin:        "",
			expectedError: `parse "": empty url`,
		},
		{
			origin:        "x",
			expectedError: `parse "x": invalid URI for request`,
		},
		{
			origin:        "ftp://solr-host.org",
			expectedError: `SetSolrURLOrigin("ftp://solr-host.org"): invalid scheme`,
		},
		{
			origin:        "http://",
			expectedError: `SetSolrURLOrigin("http://"): host is empty`,
		},
		{
			origin:        "https://",
			expectedError: `SetSolrURLOrigin("https://"): host is empty`,
		},
		{
			origin:        fakeSolrHostAndPort,
			expectedError: `SetSolrURLOrigin("` + fakeSolrHostAndPort + `"): invalid scheme`,
		},
		{
			origin:        "http://" + fakeSolrHostAndPort,
			expectedError: "",
		},
	}

	for _, testCase := range testCases {
		actualError := SetSolrURLOrigin(testCase.origin)
		var actualErrorString string
		if actualError != nil {
			actualErrorString = actualError.Error()
		}

		if testCase.expectedError == "" {
			if actualErrorString != "" {
				t.Errorf(`SetSolrURLOrigin("%s") should not have returned an error,`+
					` but it returned error "%s"`, testCase.origin, actualErrorString)
			}

			actualOrigin := GetSolrURLOrigin()
			if actualOrigin != testCase.origin {
				t.Errorf(`GetSolrURLOrigin() should return "%s", but it instead`+
					` returned "%s"`, testCase.origin, actualOrigin)
			}
		} else {
			if actualErrorString == "" {
				t.Errorf(`SetSolrURLOrigin("%s") should have returned error "%s",`+
					" but no error was returned", testCase.origin, testCase.expectedError)
			} else if actualErrorString != testCase.expectedError {
				t.Errorf(`SetSolrURLOrigin("%s") should have returned error "%s",`+
					` but instead returned error "%s"`, testCase.origin, testCase.expectedError,
					actualErrorString)
			}

			actualOrigin := GetSolrURLOrigin()
			if actualOrigin != "" {
				t.Errorf(`GetSolrURLOrigin() should return "", but it instead`+
					` returned "%s"`, actualOrigin)
			}
		}
	}
}

func clean() error {
	err := os.RemoveAll(tmpFilesDirPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(tmpFilesDirPath, 0700)
	if err != nil {
		return err
	}

	_, err = os.Create(filepath.Join(tmpFilesDirPath, ".gitkeep"))
	if err != nil {
		return err
	}

	return nil
}

func getActualFileContents(testEAD string, fileID string) (string, error) {
	actualFile := tmpFile(testEAD, fileID)

	bytes, err := os.ReadFile(actualFile)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func getExpectedPOSTRequestString(body string) string {
	return fmt.Sprintf("%s\n\n%s", getPOSTRequestHTTPHeadersString(body), body)
}

func getFileIDFromRequest(r *http.Request) (string, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	matches := idRegExp.FindStringSubmatch(string(bodyBytes))

	if len(matches) > 1 {
		return matches[1], nil
	} else {
		return "", nil
	}
}

func getPOSTRequestHTTPHeadersString(body string) string {
	// `Content-Length` should be size in bytes, not characters, so using
	// `len()` is correct, even if multi-rune characters are used in `body`.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Length
	var charLength = len(body)

	return postRequestHTTPHeadersString + strconv.Itoa(charLength)
}

func makeSolrFake(t *testing.T) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedRequest, err := httputil.DumpRequest(r, true)
			if err != nil {
				t.Errorf("httputil.DumpRequest(r) failed with error: %s", err)

				return
			}

			fileID, err := getFileIDFromRequest(r)
			if err != nil {
				t.Errorf("getFileIDFromRequest(r) failed with error: %s", err)

				return
			}

			err = writeActualSolrRequestToTmp(testEAD, fileID, string(receivedRequest))
			if err != nil {
				t.Errorf(
					"writeActualSolrRequestToTmp(testEAD, fileID, receivedRequest) failed with error: %s",
					err)

				return
			}
		}),
	)
}

func tmpFile(testEAD string, fileID string) string {
	return filepath.Join(tmpFilesDirPath, testEAD, fileID+".xml")
}

func writeActualSolrRequestToTmp(testEAD string, fileID string, actual string) error {
	tmpFile := tmpFile(testEAD, fileID)
	err := os.MkdirAll(filepath.Dir(tmpFile), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(tmpFile, []byte(actual), 0644)
}
