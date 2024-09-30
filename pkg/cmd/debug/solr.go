package debug

import (
	"github.com/spf13/cobra"
)

var eadid string
var file string
var gitCommit string
var gitRepoPath string

func init() {
	DebugCmd.AddCommand(solrCmd)
}

var solrCmd = &cobra.Command{
	Use:   "solr",
	Short: "Utilities for debugging Solr requests and inspecting the index",
}
