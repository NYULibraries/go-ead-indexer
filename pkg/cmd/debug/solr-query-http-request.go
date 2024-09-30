package debug

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	solrCmd.AddCommand(queryHTTPRequestCmd)

	queryHTTPRequestCmd.Flags().StringVarP(&eadid, "ead-id", "e",
		"", "EAD file id")
}

var queryHTTPRequestCmd = &cobra.Command{
	Use:     "query-http-request",
	Short:   "Dump the HTTP request used by the `debug solr query` command",
	Example: "go-ead-indexer debug solr query-http-request --ead-id=[EADID]",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("`%s` called with args %v", cmd.Use, args))
	},
}
