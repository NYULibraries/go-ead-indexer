package debug

import (
	"github.com/nyulibraries/go-ead-indexer/pkg/log"
	"github.com/spf13/cobra"
)

var DebugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debugging utilities",
}

var logger log.Logger // logger

func init() {
	logger = log.New()

	logger.SetLevel(log.LevelError)
}
