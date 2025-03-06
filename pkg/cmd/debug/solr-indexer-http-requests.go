package debug

import (
	"fmt"
	"github.com/nyulibraries/go-ead-indexer/pkg/debug"
	"github.com/spf13/cobra"
)

func init() {
	solrCmd.AddCommand(indexerHTTPRequestsCmd)

	indexerHTTPRequestsCmd.Flags().StringVarP(&gitCommit, "commit", "c",
		"", "hash of git commit")
	indexerHTTPRequestsCmd.Flags().StringVarP(&file, "file", "f",
		"", "path to EAD file")
	indexerHTTPRequestsCmd.Flags().StringVarP(&gitRepoPath, "git-repo", "g",
		"", "path to EAD files git repo")
}

var indexerHTTPRequestsCmd = &cobra.Command{
	Use:   "indexer-http-requests",
	Short: "Dump the HTTP POST requests to Solr used by the `index` command",
	Example: `go-ead-indexer debug solr indexer-http-requests --file=[path to EAD file]
go-ead-indexer debug solr indexer-http-requests --git-repo=[path] --commit=[hash]`,
	Run: func(cmd *cobra.Command, args []string) {
		if file != "" {
			dumpedSolrIndexerHTTPRequests, err := debug.DumpSolrIndexerHTTPRequestsForEADFile(file)
			if err != nil {
				logger.Error(fmt.Sprintf(`debug.DumpSolrIndexerHTTPRequestsForEADFile("%s")`+
					` failed with error: %s`, file, err.Error()))
			}

			fmt.Println(dumpedSolrIndexerHTTPRequests)
		} else if gitCommit != "" {
			if gitRepoPath != "" {
				dumpedSolrIndexerHTTPRequests, err :=
					debug.DumpSolrIndexerHTTPRequestsForGitCommit(gitRepoPath, gitCommit)
				if err != nil {
					logger.Error(fmt.Sprintf(`debug.DumpSolrIndexerHTTPRequestsForGitCommit("%s")`+
						` against git repo "%s" failed with error: %s`,
						gitCommit, gitRepoPath, err.Error()))
				}

				fmt.Println(dumpedSolrIndexerHTTPRequests)
			} else {
				logger.Error("Must specify --git-repo with --commit")
			}
		} else if gitRepoPath != "" {
			logger.Error("Must specify --commit with --git-repo")
		} else {
			logger.Error(fmt.Sprintf("`%s` must be called with either --file or"+
				" --commit & --git-repo", cmd.Use))
		}
	},
}
