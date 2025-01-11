package testutils

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"regexp"
	"strconv"
	"testing"
)

type ErrorResponseType string

const (
	ConnectionRefused ErrorResponseType = "connectionrefused"
	ConnectionTimeout ErrorResponseType = "connectiontimeout"

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

const errorResponseIDPrefix = "error_"

var errorResponseCounts = map[ErrorResponseType]int{}

var errorResponseTypeRegExp = regexp.MustCompile(errorResponseIDPrefix +
	"([a-z0-9]+)" + "_" + "([0-9]+)")

func MakeErrorResponseID(errorResponseType ErrorResponseType, numErrorResponsesToReturn int) string {
	if numErrorResponsesToReturn <= 0 {
		panic("`MakeErrorResponseID()` requires a positive integer for `numErrorResponsesToReturn`")
	}

	return errorResponseIDPrefix + string(errorResponseType) + "_" + strconv.Itoa(numErrorResponsesToReturn)
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

			id, err := GetID(r)
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
	if isHTTPErrorResponse(errorResponse) {
		sendErrorResponseFunction = sendHTTPErrorResponse
	} else {
		// TODO: non HTTP errors
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
			// Clear the error and send a 200 response.
			errorResponseCounts[errorResponseType] = 0
			err = send200ResponseAndWriteActualFile(w, id, receivedRequest)
		} else {
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

func isHTTPErrorResponse(errorResponse ErrorResponse) bool {
	return errorResponse.HTTPStatusCode > 0
}

func isValidSolrUpdateRequest(r *http.Request, updateURLPathAndQuery string) bool {
	var pathAndRawQuery = r.URL.Path + "?" + r.URL.RawQuery
	if pathAndRawQuery != updateURLPathAndQuery ||
		r.Method != "POST" {
		return false
	}

	return true
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
