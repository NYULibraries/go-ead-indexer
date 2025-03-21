package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nyulibraries/go-ead-indexer/pkg/cmd/debug"
	"github.com/nyulibraries/go-ead-indexer/pkg/cmd/index"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "go-ead-indexer",
	Long: "`go-ead-indexer`" + ` is the EAD file Solr indexer for Special Collections.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(debug.DebugCmd)
	rootCmd.AddCommand(index.IndexCmd)
	rootCmd.AddCommand(index.DeleteCmd)
}
