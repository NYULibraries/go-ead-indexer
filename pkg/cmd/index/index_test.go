package index

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/nyulibraries/go-ead-indexer/pkg/cmd/testutils"
	indextestutils "github.com/nyulibraries/go-ead-indexer/pkg/index/testutils"
	"github.com/nyulibraries/go-ead-indexer/pkg/log"
)

// test git repo paths
var thisPath string
var gitSourceRepoPathAbsolute string
var gitRepoTestGitRepoPathAbsolute string
var gitRepoTestGitRepoDotGitDirectory string
var gitRepoTestGitRepoHiddenGitDirectory string

// Copied this from `pkg/index` and `pkg/git` package tests in order to use the git
// repo fixture for `pkg/index` tests.  In fact, there is only one test that
// needs to use that repo: `TestIndexGitCommit_NoEADFilesInCommit`.
func init() {
	// The `filename` string is the absolute path to this source file, which should
	// be located at the root of the package directory.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("ERROR: `runtime.Caller(0)` failed")
	}

	// Get the path to the parent directory of this file.  Again, this is assuming
	// that this `init()` function is defined in a package top level file -- or
	// more precisely, that this file is in the same directory at the `testdata/`
	// directory that is referenced in the relative paths used in the functions
	// defined in this file.
	thisPath = filepath.Dir(filename)

	// Get testdata directory paths
	gitSourceRepoPathAbsolute = filepath.Join(thisPath, "testdata", "fixtures", "git-repo")

	// This could be done as a const at top level, but assigning it here to keep
	// all this path stuff in one place.
	gitRepoTestGitRepoPathAbsolute = filepath.Join(thisPath, "testdata", "fixtures", "test-git-repo")
	gitRepoTestGitRepoDotGitDirectory = filepath.Join(gitRepoTestGitRepoPathAbsolute, "dot-git")
	gitRepoTestGitRepoHiddenGitDirectory = filepath.Join(gitRepoTestGitRepoPathAbsolute, ".git")
}

func TestDelete_Cancel(t *testing.T) {
	resetDeleteArgs()

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT",
		"http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(DeleteCmd, "eadid", "fales_mss460")
	testutils.SetCmdFlag(DeleteCmd, "assume-yes", "false")

	gotStdOut, _, _ := testutils.WriteToStdinCaptureCmdStdoutStderrE(
		runDeleteCmd, DeleteCmd, []string{}, "n\n")

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut,
		fmt.Sprintf("Deletion canceled for EADID: '%s'", eadID))
}

func TestDelete_CancelAfterBadInput(t *testing.T) {
	resetDeleteArgs()

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT",
		"http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(DeleteCmd, "eadid", "fales_mss460")
	testutils.SetCmdFlag(DeleteCmd, "assume-yes", "false")

	gotStdOut, _, _ := testutils.WriteToStdinCaptureCmdStdoutStderrE(
		runDeleteCmd, DeleteCmd, []string{}, "WAFFLES\nn\n")

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, "Please enter 'y' or 'n':")
	testutils.CheckStringContains(t, gotStdOut,
		fmt.Sprintf("Deletion canceled for EADID: '%s'", eadID))
}

func TestDelete_Error(t *testing.T) {
	resetDeleteArgs()

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT",
		"http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(DeleteCmd, "eadid", "This#Is^Not!A(Valid*EADID")
	testutils.SetCmdFlag(DeleteCmd, "logging-level", "debug")
	testutils.SetCmdFlag(DeleteCmd, "assume-yes", "true")

	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runDeleteCmd,
		DeleteCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut,
		`couldn't delete data for EADID: This#Is^Not!A(Valid*EADID`)
}

func TestDelete_InitLoggerError(t *testing.T) {
	resetDeleteArgs()

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT",
		"http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(DeleteCmd, "eadid", "fales_mss460")
	testutils.SetCmdFlag(IndexCmd, "logging-level", "INVALID-LOGGING-LEVEL")
	testutils.SetCmdFlag(DeleteCmd, "assume-yes", "true")

	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runDeleteCmd,
		DeleteCmd, []string{})
	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut,
		"couldn't initialize logger: unsupported logging level:")
}

func TestDelete_InitSolrClientError(t *testing.T) {
	resetDeleteArgs()

	// ensure that the environment variable is NOT set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT", "this is not a valid url")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(DeleteCmd, "eadid", "fales_mss460")
	testutils.SetCmdFlag(DeleteCmd, "logging-level", "debug")
	testutils.SetCmdFlag(DeleteCmd, "assume-yes", "true")

	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runDeleteCmd,
		DeleteCmd, []string{})
	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut,
		`couldn't initialize Solr client: error creating Solr client: `+
			`parse \"this is not a valid url\": invalid URI for request`)
}

func TestDelete_MissingEADID(t *testing.T) {
	resetDeleteArgs()

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT",
		"http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(DeleteCmd, "eadid", "")
	testutils.SetCmdFlag(DeleteCmd, "assume-yes", "true")

	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runDeleteCmd,
		DeleteCmd, []string{})
	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, eMsgEADIDNotSet)

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
			t.Errorf("expected error message: '%s', but got '%s'", want,
				got.Error())
		}
	}
}

func TestIndex_CannotDetermineIndexingCase(t *testing.T) {
	resetIndexArgs()

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT",
		"http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(IndexCmd, "file", "")
	testutils.SetCmdFlag(IndexCmd, "git-repo", "")
	testutils.SetCmdFlag(IndexCmd, "commit", "")
	testutils.SetCmdFlag(IndexCmd, "logging-level", "debug")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd,
		IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut,
		eMsgCouldNotDetermineIndexingCase)
}

func TestIndexEAD_BadFileArgument(t *testing.T) {
	resetIndexArgs()

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT",
		"http://www.example.com:8983/solr")
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
	testutils.SetCmdFlag(IndexCmd, "file",
		filepath.Join(dir, "testdata", "fixtures", "edip", "no_file_here.xml"))
	testutils.SetCmdFlag(IndexCmd, "logging-level", "debug")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd,
		IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut,
		"EAD file does not exist: ")
}

func TestIndexEAD_Error(t *testing.T) {
	resetIndexArgs()

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT",
		"http://www.example.com:8983/solr")
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
	testutils.SetCmdFlag(IndexCmd, "file",
		filepath.Join(dir, "testdata", "fixtures", "edip", "bad_ead.xml"))
	testutils.SetCmdFlag(IndexCmd, "logging-level", "debug")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd,
		IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut,
		`couldn't index EAD file: No <ead> tag `+
			`with the expected structure was found`)
}

func TestIndexEAD_InitLoggerError(t *testing.T) {
	resetIndexArgs()

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT",
		"http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(IndexCmd, "file", "testdata/ead.xml")
	testutils.SetCmdFlag(IndexCmd, "logging-level", "INVALID-LOGGING-LEVEL")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd,
		IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut,
		"couldn't initialize logger: unsupported logging level:")
}

func TestIndexEAD_InitSolrClientError(t *testing.T) {
	resetIndexArgs()

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
	testutils.SetCmdFlag(IndexCmd, "file",
		filepath.Join(dir, "testdata", "fixtures", "edip", "bad_ead.xml"))
	testutils.SetCmdFlag(IndexCmd, "logging-level", "debug")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd,
		IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut,
		`couldn't initialize Solr client: error creating Solr client:`+
			` parse \"this is not a valid url\": invalid URI for request`)
}

func TestIndexEAD_LoggerLevelArgument(t *testing.T) {
	resetIndexArgs()

	// this value MUST BE DIFFERENT from the default logging level
	testLogLevel := "debug"
	if testLogLevel == localDefaultLogLevel {
		t.Errorf("test logging level is the same as the default logging level")
		t.FailNow()
	}

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT",
		"http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(IndexCmd, "file", "testdata/ead.xml")
	testutils.SetCmdFlag(IndexCmd, "logging-level", testLogLevel)
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd,
		IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut,
		fmt.Sprintf(`Logging level set to \"%s\"`, testLogLevel))
}

func TestIndexEAD_MissingSolrOriginEnvVariableError(t *testing.T) {
	resetIndexArgs()

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
	testutils.SetCmdFlag(IndexCmd, "file",
		filepath.Join(dir, "testdata", "fixtures", "edip", "bad_ead.xml"))
	testutils.SetCmdFlag(IndexCmd, "logging-level", "debug")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd,
		IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut,
		`couldn't initialize Solr client: `+
			`'SOLR_ORIGIN_WITH_PORT' environment variable not set`)
}

func TestIndexEAD_UnsetLoggingLevelArgument(t *testing.T) {
	resetIndexArgs()

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT",
		"http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	testutils.SetCmdFlag(IndexCmd, "file", "testdata/ead.xml")
	testutils.SetCmdFlag(IndexCmd, "logging-level", "")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd,
		IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut,
		`Logging level set to \"info\"`)
}

func TestIndexGitCommit_BadGitRepoArgument(t *testing.T) {
	resetIndexArgs()

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT",
		"http://www.example.com:8983/solr")
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
	//testutils.SetCmdFlag(IndexCmd, "file", "")
	testutils.SetCmdFlag(IndexCmd, "git-repo",
		filepath.Join(dir, "testdata", "fixtures", "this-repo-does-not-exist"))
	testutils.SetCmdFlag(IndexCmd, "commit",
		"a5ca6cca30fc08cfc13e4f1492dbfbbf3ec7cf63")
	testutils.SetCmdFlag(IndexCmd, "logging-level", "debug")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd,
		IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, "repository does not exist")
}

func TestIndexGitCommit_NoEADFilesInCommit(t *testing.T) {
	resetIndexArgs()

	// ensure that the environment variable is set
	err := os.Setenv("SOLR_ORIGIN_WITH_PORT",
		"http://www.example.com:8983/solr")
	if err != nil {
		t.Errorf("error setting environment variable: %v", err)
		t.FailNow()
	}

	// cleanup any leftovers from interrupted tests
	deleteTestGitRepo(t)
	createTestGitRepo(t)
	defer deleteTestGitRepo(t)

	testutils.SetCmdFlag(IndexCmd, "git-repo", gitRepoTestGitRepoPathAbsolute)
	testutils.SetCmdFlag(IndexCmd, "commit", indextestutils.NoEADFilesInCommitHash)
	testutils.SetCmdFlag(IndexCmd, "logging-level", "info")
	gotStdOut, _, _ := testutils.CaptureCmdStdoutStderrE(runIndexCmd,
		IndexCmd, []string{})

	if gotStdOut == "" {
		t.Errorf("expected data on StdOut but got nothing")
	}

	testutils.CheckStringContains(t, gotStdOut, wMsgNoIndexerOperationsForGitCommit)
}

func TestLocalLogLevels(t *testing.T) {
	// this is a regression test to ensure that the local log levels are still
	// valid if this test fails, the local log levels need to be updated
	loggerAvailableLevels := log.GetValidLevelOptionStrings()
	for _, level := range localLogLevels {
		if !slices.Contains(loggerAvailableLevels, level) {
			t.Errorf("local log level '%s' is not a valid logger level", level)
		}
	}
}

func createTestGitRepo(t *testing.T) {
	gitSourceRepoPathAbsoluteFS := os.DirFS(gitSourceRepoPathAbsolute)
	err := os.CopyFS(gitRepoTestGitRepoPathAbsolute, gitSourceRepoPathAbsoluteFS)
	if err != nil {
		t.Fatalf(
			`Unexpected error returned by `+
				`os.CopyFS(gitRepoTestGitRepoPathAbsolute, `+
				`gitSourceRepoPathAbsoluteFS): %s`,
			err.Error())
	}

	err = os.Rename(gitRepoTestGitRepoDotGitDirectory, gitRepoTestGitRepoHiddenGitDirectory)
	if err != nil {
		t.Fatalf(
			`Unexpected error returned by os.Rename(gitRepoTestGitRepoDotGitDirectory, `+
				`gitRepoTestGitRepoHiddenGitDirectory): %s`,
			err.Error())
	}
}

func deleteTestGitRepo(t *testing.T) {
	err := os.RemoveAll(gitRepoTestGitRepoPathAbsolute)
	if err != nil {
		t.Fatalf(
			`deleteTestGitRepo() failed with error "%s", remove %s manually`,
			err.Error(), gitRepoTestGitRepoPathAbsolute)
	}
}

func resetDeleteArgs() {
	cmd := DeleteCmd
	cmd.Flags().Set("eadid", "")
	cmd.Flags().Set("assume-yes", "")
	cmd.Flags().Set("logging-level", "")
}

func resetIndexArgs() {
	cmd := IndexCmd
	cmd.Flags().Set("file", "")
	cmd.Flags().Set("git-repo", "")
	cmd.Flags().Set("commit", "")
	cmd.Flags().Set("logging-level", "")
}
