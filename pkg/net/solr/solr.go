package solr

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const DefaultMaxRetries = 3
const DefaultTimeout = 30 * time.Second

const UpdateURLPathAndQuery = "/solr/findingaids/update?wt=json&indent=true"

var client = http.Client{
	Timeout: DefaultTimeout,
}

var maxRetries = DefaultMaxRetries

// No default Solr URL.
// We wouldn't want to corrupt the index of the default Solr server due to an
// accidental misconfiguration of an instance.
var solrURLOrigin string

func Add(xmlPostBody string) error {
	response, err := sendRequest(xmlPostBody)
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

func Commit() error {
	return nil
}

func Delete(eadID string) error {
	return nil
}

func GetMaxRetries() int {
	return maxRetries
}

func GetPOSTRequest(eadID string) error {
	return nil
}

func GetSolrURLOrigin() string {
	return solrURLOrigin
}

func SetMaxRetries(newRetries int) error {
	if newRetries < 0 {
		return fmt.Errorf("Invalid value passed to `SetMaxRetries()`: %d", newRetries)
	}

	maxRetries = newRetries

	return nil
}

func SetSolrURLOrigin(solrURLOriginArg string) error {
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

	solrURLOrigin = solrURLOriginArg

	return nil
}

func sendRequest(xmlPostBody string) (*http.Response, error) {
	response, err := client.Post(GetSolrURLOrigin()+UpdateURLPathAndQuery,
		"text/xml", bytes.NewBuffer([]byte(xmlPostBody)))
	if err != nil {
		return response, err
	}

	return response, nil
}

func setTimeout(timeoutArg time.Duration) {
	client.Timeout = timeoutArg
}
