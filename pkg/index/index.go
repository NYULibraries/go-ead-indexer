// Package index provides an interface to the EAD indexing process
//
// The SetSolrClient() function must be called before calling any of the
// indexing functions in this package.  This is because the Solr client
// is a package-level variable, and the default value is nil.
package index

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/nyulibraries/go-ead-indexer/pkg/ead"
	"github.com/nyulibraries/go-ead-indexer/pkg/ead/eadutil"
	"github.com/nyulibraries/go-ead-indexer/pkg/git"
	"github.com/nyulibraries/go-ead-indexer/pkg/log"
	"github.com/nyulibraries/go-ead-indexer/pkg/net/solr"
	"github.com/nyulibraries/go-ead-indexer/pkg/util"
)

const MessageKey = "index"
const errLoggerIsNil = "`logger` == nil"
const errSolrClientNotSet = "you must call `SetSolrClient()` before calling any indexing functions"

var sc = solr.SolrClient(nil)
var logger log.Logger
var startTime, endTime time.Time

func DeleteEADFileDataFromIndex(eadID string) error {
	logString := fmt.Sprintf("DeleteEADFileDataFromIndex(%s)", eadID)
	logDebug(logString)

	logStartTime(logString)
	defer logEndTime(logString)

	var errs []error

	// assert that the EADID is valid
	logDebug(fmt.Sprintf("eadutil.IsValidEADID(%s)", eadID))
	if !eadutil.IsValidEADID(eadID) {
		return fmt.Errorf("invalid EADID: %s", eadID)
	}

	// assert that the SolrClient has been set
	logDebug("assertSolrClientSet()")
	err := assertSolrClientSet()
	if err != nil {
		return err
	}

	logDebug(fmt.Sprintf("sc.Delete(%s)", eadID))
	err = sc.Delete(eadID)
	if err != nil {
		return appendErrIssueRollbackJoinErrs(errs, err)
	}

	// commit the change to Solr
	logDebug("sc.Commit()")
	err = sc.Commit()
	if err != nil {
		return appendErrIssueRollbackJoinErrs(errs, err)
	}

	return nil
}

func IndexEADFile(eadPath string) error {
	logString := fmt.Sprintf("IndexEADFile(%s)", eadPath)
	logDebug(logString)

	logStartTime(logString)
	defer logEndTime(logString)

	var errs []error

	// assert that the SolrClient has been set
	logDebug("assertSolrClientSet()")
	err := assertSolrClientSet()
	if err != nil {
		return appendAndJoinErrs(errs, err)
	}

	// Check if the EAD file path is absolute
	logDebug(fmt.Sprintf("filepath.IsAbs(%s)", eadPath))
	if !filepath.IsAbs(eadPath) {
		return appendAndJoinErrs(errs, fmt.Errorf("EAD file path must be absolute: %s", eadPath))
	}

	// Get the EAD's repository code
	logDebug(fmt.Sprintf("util.GetRepositoryCode(%s)", eadPath))
	repositoryCode, err := util.GetRepositoryCode(eadPath)
	if err != nil {
		return appendAndJoinErrs(errs, err)
	}

	// Read the EAD file
	logDebug(fmt.Sprintf("os.ReadFile(%s)", eadPath))
	eadXML, err := os.ReadFile(eadPath)
	if err != nil {
		return appendAndJoinErrs(errs, err)
	}

	// Parse the EAD file
	//logDebug(fmt.Sprintf("ead.New(%s, (XML for %s))", repositoryCode, eadPath))
	logDebug(fmt.Sprintf("ead.New(%s, %s)", repositoryCode, eadXML))
	EAD, err := ead.New(repositoryCode, string(eadXML))
	if err != nil {
		return appendAndJoinErrs(errs, err)
	}

	// Delete the data for this EAD from Solr
	logDebug(fmt.Sprintf("sc.Delete(%s)", EAD.CollectionDoc.Parts.EADID.Values[0]))
	err = sc.Delete(EAD.CollectionDoc.Parts.EADID.Values[0])
	if err != nil {
		return appendErrIssueRollbackJoinErrs(errs, err)
	}

	// Add the EAD Collection-level document to Solr
	xmlPostBody := EAD.CollectionDoc.SolrAddMessage.String()
	logDebug(fmt.Sprintf("collection-level: sc.Add(%s)", xmlPostBody))

	err = sc.Add(xmlPostBody)
	if err != nil {
		return appendErrIssueRollbackJoinErrs(errs, err)
	}

	// Add the EAD Component-level documents to Solr
	if EAD.Components != nil {
		for _, component := range *EAD.Components {
			xmlPostBody = component.SolrAddMessage.String()
			logDebug(fmt.Sprintf("component-level: sc.Add(%s)", xmlPostBody))

			err = sc.Add(xmlPostBody)
			if err != nil {
				logDebug("error: " + err.Error())
				errs = append(errs, err)
			}
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
	logDebug("sc.Commit()")
	err = sc.Commit()
	if err != nil {
		return appendErrIssueRollbackJoinErrs(errs, err)
	}

	return nil
}

func IndexGitCommit(repoPath, commit string) (int, error) {
	numIndexOperations := 0

	logString := fmt.Sprintf("IndexGitCommit(%s, %s)", repoPath, commit)
	logDebug(logString)

	// assert that the SolrClient has been set
	logDebug("assertSolrClientSet()")
	err := assertSolrClientSet()
	if err != nil {
		return numIndexOperations, err
	}

	// checkout the git commit
	logDebug(fmt.Sprintf("git.CheckoutMergeReset(%s, %s)", repoPath, commit))
	err = git.CheckoutMergeReset(repoPath, commit)
	if err != nil {
		return numIndexOperations, err
	}

	// get the list of EAD files and their operations
	logDebug(fmt.Sprintf("git.ListEADFilesForCommit(%s, %s)", repoPath, commit))
	operations, err := git.ListEADFilesForCommit(repoPath, commit)
	if err != nil {
		return numIndexOperations, err
	}

	numIndexOperations = len(operations)

	for _, eadFileRelativePath := range slices.Sorted(maps.Keys(operations)) {
		operation := operations[eadFileRelativePath]

		switch operation {
		case git.Add:
			err = IndexEADFile(filepath.Join(repoPath, eadFileRelativePath))
			if err != nil {
				return numIndexOperations, err
			}

		case git.Delete:
			eadID, err := eadutil.EADPathToEADID(eadFileRelativePath)
			if err != nil {
				return numIndexOperations, err
			}

			err = DeleteEADFileDataFromIndex(eadID)
			if err != nil {
				return numIndexOperations, err
			}

		default:
			return numIndexOperations, fmt.Errorf("unknown operation: %s", operation)
		}
	}

	return numIndexOperations, nil
}

func InitLogger(l log.Logger) error {
	logger = l
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

func logStartTime(s string) {
	startTime = time.Now()
	logInfo(fmt.Sprintf("%s started at %s", s, startTime))
}

func logDebug(s string) {
	if logger == nil {
		_, _ = fmt.Fprintln(os.Stderr, "logDebug() error: "+errLoggerIsNil)

		return
	}
	logger.Debug(MessageKey, s)
}

func logEndTime(s string) {
	endTime = time.Now()
	logInfo(fmt.Sprintf("%s ended at %s", s, endTime))
	logInfo(fmt.Sprintf("%s duration: %s", s, endTime.Sub(startTime)))
}

func logInfo(s string) {
	if logger == nil {
		_, _ = fmt.Fprintln(os.Stderr, "logInfo() error: "+errLoggerIsNil)

		return
	}
	logger.Info(MessageKey, s)
}
