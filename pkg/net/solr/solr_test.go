package solr

import (
	"flag"
	eadtestutils "go-ead-indexer/pkg/ead/testutils"
	"go-ead-indexer/pkg/net/solr/testutils"
	"go-ead-indexer/pkg/util"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var fakeSolrServer *httptest.Server

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func TestAdd(t *testing.T) {
	// Have to pass in `UpdateURLPathAndQuery` to `testutils` sub-package, which
	// can't import its own parent package.
	fakeSolrServer = testutils.MakeSolrFake(UpdateURLPathAndQuery, t)
	defer fakeSolrServer.Close()

	err := SetSolrURLOrigin(fakeSolrServer.URL)
	if err != nil {
		t.Fatalf(`Setup of Solr fake failed with error: %s`, err)
	}

	t.Run("Do not retry indefinitely", testAdd_doNotRetryIndefinitely)
	t.Run("Never retry certain HTTP errors", testAdd_neverRetryCertainHTTPErrors)
	// TODO: Re-enable once these pass
	//t.Run("Retry certain HTTP errors", testAdd_retryCertainHTTPErrors)
	//t.Run("Retry connection refused errors", testAdd_retryConnectionRefused)
	//t.Run("Retry connection timeout errors", testAdd_retryConnectionTimeouts)
	t.Run("Successfully add", testAdd_successAdds)
}

func testAdd_doNotRetryIndefinitely(t *testing.T) {
	testutils.ResetErrorResponseCounts()

	const expectedError = `HTTP/1.1 408 Request Timeout
Transfer-Encoding: chunked
Content-Type: text/plain;charset=UTF-8

60
{"responseHeader":{"status":408,"QTime":0},"error":{"msg":"[http408requesttimeout]","code":408}}
0

`

	// Have Solr fake error out more times than `Add()` will retry.
	id, postBody := testutils.MakeErrorResponseIDAndPostBody(
		testutils.HTTP408RequestTimeout, GetRetries()+1)

	err := Add(postBody)

	if err == nil {
		t.Errorf(`Expected Add() for id="%s" to return an error, but no error was returned`,
			id)

		return
	}

	// The returned error contains carriage returns, which would be a pain
	// to get into the copy/pasted values above, so we just remove it from
	// actual values before comparison.  Not using golden files for these
	// because they very likely will never change.
	massagedActualError := strings.ReplaceAll(err.Error(), "\r", "")
	if massagedActualError != expectedError {
		t.Errorf(`Expected request for id="%s" to return error "%s"`+
			` but got error "%s"`, id, expectedError, massagedActualError)
	}
}

// Test that `Add()` will not attempt to retry certain errors which are not worth
// retrying.
func testAdd_neverRetryCertainHTTPErrors(t *testing.T) {
	testutils.ResetErrorResponseCounts()

	const expectedErrorHTTP400BadRequest = `HTTP/1.1 400 Bad Request
Transfer-Encoding: chunked
Content-Type: text/plain;charset=UTF-8

5f
{"responseHeader":{"status":400,"QTime":0},"error":{"msg":"missing content stream","code":400}}
0

`
	const expectedErrorHTTPE401Unauthorized = `HTTP/1.1 401 Unauthorized
Transfer-Encoding: chunked
Content-Type: text/plain;charset=UTF-8

5e
{"responseHeader":{"status":401,"QTime":0},"error":{"msg":"[http401unauthorized]","code":401}}
0

`
	const expectedErrorHTTP403Forbidden = `HTTP/1.1 403 Forbidden
Transfer-Encoding: chunked
Content-Type: text/plain;charset=UTF-8

5b
{"responseHeader":{"status":403,"QTime":0},"error":{"msg":"[http403forbidden]","code":403}}
0

`
	const expectedErrorHTTP404NotFound = `HTTP/1.1 404 Not Found
Transfer-Encoding: chunked
Content-Type: text/plain;charset=UTF-8

56a
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=ISO-8859-1"/>
<title>Error 404 Not Found</title>
</head>
<body><h2>HTTP ERROR 404</h2>
<p>Problem accessing /solr/nonexistent-path.. Reason:
<pre>    Not Found</pre></p><hr /><i><small>Powered by Jetty://</small></i><br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                

</body>
</html>
0

`
	const expectedErrorHTTP405MethodNotAllowed = `HTTP/1.1 405 Method Not Allowed
Transfer-Encoding: chunked
Content-Type: text/plain;charset=UTF-8

5ac
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=ISO-8859-1"/>
<title>Error 405 HTTP method POST is not supported by this URL</title>
</head>
<body><h2>HTTP ERROR 405</h2>
<p>Problem accessing /solr/admin.html.. Reason:
<pre>    HTTP method POST is not supported by this URL</pre></p><hr /><i><small>Powered by Jetty://</small></i><br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                

</body>
</html>
0

`
	testCases := []struct {
		errorResponseType testutils.ErrorResponseType
		expectedError     string
	}{
		{
			errorResponseType: testutils.HTTP400BadRequest,
			expectedError:     expectedErrorHTTP400BadRequest,
		},
		{
			errorResponseType: testutils.HTTP401Unauthorized,
			expectedError:     expectedErrorHTTPE401Unauthorized,
		},
		{
			errorResponseType: testutils.HTTP403Forbidden,
			expectedError:     expectedErrorHTTP403Forbidden,
		},
		{
			errorResponseType: testutils.HTTP404NotFound,
			expectedError:     expectedErrorHTTP404NotFound,
		},
		{
			errorResponseType: testutils.HTTP405HTTPMethodNotAllowed,
			expectedError:     expectedErrorHTTP405MethodNotAllowed,
		},
	}

	for _, testCase := range testCases {
		// Set `numErrorResponsesToReturn` to 1 because `Add()` should never retry
		// these kinds of errors.
		id, postBody := testutils.MakeErrorResponseIDAndPostBody(
			testCase.errorResponseType, 1)

		err := Add(postBody)

		if err == nil {
			t.Errorf(`Expected Add() for id="%s" to return an error, but no error was returned`,
				id)

			continue
		}

		// The returned error contains carriage returns, which would be a pain
		// to get into the copy/pasted values above, so we just remove it from
		// actual values before comparison.  Not using golden files for these
		// because they very likely will never change.
		massagedActualError := strings.ReplaceAll(err.Error(), "\r", "")
		if massagedActualError != testCase.expectedError {
			t.Errorf(`Expected request for id="%s" to return error "%s", `+
				` but got error "%s"`, id, testCase.expectedError, err.Error())
		}
	}
}

func testAdd_retryCertainHTTPErrors(t *testing.T) {
	testutils.ResetErrorResponseCounts()
	errorResponseTypes := []testutils.ErrorResponseType{
		testutils.HTTP408RequestTimeout,
		testutils.HTTP500InternalServerError,
		testutils.HTTP502BadGateway,
		testutils.HTTP503ServiceUnavailable,
		testutils.HTTP504GatewayTimeout,
	}

	for _, errorResponseType := range errorResponseTypes {
		id, postBody := testutils.MakeErrorResponseIDAndPostBody(
			errorResponseType, GetRetries())

		err := Add(postBody)

		if err != nil {
			t.Errorf(`Expected request for id="%s" to succeed, but it failed with error "%s"`,
				id, err.Error())
		}
	}
}

func testAdd_retryConnectionRefused(t *testing.T) {
	testutils.ResetErrorResponseCounts()

	// Set Solr origin to the address of an unused port on localhost.
	connectionRefusedOrigin := "http://" + util.GetUnusedLocalhostNetworkAddress()
	err := SetSolrURLOrigin(connectionRefusedOrigin)
	if err != nil {
		t.Fatalf(`Setup of Solr fake failed with error: %s`, err)
	}

	// TODO: After we have decided upon and implemented the retries, set this
	// timer to execute its function before the first retry interval passes.
	//
	// Switch to Solr fake before `Add()` is done with its retries.
	// This test was initially itself tested for correctness by having `Add()`
	// do several consecutive POST requests with no break in between them.
	// It's worth noting that one millisecond was the maximum number of
	// milliseconds that could be used to get this test to pass with four POST
	// requests in a row.  Even two milliseconds gave the `Add()` too much time
	// to do attempt all the retries.
	time.AfterFunc(1*time.Millisecond, func() {
		err := SetSolrURLOrigin(fakeSolrServer.URL)
		if err != nil {
			t.Fatalf(`Setup of Solr fake failed with error: %s`, err)
		}
	})

	postBody := `<?xml version="1.0" encoding="UTF-8"?>
<add>
  <doc>
    <field name="id">connectionrefused_0</field>
  </doc>
</add>`

	err = Add(postBody)
	if err != nil {
		t.Errorf(`Expected simulated retry of connection refused request`+
			` to succeed, but it failed with error "%s"`, err.Error())
	}
}

func testAdd_retryConnectionTimeouts(t *testing.T) {
	testutils.ResetErrorResponseCounts()

	setTimeout(testutils.ConnectionTimeoutDuration)

	id, postBody := testutils.MakeErrorResponseIDAndPostBody(
		testutils.ConnectionTimeout, GetRetries())

	err := Add(postBody)
	if err != nil {
		t.Errorf(`Expected request for id="%s" to succeed, but it failed with error "%s"`,
			id, err.Error())
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
	massagedActualRequest := testutils.MassagedGoHTTPClientRequest(actualRequest)

	expectedRequest := testutils.GetExpectedPOSTRequestString(postBody)
	diff := util.DiffStrings("expected", expectedRequest,
		"actual", massagedActualRequest)
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
	t.Run("Errors", testSetSolrURLOrigin_errors)
	t.Run("Successfully set URL origin", testSetSolrURLOrigin_normal)
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
