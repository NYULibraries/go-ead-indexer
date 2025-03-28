package index

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/nyulibraries/go-ead-indexer/pkg/cmd/testutils"
	"github.com/nyulibraries/go-ead-indexer/pkg/log"
)

func TestDelete_Cancel(t *testing.T) {
	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT", "http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(DeleteCmd, "eadid", "fales_mss460")
	testutils.SetCmdFlag(DeleteCmd, "assume-yes", "false")

	gotStdOut, _, _ := testutils.WriteToStdinCaptureCmdStdoutStderrE(runDeleteCmd, DeleteCmd, []string{}, "n\n")

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, fmt.Sprintf("Deletion canceled for EADID: '%s'", eadID))
}

func TestDelete_CancelAfterBadInput(t *testing.T) {
	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT", "http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(DeleteCmd, "eadid", "fales_mss460")
	testutils.SetCmdFlag(DeleteCmd, "assume-yes", "false")

	gotStdOut, _, _ := testutils.WriteToStdinCaptureCmdStdoutStderrE(runDeleteCmd, DeleteCmd, []string{}, "WAFFLES\nn\n")

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, "Please enter 'y' or 'n':")
	testutils.CheckStringContains(t, gotStdOut, fmt.Sprintf("Deletion canceled for EADID: '%s'", eadID))
}

func TestDelete_Error(t *testing.T) {
	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT", "http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(DeleteCmd, "eadid", "This#Is^Not!A(Valid*EADID")
	testutils.SetCmdFlag(DeleteCmd, "logging-level", "debug")
	testutils.SetCmdFlag(DeleteCmd, "assume-yes", "true")

	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runDeleteCmd, DeleteCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, "ERROR: couldn't delete data for EADID: This#Is^Not!A(Valid*EADID error:")
}

func TestIndex_ArgumentValidation(t *testing.T) {
	var want string
	var got error
	var args []string

	fileFlag := "file"
	gitRepoFlag := "git-repo"
	gitCommitFlag := "commit"

	// set up arguments
	dir, err := testutils.GetCallingFileDirPath()
	if err != nil {
		t.Errorf("error getting calling file directory: %v", err)
		t.FailNow()
	}
	file := filepath.Join(dir, "testdata", "fixtures", "edip", "mos_2024.xml")
	gitRepoPath := filepath.Join(dir, "testdata", "fixtures", "edip")
	gitCommit := "a5ca6cca30fc08cfc13e4f1492dbfbbf3ec7cf63"

	scenarios := []struct {
		File        string
		GitRepoPath string
		GitCommit   string
		Want        string
	}{
		{"", "", "", eMsgNeedOneButNotBothFileAndGitRepo},                   // fail: neither the file nor the git-repo flag is set
		{"", "", gitCommit, eMsgNeedOneButNotBothFileAndGitRepo},            // fail: only the commit flag is set
		{"", gitRepoPath, "", eMsgMissingCommitOrGitRepo},                   // fail: only the git-repo flag is set
		{"", gitRepoPath, gitCommit, ""},                                    // pass: the git-repo and commit flags are set
		{file, "", "", ""},                                                  // pass: only the file flag is set
		{file, "", gitCommit, eMsgCommitOnlyWithGitRepo},                    // fail: both file and commit flags are set
		{file, gitRepoPath, "", eMsgNeedOneButNotBothFileAndGitRepo},        // fail: both file and git-repo flags are set
		{file, gitRepoPath, gitCommit, eMsgNeedOneButNotBothFileAndGitRepo}, // fail: all three flags are set
	}

	for _, scenario := range scenarios {
		// set the flags
		testutils.SetCmdFlag(IndexCmd, fileFlag, scenario.File)
		testutils.SetCmdFlag(IndexCmd, gitRepoFlag, scenario.GitRepoPath)
		testutils.SetCmdFlag(IndexCmd, gitCommitFlag, scenario.GitCommit)

		want = scenario.Want
		got = indexCheckArgs(IndexCmd, args)

		switch {
		case want == "" && got != nil:
			t.Errorf("expected no error but got: %v", got)
		case want != "" && got == nil:
			t.Errorf("expected an error but got nothing")
		case (want != "" && got != nil) && (got.Error() != want):
			t.Errorf("expected error message: '%s', but got '%s'", want, got.Error())
		}
	}
}

func TestIndexEAD_BadFileArgument(t *testing.T) {
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

func TestIndexEAD_Error(t *testing.T) {
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

func TestIndexEAD_InitLoggerError(t *testing.T) {
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

	testutils.CheckStringContains(t, gotStdOut, "ERROR: couldn't initialize logger: ERROR: unsupported logging level:")
}

func TestIndexEAD_InitSolrClientError(t *testing.T) {
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

func TestIndexEAD_LoggerLevelArgument(t *testing.T) {
	testLogLevel := "debug" // this value MUST BE DIFFERENT from the default logging level

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT", "http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(IndexCmd, "file", "testdata/ead.xml")
	testutils.SetCmdFlag(IndexCmd, "logging-level", testLogLevel)
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd, IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, fmt.Sprintf(`Logging level set to \"%s\"`, testLogLevel))
}

func TestIndexEAD_MissingSolrOriginEnvVariableError(t *testing.T) {
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

func TestIndexEAD_UnsetLoggingLevelArgument(t *testing.T) {
	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT", "http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(IndexCmd, "file", "testdata/ead.xml")
	testutils.SetCmdFlag(IndexCmd, "logging-level", "")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd, IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, `Logging level set to \"info\"`)
}

func TestLocalLogLevels(t *testing.T) {
	// this is a regression test to ensure that the local log levels are still valid
	// if this test fails, the local log levels need to be updated
	loggerAvailableLevels := log.GetValidLevelOptionStrings()
	for _, level := range localLogLevels {
		if !slices.Contains(loggerAvailableLevels, level) {
			t.Errorf("local log level '%s' is not a valid logger level", level)
		}
	}
}
