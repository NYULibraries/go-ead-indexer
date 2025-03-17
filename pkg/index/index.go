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
	"github.com/nyulibraries/go-ead-indexer/pkg/ead/eadutil"
	"github.com/nyulibraries/go-ead-indexer/pkg/git"
	"github.com/nyulibraries/go-ead-indexer/pkg/net/solr"
	"github.com/nyulibraries/go-ead-indexer/pkg/util"
)

const MessageKey = "index"
const errSolrClientNotSet = "you must call `SetSolrClient()` before calling any indexing functions"

var sc = solr.SolrClient(nil)

func DeleteEADFileDataFromIndex(eadID string) error {
	var errs []error

	// assert that the EADID is valid
	if !eadutil.IsValidEADID(eadID) {
		return fmt.Errorf("invalid EADID: %s", eadID)
	}

	// assert that the SolrClient has been set
	err := assertSolrClientSet()
	if err != nil {
		return err
	}

	err = sc.Delete(eadID)
	if err != nil {
		return appendErrIssueRollbackJoinErrs(errs, err)
	}
	return nil
}

func IndexEADFile(eadPath string) error {

	var errs []error

	// assert that the SolrClient has been set
	err := assertSolrClientSet()
	if err != nil {
		return appendAndJoinErrs(errs, err)
	}

	// Check if the EAD file path is absolute
	if !filepath.IsAbs(eadPath) {
		return appendAndJoinErrs(errs, fmt.Errorf("EAD file path must be absolute: %s", eadPath))
	}

	// Get the EAD's repository code
	repositoryCode, err := util.GetRepositoryCode(eadPath)
	if err != nil {
		return appendAndJoinErrs(errs, err)
	}

	// Read the EAD file
	eadXML, err := os.ReadFile(eadPath)
	if err != nil {
		return appendAndJoinErrs(errs, err)
	}

	// Parse the EAD file
	EAD, err := ead.New(repositoryCode, string(eadXML))
	if err != nil {
		return appendAndJoinErrs(errs, err)
	}

	// Delete the data for this EAD from Solr
	err = sc.Delete(EAD.CollectionDoc.Parts.EADID.Values[0])
	if err != nil {
		return appendErrIssueRollbackJoinErrs(errs, err)
	}

	// Add the EAD Collection-level document to Solr
	xmlPostBody := EAD.CollectionDoc.SolrAddMessage.String()
	err = sc.Add(xmlPostBody)
	if err != nil {
		return appendErrIssueRollbackJoinErrs(errs, err)
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
		// NOTE: in this scenario, there isn't a new error,
		// but we still want to take advantage of the rollback functionality,
		// so we pass "nil" as the error
		return appendErrIssueRollbackJoinErrs(errs, nil)
	}

	// commit the documents to Solr
	err = sc.Commit()
	if err != nil {
		return appendErrIssueRollbackJoinErrs(errs, err)
	}

	return nil
}

func IndexGitCommit(repoPath, commit string) error {
	var errs []error

	// assert that the SolrClient has been set
	err := assertSolrClientSet()
	if err != nil {
		return appendAndJoinErrs(errs, err)
	}

	// get the list of EAD files and their operations
	indexerOperations, err := git.ListEADFilesForCommit(repoPath, commit)
	if err != nil {
		return appendAndJoinErrs(errs, err)
	}

	// iterate over the EAD files and their operations
	for eadPath, operation := range indexerOperations {
		switch operation {
		case git.Add:
			err = IndexEADFile(eadPath)
			if err != nil {
				return appendErrIssueRollbackJoinErrs(errs, err)
			}

		case git.Delete:
			eadID, err := git.EADPathToEADID(eadPath)
			if err != nil {
				return appendAndJoinErrs(errs, err)
			}

			err = DeleteEADFileDataFromIndex(eadID)
			if err != nil {
				return appendErrIssueRollbackJoinErrs(errs, err)
			}

		default:
			return appendAndJoinErrs(errs, fmt.Errorf("unknown operation: %s", operation))
		}
	}

	return nil
}

func SetSolrClient(solrClient solr.SolrClient) {
	sc = solrClient
}

func appendAndJoinErrs(errs []error, err error) error {
	errs = append(errs, err)
	return errors.Join(errs...)
}

func appendErrIssueRollbackJoinErrs(errs []error, err error) error {
	errs = append(errs, err)
	err = sc.Rollback()
	if err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func assertSolrClientSet() error {
	if sc == nil {
		return errors.New(errSolrClientNotSet)
	}

	if sc.GetSolrURLOrigin() == "" {
		return errors.New("the SolrClient URL origin is not set")
	}

	return nil
}
