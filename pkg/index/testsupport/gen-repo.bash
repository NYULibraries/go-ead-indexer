#!/bin/bash
set -uo pipefail

# This script generates a git repo test fixture for use in the pkg/index IndexGitCommit() tests.
# (Please see https://jira.nyu.edu/browse/DLFA-230 for details.)
#
# The script creates a directory named 'git-repo', then generates and commits various files.
# Finally, the script renames the 'git-repo/.git' directory to 'git-repo/dot-git'.
# The git-repo directory can be moved into the pkg/git/testdata/fixtures directory
# for use in git pkg tests.

# Commit history replicated in repo (NOTE: commit hashes WILL differ)
# 00fd44a8e69285cf3789be3e7bc0e4e88d5f6dd8 2025-03-26 13:18:07 -0400 | Updating nyuad/ad_mc_019.xml, Deleting file tamwag/tam_143.xml EADID='tam_143', Updating edip/mos_2024.xml, Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Updating akkasah/ad_mc_030.xml (HEAD -> main) [jgpawletko]
# 5fec61740cb7e4f05bbfa77548b42be2003e278b 2025-03-26 13:18:07 -0400 | Updating tamwag/tam_143.xml, Updating cbh/arc_212_plymouth_beecher.xml [jgpawletko]
# d8144b3136ef4a9abf0613a1302606644f90bd6c 2025-03-26 13:18:07 -0400 | Deleting file fales/mss_420.xml EADID='mss_420', Updating fales/mss_420.xml [jgpawletko]
# 0afcf14e99bbd6e158f486090877fbd50370494c 2025-03-26 13:18:07 -0400 | Updating fales/mss_420.xml [jgpawletko]
# aee0af16b6d92444326eea4847893844f3ca59ae 2025-03-26 13:18:07 -0400 | Deleting file fales/mss_460.xml EADID='mss_460' [jgpawletko]
# e24960b1dc934c628d4475cb4537f7e21f54032c 2025-03-26 13:18:07 -0400 | Updating fales/mss_460.xml [jgpawletko]
# 0fcdd54abaeb3b2f15b50f8eb5ef903ba2231896 2025-03-26 13:18:07 -0400 | Deleting file akkasah/ad_mc_030.xml EADID='ad_mc_030', Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Deleting file edip/mos_2024.xml EADID='mos_2024', Deleting file fales/mss_420.xml EADID='mss_420', Deleting file fales/mss_460.xml EADID='mss_460', Deleting file nyhs/ms256_harmon_hendricks_goldstone.xml EADID='ms256_harmon_hendricks_goldstone', Deleting file nyhs/ms347_foundling_hospital.xml EADID='ms347_foundling_hospital', Deleting file nyuad/ad_mc_019.xml EADID='ad_mc_019', Deleting file tamwag/tam_143.xml EADID='tam_143' [jgpawletko]
# 7fdb03f4ab09f0eddf9b3c0e77ba50f5d036b2e9 2025-03-26 13:18:07 -0400 | Updating akkasah/ad_mc_030.xml, Updating cbh/arc_212_plymouth_beecher.xml, Updating edip/mos_2024.xml, Updating fales/mss_420.xml, Updating fales/mss_460.xml, Updating nyhs/ms256_harmon_hendricks_goldstone.xml, Updating nyhs/ms347_foundling_hospital.xml, Updating nyuad/ad_mc_019.xml, Updating tamwag/tam_143.xml [jgpawletko]

err_exit() {
    echo "$@" 1>&2
    exit 1
}


# Plan:
# create git-repo directory
# init git repo
# copy     all
#   add    all
#   delete all
# copy     1
#   add    1
#   delete 1
# copy     A
#   add    A
#   delete A
#   recopy A
#   modify A
#   add    A
# copy     2
#   add    2
# copy     3
#   add 1, delete 1, add 1, delete 1, add 1
# copy     A
#   add    A
#   delete A
#   recopy A
#   modify A
#   add    A
#

#------------------------------------------------------------------------------
# VARIABLES
#------------------------------------------------------------------------------
SCRIPT_ROOT=$(dirname "$(realpath "$0")") || err_exit "Failed to get script root"
REPO_NAME="git-repo"
REPO_ROOT="${SCRIPT_ROOT}/${REPO_NAME}"
EAD_FILE_ROOT="${SCRIPT_ROOT}/../../ead/testdata/fixtures/ead-files"
EAD_FILES='akkasah/ad_mc_030.xml
cbh/arc_212_plymouth_beecher.xml
edip/mos_2024.xml
fales/mss_420.xml
fales/mss_460.xml
nyhs/ms256_harmon_hendricks_goldstone.xml
nyhs/ms347_foundling_hospital.xml
nyuad/ad_mc_019.xml
tamwag/tam_143.xml
'

COMMIT_HASHES_GO_FILENAME="commit-hashes.go"
COMMIT_HASHES_GO_FILEPATH="${SCRIPT_ROOT}/${COMMIT_HASHES_GO_FILENAME}"

# EAD files for various scenarios, committed in reverse alphabetical order
ADD_DELETE_MODIFY_ADD_ONE_EAD='fales/mss_420.xml'
ADD_AND_DELETE_ONE_EAD='fales/mss_460.xml'
ADD_AND_DELETE_TWO_EADS_NUM_01='tamwag/tam_143.xml'
ADD_AND_DELETE_TWO_EADS_NUM_02='cbh/arc_212_plymouth_beecher.xml'
ADD_THREE_EADS_NUM_01='nyuad/ad_mc_019.xml'
ADD_THREE_EADS_NUM_02='edip/mos_2024.xml'
ADD_THREE_EADS_NUM_03='akkasah/ad_mc_030.xml'

# Non-EAD file
README='README.md'

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

cp_ead() {
    local f

    f="$1"
    dir=$(dirname "$f")
    tgt_dir="${REPO_ROOT}/${dir}"
    mkdir -p "${tgt_dir}" || err_exit "problem creating '$tgt_dir'"
    cp "${EAD_FILE_ROOT}/${f}"  "${tgt_dir}" || err_exit "problem copying EAD '$f'"
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

if [[ ! -d "$EAD_FILE_ROOT" ]]; then
    err_exit "EAD directory '$EAD_FILE_ROOT' does not exist. Please locate or set up the EAD directory hierarchy."
fi


echo "------------------------------------------------------------------------------"
echo "setting up git repository"
echo "------------------------------------------------------------------------------"
mkdir -pv "$REPO_ROOT" || err_exit "Failed to create git-repo directory"
pushd "$REPO_ROOT" &>/dev/null || err_exit "Failed to change directory to ${REPO_ROOT}"
git init . || err_exit "Failed to init git repo"


echo "------------------------------------------------------------------------------"
echo "setting up 'add all' commit"
echo "------------------------------------------------------------------------------"
commit_str=""
for f in $EAD_FILES; do
    cp_ead "$f"
    add_file "$f"
done
strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables AddAllHash

echo "------------------------------------------------------------------------------"
echo "setting up 'delete all' commit"
echo "------------------------------------------------------------------------------"
commit_str=""
for f in $EAD_FILES; do
    rm_file "$f"
done
strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables DeleteAllHash


echo "------------------------------------------------------------------------------"
echo "setting up 'add one' and 'delete one' commits"
echo "------------------------------------------------------------------------------"
commit_str=""
cp_ead "$ADD_AND_DELETE_ONE_EAD"
add_file "$ADD_AND_DELETE_ONE_EAD"
strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables AddOneHash

commit_str=""
rm_file "$ADD_AND_DELETE_ONE_EAD"
strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables DeleteOneHash


echo "------------------------------------------------------------------------------"
echo "setting up 'add A', 'delete A', 'modify A', 'add A' commit"
echo "------------------------------------------------------------------------------"
# add the file to the repo
commit_str=""
cp_ead "$ADD_DELETE_MODIFY_ADD_ONE_EAD"
add_file "$ADD_DELETE_MODIFY_ADD_ONE_EAD"
strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
# Note there is no `const` declaration for this commit.
update_commit_hash_go_file_variables

# git rm the file but do not commit yet
commit_str=""
rm_file "$ADD_DELETE_MODIFY_ADD_ONE_EAD"
# copy, modify, and add the file back to the repo
cp_ead "$ADD_DELETE_MODIFY_ADD_ONE_EAD"
echo "   " >> "${REPO_ROOT}/${ADD_DELETE_MODIFY_ADD_ONE_EAD}"
add_file "$ADD_DELETE_MODIFY_ADD_ONE_EAD"
# commit the changes
strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables DeleteModifyAddHash


echo "------------------------------------------------------------------------------"
echo "setting up 'add two' commit"
echo "------------------------------------------------------------------------------"
commit_str=""
for f in $ADD_AND_DELETE_TWO_EADS_NUM_01 $ADD_AND_DELETE_TWO_EADS_NUM_02; do
    cp_ead "$f"
    add_file "$f"
done

strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables AddTwoHash


echo "------------------------------------------------------------------------------"
echo "setting up 'add three and delete two' commit"
echo "------------------------------------------------------------------------------"
commit_str=""
# INTERLEAVE ADD AND DELETE OPERATIONS
cp_ead "$ADD_THREE_EADS_NUM_01"
add_file "$ADD_THREE_EADS_NUM_01"

rm_file "$ADD_AND_DELETE_TWO_EADS_NUM_01"

cp_ead "$ADD_THREE_EADS_NUM_02"
add_file "$ADD_THREE_EADS_NUM_02"

rm_file "$ADD_AND_DELETE_TWO_EADS_NUM_02"

cp_ead "$ADD_THREE_EADS_NUM_03"
add_file "$ADD_THREE_EADS_NUM_03"

strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables AddThreeDeleteTwoHash


echo "------------------------------------------------------------------------------"
echo "setting up 'no EAD files to index'"
echo "------------------------------------------------------------------------------"
commit_str=""

echo 'README.md' > $README
add_file "$README"

strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"
update_commit_hash_go_file_variables NoEADFilesInCommitHash


# Need to do this to prevent https://jira.nyu.edu/browse/DLFA-276 bug:
# "`git\.CheckoutMergeReset` will silently check out a default commit if `commitHash` is not a valid commit hash string"
echo "------------------------------------------------------------------------------"
echo "setting branch name to 'master' (see https://jira.nyu.edu/browse/DLFA-276)"
echo "------------------------------------------------------------------------------"
git branch -m master || err_exit "problem renaming branch to master"


# generate log information for the developer to use in updating tests:
echo "------------------------------------------------------------------------------"
echo "listing commit history so that hashes can be used in tests"
echo "------------------------------------------------------------------------------"
git log --pretty=format:"%H %ad | %s%d [%an]" --date=iso
echo "------------------------------------------------------------------------------"

popd &>/dev/null || err_exit "Failed to popd after git operations"

echo "------------------------------------------------------------------------------"
echo "renaming .git to dot-git"
echo "------------------------------------------------------------------------------"
mv -nv "${REPO_ROOT}/.git" "${REPO_ROOT}/dot-git" &>/dev/null || err_exit "Failed to rename git-repo/.git to git-repo/dot-git"

echo "------------------------------------------------------------------------------"
echo "NEXT STEPS:"
echo "1. move ${REPO_NAME} to pkg/index/testdata/fixtures/"
echo "2. move ${COMMIT_HASHES_GO_FILENAME} to pkg/index/testutils/"
echo "3. Run the git pkg tests"
echo "------------------------------------------------------------------------------"

echo "------------------------------------------------------------------------------"
echo "updating $commit_history_from_test_fixture_code_comment"
echo "------------------------------------------------------------------------------"
cat << EOF > $COMMIT_HASHES_GO_FILEPATH
package testutils

// ------------------------------------------------------------------------------
// git repo fixture constants shared by cmd/index and pkg/index tests
// ------------------------------------------------------------------------------

/*
	# Commit history from test fixture
EOF

echo -en "$commit_history_from_test_fixture_code_comment*/\n\n" >> $COMMIT_HASHES_GO_FILEPATH

echo -en "${commit_hash_constants}\n" >> $COMMIT_HASHES_GO_FILEPATH

exit 0
