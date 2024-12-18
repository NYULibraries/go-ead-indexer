package testutils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"regexp"
	"strconv"
	"testing"
)

type ErrorResponseType int

const NotAnError ErrorResponseType = -1

const (
	ConnectionAborted ErrorResponseType = iota + 1
	ConnectionRefused
	ConnectionReset
	ConnectionTimeout

	HTTP302Found
	HTTP400BadRequest
	HTTP401Unauthorized
	HTTP403Forbidden
	HTTP404NotFound
	HTTP405HTTPMethodNotAllowed
	HTTP408RequestTimeout
	HTTP500InternalServerError
	HTTP502BadGateway
	HTTP503ServiceUnavailable
	HTTP504GatewayTimeout
)

const errorResponseIDPrefix = "error_"

var errorResponseTypeRegExp = regexp.MustCompile(errorResponseIDPrefix + "([a0-9]+)")

func MakeErrorResponseID(errorResponseType ErrorResponseType) string {
	errorResponseTypeString := strconv.Itoa(int(errorResponseType))

	return errorResponseIDPrefix + errorResponseTypeString
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

			errorResponseType := getErrorResponseType(id)

			if errorResponseType == NotAnError {
				err = writeActualSolrRequestToTmp(TestEAD, id, string(receivedRequest))
				if err != nil {
					t.Errorf(
						"writeActualSolrRequestToTmp(TestEAD, fileID, receivedRequest) failed with error: %s",
						err)

					return
				}
			} else {
				switch errorResponseType {
				case ConnectionAborted:
				case ConnectionRefused:
				case ConnectionReset:
				case ConnectionTimeout:
				case HTTP302Found:
				case HTTP400BadRequest:
				case HTTP401Unauthorized:
				case HTTP403Forbidden:
				case HTTP404NotFound:
				case HTTP405HTTPMethodNotAllowed:
				case HTTP408RequestTimeout:
				case HTTP500InternalServerError:
				case HTTP502BadGateway:
				case HTTP503ServiceUnavailable:
				case HTTP504GatewayTimeout:
				default:
					t.Fatalf(fmt.Sprintf("Unrecognized `ErrorResponseType`: %s",
						errorResponseType))
				}
			}
		}),
	)
}

func getErrorResponseType(id string) ErrorResponseType {
	matches := errorResponseTypeRegExp.FindStringSubmatch(id)

	if len(matches) > 1 {
		value, err := strconv.Atoi(matches[1])
		if err != nil {
			panic(fmt.Sprintf(
				"Check regular expression `errorResponseTypeRegExp`!  Error: %s",
				err))
		}

		errorResponseType := ErrorResponseType(value)

		return errorResponseType
	} else {
		return NotAnError
	}
}

func isValidSolrUpdateRequest(r *http.Request, updateURLPathAndQuery string) bool {
	var pathAndRawQuery = r.URL.Path + "?" + r.URL.RawQuery
	if pathAndRawQuery != updateURLPathAndQuery ||
		r.Method != "POST" {
		return false
	}

	return true
}
