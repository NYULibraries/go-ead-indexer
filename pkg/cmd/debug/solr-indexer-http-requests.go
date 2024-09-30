package debug

import (
	"fmt"
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
		fmt.Println(fmt.Sprintf("`%s` called with args %v", cmd.Use, args))
	},
}
