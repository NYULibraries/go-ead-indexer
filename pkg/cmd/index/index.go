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
var loggingLevel string // logging level
var logger log.Logger   // logger

// This init() function contains a subset of the full 'index' command functionality
func init() {
	IndexCmd.Flags().StringVarP(&file, "file", "f",
		"", "path to EAD file")
	IndexCmd.Flags().StringVarP(&loggingLevel, "logging-level", "l",
		log.DefaultLevelStringOption,
		"Sets logging level: "+strings.Join(log.GetValidLevelOptionStrings(), ", ")+"")
}

var IndexCmd = &cobra.Command{
	Use:     "index",
	Short:   "Index EAD file",
	Example: `go-ead-indexer index --file=[path to EAD file] --logging-level="debug"`,
	Run:     runIndexCmd,
}

// runIndexCmd is the main function for the 'index' command
// It initializes the logger and Solr client, then indexes the EAD file
// It exits with a fatal error if any of these steps fail
// It logs a message when the EAD file is successfully indexed
func runIndexCmd(cmd *cobra.Command, args []string) {

	err := initLogger()
	if err != nil {
		logger.Fatal("ERROR: couldn't initialize logger", err)
	}

	if file == "" {
		logger.Fatal("ERROR: EAD file path not set")
	}

	err = initSolrClient()
	if err != nil {
		logger.Fatal("ERROR: couldn't initialize Solr client", err)
	}

	// index EAD file
	err = index.IndexEADFile(file)
	if err != nil {
		logger.Fatal("ERROR: couldn't index EAD file", err)
	}

	logger.Info(index.MessageKey, fmt.Sprintf("SUCCESS: indexed EAD file: %s", file))
}

// initLogger initializes the logger in the pkg/cmd/index package
func initLogger() error {
	logger = log.New()
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
