package ead

import (
	"errors"
	"flag"
	"fmt"
	"go-ead-indexer/pkg/ead/collectiondoc"
	"go-ead-indexer/pkg/ead/component"
	"go-ead-indexer/pkg/ead/eadutil"
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

	err := clean()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestCollectionDocSolrAddMessage(t *testing.T) {
	testEADs := testutils.GetTestEADs()

	for _, testEAD := range testEADs {
		t.Run(testEAD, func(t *testing.T) {
			eadXML, err := testutils.GetEADFixtureValue(testEAD)
			if err != nil {
				t.Fatal(err)
			}

			repositoryCode := testutils.ParseRepositoryCode(testEAD)
			eadToTest, err := New(repositoryCode, eadXML)
			if err != nil {
				t.Fatal(err)
			}

			eadID := testutils.ParseEADID(testEAD)
			testCollectionDocSolrAddMessage(testEAD, eadID,
				eadToTest.CollectionDoc.SolrAddMessage, t)
		})
	}
}

func TestComponentDocSolrAddMessage(t *testing.T) {
	testEADs := testutils.GetTestEADs()

	for _, testEAD := range testEADs {
		t.Run(testEAD, func(t *testing.T) {
			eadXML, err := testutils.GetEADFixtureValue(testEAD)
			if err != nil {
				t.Fatal(err)
			}

			repositoryCode := testutils.ParseRepositoryCode(testEAD)
			eadToTest, err := New(repositoryCode, eadXML)
			if err != nil {
				t.Fatal(err)
			}

			componentIDs := []string{}
			for _, component := range *eadToTest.Components {
				componentIDs = append(componentIDs, component.ID)
				testComponentSolrAddMessage(testEAD, component.ID,
					component.SolrAddMessage, t)
			}

			testNoMissingComponents(testEAD, componentIDs, t)
		})
	}
}

func TestNewWithBadEADXML(t *testing.T) {
	testCases := []struct {
		eadXML        string
		expectedError string
	}{
		{
			eadXML:        "",
			expectedError: errorNoEADTagWithExpectedStructureFound,
		},
		{
			eadXML:        "/path/to/ead/file/edip/mos_2024.xml",
			expectedError: errorNoEADTagWithExpectedStructureFound,
		},
	}

	const repositoryCode = "edip"

	for _, testCase := range testCases {
		_, err := New(repositoryCode, testCase.eadXML)
		if err == nil {
			t.Errorf(`Expected an error to be returned by New("%s", "%s"), but no error was returned.`,
				repositoryCode, testCase.eadXML)

			continue
		}

		errorString := err.Error()
		if errorString != testCase.expectedError {
			t.Errorf(`Expected error "%s" to be returned by New("%s", "%s"),`+
				` but got: "%s"`,
				testCase.expectedError, repositoryCode, testCase.eadXML, errorString)
		}
	}
}

// We don't have an official repository code format, but there is a comprehensive
// list of repository codes:
// https://jira.nyu.edu/browse/FADESIGN-65
func TestNewWithInvalidRepositoryCode(t *testing.T) {
	testCases := []struct {
		repositoryCode string
		expectedError  string
	}{
		{
			repositoryCode: "",
			expectedError:  `Invalid repository code: ""`,
		},
		{
			repositoryCode: "!#$",
			expectedError:  `Invalid repository code: "!#$"`,
		},
	}

	eadXML, err := testutils.GetEADFixtureValue("edip/mos_2024")
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		_, err := New(testCase.repositoryCode, eadXML)
		if err == nil {
			t.Errorf(`Expected an error to be returned by New("%s", "..."), but no error was returned.`,
				testCase.repositoryCode)

			continue
		}

		errorString := err.Error()
		if errorString != testCase.expectedError {
			t.Errorf(`Expected error "%s" to be returned by New("%s", "..."),`+
				` but got: "%s"`,
				testCase.expectedError, testCase.repositoryCode, errorString)
		}
	}
}

func clean() error {
	err := os.RemoveAll(tmpFilesDirPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(tmpFilesDirPath, 0700)
	if err != nil {
		return err
	}

	_, err = os.Create(filepath.Join(tmpFilesDirPath, ".gitkeep"))
	if err != nil {
		return err
	}

	return nil
}

func testCollectionDocSolrAddMessage(testEAD string, fileID string,
	solrAddMessage collectiondoc.SolrAddMessage, t *testing.T) {
	testSolrAddMessageXML(testEAD, fileID, fmt.Sprintf("%s", solrAddMessage), t)
}

func testComponentSolrAddMessage(testEAD string, fileID string,
	solrAddMessage component.SolrAddMessage, t *testing.T) {
	testSolrAddMessageXML(testEAD, fileID, fmt.Sprintf("%s", solrAddMessage), t)
}

func testNoMissingComponents(testEAD string, componentIDs []string, t *testing.T) {
	missingComponents := []string{}

	goldenFileIDs := testutils.GetGoldenFileIDs(testEAD)
	goldenFileIDs = slices.DeleteFunc(goldenFileIDs, func(goldenFileID string) bool {
		return goldenFileID == testutils.ParseEADID(testEAD)
	})

	for _, goldenFileID := range goldenFileIDs {
		if !slices.Contains(componentIDs, goldenFileID) {
			missingComponents = append(missingComponents, goldenFileID)
		}
	}

	if len(missingComponents) > 0 {
		slices.SortStableFunc(missingComponents, func(a string, b string) int {
			return strings.Compare(a, b)
		})
		failMessage := fmt.Sprintf("`EAD.Components` for testEAD %s is missing the following component IDs:\n%s",
			testEAD, strings.Join(missingComponents, "\n"))
		t.Errorf(failMessage)
	}
}

func testSolrAddMessageXML(testEAD string, fileID string,
	actualValue string, t *testing.T) {
	if *updateGoldenFiles {
		err := testutils.UpdateGoldenFile(testEAD, fileID, actualValue)
		if err != nil {
			t.Fatalf("Error updating golden file: %s", err)
		}
	}

	goldenValue, err := testutils.GetGoldenFileValue(testEAD, fileID)
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
		err := writeActualSolrXMLToTmp(testEAD, fileID, actualValue)
		if err != nil {
			t.Fatalf("Error writing actual temp file for test case \"%s/%s\": %s",
				testEAD, fileID, err)
		}

		goldenFile := testutils.GoldenFilePath(testEAD, fileID)
		actualFile := tmpFile(testEAD, fileID)
		goldenValue, err := testutils.GetGoldenFileValue(testEAD, fileID)
		if err != nil {
			t.Fatalf("Error fetching golden value for %s: %s\n"+
				"Manually diff these files to determine the reasons for test failure.",
				goldenFile, err)
		}
		actualValue, err := os.ReadFile(tmpFile(testEAD, fileID))
		if err != nil {
			t.Fatalf("Error fetching actual value for %s: %s\n"+
				"Manually diff these files to determine the reasons for test failure.",
				actualFile, err)
		}

		// Diff the prettified XML and not the minified XML in the files,
		// to make it easier for a human to see the discrepancies.
		diff := util.DiffStrings("golden XML (prettified -- NOT the original)",
			eadutil.PrettifySolrAddMessageXML(goldenValue),
			"actual (prettified -- NOT the original)",
			eadutil.PrettifySolrAddMessageXML(string(actualValue)))
		t.Errorf("golden and actual values for %s do not match:\n%s"+
			fileID, diff)
	}
}

func tmpFile(testEAD string, fileID string) string {
	return filepath.Join(tmpFilesDirPath, testEAD, fileID)
}

func writeActualSolrXMLToTmp(testEAD string, fileID string, actual string) error {
	tmpFile := tmpFile(testEAD, fileID)
	err := os.MkdirAll(filepath.Dir(tmpFile), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(tmpFile, []byte(actual), 0644)
}
