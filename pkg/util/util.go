package util

import (
	"errors"
	"fmt"
	"go-ead-indexer/pkg/util/diff"
	"net"
	"os"
	"runtime"
	"strings"
)

// Replicate https://github.com/NYULibraries/ead_indexer/blob/a367ab8cc791376f0d8a287cbcd5b6ee43d5c04f/lib/ead_indexer/behaviors.rb#L137
// TODO: If it turns out that we don't need to preserve order in any of our uniq'ed
// slices, remove this.
func CompactStringSlicePreserveOrder(stringSlice []string) []string {
	compactedSlice := []string{}
	seen := map[string]bool{}
	for _, element := range stringSlice {
		if seen[element] {
			continue
		}
		compactedSlice = append(compactedSlice, element)
		seen[element] = true
	}

	return compactedSlice
}

func DiffFiles(path1 string, path2 string) (string, error) {
	bytes1, err := os.ReadFile(path1)
	if err != nil {
		return "", err
	}

	bytes2, err := os.ReadFile(path2)
	if err != nil {
		return "", err
	}

	diffBytes := diff.Diff(path1, bytes1, path2, bytes2)

	return string(diffBytes), nil
}

func DiffStrings(label1 string, string1 string, label2, string2 string) string {
	diffString := string(diff.Diff(label1, []byte(string1), label2, []byte(string2)))

	return diffString
}

// Return package name and function name for caller on the callstack, selected
// by `skip` which is the number of stack frames to ascend, where 0 identifies
// the caller of this utility function.  See `runtime.Caller` code and/or
// documentations for more details.
func GetCallerFunctionName(skip int) (string, string, error) {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "", "", errors.New(fmt.Sprintf("runtime.Caller(%d) failed", skip))
	}

	fullyQualifiedFunctionName := runtime.FuncForPC(pc).Name()
	dotSeparator := strings.LastIndexByte(fullyQualifiedFunctionName, '.')

	return fullyQualifiedFunctionName[:dotSeparator],
		fullyQualifiedFunctionName[dotSeparator+1:],
		nil
}

// This is the method that Go itself uses: see net/http/httptest/server.go:
// https://github.com/golang/go/blob/69234ded30614a471c35cef5d87b0e0d3c136cd9/src/net/http/httptest/server.go#L60-L75
// See also:
// "Is it possible to connect to TCP port 0?"
// https://unix.stackexchange.com/questions/180492/is-it-possible-to-connect-to-tcp-port-0
func GetUnusedLocalhostNetworkAddress() string {
	// Based on https://github.com/golang/go/blob/69234ded30614a471c35cef5d87b0e0d3c136cd9/src/net/http/httptest/server.go#L68-L73
	throwawayListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if throwawayListener, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			panic(fmt.Sprintf("httptest: failed to listen on a port: %v", err))
		}
	}
	defer throwawayListener.Close()

	return throwawayListener.Addr().String()
}

// Based on: https://stackoverflow.com/questions/18594330/what-is-the-best-way-to-test-for-an-empty-string-in-go
func IsNonEmptyString(value string) bool {
	return len(strings.TrimSpace(value)) > 0
}
