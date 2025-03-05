package index

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nyulibraries/go-ead-indexer/pkg/cmd/testutils"
)

func TestInitLoggerError(t *testing.T) {
	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT", "http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(IndexCmd, "file", "testdata/ead.xml")
	testutils.SetCmdFlag(IndexCmd, "logging-level", "INVALID-LOGGING-LEVEL")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd, IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, "ERROR: couldn't initialize logger: ERROR: couldn't set log level:")
}

func TestLoggerLevelArgument(t *testing.T) {
	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT", "http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(IndexCmd, "file", "testdata/ead.xml")
	testutils.SetCmdFlag(IndexCmd, "logging-level", "none")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd, IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, `Logging level set to \"none\"`)
}

func TestMissingFileArgument(t *testing.T) {
	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT", "http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	// clear the file flag
	testutils.SetCmdFlag(IndexCmd, "file", "")
	testutils.SetCmdFlag(IndexCmd, "logging-level", "debug")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd, IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, "ERROR: EAD file path not set")
}

func TestBadFileArgument(t *testing.T) {
	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT", "http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	dir, err := testutils.GetCallingFileDirPath()
	if err != nil {
		t.Errorf("error getting calling file directory: %v", err)
		t.FailNow()
	}

	// set the file flag to a non-existent file
	testutils.SetCmdFlag(IndexCmd, "file", filepath.Join(dir, "testdata", "fixtures", "edip", "no_file_here.xml"))
	testutils.SetCmdFlag(IndexCmd, "logging-level", "debug")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd, IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, "ERROR: EAD file does not exist: ")
}

func TestMissingSolrOriginEnvVariableError(t *testing.T) {
	// ensure that the environment variable is NOT set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT", "")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	dir, err := testutils.GetCallingFileDirPath()
	if err != nil {
		t.Errorf("error getting calling file directory: %v", err)
		t.FailNow()
	}

	// set the file flag to an existing file
	testutils.SetCmdFlag(IndexCmd, "file", filepath.Join(dir, "testdata", "fixtures", "edip", "bad_ead.xml"))
	testutils.SetCmdFlag(IndexCmd, "logging-level", "debug")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd, IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, "ERROR: couldn't initialize Solr client: 'SOLR_ORIGIN_WITH_PORT' environment variable not set")
}

func TestInitSolrClientError(t *testing.T) {
	// ensure that the environment variable is NOT set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT", "this is not a valid url")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	dir, err := testutils.GetCallingFileDirPath()
	if err != nil {
		t.Errorf("error getting calling file directory: %v", err)
		t.FailNow()
	}

	// set the file flag to an existing file
	testutils.SetCmdFlag(IndexCmd, "file", filepath.Join(dir, "testdata", "fixtures", "edip", "bad_ead.xml"))
	testutils.SetCmdFlag(IndexCmd, "logging-level", "debug")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd, IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, `ERROR: couldn't initialize Solr client: error creating Solr client: parse \"this is not a valid url\": invalid URI for request`)
}

func TestIndexingError(t *testing.T) {
	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT", "http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	dir, err := testutils.GetCallingFileDirPath()
	if err != nil {
		t.Errorf("error getting calling file directory: %v", err)
		t.FailNow()
	}

	// set the file flag to an existing file that is not a valid EAD file
	testutils.SetCmdFlag(IndexCmd, "file", filepath.Join(dir, "testdata", "fixtures", "edip", "bad_ead.xml"))
	testutils.SetCmdFlag(IndexCmd, "logging-level", "debug")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd, IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, "ERROR: couldn't index EAD file: No <ead> tag with the expected structure was found")
}
