// Package index provides an interface to the EAD indexing process
//
// The SetSolrClient() function must be called before calling any of the
// indexing functions in this package.  This is because the Solr client
// is a package-level variable, and the default value is nil.
package index

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nyulibraries/go-ead-indexer/pkg/ead"
	"github.com/nyulibraries/go-ead-indexer/pkg/net/solr"
	"github.com/nyulibraries/go-ead-indexer/pkg/util"
)

const MessageKey = "index"

var sc = solr.SolrClient(nil)

func SetSolrClient(solrClient solr.SolrClient) {
	sc = solrClient
}

const errSolrClientNotSet = "you must call `SetSolrClient()` before calling any indexing functions"

func assertSolrClientSet() error {
	if sc == nil {
		return errors.New(errSolrClientNotSet)
	}

	if sc.GetSolrURLOrigin() == "" {
		return errors.New("the SolrClient URL origin is not set")
	}

	return nil
}

func IndexEADFile(eadPath string) error {

	var errs []error

	// assert that the SolrClient has been set
	err := assertSolrClientSet()
	if err != nil {
		errs = append(errs, err)
		return errors.Join(errs...)
	}

	// Check if the EAD file path is absolute
	if !filepath.IsAbs(eadPath) {
		errs = append(errs, fmt.Errorf("EAD file path must be absolute: %s", eadPath))
		return errors.Join(errs...)
	}

	// Get the EAD's repository code
	repoCode, err := util.GetRepositoryCode(eadPath)
	if err != nil {
		errs = append(errs, err)
		return errors.Join(errs...)
	}

	// Read the EAD file
	eadXML, err := os.ReadFile(eadPath)
	if err != nil {
		errs = append(errs, err)
		return errors.Join(errs...)
	}

	// Parse the EAD file
	EAD, err := ead.New(repoCode, string(eadXML))
	if err != nil {
		errs = append(errs, err)
		return errors.Join(errs...)
	}

	// Delete the data for this EAD from Solr
	err = sc.Delete(EAD.CollectionDoc.Parts.EADID.Values[0])
	if err != nil {
		sc.Rollback()
		errs = append(errs, err)
		return errors.Join(errs...)
	}

	// Add the EAD Collection-level document to Solr
	xmlPostBody := EAD.CollectionDoc.SolrAddMessage.String()
	err = sc.Add(string(xmlPostBody))
	if err != nil {
		sc.Rollback()
		errs = append(errs, err)
		return errors.Join(errs...)
	}

	// Add the EAD Component-level documents to Solr
	for _, component := range *EAD.Components {
		xmlPostBody = component.SolrAddMessage.String()

		err = sc.Add(string(xmlPostBody))
		if err != nil {
			errs = append(errs, err)
		}
	}

	// Rollback if there were any errors during the component-level indexing
	if errs != nil {
		sc.Rollback()
		return errors.Join(errs...)
	}

	// commit the documents to Solr
	err = sc.Commit()
	if err != nil {
		sc.Rollback()
		errs = append(errs, err)
		return errors.Join(errs...)
	}

	return nil
}
