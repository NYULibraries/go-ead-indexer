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
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f(cmd, args)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out)
}

func CaptureCmdOutputE(f func(*cobra.Command, []string) error, cmd *cobra.Command, args []string) (string, error) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := f(cmd, args)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out), err
}

func CaptureCmdStdoutStderr(f func(*cobra.Command, []string), cmd *cobra.Command, args []string) (string, string) {
	rescueStdout := os.Stdout
	rescueStderr := os.Stderr

	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout = stdoutW
	os.Stderr = stderrW

	f(cmd, args)

	stdoutW.Close()
	stderrW.Close()

	stdout, _ := io.ReadAll(stdoutR)
	stderr, _ := io.ReadAll(stderrR)
	os.Stdout = rescueStdout
	os.Stderr = rescueStderr

	return string(stdout), string(stderr)
}

func CaptureCmdStdoutStderrE(f func(*cobra.Command, []string) error, cmd *cobra.Command, args []string) (string, string, error) {
	rescueStdout := os.Stdout
	rescueStderr := os.Stderr

	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout = stdoutW
	os.Stderr = stderrW

	err = f(cmd, args)

	stdoutW.Close()
	stderrW.Close()

	stdout, _ := io.ReadAll(stdoutR)
	stderr, _ := io.ReadAll(stderrR)
	os.Stdout = rescueStdout
	os.Stderr = rescueStderr

	return string(stdout), string(stderr), err
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
