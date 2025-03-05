package testutils

import (
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func CaptureCmdOutput(f func(*cobra.Command, []string), cmd *cobra.Command, args []string) string {
	// save the current stdout
	rescueStdout := os.Stdout

	// create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	// map stdout to the write end of the pipe
	os.Stdout = w

	// execute the command
	f(cmd, args)

	// close the write end of the pipe
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

	// read from the read end of the pipe
	out, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	// restore stdout
	os.Stdout = rescueStdout
	return string(out)
}

func CaptureCmdOutputE(f func(*cobra.Command, []string) error, cmd *cobra.Command, args []string) (string, error) {
	// save the current stdout
	rescueStdout := os.Stdout

	// create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	// map stdout to the write end of the pipe
	os.Stdout = w

	// execute the command
	cmdErr := f(cmd, args)

	// close the write end of the pipe
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

	// read from the read end of the pipe
	out, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	// restore stdout
	os.Stdout = rescueStdout

	return string(out), cmdErr
}

func CaptureCmdStdoutStderr(f func(*cobra.Command, []string), cmd *cobra.Command, args []string) (string, string) {
	// save the current stdout and stderr
	rescueStdout := os.Stdout
	rescueStderr := os.Stderr

	// create a pipe for stdout
	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	// create a pipe for stderr
	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	// map stdout and stderr to the write ends of the pipes
	os.Stdout = stdoutW
	os.Stderr = stderrW

	// execute the command
	f(cmd, args)

	// close the write end of the stdout pipe
	err = stdoutW.Close()
	if err != nil {
		log.Fatal(err)
	}

	// close the write end of the stderr pipe
	err = stderrW.Close()
	if err != nil {
		log.Fatal(err)
	}

	// read from the read ends of the stdout pipe
	stdout, err := io.ReadAll(stdoutR)
	if err != nil {
		log.Fatal(err)
	}

	// read from the read ends of the stderr pipe
	stderr, err := io.ReadAll(stderrR)
	if err != nil {
		log.Fatal(err)
	}

	// restore stdout and stderr
	os.Stdout = rescueStdout
	os.Stderr = rescueStderr

	return string(stdout), string(stderr)
}

func CaptureCmdStdoutStderrE(f func(*cobra.Command, []string) error, cmd *cobra.Command, args []string) (string, string, error) {

	// save the current stdout and stderr
	rescueStdout := os.Stdout
	rescueStderr := os.Stderr

	// create a pipe for stdout
	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	// create a pipe for stderr
	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	// map stdout and stderr to the write ends of the pipes
	os.Stdout = stdoutW
	os.Stderr = stderrW

	// execute the command
	cmdErr := f(cmd, args)

	// close the write end of the stdout pipe
	err = stdoutW.Close()
	if err != nil {
		log.Fatal(err)
	}

	// close the write end of the stderr pipe
	err = stderrW.Close()
	if err != nil {
		log.Fatal(err)
	}

	// read from the read end of the stdout pipe
	stdout, err := io.ReadAll(stdoutR)
	if err != nil {
		log.Fatal(err)
	}

	// read from the read end of the stderr pipe
	stderr, err := io.ReadAll(stderrR)
	if err != nil {
		log.Fatal(err)
	}

	// restore stdout and stderr
	os.Stdout = rescueStdout
	os.Stderr = rescueStderr

	return string(stdout), string(stderr), cmdErr
}

func SetCmdFlag(f *cobra.Command, flag string, val string) {
	f.Flags().Set(flag, val)
}

// copied from Cobra test code:
// https://github.com/spf13/cobra/blob/40d34bca1bffe2f5e84b18d7fd94d5b3c02275a6/command_test.go#L49
func CheckStringContains(t *testing.T, got, expected string) {
	if !strings.Contains(got, expected) {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", expected, got)
	}
}
