package solr

import (
	"bytes"
	"flag"
	"fmt"
	eadtestutils "go-ead-indexer/pkg/ead/testutils"
	"go-ead-indexer/pkg/net/solr/testutils"
	"go-ead-indexer/pkg/util"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func TestAdd(t *testing.T) {
	// TODO: Re-enable once these are fully implemented.
	//testAdd_failAdds(t)
	//testAdd_retryFailAdds(t)
	//testAdd_retrySuccessAdds(t)
	// TODO: Re-enable once these pass.
	//testAdd_successAdds(t)
}

func errorResponseXMLPostBody(id string) string {
	postBody := []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<add>
  <doc>
    <field name="id">%s</field>
  </doc>
</add>`, id))

	return bytes.NewBuffer(postBody).String()
}

// Test requests which trigger error responses for which retries are not attempted.
func testAdd_failAdds(t *testing.T) {
	// Have to pass in `UpdateURLPathAndQuery` to `testutils` sub-package, which
	// can't import its own parent package.
	fakeSolrServer := testutils.MakeSolrFake(UpdateURLPathAndQuery, t)
	defer fakeSolrServer.Close()

	err := SetSolrURLOrigin(fakeSolrServer.URL)
	if err != nil {
		t.Fatalf(`Setup of Solr fake failed with error: %s`, err)
	}

	testCases := []struct {
		errorResponseType testutils.ErrorResponseType
		expectedError     string
	}{
		{
			errorResponseType: testutils.HTTP400BadRequest,
			expectedError:     "",
		},
		{
			errorResponseType: testutils.HTTP401Unauthorized,
			expectedError:     "",
		},
		{
			errorResponseType: testutils.HTTP403Forbidden,
			expectedError:     "",
		},
		{
			errorResponseType: testutils.HTTP404NotFound,
			expectedError:     "",
		},
		{
			errorResponseType: testutils.HTTP405HTTPMethodNotAllowed,
			expectedError:     "",
		},
		{
			errorResponseType: testutils.HTTP403Forbidden,
			expectedError:     "",
		},
	}

	for _, testCase := range testCases {
		id := testutils.MakeErrorResponseID(testCase.errorResponseType)

		err := Add(errorResponseXMLPostBody(id))

		if err == nil {
			t.Errorf(`Expected Add() for id="%s" to return an error, but no error was returned`,
				id)

			continue
		}

		if err.Error() != testCase.expectedError {
			t.Errorf(`Expected request for id="%s" to return error "%s", `+
				` but got error "%s"`, id, testCase.expectedError, err.Error())
		}
	}
}

func testAdd_retryFailAdds(t *testing.T) {
	const expectedError = ""

	id := testutils.MakeErrorResponseID(testutils.ConnectionTimeoutPermanent)

	err := Add(errorResponseXMLPostBody(id))

	if err == nil {
		t.Errorf(`Expected Add() for id="%s" to return an error, but no error was returned`,
			id)

		return
	}

	if err.Error() != expectedError {
		t.Errorf(`Expected request for id="%s" to return error "%s", `+
			` but got error "%s"`, id, expectedError, err.Error())
	}
}

func testAdd_retrySuccessAdds(t *testing.T) {
	errorResponseTypes := []testutils.ErrorResponseType{
		testutils.ConnectionAborted,
		testutils.ConnectionRefused,
		testutils.ConnectionReset,
		testutils.ConnectionTimeout,
		testutils.HTTP408RequestTimeout,
		testutils.HTTP500InternalServerError,
		testutils.HTTP502BadGateway,
		testutils.HTTP503ServiceUnavailable,
		testutils.HTTP504GatewayTimeout,
	}

	for _, errorResponseType := range errorResponseTypes {
		id := testutils.MakeErrorResponseID(errorResponseType)

		err := Add(errorResponseXMLPostBody(id))

		if err != nil {
			t.Errorf(`Expected request for id="%s" to succeed, but it failed with error "%s"`,
				id, err.Error())
		}
	}
}

func testAdd_successAdd(goldenFileID string, t *testing.T) {
	postBody, err := eadtestutils.GetGoldenFileValue(testutils.TestEAD, goldenFileID)
	if err != nil {
		t.Fatalf("eadtestutils.GetGoldenFileValue(testutils.TestEAD, goldenFileID) failed with error: %s", err)
	}

	err = Add(postBody)
	if err != nil {
		t.Fatalf("Expected no error for %s, got: %s", goldenFileID, err)
	}

	actualRequest, err := testutils.GetActualFileContents(testutils.TestEAD, goldenFileID)
	if err != nil {
		t.Fatalf("testutils.getActualFileContents(testutils.TestEAD, goldenFileID) failed with error: %s", err)
	}

	expectedRequest := testutils.GetExpectedPOSTRequestString(postBody)
	diff := util.DiffStrings("expected", expectedRequest,
		"actual", actualRequest)
	if diff != "" {
		t.Errorf(`%s fail: actual request does not match expected: %s`,
			goldenFileID, diff)
	}
}

func testAdd_successAdds(t *testing.T) {
	err := testutils.Clean()
	if err != nil {
		t.Errorf("clean() failed with error: %s", err)
	}

	// Have to pass in `UpdateURLPathAndQuery` to `testutils` sub-package, which
	// can't import its own parent package.
	fakeSolrServer := testutils.MakeSolrFake(UpdateURLPathAndQuery, t)
	defer fakeSolrServer.Close()

	err = SetSolrURLOrigin(fakeSolrServer.URL)
	if err != nil {
		t.Fatalf(`Setup of Solr fake failed with error: %s`, err)
	}

	goldenFileIDs := eadtestutils.GetGoldenFileIDs(testutils.TestEAD)
	for _, goldenFileID := range goldenFileIDs {
		testAdd_successAdd(goldenFileID, t)
	}
}

func TestGetRetries(t *testing.T) {
	actual := GetRetries()
	if actual != DefaultRetries {
		t.Errorf(`Expected %d, got %d`, DefaultRetries, actual)
	}
}

func TestSetRetries(t *testing.T) {
	testSetRetries_badInput(t)
	testSetRetries_normal(t)
}

func testSetRetries_badInput(t *testing.T) {
	err := SetRetries(-1)
	if err == nil {
		t.Error("Expected `SetRetries(-1)` to return an error, but no error was returned")
	}
}

func testSetRetries_normal(t *testing.T) {
	newRetries := 999
	err := SetRetries(newRetries)
	if err != nil {
		t.Errorf("`SetRetries(%d)` returned an error: %s",
			newRetries, err.Error())
	}

	if GetRetries() != newRetries {
		t.Errorf("Expected `GetRetries()` to return %d, but it returned %d",
			newRetries, retries)
	}
}

func TestSetSolrURLOrigin(t *testing.T) {
	testSetSolrURLOrigin_errors(t)
	testSetSolrURLOrigin_normal(t)
}

func testSetSolrURLOrigin_errors(t *testing.T) {
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
			expectedError: `SetSolrURLOrigin("https://"): https is not currently supported`,
		},
		{
			origin:        testutils.FakeSolrHostAndPort,
			expectedError: `SetSolrURLOrigin("` + testutils.FakeSolrHostAndPort + `"): invalid scheme`,
		},
	}

	initialSolrURLOrigin := GetSolrURLOrigin()

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
			if actualOrigin != initialSolrURLOrigin {
				t.Errorf("`GetSolrURLOrigin()` should have returned the"+
					` unchanged initial value "%s", but it instead returned "%s"`,
					initialSolrURLOrigin, actualOrigin)
			}
		}
	}
}

func testSetSolrURLOrigin_normal(t *testing.T) {
	testCases := []struct {
		origin        string
		expectedError string
	}{
		{
			origin:        "http://" + testutils.FakeSolrHostAndPort,
			expectedError: "",
		},
	}

	for _, testCase := range testCases {
		err := SetSolrURLOrigin(testCase.origin)
		if err != nil {
			t.Errorf(`SetSolrURLOrigin("%s") should not have returned an error,`+
				` but it returned error "%s"`, testCase.origin, err.Error())
		}

		actualOrigin := GetSolrURLOrigin()
		if actualOrigin != testCase.origin {
			t.Errorf(`GetSolrURLOrigin() should have returned "%s",`+
				` but it instead returned "%s"`, testCase.origin, actualOrigin)
		}
	}
}
