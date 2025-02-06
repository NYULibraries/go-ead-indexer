package testutils

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type SolrClientMock struct {
	fileDir          string            // Directory containing the "golden files" for the test
	goldenFileHashes map[string]string // Hashes of the golden files
}

func GetSolrClientMock() *SolrClientMock {
	return &SolrClientMock{
		fileDir:          "",
		goldenFileHashes: make(map[string]string),
	}
}

func (sc *SolrClientMock) Add(xmlPostBody string) error {
	return sc.updateHash(xmlPostBody)
}

func (sc *SolrClientMock) Commit() error {
	commitXML := `<?xml version="1.0" encoding="UTF-8"?><commit/>`
	return sc.updateHash(commitXML)
}

func (sc *SolrClientMock) Delete(string) error {
	return nil
}

func (sc *SolrClientMock) GetPostRequest(string) (*http.Request, error) {
	return nil, nil
}

func (sc *SolrClientMock) GetSolrURLOrigin() string {
	return ""
}

func (sc *SolrClientMock) Reset() {
	// reset the solr client mock
	clear(sc.goldenFileHashes)
	sc.fileDir = ""
}

func (sc *SolrClientMock) IsComplete() bool {
	if len(sc.goldenFileHashes) != 0 {
		return false
	}
	return true
}

func (sc *SolrClientMock) SetupMock(goldenFileDir string) error {
	// assumes all files in the directory are golden files
	// and that all files will be consumed by the test
	sc.fileDir = goldenFileDir

	files, err := os.ReadDir(goldenFileDir)
	if err != nil {
		return err
	}

	// load the golden file hashes map
	for _, file := range files {
		filePath := filepath.Join(goldenFileDir, file.Name())

		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}
		sc.goldenFileHashes[string(h.Sum(nil))] = filePath
	}

	return nil
}

func (sc *SolrClientMock) updateHash(xmlPostBody string) error {
	h := md5.New()
	io.WriteString(h, xmlPostBody)

	hash := string(h.Sum(nil))
	if _, ok := sc.goldenFileHashes[hash]; !ok {
		return fmt.Errorf("hash '%s' not found in golden file hashes".hash)
	}
	// remove the hash from the golden file hash map
	delete(sc.goldenFileHashes, hash)
	return nil
}
