package git

import (
	"errors"
	"fmt"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gitdiff "github.com/go-git/go-git/v5/plumbing/format/diff"
	"strings"
)

type IndexerOperation string

const (
	Add     IndexerOperation = "add"
	Delete  IndexerOperation = "delete"
	Unknown IndexerOperation = "unknown"
)

const errNotAValidCommitHashStringTemplate = `"%s" is not a valid commit hash string`

// CheckoutMergeReset checks out a commit hash in a git repository.
//
// WARNING:
// This function uses the default "gogit.CheckoutOptions{Keep: false}".
//
// This means that:
// 1.) if there are any files under version control with uncommitted changes,
// the checkout will FAIL
//
// 2.) if there are any files in the git repo directory hierarchy that are
// NOT under version control, THOSE FILES WILL BE DELETED ON CHECKOUT!
func CheckoutMergeReset(repoPath string, commitHash string) error {
	// We need to test this before the `worktree.Checkout()` call below, because if
	// an invalid commit hash string is passed to `plumbing.NewHash()` the result
	// is a zero hash, and in such cases undesirable behavior can result.  We cannot
	// rely on `worktree.NewHash()` error handling to do what we would like.
	// For details, see https://jira.nyu.edu/browse/DLFA-276.
	if !plumbing.IsHash(commitHash) {
		return fmt.Errorf(errNotAValidCommitHashStringTemplate, commitHash)
	}

	repo, err := gogit.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Note if both Hash and Branch are empty, the Branch defaults to "master".
	// Source: https://github.com/go-git/go-git/blob/v5.13.1/options.go#L353
	// This includes if `plumbing.NewHash(commitHash)` returns an empty Hash.
	// If Hash is empty, undesirable things will happen, which is why we test
	// `commitHash` before this call.
	// For details, see https://jira.nyu.edu/browse/DLFA-276.
	err = worktree.Checkout(&gogit.CheckoutOptions{
		// `plumbing.NewHash()` will return an empty hash if `commitHash` is not a valid hexadecimal string.
		Hash: plumbing.NewHash(commitHash),
	})
	if err != nil {
		return fmt.Errorf("problem checking out hash '%s', error: '%s'",
			commitHash, err.Error())
	}

	return nil
}

// TODO: improve filtering out of files we don't want to accidentally process.
// See long comment before filter helper function definition.
func ListEADFilesForCommit(repoPath string,
	thisCommitHashString string) (map[string]IndexerOperation, error) {
	// This is ahelper function to prevent accidental inclusion of README.md and
	// .circleci/config.yml files.  See https://nyu.atlassian.net/browse/DLFA-302.
	//
	// Ideally, we want to only include EAD files that we would like to
	// index.  The rules for this might require some thought, as we've never
	// strictly defined what constitutes a valid filepath.	For example, do we
	// want to exclude:
	//   * EAD files that are not in one of the repository code subdirectories?
	//   * EAD files that are named correctly and are placed in the correct
	//     repository code directories, but are not formed correctly?
	//     E.g. empty or truncated?  If they were published in a proper directory
	//     with a proper name by the publisher and GT, should the indexer be
	//     second-guessing?  This is a question to consider because we might
	//     for example think that doing a string test for an expected tag might
	//     be a cheap way to test for inclusion.
	//   * EAD files that don't have extension ".xml", like ".XML"?
	//
	// Keep in mind that if we let through files that are not valid EAD2002
	// files it's not the end of the world, because they will never make it past
	// the parsing step.  Likewise, files located in the wrong place or named
	// with an unexpected extension could also have trouble making it through
	// various processing steps.
	//
	// At the moment, we don't keep a list of valid repository codes, but if we
	// did, we could filter by ensuring all included filepaths were of the form:
	// `<valid repository code>/<valid EAD ID>.xml`.
	//
	// For now, we just test for ".xml" extension.  We currently have no XML
	// files in the repo that are not actual EAD files.  We could conceivably
	// have  some in the future though, for example if we switch to a new
	// CI/CD solution which uses XML configuration files, and if that ever ends
	// up being the case, we would need to enhance this function.
	//
	// Since this is a very context specific function, we are not putting it in
	// the `util` package, for example as `util.IsEADFile()`, and for the same
	// reason we don't even necessarily want it to have package level scope.
	isValidFilepath := func(filepath string) bool {
		return strings.HasSuffix(filepath, ".xml")
	}

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
		return nil,
			fmt.Errorf("problem getting commit object for commit hash %s: %s",
				thisCommitHash, err)
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

			if isValidFilepath(file.Name) {
				operations[file.Name] = Add
			}
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

		if !isValidFilepath(k) {
			continue
		}

		if v == Unknown {
			// unable to determine the type of change
			errs = append(errs,
				fmt.Errorf("unable to determine file transition: Commits: "+
					"commit '%s', parent: '%s', Files: from '%s', to '%s'",
					thisCommitHashString,
					parentHash.String(),
					getPath(from),
					getPath(to)))
			continue
		}

		operations[k] = v
	}
	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	return operations, nil
}

func getPath(f gitdiff.File) string {
	if f == nil {
		return ""
	}
	return f.Path()
}

func classifyFileChange(from, to gitdiff.File) (string, IndexerOperation) {
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
