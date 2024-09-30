package debug

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	solrCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringVarP(&eadid, "ead-id", "e",
		"", "EAD file id")
}

var queryCmd = &cobra.Command{
	Use:     "query",
	Short:   "Query the Solr index",
	Example: "go-ead-indexer debug solr query --ead-id=[EADID]",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("`%s` called with args %v", cmd.Use, args))
	},
}
