#!/bin/bash
set -uo pipefail

err_exit() {
    echo "$@" 1>&2
    exit 1
}


#------------------------------------------------------------------------------
# VARIABLES
#------------------------------------------------------------------------------

COMMIT_HASHES_GO_FILENAME=commit-hashes.go
GIT_REPO_DIRNAME=git-repo

SCRIPT_ROOT=$(dirname "$(realpath "$0")") || err_exit "Failed to get script root"
GENERATE_REPO_SCRIPT="${SCRIPT_ROOT}/gen-repo.bash"

SOURCE_COMMIT_HASHES_GO_FILEPATH="${SCRIPT_ROOT}/${COMMIT_HASHES_GO_FILENAME}"
TARGET_COMMIT_HASHES_GO_FILEPATH="${SCRIPT_ROOT}/../${COMMIT_HASHES_GO_FILENAME}"

SOURCE_GIT_REPO_DIRPATH="${SCRIPT_ROOT}/${GIT_REPO_DIRNAME}"
TARGET_GIT_REPO_DIRPATH="${SCRIPT_ROOT}/../testdata/fixtures/${GIT_REPO_DIRNAME}"


#------------------------------------------------------------------------------
# MAIN
#------------------------------------------------------------------------------

# Clean out any previously generated fixtures that haven't been moved into place.
# The generate repo script will abort if this isn't done.
rm -fr $SOURCE_COMMIT_HASHES_GO_FILEPATH $SOURCE_GIT_REPO_DIRPATH || \
    err_exit "Failed to delete: ${SOURCE_COMMIT_HASHES_GO_FILEPATH} ${SOURCE_GIT_REPO_DIRPATH}"

$GENERATE_REPO_SCRIPT || err_exit "Failed to run ${GENERATE_REPO_SCRIPT}"

# Clean out current fixtures.
rm -fr $TARGET_COMMIT_HASHES_GO_FILEPATH $TARGET_GIT_REPO_DIRPATH || \
    err_exit "Failed to delete: ${TARGET_COMMIT_HASHES_GO_FILEPATH} ${TARGET_GIT_REPO_DIRPATH}"

# Move new fixtures into place.
mv $SOURCE_COMMIT_HASHES_GO_FILEPATH $TARGET_COMMIT_HASHES_GO_FILEPATH || \
    err_exit "Failed to move ${SOURCE_COMMIT_HASHES_GO_FILEPATH} to ${TARGET_COMMIT_HASHES_GO_FILEPATH}"
mv $SOURCE_GIT_REPO_DIRPATH $TARGET_GIT_REPO_DIRPATH || \
    err_exit "Failed to move ${SOURCE_GIT_REPO_DIRPATH} to ${TARGET_GIT_REPO_DIRPATH}"

exit 0
