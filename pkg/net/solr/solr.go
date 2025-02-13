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

type SolrClient interface {
	Add(string) error
	Commit() error
	Delete(string) error
	GetPostRequest(string) (*http.Request, error)
	GetSolrURLOrigin() string
}

type solrClient struct {
	backoffInitialInterval time.Duration
	backoffMultiplier      time.Duration
	client                 http.Client
	urlOrigin              string
}

// No default Solr URL.
// We wouldn't want to corrupt the index of the default Solr server due to an
// accidental misconfiguration of an instance.
const DefaultBackoffInitialInterval = 1 * time.Second
const DefaultBackoffMultiplier = 4
const DefaultTimeout = 30 * time.Second

const UpdateURLPathAndQuery = "/solr/findingaids/update?wt=json&indent=true"

var maxRetries = 3

func NewSolrClient(urlOrigin string) (SolrClient, error) {
	solrClient, err := newSolrClient(urlOrigin)

	return &solrClient, err
}

// This is used by the tests, which require access to private `solrClient`
// data and methods.
func newSolrClient(urlOrigin string) (solrClient, error) {
	solrClient := solrClient{
		backoffInitialInterval: DefaultBackoffInitialInterval,
		backoffMultiplier:      DefaultBackoffMultiplier,
		client: http.Client{
			Timeout: DefaultTimeout,
		},
	}

	err := solrClient.setSolrURLOrigin(urlOrigin)

	return solrClient, err
}

func (sc *solrClient) Add(xmlPostBody string) error {
	return sc.solrRequest(xmlPostBody)
}

func (sc *solrClient) Commit() error {
	xmlPostBody := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<commit/>
`)

	return sc.solrRequest(xmlPostBody)
}

func (sc *solrClient) Delete(eadID string) error {
	xmlPostBody := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<delete>
  <query>ead_ssi:"%s"</query>
</delete>
`, eadID)

	return sc.solrRequest(xmlPostBody)
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

func (sc *solrClient) sendRequest(xmlPostBody string) (*http.Response, error) {
	request, err := sc.GetPostRequest(xmlPostBody)
	if err != nil {
		return nil, err
	}

	var response *http.Response
	numRetries := getMaxRetries()
	sleepInterval := sc.backoffInitialInterval
	for i := 0; i < 1+numRetries; i++ {
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

func (sc *solrClient) setSolrURLOrigin(solrURLOriginArg string) error {
	parsedURL, err := url.ParseRequestURI(solrURLOriginArg)
	if err != nil {
		return err
	}

	// Are the servers going to eventually be HTTPS?
	if parsedURL.Scheme != "http" {
		if parsedURL.Scheme == "https" {
			return errors.New(fmt.Sprintf(`setSolrURLOrigin("%s"): https is not currently supported`,
				solrURLOriginArg))
		} else {
			return errors.New(fmt.Sprintf(`setSolrURLOrigin("%s"): invalid scheme`,
				solrURLOriginArg))
		}
	}

	if parsedURL.Host == "" {
		return errors.New(fmt.Sprintf(`setSolrURLOrigin("%s"): host is empty`,
			solrURLOriginArg))
	}

	sc.urlOrigin = solrURLOriginArg

	return nil
}

func (sc *solrClient) setTimeout(timeoutArg time.Duration) {
	sc.client.Timeout = timeoutArg
}

func (sc *solrClient) solrRequest(xmlPostBody string) error {
	response, err := sc.sendRequest(xmlPostBody)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		// NOTE: some extra characters appear in the dumped response body we
		// include in the returned error.  These are chunked encoding sizes,
		// according to this discussion:
		// "http resp.Write & httputil.DumpResponse include extra text with body"
		// https://groups.google.com/g/golang-nuts/c/LCoPQOpDvx4?pli=1
		// Confirmed this by removing the "Transfer-Encoding: chunked" HTTP header
		// from the Solr fake responses, which did away with the extra characters
		// and added the Content-Length header.
		dumpedResponse, dumpResponseError := httputil.DumpResponse(response, true)
		if dumpResponseError != nil {
			return dumpResponseError
		}

		return errors.New(string(dumpedResponse))
	}

	return nil
}

func getMaxRetries() int {
	return maxRetries
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
