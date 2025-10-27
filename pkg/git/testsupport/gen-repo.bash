#!/bin/bash
set -uo pipefail

# This script generates a git repo test fixture for use in the git package tests.
# The script creates a directory named 'git-repo', then generates and commits various files.
# Finally, the script renames the 'git-repo/.git' directory to 'git-repo/dot-git'.
# The git-repo directory can be moved into the pkg/git/testdata/fixtures directory
# for use in git pkg tests.

err_exit() {
    echo "$@" 1>&2
    exit 1
}


#------------------------------------------------------------------------------
# VARIABLES
#------------------------------------------------------------------------------
SCRIPT_ROOT=$(dirname "$(realpath "$0")") || err_exit "Failed to get script root"
REPO_NAME="git-repo"
REPO_ROOT="${SCRIPT_ROOT}/${REPO_NAME}"

COMMIT_HASHES_GO_FILENAME="commit-hashes.go"
COMMIT_HASHES_GO_FILEPATH="${SCRIPT_ROOT}/${COMMIT_HASHES_GO_FILENAME}"

# String variables for writing out the new commit-hashes.go file.
commit_history_from_test_fixture_code_comment=''
commit_hash_constants="// hashes from the git-repo fixture (in order of commits)\n"

#------------------------------------------------------------------------------
# FUNCTIONS
#------------------------------------------------------------------------------
add_file() {
    local file
    file="$1"
    git add "$file"  || err_exit "Failed to add '$file' to git repo"
    commit_str+="Updating $file, "
}

rename_file() {
    local filename1 filename2
    filename1="$1"
    filename2="$2"
    git mv "$filename1" "$filename2"  || err_exit "Failed to rename '$filename1' to '$filename2'"
    commit_str+="Renaming $filename1 -> $filename2, "
}

rm_file() {
    local file eadid
    file="$1"
    git rm "$file"  || err_exit "Failed to rm '$file' from git repo"
    eadid=$(echo "$file" | cut -d/ -f2 | cut -d\. -f1)
    commit_str+="Deleting file ${1} EADID='${eadid}', "
}

strip_commit_str_trailing_comma_space() {
    commit_str=$(echo "$commit_str" | sed -e 's/, $//')
}

update_commit_hash_go_file_variables() {
    current_commit=$(git rev-parse HEAD)

    commit_history_from_test_fixture_code_comment="\t$current_commit $commit_str\n$commit_history_from_test_fixture_code_comment"

    if [ "$#" -gt 0 ]; then
        commit_hash_constants="${commit_hash_constants}const ${1} = \"$current_commit\"\n"
    fi
}

#------------------------------------------------------------------------------
# MAIN
#------------------------------------------------------------------------------
if [[ -d "$REPO_ROOT" ]]; then
    err_exit "'$REPO_ROOT' directory already exists. Please remove it before running this script."
fi

if [[ -f "$COMMIT_HASHES_GO_FILEPATH" ]]; then
    err_exit "'$COMMIT_HASHES_GO_FILEPATH' file already exists. Please remove it before running this script."
fi

echo "------------------------------------------------------------------------------"
echo "creating directory hierarchy and test files"
echo "------------------------------------------------------------------------------"
mkdir -p "$REPO_ROOT/.circleci" "$REPO_ROOT/archives" "$REPO_ROOT/fales" "$REPO_ROOT/tamwag"

pushd "$REPO_ROOT" &>/dev/null || err_exit "Failed to change directory to ${REPO_ROOT}"
echo 'README.md' > README.md
popd &>/dev/null || err_exit "Failed to popd after creating README.md file"

pushd "$REPO_ROOT/.circleci" &>/dev/null || err_exit "Failed to change directory to ${REPO_ROOT}/.circleci"
echo 'config.yml' > config.yml
popd &>/dev/null || err_exit "Failed to popd after creating .circleci/config.yml file"

pushd "$REPO_ROOT/archives" &>/dev/null || err_exit "Failed to change directory to ${REPO_ROOT}/archives"
for e in 'mc_1' 'cap_1' ; do
    echo "$e" > "${e}.xml"
done
popd &>/dev/null || err_exit "Failed to popd after creating archives files"

pushd "$REPO_ROOT/fales" &>/dev/null || err_exit "Failed to change directory to ${REPO_ROOT}/fales"
for i in {1..5}; do
    echo "mss_00${i}" > "mss_00${i}.xml"
done
popd &>/dev/null || err_exit "Failed to popd after creating fales files"

pushd "$REPO_ROOT/tamwag" &>/dev/null || err_exit "Failed to change directory to ${REPO_ROOT}/tamwag"
for i in {1..2}; do
    echo "aia_00${i}" > "aia_00${i}.xml"
done
popd &>/dev/null || err_exit "Failed to popd after creating tamwag files"

pushd "$REPO_ROOT" &>/dev/null || err_exit "Failed to change directory to ${REPO_ROOT}"

echo "------------------------------------------------------------------------------"
echo "setting up git repository"
echo "------------------------------------------------------------------------------"
git init .

# It is required that this test repo initial commit have a mix of at least one
# EAD file to index and one non-EAD file, with at least one EAD file sorting
# lexicographically after at least one non-EAD file.
# (Joe established that `go-git` sorts files affected in a commit
# lexicographically-- see this Jira comment:
# https://nyu.atlassian.net/browse/DLFA-222?focusedCommentId=52815).
# The reason for this requirement is that when this initial test repo commit was
# originally added for https://nyu.atlassian.net/browse/DLFA-302, it had only
# the README.md and .circleci/config.yml file in it, and the appropriate test
# was added to ensure that commits with no EAD files to index did not get
# processed.  However, the test for the bugfix itself had a bug:
# https://github.com/NYULibraries/go-ead-indexer/blob/8e4495f8130f9d155d642ffa4cc2dba935c277e5/pkg/git/git.go#L143
# ...which only manifested an error if the initial test repo commit contained an
# EAD file to index after a non-EAD file, which unfortunately wasn't the case.
# The bug has been fixed and we have updated this test repo initial commit
# accordingly.  This condition must be maintained going forward.
# For details, see this comment in the bug ticket:
# https://nyu.atlassian.net/browse/DLFA-302?focusedCommentId=222897
for f in fales/mss_001.xml README.md .circleci/config.yml; do
    add_file "$f"
done
# Override the `commit_str` processing done by `add_file`.
commit_str='Initial commit of fales/mss_001.xml, README.md, and .circle/config.yml'
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables Commit1Hash

commit_str=""
for f in fales/mss_002.xml fales/mss_003.xml; do
    add_file "$f"
done
strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables Commit2Hash

commit_str=""
for f in archives/cap_1.xml fales/mss_004.xml tamwag/aia_001.xml; do
    add_file "$f"
done
strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables Commit3Hash

commit_str=""
add_file archives/mc_1.xml
rm_file fales/mss_002.xml
add_file fales/mss_005.xml
add_file tamwag/aia_002.xml
strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables Commit4Hash

commit_str=""
echo "mss_001 update" > fales/mss_001.xml
add_file fales/mss_001.xml
strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables Commit5Hash

echo 'README.md update' > README.md
echo 'config.yml update' > .circleci/config.yml
for f in README.md .circleci/config.yml; do
    add_file "$f"
done
# Override the `commit_str` processing done by `add_file`.
commit_str='Updating README.md with [whatever] and .circleci/config.yml with [whatever]'
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables Commit6Hash

echo 'README.md [had to do something special to archives/mc_1.xml' > README.md
add_file README.md
echo 'Do something special to mc_1.xml' > archives/mc_1.xml
add_file archives/mc_1.xml
# Override the `commit_str` processing done by `add_file`.
commit_str='[Do something special to] archives/mc_1.xml and add note about it to README.md'
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables Commit7Hash

commit_str=""
rename_file archives/cap_1.xml archives/cap_001.xml
rename_file archives/mc_1.xml archives/mc_001.xml
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables Commit8Hash

# Rename an EAD file with new suffix to prevent it from being recognized as an
# EAD file that needs to be indexed.
commit_str=""
rename_file archives/cap_001.xml archives/cap_001.xml.temporarily-disabled
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables Commit9Hash

# Restore the original name of the temporarily disabled EAD file.
commit_str=""
rename_file archives/cap_001.xml.temporarily-disabled archives/cap_001.xml
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables Commit10Hash

# Need to do this to prevent https://jira.nyu.edu/browse/DLFA-276 bug:
# "`git.CheckoutMergeReset` will silently check out a default commit if `commitHash` is not a valid commit hash string"
echo "------------------------------------------------------------------------------"
echo "setting branch name to 'master' (see https://jira.nyu.edu/browse/DLFA-276)"
echo "------------------------------------------------------------------------------"
git branch -m master || err_exit "problem renaming branch to master"

# generate log information for the developer to use in updating tests:
echo "------------------------------------------------------------------------------"
echo "listing commit history so that hashes can be used in tests"
echo "------------------------------------------------------------------------------"
git log --pretty=format:"%H %ad | %s%d" --date=iso
echo "------------------------------------------------------------------------------"

popd &>/dev/null || err_exit "Failed to popd after git operations"

echo "------------------------------------------------------------------------------"
echo "renaming .git to dot-git"
echo "------------------------------------------------------------------------------"
# NOTE: you MUST include the trailing /. or the .git directory will not be included in the tarball
mv -nv "$REPO_ROOT/.git" "$REPO_ROOT/dot-git" &>/dev/null || err_exit "Failed to rename ${REPO_ROOT}/.git to ${REPO_ROOT}/dot-git"

echo "------------------------------------------------------------------------------"
echo "NEXT STEPS:"
echo "1. move ${REPO_ROOT} to pkg/git/testdata/fixtures"
echo "2. Update the git pkg test scenarios with the new commit hash values"
echo "3. Run the git pkg tests"
echo "------------------------------------------------------------------------------"

echo "------------------------------------------------------------------------------"
echo "updating $commit_history_from_test_fixture_code_comment"
echo "------------------------------------------------------------------------------"
cat << EOF > $COMMIT_HASHES_GO_FILEPATH
// Code generated by pkg/index/testsupport/gen-repo.bash. DO NOT EDIT.

package git

// ------------------------------------------------------------------------------
// git repo fixture constants used by pkg/git tests
// ------------------------------------------------------------------------------

/*
	# Commit history from test fixture
EOF

echo -en "$commit_history_from_test_fixture_code_comment*/\n\n" >> $COMMIT_HASHES_GO_FILEPATH

echo -en "${commit_hash_constants}" >> $COMMIT_HASHES_GO_FILEPATH

exit 0
