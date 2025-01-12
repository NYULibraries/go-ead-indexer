package solr

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

const UpdateURLPathAndQuery = "/solr/findingaids/update?wt=json&indent=true"

const DefaultRetries = 3

var retries = DefaultRetries

// No default Solr URL.
// We wouldn't want to corrupt the index of the default Solr server due to an
// accidental misconfiguration of an instance.
var solrURLOrigin string

// TODO: Obviously replace this fake stuff after `TestAdd()` is completed.  There
// needs to be actual files written out for the test to do diff'ing against.
func Add(xmlPostBody string) error {
	var idRegExp = regexp.MustCompile(`<field name="id">([a-z0-9_-]+)</field>`)

	matches := idRegExp.FindStringSubmatch(xmlPostBody)
	if len(matches) < 2 {
		return errors.New("No id found")
	}

	id := matches[1]

	postBody := []byte(`<field name="id">` + id + "</field>")
	postBodyBuffer := bytes.NewBuffer(postBody)
	response, err := http.Post(GetSolrURLOrigin()+UpdateURLPathAndQuery,
		"text/xml", postBodyBuffer)
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

func GetPOSTRequest(eadID string) error {
	return nil
}

func GetRetries() int {
	return retries
}

func GetSolrURLOrigin() string {
	return solrURLOrigin
}

func SetRetries(newRetries int) error {
	if newRetries < 0 {
		return fmt.Errorf("Invalid value passed to `SetRetries()`: %d", newRetries)
	}

	retries = newRetries

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
