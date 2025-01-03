package testutils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"regexp"
	"testing"
)

type ErrorResponseType string

const NotAnError ErrorResponseType = "notanerror"

const (
	ConnectionAborted ErrorResponseType = "connectionaborted"
	ConnectionRefused ErrorResponseType = "connectionrefused"
	ConnectionReset   ErrorResponseType = "connectionreset"
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

	ConnectionTimeoutPermanent ErrorResponseType = "connectiontimeoutpermanent"
)

const errorResponseIDPrefix = "error_"

var errorResponseTypeRegExp = regexp.MustCompile(errorResponseIDPrefix + "([a-z0-9]+)")

func MakeErrorResponseID(errorResponseType ErrorResponseType) string {
	return errorResponseIDPrefix + string(errorResponseType)
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
				errorResponseType := getErrorResponseType(id)

				switch errorResponseType {
				case ConnectionAborted:
				case ConnectionRefused:
				case ConnectionReset:
				case ConnectionTimeout:
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
				case ConnectionTimeoutPermanent:
				default:
					t.Fatalf(fmt.Sprintf("Unrecognized `ErrorResponseType`: %s",
						errorResponseType))
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

func getErrorResponseType(id string) ErrorResponseType {
	matches := errorResponseTypeRegExp.FindStringSubmatch(id)

	if len(matches) > 1 {
		errorResponseType := ErrorResponseType(matches[1])

		return errorResponseType
	} else {
		return NotAnError
	}
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

func send200Response(w http.ResponseWriter) error {
	w.Header().Add("Content-Type", "text/plain;charset=UTF-8")
	w.Header().Add("Transfer-Encoding", "chunked")
	// Suppress automatic header.
	w.Header()["Date"] = nil

	_, err := w.Write([]byte(`{
  "responseHeader":{
    "status":0,
    "QTime":0}}`))

	return err
}

func send200ResponseAndWriteActualFile(w http.ResponseWriter, id string, receivedRequest []byte) error {
	err := send200Response(w)
	if err != nil {
		return err
	}

	return writeActualSolrRequestToTmp(TestEAD, id, string(receivedRequest))
}
