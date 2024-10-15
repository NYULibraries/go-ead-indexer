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

			eadToTest, err := New(eadXML)
			if err != nil {
				t.Fatal(err)
			}

			testSolrAddMessage(eadID, eadID, eadToTest.Collection.SolrAddMessage, t)

			componentIDs := []string{}
			for _, component := range *eadToTest.Components {
				componentIDs = append(componentIDs, component.ID)
				testSolrAddMessage(eadID, component.ID, component.SolrAddMessage, t)
			}

			testNoMissingComponents(eadID, componentIDs, t)
		})
	}
}

func TestReplaceMARCSubfieldDemarcators(t *testing.T) {
	// To see where some of these real life examples came from:
	// https://jira.nyu.edu/browse/DLFA-229?focusedCommentId=10153922&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10153922
	testCases := []struct {
		in  string
		out string
	}{
		{
			"",
			"",
		},
		{
			"Laundry industry |z New York (State) |z New York.",
			"Laundry industry -- New York (State) -- New York.",
		},
		{
			"China |x History |x Tiananmen Square Incident, 1989",
			"China -- History -- Tiananmen Square Incident, 1989",
		},
		{
			"Labor Unions |z United States |y 1980-1990.",
			"Labor Unions -- United States -- 1980-1990.",
		},
		{
			"Elections |z United States |x History |y 20th century |v Statistics.",
			"Elections -- United States -- History -- 20th century -- Statistics.",
		},
		{
			"Randall, Margaret, |d 1936-",
			"Randall, Margaret, -- 1936-",
		},
		{
			"General strikes |Z New York (State) |z Kings County",
			"General strikes -- New York (State) -- Kings County",
		},
		{
			"Theaters |x Employees |X Labor unions |z United States.",
			"Theaters -- Employees -- Labor unions -- United States.",
		},
		{
			"France. |t Constitution (1958).",
			"France. -- Constitution (1958).",
		},
		{
			"United States. Congress. House. |b Committee on Education and Labor. |b Select Subcommittee on Education",
			"United States. Congress. House. -- Committee on Education and Labor. -- Select Subcommittee on Education",
		},
		{
			"Wagner, Richard, 1813-1883. |t Operas. |k Selections",
			"Wagner, Richard, 1813-1883. -- Operas. -- Selections",
		},
		{
			"Germany. |t Treaties, etc. |g Soviet Union, |d 1939 Aug. 23.",
			"Germany. -- Treaties, etc. -- Soviet Union, -- 1939 Aug. 23.",
		},
		{
			"DO | NOT || CHANGE",
			"DO | NOT || CHANGE",
		},
		// TODO: fix the bug we've intentionally preserved in MARC subfield demarcation
		// replacement.  For details, see:
		//
		//   - https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10154897&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10154897
		//   - https://jira.nyu.edu/browse/DLFA-229?focusedCommentId=10153922&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10153922
		//
		// Once that is done, we can uncomment these tests, which currently fail.
		//{
		//	"Violence: Recode / UNDER|STAND",
		//	"Violence: Recode / UNDER|STAND",
		//},
		//{
		//	"85-2126 | John Hans[e|o]n (from Box 4 of 6)",
		//	"85-2126 | John Hans[e|o]n (from Box 4 of 6)",
		//},
	}

	for _, testCase := range testCases {
		actual := replaceMARCSubfieldDemarcators(testCase.in)
		if actual != testCase.out {
			t.Errorf(`Expected output string "%s" for input string "%s", got "%s""`,
				testCase.out, testCase.in, actual)
		}
	}
}

func testNoMissingComponents(eadID string, componentIDs []string, t *testing.T) {
	missingComponents := []string{}

	goldenFileIDs := testutils.GetGoldenFileIDs(eadID)
	goldenFileIDs = slices.DeleteFunc(goldenFileIDs, func(goldenFileID string) bool {
		return goldenFileID == eadID
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
		failMessage := fmt.Sprintf("`EAD.Components` for eadID %s is missing the following component IDs:\n%s",
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
