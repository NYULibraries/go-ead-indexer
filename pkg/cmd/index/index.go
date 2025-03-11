package index

import (
	"fmt"
	"os"
	"strings"

	"github.com/nyulibraries/go-ead-indexer/pkg/index"
	"github.com/spf13/cobra"

	"github.com/nyulibraries/go-ead-indexer/pkg/log"
	"github.com/nyulibraries/go-ead-indexer/pkg/net/solr"
)

//------------------------------------------------------------------------------
// THIS PACKAGE IS A WORK IN PROGRESS!!!
//------------------------------------------------------------------------------

// environment variable that holds the Solr origin with port information
const originEnvVar = "SOLR_ORIGIN_WITH_PORT"

var file string         // EAD file to be indexed
var eadID string        // EADID value of EAD data to delete
var loggingLevel string // logging level
var logger log.Logger   // logger

// This init() function contains a subset of the full 'index' command functionality
func init() {
	IndexCmd.Flags().StringVarP(&file, "file", "f", "", "path to EAD file")
	IndexCmd.Flags().StringVarP(&loggingLevel, "logging-level", "l",
		log.DefaultLevelStringOption,
		"Sets logging level: "+strings.Join(log.GetValidLevelOptionStrings(), ", ")+"")

	DeleteCmd.Flags().StringVarP(&eadID, "eadid", "e", "", "EADID value of EAD data to delete")
	DeleteCmd.Flags().StringVarP(&loggingLevel, "logging-level", "l",
		log.DefaultLevelStringOption,
		"Sets logging level: "+strings.Join(log.GetValidLevelOptionStrings(), ", ")+"")
}

var DeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete data by EADID",
	Long:    "Delete data from the index using the EADID",
	Example: `go-ead-indexer delete --eadid=[EADID] --logging-level="debug"`,
	RunE:    runDeleteCmd,
}

var IndexCmd = &cobra.Command{
	Use:     "index",
	Short:   "Index EAD file",
	Example: `go-ead-indexer index --file=[path to EAD file] --logging-level="debug"`,
	RunE:    runIndexCmd,
}

// runDeleteCmd is the main function for the 'delete' verb
// It initializes the logger and Solr client, then deletes the data by EADID
// It exits with a fatal error if any of these steps fail
// It logs a message when the EAD data is successfully deleted
func runDeleteCmd(cmd *cobra.Command, args []string) error {

	// initialize logger
	err := initLogger()
	if err != nil {
		emsg := fmt.Sprintf("ERROR: couldn't initialize logger: %s", err)
		return logAndReturnError(emsg)
	}

	// check if EAD file path is set
	if eadID == "" {
		emsg := "ERROR: EADID is not set"
		return logAndReturnError(emsg)
	}

	// initialize Solr client
	err = initSolrClient()
	if err != nil {
		emsg := fmt.Sprintf("ERROR: couldn't initialize Solr client: %s", err)
		return logAndReturnError(emsg)
	}

	// delete data associated with EADID
	err = index.DeleteEADFileDataFromIndex(eadID)
	if err != nil {
		emsg := fmt.Sprintf("ERROR: couldn't delete data for EADID: %s error: %s", eadID, err)
		return logAndReturnError(emsg)
	}

	// log success message
	logger.Info(index.MessageKey, fmt.Sprintf("SUCCESS: deleted data for EADID: %s", eadID))
	return nil
}

// runIndexCmd is the main function for the 'index' command
// It initializes the logger and Solr client, then indexes the EAD file
// It exits with a fatal error if any of these steps fail
// It logs a message when the EAD file is successfully indexed
func runIndexCmd(cmd *cobra.Command, args []string) error {

	// initialize logger
	err := initLogger()
	if err != nil {
		emsg := fmt.Sprintf("ERROR: couldn't initialize logger: %s", err)
		return logAndReturnError(emsg)
	}

	// check if EAD file path is set
	if file == "" {
		emsg := "ERROR: EAD file path not set"
		return logAndReturnError(emsg)
	}

	// check that the EAD file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		emsg := fmt.Sprintf("ERROR: EAD file does not exist: %s", file)
		return logAndReturnError(emsg)
	}

	// initialize Solr client
	err = initSolrClient()
	if err != nil {
		emsg := fmt.Sprintf("ERROR: couldn't initialize Solr client: %s", err)
		return logAndReturnError(emsg)
	}

	// index EAD file
	err = index.IndexEADFile(file)
	if err != nil {
		emsg := fmt.Sprintf("ERROR: couldn't index EAD file: %s", err)
		return logAndReturnError(emsg)
	}

	// log success message
	logger.Info(index.MessageKey, fmt.Sprintf("SUCCESS: indexed EAD file: %s", file))
	return nil
}

// initLogger initializes the logger in the pkg/cmd/index package
func initLogger() error {

	logger = log.New()

	// set logging level if it was not specified on the command line
	if loggingLevel == "" {
		loggingLevel = log.DefaultLevelStringOption
	}

	normalizedLogLevel := strings.ToLower(loggingLevel)
	err := logger.SetLevelByString(normalizedLogLevel)
	if err != nil {
		return fmt.Errorf("ERROR: couldn't set log level: %s", err)
	}

	logger.Info(index.MessageKey, fmt.Sprintf("Logging level set to \"%s\"", normalizedLogLevel))
	return nil
}

// initSolrClient initializes the Solr client in the pkg/index package
func initSolrClient() error {
	solrOrigin := os.Getenv(originEnvVar)
	if solrOrigin == "" {
		return fmt.Errorf("'%s' environment variable not set", originEnvVar)
	}

	sc, err := solr.NewSolrClient(solrOrigin)
	if err != nil {
		return fmt.Errorf("error creating Solr client: %s", err)
	}

	index.SetSolrClient(sc)

	return nil
}

func logAndReturnError(emsg string) error {
	logger.Error(index.MessageKey, emsg)
	return fmt.Errorf("%s", emsg)
}

//------------------------------------------------------------------------------
// THIS PACKAGE IS A WORK IN PROGRESS!!!
//
// THE FOLLOWING CODE HAS BEEN COMMENTED OUT,
// BUT REPRESENTS THE FULL FUNCTIONALITY OF THE 'index' COMMAND
//------------------------------------------------------------------------------
// git commit hash and path to EAD files git repo
// var gitCommit string
// var gitRepoPath string
//
// the following init() function contains the full implementation of the 'index' command
// func init() {
// 	IndexCmd.Flags().StringVarP(&gitCommit, "commit", "c",
// 		"", "hash of git commit")
// 	IndexCmd.Flags().StringVarP(&file, "file", "f",
// 		"", "path to EAD file")
// 	IndexCmd.Flags().StringVarP(&gitRepoPath, "git-repo", "g",
// 		"", "path to EAD files git repo")
// 	IndexCmd.Flags().StringVarP(&loggingLevel, "logging-level", "l",
// 		log.DefaultLevelStringOption,
// 		"Sets logging level: "+strings.Join(log.GetValidLevelOptionStrings(), ", ")+"")
// }
// var IndexCmd = &cobra.Command{
// 	Use:     "index",
// 	Short:   "Index EAD file",
// 	Example: `go-ead-indexer index --file=[path to EAD file] --logging-level="debug"`,
// 	go-ead-indexer index --git-repo=[path] --commit=[hash]
// 	go-ead-indexer index --git-repo=[path] --commit=[hash] --logging-level="debug"`,
// 	Run: runIndexCmd,
// }
