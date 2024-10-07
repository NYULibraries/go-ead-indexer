package ead

import (
	"errors"
	"flag"
	"fmt"
	"go-ead-indexer/pkg/ead/testutils"
	"go-ead-indexer/pkg/util"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

var tmpFilesDirPath = filepath.Join("testdata", "tmp", "actual")

var updateGoldenFiles = flag.Bool("update-golden-files", false, "update the golden files")

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func TestParseSolrAddMessages(t *testing.T) {
	eadIDs := testutils.GetTestEADIDs()

	for _, eadID := range eadIDs {
		t.Run(eadID, func(t *testing.T) {
			eadXML, err := testutils.GetEADFixtureValue(eadID)
			if err != nil {
				t.Fatal(err)
			}

			solrAddMessages, err := ParseSolrAddMessages(eadXML)
			if err != nil {
				t.Fatal(err)
			}

			testSolrAddMessage(eadID, eadID, solrAddMessages.Collection, t)

			for _, componentRequestBody := range *solrAddMessages.Components {
				testSolrAddMessage(eadID, componentRequestBody.ID, componentRequestBody.Message, t)
			}

			testNoMissingComponents(eadID, solrAddMessages, t)
		})
	}
}

func testNoMissingComponents(eadID string, solrAddMessages SolrAddMessages, t *testing.T) {
	missingComponents := []string{}

	goldenFileIDs := testutils.GetGoldenFileIDs(eadID)
	goldenFileIDs = slices.DeleteFunc(goldenFileIDs, func(goldenFileID string) bool {
		return goldenFileID == eadID
	})

	actualFileIDs := []string{}
	for _, componentAddMessage := range *solrAddMessages.Components {
		actualFileIDs = append(actualFileIDs, componentAddMessage.ID)
	}

	for _, goldenFileID := range goldenFileIDs {
		if !slices.Contains(actualFileIDs, goldenFileID) {
			missingComponents = append(missingComponents, goldenFileID)
		}
	}

	if len(missingComponents) > 0 {
		slices.SortStableFunc(missingComponents, func(a string, b string) int {
			return strings.Compare(a, b)
		})
		failMessage := fmt.Sprintf("`SolrAddMessages.Components` for eadID %s is missing messages for these fileIDs:\n%s",
			eadID, strings.Join(missingComponents, "\n"))
		t.Errorf(failMessage)
	}
}

func testSolrAddMessage(eadID string, fileID string, actualValue string, t *testing.T) {
	if *updateGoldenFiles {
		err := testutils.UpdateGoldenFile(eadID, fileID, actualValue)
		if err != nil {
			t.Fatalf("Error updating golden file: %s", err)
		}
	}

	goldenValue, err := testutils.GetGoldenFileValue(eadID, fileID)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// This is a test fail, not a fatal test execution error.
			// A missing golden file means that a Solr add message was created
			// for a component that shouldn't exist.
			t.Errorf("No golden file exists for \"%s\": %s",
				fileID, err)

			return
		} else {
			t.Fatalf("Error retrieving golden value for \"%s\": %s",
				fileID, err)
		}
	}

	if actualValue != goldenValue {
		err := writeActualSolrXMLToTmp(eadID, fileID, actualValue)
		if err != nil {
			t.Fatalf("Error writing actual temp file for test case \"%s/%s\": %s",
				eadID, fileID, err)
		}

		goldenFile := testutils.GoldenFilePath(eadID, fileID)
		actualFile := tmpFile(eadID, fileID)
		diff, err := util.DiffFiles(goldenFile, actualFile)
		if err != nil {
			t.Fatalf("Error diff'ing %s vs. %s: %s\n"+
				"Manually diff these files to determine the reasons for test failure.",
				goldenFile, actualFile, err)
		}

		t.Errorf("golden and actual values for %s do not match:\n%s\n",
			fileID, diff)
	}
}

func tmpFile(eadID string, fileID string) string {
	return filepath.Join(tmpFilesDirPath, eadID, fileID)
}

func writeActualSolrXMLToTmp(eadID string, fileID string, actual string) error {
	tmpFile := tmpFile(eadID, fileID)
	err := os.MkdirAll(filepath.Dir(tmpFile), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(tmpFile, []byte(actual), 0644)
}
