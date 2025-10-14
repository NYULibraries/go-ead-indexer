#!/bin/bash
set -uo pipefail

# This script generates a git repo test fixture for use in the git package tests.
# The script creates a directory named 'git-repo', then generates and commits various files.
# Finally, the script renames the 'git-repo/.git' directory to 'git-repo/dot-git'.
# The git-repo directory can be moved into the pkg/git/testdata/fixtures directory
# for use in git pkg tests.

# Commit history replicated in repo
# 8adb38c5f05fce5ef8ef9b97b4721b5a962057ea 2025-10-06 17:50:33 -0400 | [Do something special to] archives/mc_1.xml and add note about it to README.md (HEAD -> master)
# e6af7e810b8002761077a943689529405d558697 2025-10-06 17:50:33 -0400 | Updating README.md with [whatever] and .circleci/config.yml with [whatever]
# df2bfddf4a599e4a24373320e91366df90dc708a 2025-10-06 17:50:33 -0400 | Updating file fales/mss_001.xml
# f34cf26e0a8c70511b7941921ee5016c4fcf3fce 2025-10-06 17:50:33 -0400 | Updating file archives/mc_1.xml, Deleting file fales/mss_002.xml EADID='mss_002', Updating file fales/mss_005.xml, Updating file tamwag/aia_002.xml
# 2a5cc008d17384ab183dba69190251e0503fa315 2025-10-06 17:50:33 -0400 | Updating file archives/cap_1.xml, Updating file fales/mss_004.xml, Updating file tamwag/aia_001.xml
# 80301c37ccc2998fd2a8b021a731296273d37467 2025-10-06 17:50:33 -0400 | Updating file fales/mss_002.xml, Updating file fales/mss_003.xml
# 3c20e78557fbf11e77b7fb9e551b7c1b2d508261 2025-10-06 17:50:32 -0400 | Initial commit of fales/mss_001.xml, README.md, and .circle/config.yml

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

rm_file() {
    local file eadid
    file="$1"
    git rm "$file"  || err_exit "Failed to rm '$file' from git repo"
    eadid=$(echo "$file" | cut -d/ -f2 | cut -d\. -f1)
    commit_str+="Deleting file ${1} EADID='${eadid}', "
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
git add fales/mss_001.xml README.md .circleci/config.yml
git commit -m "Initial commit of fales/mss_001.xml, README.md, and .circle/config.yml"

git add fales/mss_001.xml
git commit -m "Updating file fales/mss_001.xml"

git add fales/mss_002.xml fales/mss_003.xml
git commit -m "Updating file fales/mss_002.xml, Updating file fales/mss_003.xml"

git add archives/cap_1.xml fales/mss_004.xml tamwag/aia_001.xml
git commit -m "Updating file archives/cap_1.xml, Updating file fales/mss_004.xml, Updating file tamwag/aia_001.xml"

git add archives/mc_1.xml
git rm fales/mss_002.xml
git add fales/mss_005.xml
git add tamwag/aia_002.xml
git commit -m "Updating file archives/mc_1.xml, Deleting file fales/mss_002.xml EADID='mss_002', Updating file fales/mss_005.xml, Updating file tamwag/aia_002.xml"

echo "mss_001 update" > fales/mss_001.xml
git add fales/mss_001.xml
git commit -m 'Updating file fales/mss_001.xml'

echo 'README.md update' > README.md
echo 'config.yml update' > .circleci/config.yml
git add README.md
git add .circleci/config.yml
git commit -m 'Updating README.md with [whatever] and .circleci/config.yml with [whatever]'

echo 'README.md [had to do something special to archives/mc_1.xml' > README.md
git add README.md
echo 'Do something special to mc_1.xml' > archives/mc_1.xml
git add archives/mc_1.xml
git commit -m '[Do something special to] archives/mc_1.xml and add note about it to README.md'

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

exit 0
