package index

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
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
const eMsgNeedOneButNotBothFileAndGitRepo = "one, but not both, of --file or --git-repo arguments must be specified"
const eMsgCommitOnlyWithGitRepo = "the --commit argument can only be used with the --git-repo argument"
const eMsgMissingCommitOrGitRepo = "missing argument: the --git-repo argument must be used with the --commit argument"

// log levels used by this package, in increasing order of severity
var localLogLevels = []string{"debug", "info", "error"}

var file string         // EAD file to be indexed
var gitCommit string    // commit to index
var gitRepoPath string  // path to EAD files git repo
var eadID string        // EADID value of EAD data to delete
var assumeYes bool      // flag to disable interactive mode
var loggingLevel string // logging level
var logger log.Logger   // logger

// This init() function contains a subset of the full 'index' command functionality
func init() {
	IndexCmd.Flags().StringVarP(&gitCommit, "commit", "c",
		"", "hash of git commit")
	IndexCmd.Flags().StringVarP(&file, "file", "f", "",
		"path to EAD file")
	IndexCmd.Flags().StringVarP(&gitRepoPath, "git-repo", "g", "",
		"path to EAD files git repo")
	IndexCmd.Flags().StringVarP(&loggingLevel, "logging-level", "l",
		log.DefaultLevelStringOption,
		"Sets logging level: "+strings.Join(localLogLevels, ", ")+"")

	DeleteCmd.Flags().StringVarP(&eadID, "eadid", "e", "",
		"EADID value of EAD data to delete")
	DeleteCmd.Flags().BoolVarP(&assumeYes, "assume-yes", "y", false,
		"disable interactive mode")
	DeleteCmd.Flags().StringVarP(&loggingLevel, "logging-level", "l",
		log.DefaultLevelStringOption,
		"Sets logging level: "+strings.Join(localLogLevels, ", ")+"")
}

var DeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete data by EADID",
	Long:    "Delete data from the index using the EADID",
	Example: `go-ead-indexer delete --eadid=[EADID] --logging-level="debug --assume-yes"`,
	RunE:    runDeleteCmd,
}

var IndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index EAD file or commit",
	Example: `go-ead-indexer index --file=[path to EAD file] --logging-level="debug"
	go-ead-indexer index --git=[path] --commit=[hash] --logging-level="error"`,
	Args: indexCheckArgs,
	RunE: runIndexCmd,
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

	// request confirmation if interactive mode is enabled
	if !assumeYes {
		confirmed := confirmDelete(eadID)
		if !confirmed {
			msg := fmt.Sprintf("Deletion canceled for EADID: '%s'", eadID)
			logger.Info(index.MessageKey, msg)
			fmt.Println(msg)
			return nil
		}
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
	if _, err := os.Stat(file); errors.Is(err, fs.ErrNotExist) {
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
	if !slices.Contains(localLogLevels, normalizedLogLevel) {
		return fmt.Errorf("ERROR: unsupported logging level: '%s'. Supported levels are: %s", normalizedLogLevel, strings.Join(localLogLevels, ", "))
	}

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

func confirmDelete(eadID string) bool {
	var response string

	fmt.Printf("Are you sure you want to delete data for EADID: %s? (y/n): ", eadID)
	fmt.Scanln(&response)

	lowercaseResponse := strings.ToLower(response)
	for lowercaseResponse != "y" && lowercaseResponse != "n" {
		fmt.Printf("Please enter 'y' or 'n': ")
		fmt.Scanln(&response)
		lowercaseResponse = strings.ToLower(response)
	}

	return lowercaseResponse == "y"
}

func indexCheckArgs(cmd *cobra.Command, args []string) error {
	if (file == "" && gitRepoPath == "") ||
		(file != "" && gitRepoPath != "") {
		return fmt.Errorf("%s", eMsgNeedOneButNotBothFileAndGitRepo)
	}

	if file != "" && gitCommit != "" {
		return fmt.Errorf("%s", eMsgCommitOnlyWithGitRepo)
	}

	if (gitRepoPath != "" && gitCommit == "") ||
		(gitRepoPath == "" && gitCommit != "") {
		return fmt.Errorf("%s", eMsgMissingCommitOrGitRepo)
	}

	// arguments are OK so disable Cobra's usage output on error
	cmd.SilenceUsage = true

	return nil
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
