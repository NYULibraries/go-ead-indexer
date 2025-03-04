package index

import (
	"os"
	"testing"

	"github.com/nyulibraries/go-ead-indexer/pkg/cmd/testutils"
)

func TestInitLoggerError(t *testing.T) {
	// ensure that the environment variable is NOT set
	os.Setenv("SOLR_ORIGIN_WITH_PORT", "")

	testutils.SetCmdFlag(IndexCmd, "file", "testdata/ead.xml")
	testutils.SetCmdFlag(IndexCmd, "logging-level", "debug")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd, IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, "ERROR: couldn't initialize Solr client: 'SOLR_ORIGIN_WITH_PORT' environment variable not set")
}
