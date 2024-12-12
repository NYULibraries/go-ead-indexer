package solr

import (
	"flag"
	eadtestutils "go-ead-indexer/pkg/ead/testutils"
	"go-ead-indexer/pkg/net/solr/testutils"
	"go-ead-indexer/pkg/util"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func TestAdd(t *testing.T) {
	testAdd_successfulAdds(t)
}

func testAdd_successfulAdd(goldenFileID string, t *testing.T) {
	postBody, err := eadtestutils.GetGoldenFileValue(testutils.TestEAD, goldenFileID)
	if err != nil {
		t.Fatalf("eadtestutils.GetGoldenFileValue(testutils.TestEAD, goldenFileID) failed with error: %s", err)
	}

	err = Add(postBody)
	if err != nil {
		t.Fatalf("Expected no error for %s, got: %s", goldenFileID, err)
	}

	actualRequest, err := testutils.GetActualFileContents(testutils.TestEAD, goldenFileID)
	if err != nil {
		t.Fatalf("testutils.getActualFileContents(testutils.TestEAD, goldenFileID) failed with error: %s", err)
	}

	expectedRequest := testutils.GetExpectedPOSTRequestString(postBody)
	diff := util.DiffStrings("expected", expectedRequest,
		"actual", actualRequest)
	if diff != "" {
		t.Errorf(`%s fail: actual request does not match expected: %s`,
			goldenFileID, diff)
	}
}

func testAdd_successfulAdds(t *testing.T) {
	err := testutils.Clean()
	if err != nil {
		t.Errorf("clean() failed with error: %s", err)
	}

	fakeSolrServer := testutils.MakeSolrFake(t)
	defer fakeSolrServer.Close()

	err = SetSolrURLOrigin(fakeSolrServer.URL)
	if err != nil {
		t.Fatalf(`Setup of Solr fake failed with error: %s`, err)
	}

	goldenFileIDs := eadtestutils.GetGoldenFileIDs(testutils.TestEAD)
	for _, goldenFileID := range goldenFileIDs {
		testAdd_successfulAdd(goldenFileID, t)
	}
}

func TestSetSolrURL(t *testing.T) {
	testCases := []struct {
		origin        string
		expectedError string
	}{
		{
			origin:        "",
			expectedError: `parse "": empty url`,
		},
		{
			origin:        "x",
			expectedError: `parse "x": invalid URI for request`,
		},
		{
			origin:        "ftp://solr-host.org",
			expectedError: `SetSolrURLOrigin("ftp://solr-host.org"): invalid scheme`,
		},
		{
			origin:        "http://",
			expectedError: `SetSolrURLOrigin("http://"): host is empty`,
		},
		{
			origin:        "https://",
			expectedError: `SetSolrURLOrigin("https://"): host is empty`,
		},
		{
			origin:        testutils.FakeSolrHostAndPort,
			expectedError: `SetSolrURLOrigin("` + testutils.FakeSolrHostAndPort + `"): invalid scheme`,
		},
		{
			origin:        "http://" + testutils.FakeSolrHostAndPort,
			expectedError: "",
		},
	}

	for _, testCase := range testCases {
		actualError := SetSolrURLOrigin(testCase.origin)
		var actualErrorString string
		if actualError != nil {
			actualErrorString = actualError.Error()
		}

		if testCase.expectedError == "" {
			if actualErrorString != "" {
				t.Errorf(`SetSolrURLOrigin("%s") should not have returned an error,`+
					` but it returned error "%s"`, testCase.origin, actualErrorString)
			}

			actualOrigin := GetSolrURLOrigin()
			if actualOrigin != testCase.origin {
				t.Errorf(`GetSolrURLOrigin() should return "%s", but it instead`+
					` returned "%s"`, testCase.origin, actualOrigin)
			}
		} else {
			if actualErrorString == "" {
				t.Errorf(`SetSolrURLOrigin("%s") should have returned error "%s",`+
					" but no error was returned", testCase.origin, testCase.expectedError)
			} else if actualErrorString != testCase.expectedError {
				t.Errorf(`SetSolrURLOrigin("%s") should have returned error "%s",`+
					` but instead returned error "%s"`, testCase.origin, testCase.expectedError,
					actualErrorString)
			}

			actualOrigin := GetSolrURLOrigin()
			if actualOrigin != "" {
				t.Errorf(`GetSolrURLOrigin() should return "", but it instead`+
					` returned "%s"`, actualOrigin)
			}
		}
	}
}
