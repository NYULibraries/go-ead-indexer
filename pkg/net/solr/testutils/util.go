package testutils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const FakeSolrHostAndPort = "fake-solr-host.library.nyu.edu:8080"
const TestEAD = "appdev/mos_2021"

var idRegExp = regexp.MustCompile(`<field name="id">([a-z0-9_-]+)</field>`)

var tmpFilesDirPath = filepath.Join("testdata", "tmp", "actual")

var postRequestHTTPHeadersString = fmt.Sprintf(
	`POST /solr/findingaids/update?wt=json HTTP/1.1
Content-Type: text/xml
Connection: close
Host: %s
Content-Length: `, FakeSolrHostAndPort)

func Clean() error {
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

func GetActualFileContents(testEAD string, fileID string) (string, error) {
	actualFile := tmpFile(testEAD, fileID)

	bytes, err := os.ReadFile(actualFile)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func GetExpectedPOSTRequestString(body string) string {
	return fmt.Sprintf("%s\n\n%s", getPOSTRequestHTTPHeadersString(body), body)
}

func GetID(r *http.Request) (string, error) {
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
