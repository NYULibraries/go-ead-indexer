package git

import (
	"errors"
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
)

type IndexerOperation string

const (
	Add     IndexerOperation = "add"
	Delete  IndexerOperation = "delete"
	Unknown IndexerOperation = "unknown"
)

func ListEADFilesForCommit(repoPath string, thisCommitHashString string) (map[string]IndexerOperation, error) {

	operations := make(map[string]IndexerOperation)

	// Opens an already existing repository.
	repo, err := gogit.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	// Get the commit object by hash
	thisCommitHash := plumbing.NewHash(thisCommitHashString)
	thisCommit, err := repo.CommitObject(thisCommitHash)
	if err != nil {
		return nil, fmt.Errorf("problem getting commit object for commit hash %s: %s", thisCommitHash, err)
	}

	// handle the initial commit case
	if len(thisCommit.ParentHashes) == 0 {
		// get the tree and list the files
		tree, err := thisCommit.Tree()
		if err != nil {
			return nil, err
		}
		files := tree.Files()

		for {
			file, err := files.Next()
			if err != nil {
				break
			}
			operations[file.Name] = Add
		}
		return operations, nil
	}

	// Get the parent commit
	parentHash := thisCommit.ParentHashes[0]
	parentCommit, err := repo.CommitObject(parentHash)
	if err != nil {
		return nil, err
	}

	// Get the changes between the two commits
	patch, err := parentCommit.Patch(thisCommit)
	if err != nil {
		return nil, err
	}

	var errs []error
	for _, fileChange := range patch.FilePatches() {
		from, to := fileChange.Files()
		k, v := classifyFileChange(from, to)
		if v == Unknown {
			// unable to determine the type of change
			errs = append(errs, fmt.Errorf("unable to determine file transition: Commits: commit '%s', parent: '%s', Files: from '%s', to '%s'", thisCommitHashString, parentHash.String(), getPath(from), getPath(to)))
			continue
		}

		operations[k] = v
	}
	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	return operations, nil
}

func getPath(f diff.File) string {
	if f == nil {
		return ""
	}
	return f.Path()
}

func classifyFileChange(from, to diff.File) (string, IndexerOperation) {
	/*
		add    --> from.Path() is nil &&
				     to.Path() is not nil

		update --> from.Path() == to.Path()

		delete --> from.Path() is not nil &&
				     to.Path() is nil
	*/
	switch {
	case from == nil && to == nil:
		// this shouldn't happen
		return "", Unknown
	case from == nil && to != nil:
		return to.Path(), Add
	case from != nil && to == nil:
		return from.Path(), Delete
	case from.Path() == to.Path():
		return to.Path(), Add
	default:
		return "", Unknown
	}
}
