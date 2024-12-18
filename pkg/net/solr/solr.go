package solr

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

const UpdateURLPathAndQuery = "/solr/findingaids/update?wt=json&indent=true"

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
	responseBody := bytes.NewBuffer(postBody)
	_, _ = http.Post(GetSolrURLOrigin()+UpdateURLPathAndQuery,
		"text/xml", responseBody)

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

func GetSolrURLOrigin() string {
	return solrURLOrigin
}

func SetSolrURLOrigin(solrURLOriginArg string) error {
	parsedURL, err := url.ParseRequestURI(solrURLOriginArg)
	if err != nil {
		return err
	}

	// Are the servers going to eventually be HTTPS?
	if parsedURL.Scheme != "http" {
		if parsedURL.Scheme != "https" {
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
