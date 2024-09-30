package debug

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	solrCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVarP(&file, "file", "f",
		"", "path to EAD file")
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify that an EAD file or git commit has been correctly indexed",
	Example: `go-ead-indexer debug solr verify --file=[path to EAD file]
go-ead-indexer debug solr verify --git-repo=[path] --commit=[hash]`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("`%s` called with args %v", cmd.Use, args))
	},
}
