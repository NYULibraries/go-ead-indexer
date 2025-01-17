package testutils

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"regexp"
	"strconv"
	"testing"
	"time"
)

type ErrorResponseType string

// NOTE: It is challenging reliably simulate network errors such as connection refused,
// connection reset, and connection timeout.  It would be a lot of work to write
// tests for retry of those kinds of errors, so for now we don't.
// Even something that should be straightforward to simulate like connection
// refused proved to be difficult due to the limitations of net/http/httptest
// Server:
// * `Server` can't be stopped and restarted because `Start()` only works on
// an unstarted server.
// * Setting the URL on an unstarted server and then calling `Start()` doesn't
// work because `Start()` will panic with error "Server already started" if URL is
// not empty.
// * The URL of `Server()` is empty until it is started (needs to find an unused
// port), so there is no way to set the URL origin on the Solr client in advance.
// The `Start()` code sets URL from `.Listener.Addr().String()`, but when this
// is used in advance to set the Solr client URL origin it causes the test to
// spin.  Loading the URL into a browser results in a white page with spinner.
const (
	ContextDeadlineExceeded ErrorResponseType = "connectiontimeout"

	HTTP400BadRequest           ErrorResponseType = "http400badrequest"
	HTTP401Unauthorized         ErrorResponseType = "http401unauthorized"
	HTTP403Forbidden            ErrorResponseType = "http403forbidden"
	HTTP404NotFound             ErrorResponseType = "http404notfound"
	HTTP405HTTPMethodNotAllowed ErrorResponseType = "http405httpmethodnotallowed"
	HTTP408RequestTimeout       ErrorResponseType = "http408requesttimeout"
	HTTP500InternalServerError  ErrorResponseType = "http500internalservererror"
	HTTP502BadGateway           ErrorResponseType = "http502badgateway"
	HTTP503ServiceUnavailable   ErrorResponseType = "http503serviceunavailable"
	HTTP504GatewayTimeout       ErrorResponseType = "http504gatewaytimeout"
)

const ContextDeadlineExceededErrorResponseDuration = 1 * time.Second

const errorResponseIDPrefix = "error_"

const errorsTurnedOff = -1

var errorResponseCounts = map[ErrorResponseType]int{}

var errorResponseTypeRegExp = regexp.MustCompile(errorResponseIDPrefix +
	"([a-z0-9]+)" + "_" + "([0-9]+)")

func MakeErrorResponseIDAndPostBody(errorResponseType ErrorResponseType,
	numErrorResponsesToReturn int) (string, string) {
	id := makeErrorResponseID(errorResponseType, numErrorResponsesToReturn)
	postBody := []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<add>
  <doc>
    <field name="id">%s</field>
  </doc>
</add>`, id))

	return id, bytes.NewBuffer(postBody).String()
}

// Need to pass in `updateURLPathAndQuery` because can't use `UpdateURLPathAndQuery`
// from `solr` package directly because importing `solr` throws an import cycle
// compile error.
func MakeSolrFake(updateURLPathAndQuery string, t *testing.T) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isValidSolrUpdateRequest(r, updateURLPathAndQuery) {
				t.Fatal("Solr fake received an invalid Solr request from the test code.")
			}

			receivedRequest, err := httputil.DumpRequest(r, true)
			if err != nil {
				t.Errorf("httputil.DumpRequest(r) failed with error: %s", err)

				return
			}

			id, err := GetID(receivedRequest)
			if err != nil {
				t.Errorf("GetID(r) failed with error: %s", err)

				return
			}

			if isErrorResponseID(id) {
				err := handleErrorResponse(w, id, receivedRequest)
				if err != nil {
					t.Errorf(
						"handleErrorResponse() failed with error: %s",
						err)

					return
				}
			} else {
				err := send200ResponseAndWriteActualFile(w, id, receivedRequest)
				if err != nil {
					t.Errorf(
						"send200ResponseAndWriteActualFile() failed with error: %s",
						err)

					return
				}
			}
		}),
	)
}

func ResetErrorResponseCounts() {
	errorResponseCounts = map[ErrorResponseType]int{}
}

func getErrorResponse(id string) (ErrorResponse, error) {
	matches := errorResponseTypeRegExp.FindStringSubmatch(id)

	if len(matches) > 2 {
		errorResponseType := ErrorResponseType(matches[1])
		numRetriesRequired, err := strconv.Atoi(matches[2])
		// An error should only be possible if `errorResponseTypeRegExp` is buggy,
		// or if `MakeErrorResponseID()` does not limit the error count to int values.
		if err != nil {
			panic(err)
		}

		errorResponse, ok := errorResponseMap[errorResponseType]
		if !ok {
			return errorResponse, errors.New(
				fmt.Sprintf(`No ErrorResponse found for ID "%s"`, id))
		}
		errorResponse.NumRetriesRequired = numRetriesRequired
		errorResponse.Type = errorResponseType

		return errorResponse, nil
	} else {
		return ErrorResponse{}, errors.New(`"%s" is not a valid ErrorResponseType ID`)
	}
}

func handleErrorResponse(w http.ResponseWriter, id string, receivedRequest []byte) error {
	errorResponse, err := getErrorResponse(id)
	if err != nil {
		return err
	}

	errorResponseType := errorResponse.Type
	numRetriesRequired := errorResponse.NumRetriesRequired

	var sendErrorResponseFunction func(http.ResponseWriter, ErrorResponse) error
	if errorResponseType == ContextDeadlineExceeded {
		sendErrorResponseFunction = sendConnectionTimeoutResponse
	} else {
		sendErrorResponseFunction = sendHTTPErrorResponse
	}

	// Check the number of times this error response type has been sent, and
	// response accordingly to this current request.
	if _, ok := errorResponseCounts[errorResponseType]; !ok {
		// This is the first occurrence.  Start the count, and send the
		// error response.
		errorResponseCounts[errorResponseType] = 1
		err = sendErrorResponseFunction(w, errorResponse)
	} else {
		currentCount := errorResponseCounts[errorResponseType]
		if currentCount == numRetriesRequired {
			// Send a 200 response, and don't response with an error for this
			// error response type anymore.
			err = send200Response(w)
			errorResponseCounts[errorResponseType] = errorsTurnedOff
		} else if currentCount == errorsTurnedOff {
			// The error responses have been used up.  Send a 200 response.
			err = send200Response(w)
		} else {
			// We've not used up the errors yet.
			// Increment the error count and send an error response.
			errorResponseCounts[errorResponseType] += 1
			err = sendErrorResponseFunction(w, errorResponse)
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func isErrorResponseID(id string) bool {
	return errorResponseTypeRegExp.MatchString(id)
}

func isValidSolrUpdateRequest(r *http.Request, updateURLPathAndQuery string) bool {
	var pathAndRawQuery = r.URL.Path + "?" + r.URL.RawQuery
	if pathAndRawQuery != updateURLPathAndQuery ||
		r.Method != "POST" {
		return false
	}

	return true
}

func makeErrorResponseID(errorResponseType ErrorResponseType, numErrorResponsesToReturn int) string {
	if numErrorResponsesToReturn <= 0 {
		panic("`MakeErrorResponseID()` requires a positive integer for `numErrorResponsesToReturn`")
	}

	return errorResponseIDPrefix + string(errorResponseType) + "_" + strconv.Itoa(numErrorResponsesToReturn)
}

func send200Response(w http.ResponseWriter) error {
	return sendResponse(w, http.StatusOK, `{
  "responseHeader":{
    "status":0,
    "QTime":0}}`)
}

func send200ResponseAndWriteActualFile(w http.ResponseWriter, id string, receivedRequest []byte) error {
	err := send200Response(w)
	if err != nil {
		return err
	}

	return writeActualSolrRequestToTmp(TestEAD, id, string(receivedRequest))
}

func sendConnectionTimeoutResponse(w http.ResponseWriter, errorResponse ErrorResponse) error {
	time.Sleep(ContextDeadlineExceededErrorResponseDuration)

	return nil
}

func sendHTTPErrorResponse(w http.ResponseWriter, errorResponse ErrorResponse) error {
	return sendResponse(w, errorResponse.HTTPStatusCode, errorResponse.ResponseBody)
}

func sendResponse(w http.ResponseWriter, statusCode int, body string) error {
	w.Header().Add("Content-Type", "text/plain;charset=UTF-8")
	w.Header().Add("Transfer-Encoding", "chunked")
	// Suppress automatic header.
	w.Header()["Date"] = nil

	w.WriteHeader(statusCode)

	_, err := w.Write([]byte(body))

	return err
}
