package testutils

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

func MakeSolrFake(t *testing.T) *httptest.Server {
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

			err = writeActualSolrRequestToTmp(TestEAD, fileID, string(receivedRequest))
			if err != nil {
				t.Errorf(
					"writeActualSolrRequestToTmp(TestEAD, fileID, receivedRequest) failed with error: %s",
					err)

				return
			}
		}),
	)
}
