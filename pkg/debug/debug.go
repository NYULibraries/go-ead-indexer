package debug

import (
	"encoding/json"
	"errors"
	"github.com/nyulibraries/go-ead-indexer/pkg/ead"
	"github.com/nyulibraries/go-ead-indexer/pkg/git"
	"github.com/nyulibraries/go-ead-indexer/pkg/net/solr"
	"github.com/nyulibraries/go-ead-indexer/pkg/util"
	"io/fs"
	"net/http/httputil"
	"os"
	"path"
	"path/filepath"
)

type dumpedSolrIndexerHTTPRequestsForEADFile struct {
	CollectionDoc string            `json:"collectiondoc"`
	Components    map[string]string `json:"components"`
}

func DumpSolrIndexerHTTPRequestsForEADFile(eadFile string) (string, error) {
	dumpedHTTPRequests, err :=
		getDumpedSolrIndexerHTTPRequestsForEADFile(eadFile)
	if err != nil {
		return "", err
	}

	dumpedHTTPRequestsJSONBytes, err :=
		json.MarshalIndent(dumpedHTTPRequests, "", "    ")
	if err != nil {
		return "", err
	}

	return string(dumpedHTTPRequestsJSONBytes), nil
}

func DumpSolrIndexerHTTPRequestsForGitCommit(repoPath string, commit string) (string, error) {
	var repoPathAbsolute string
	if filepath.IsAbs(repoPath) {
		repoPathAbsolute = repoPath
	} else {
		var err error
		repoPathAbsolute, err = filepath.Abs(repoPath)
		if err != nil {
			return "", err
		}
	}

	eadFilesForCommit, err := git.ListEADFilesForCommit(repoPathAbsolute, commit)
	if err != nil {
		return "", err
	}

	dumpedSolrIndexerHTTPRequests := map[string]dumpedSolrIndexerHTTPRequestsForEADFile{}
	for eadFileRelativePath, _ := range eadFilesForCommit {
		if eadFilesForCommit[eadFileRelativePath] == git.Add {
			eadFileAbsolutePath := path.Join(repoPathAbsolute, eadFileRelativePath)
			dumpedHTTPRequests, err :=
				getDumpedSolrIndexerHTTPRequestsForEADFile(eadFileAbsolutePath)
			if err != nil {
				return "", err
			}

			dumpedSolrIndexerHTTPRequests[eadFileRelativePath] = dumpedHTTPRequests
		}
	}

	dumpedHTTPRequestsJSONBytes, err :=
		json.MarshalIndent(dumpedSolrIndexerHTTPRequests, "", "    ")
	if err != nil {
		return "", err
	}

	return string(dumpedHTTPRequestsJSONBytes), nil
}

func getDumpedSolrIndexerHTTPRequestsForEADFile(eadFile string) (dumpedSolrIndexerHTTPRequestsForEADFile, error) {
	eadXML, err := os.ReadFile(eadFile)
	if errors.Is(err, fs.ErrNotExist) {
		return dumpedSolrIndexerHTTPRequestsForEADFile{}, err
	}

	repositoryCode, err := util.GetRepositoryCode(eadFile)
	if err != nil {
		return dumpedSolrIndexerHTTPRequestsForEADFile{}, err
	}

	eadObject, err := ead.New(repositoryCode, string(eadXML))
	if err != nil {
		return dumpedSolrIndexerHTTPRequestsForEADFile{}, err
	}

	sc, err := solr.NewSolrClient("http://example.com")
	if err != nil {
		return dumpedSolrIndexerHTTPRequestsForEADFile{}, err
	}

	dumpedHTTPRequests := dumpedSolrIndexerHTTPRequestsForEADFile{}
	dumpedHTTPRequests.Components = map[string]string{}

	postRequest, err :=
		sc.GetPostRequest(eadObject.CollectionDoc.SolrAddMessage.String())
	if err != nil {
		return dumpedSolrIndexerHTTPRequestsForEADFile{}, err
	}

	dumpedPostRequest, err := httputil.DumpRequest(postRequest, true)
	if err != nil {
		return dumpedSolrIndexerHTTPRequestsForEADFile{}, err
	}

	dumpedHTTPRequests.CollectionDoc = string(dumpedPostRequest)

	for _, component := range *eadObject.Components {
		postRequest, err :=
			sc.GetPostRequest(eadObject.CollectionDoc.SolrAddMessage.String())
		if err != nil {
			return dumpedSolrIndexerHTTPRequestsForEADFile{}, err
		}

		dumpedPostRequest, err := httputil.DumpRequest(postRequest, true)
		if err != nil {
			return dumpedSolrIndexerHTTPRequestsForEADFile{}, err
		}

		dumpedHTTPRequests.Components[component.ID] = string(dumpedPostRequest)
	}

	return dumpedHTTPRequests, nil
}
