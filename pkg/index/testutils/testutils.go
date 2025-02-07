package testutils

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type SolrClientMock struct {
	fileDir          string            // Directory containing the "golden files" for the test
	GoldenFileHashes map[string]string // Hashes of the golden files
}

func GetSolrClientMock() *SolrClientMock {
	return &SolrClientMock{
		fileDir:          "",
		GoldenFileHashes: make(map[string]string),
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
	clear(sc.GoldenFileHashes)
	sc.fileDir = ""
}

func (sc *SolrClientMock) IsComplete() bool {
	return len(sc.GoldenFileHashes) == 0
}

// func (sc *SolrClientMock) SetupMock(goldenFileDir, suffix string) error {
// 	// assumes all files in the directory are golden files
// 	// and that all files will be consumed by the test
// 	sc.fileDir = goldenFileDir

// 	files, err := os.ReadDir(goldenFileDir)
// 	if err != nil {
// 		return err
// 	}

// 	// load the golden file hashes map
// 	for _, file := range files {
// 		// skip non-matching files
// 		if !strings.HasSuffix(file.Name(), suffix) {
// 			continue
// 		}

// 		filePath := filepath.Join(goldenFileDir, file.Name())

// 		f, err := os.Open(filePath)
// 		if err != nil {
// 			return err
// 		}
// 		defer f.Close()

// 		h := md5.New()
// 		if _, err := io.Copy(h, f); err != nil {
// 			return err
// 		}
// 		sc.GoldenFileHashes[string(h.Sum(nil))] = filePath
// 	}

// 	return nil
// }

func (sc *SolrClientMock) SetupMock(goldenFileDir string) error {
	// assumes all files in the directory are golden files
	// and that all files will be consumed by the test
	sc.fileDir = goldenFileDir

	files, err := os.ReadDir(goldenFileDir)
	if err != nil {
		return err
	}

	// load the golden file hashes map
	h := md5.New()
	for _, file := range files {
		filePath := filepath.Join(goldenFileDir, file.Name())

		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(h, f); err != nil {
			return err
		}

		sum := formattedHashSum(h)

		// look for collisions
		if sc.GoldenFileHashes[sum] != "" {
			return fmt.Errorf("duplicate hash '%s' found in golden file hashes for file: %s, file already in hash: %s", sum, filePath, sc.GoldenFileHashes[sum])
		}
		// no collision, add the hash to the golden file hash map
		sc.GoldenFileHashes[sum] = filePath
		h.Reset()
	}

	return nil
}

func (sc *SolrClientMock) updateHash(xmlPostBody string) error {
	h := md5.New()
	io.WriteString(h, xmlPostBody)

	hash := formattedHashSum(h)
	if _, ok := sc.GoldenFileHashes[hash]; !ok {
		return fmt.Errorf("hash '%s' not found in golden file hashes", hash)
	}
	// remove the hash from the golden file hash map
	delete(sc.GoldenFileHashes, hash)
	return nil
}

func formattedHashSum(h hash.Hash) string {
	return fmt.Sprintf("%x", h.Sum(nil))
}
