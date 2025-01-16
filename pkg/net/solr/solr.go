package solr

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"syscall"
	"time"
)

type solrClient struct {
	backoffInitialInterval time.Duration
	backoffMultiplier      time.Duration
	client                 http.Client
	maxRetries             int
	urlOrigin              string
}

// No default Solr URL.
// We wouldn't want to corrupt the index of the default Solr server due to an
// accidental misconfiguration of an instance.
const DefaultBackoffInitialInterval = 1 * time.Second
const DefaultBackoffMultiplier = 4
const DefaultMaxRetries = 3
const DefaultTimeout = 30 * time.Second

const UpdateURLPathAndQuery = "/solr/findingaids/update?wt=json&indent=true"

func (sc *solrClient) Add(xmlPostBody string) error {
	response, err := sc.sendRequest(xmlPostBody)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		// Some extra characters appear in the dumped response body.  See:
		// "http resp.Write & httputil.DumpResponse include extra text with body"
		// https://groups.google.com/g/golang-nuts/c/LCoPQOpDvx4?pli=1
		//
		// To test this, removed "Transfer-Encoding: chunked" HTTP header
		// from the Solr fake responses, and extra characters no longer appeared
		// (and Content-Length header was automatically added).
		dumpedResponse, dumpResponseError := httputil.DumpResponse(response, true)
		if dumpResponseError != nil {
			return dumpResponseError
		}

		return errors.New(string(dumpedResponse))
	}

	return nil
}

func (sc *solrClient) Commit() error {
	return nil
}

func (sc *solrClient) Delete(eadID string) error {
	return nil
}

func (sc *solrClient) GetMaxRetries() int {
	return sc.maxRetries
}

func (sc *solrClient) GetPostRequest(xmlPostBody string) (*http.Request, error) {
	postRequest, err := http.NewRequest(http.MethodPost,
		sc.GetSolrURLOrigin()+UpdateURLPathAndQuery,
		bytes.NewReader([]byte(xmlPostBody)))
	if err != nil {
		return postRequest, err
	}

	postRequest.Header.Set("Content-Type", "text/xml")

	return postRequest, nil
}

func (sc *solrClient) GetSolrURLOrigin() string {
	return sc.urlOrigin
}

func (sc *solrClient) SetMaxRetries(newRetries int) error {
	if newRetries < 0 {
		return fmt.Errorf("Invalid value passed to `SetMaxRetries()`: %d", newRetries)
	}

	sc.maxRetries = newRetries

	return nil
}

func (sc *solrClient) SetSolrURLOrigin(solrURLOriginArg string) error {
	parsedURL, err := url.ParseRequestURI(solrURLOriginArg)
	if err != nil {
		return err
	}

	// Are the servers going to eventually be HTTPS?
	if parsedURL.Scheme != "http" {
		if parsedURL.Scheme == "https" {
			return errors.New(fmt.Sprintf(`SetSolrURLOrigin("%s"): https is not currently supported`,
				solrURLOriginArg))
		} else {
			return errors.New(fmt.Sprintf(`SetSolrURLOrigin("%s"): invalid scheme`,
				solrURLOriginArg))
		}
	}

	if parsedURL.Host == "" {
		return errors.New(fmt.Sprintf(`SetSolrURLOrigin("%s"): host is empty`,
			solrURLOriginArg))
	}

	sc.urlOrigin = solrURLOriginArg

	return nil
}

func (sc *solrClient) sendRequest(xmlPostBody string) (*http.Response, error) {
	request, err := sc.GetPostRequest(xmlPostBody)
	if err != nil {
		return nil, err
	}

	var response *http.Response
	numRetries := sc.GetMaxRetries()
	sleepInterval := sc.backoffInitialInterval
	for i := 0; i < numRetries+1; i++ {
		response, err = sc.client.Do(request)
		if err != nil && !isRetryableError(err) {
			break
		}

		if response != nil {
			if response.StatusCode == http.StatusOK ||
				!isRetryableHTTPError(response.StatusCode) {
				break
			}
		}

		// Restore POST body of request for next try.
		request.Body = io.NopCloser(bytes.NewBuffer([]byte(xmlPostBody)))

		// Wait.
		time.Sleep(sleepInterval)
		sleepInterval = sleepInterval * sc.backoffMultiplier
	}

	return response, err
}

func (sc *solrClient) setTimeout(timeoutArg time.Duration) {
	sc.client.Timeout = timeoutArg
}

func isRetryableError(err error) bool {
	var syscallErrno syscall.Errno
	switch {
	case errors.As(err, &syscallErrno):
		switch {
		case errors.Is(err, syscall.ECONNREFUSED),
			errors.Is(err, syscall.ECONNRESET),
			errors.Is(err, syscall.ETIMEDOUT):
			return true
		default:
			return false
		}
	case errors.Is(err, context.DeadlineExceeded):
		return true
	default:
		return false
	}
}

func isRetryableHTTPError(statusCode int) bool {
	switch statusCode {
	case http.StatusBadGateway,
		http.StatusGatewayTimeout,
		http.StatusInternalServerError,
		http.StatusRequestTimeout,
		http.StatusServiceUnavailable:
		return true

	default:
		return false
	}
}
