package index

import (
	"fmt"
	"github.com/nyulibraries/go-ead-indexer/pkg/index"
	"github.com/spf13/cobra"
	"strings"

	"github.com/nyulibraries/go-ead-indexer/pkg/log"
)

var file string
var gitCommit string
var gitRepoPath string
var loggingLevel string

func init() {
	IndexCmd.Flags().StringVarP(&gitCommit, "commit", "c",
		"", "hash of git commit")
	IndexCmd.Flags().StringVarP(&file, "file", "f",
		"", "path to EAD file")
	IndexCmd.Flags().StringVarP(&gitRepoPath, "git-repo", "g",
		"", "path to EAD files git repo")
	IndexCmd.Flags().StringVarP(&loggingLevel, "logging-level", "l",
		log.DefaultLevelStringOption,
		"Sets logging level: "+strings.Join(log.GetValidLevelOptionStrings(), ", ")+"")
}

var IndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index EAD file",
	Example: `go-ead-indexer index --file=[path to EAD file]
go-ead-indexer index --git-repo=[path] --commit=[hash]
go-ead-indexer index --git-repo=[path] --commit=[hash] --logging-level="debug"`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("`%s` called with args %v", cmd.Use, args))

		logger := log.New()

		normalizedLogLevel := strings.ToLower(loggingLevel)
		err := logger.SetLevelByString(normalizedLogLevel)
		if err != nil {
			logger.Fatal("ERROR: couldn't set log level", err)
		}

		logger.Info(index.MessageKey, fmt.Sprintf("Logging level set to \"%s\"", normalizedLogLevel))
	},
}
