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
# e5c5336b63b109b68c495bbfea94d30ecbc1ef67 2025-03-24 19:53:31 -0400 | Updating akkasah/ad_mc_030.xml, Updating cbh/arc_212_plymouth_beecher.xml, Updating edip/mos_2024.xml, Deleting file nyuad/ad_mc_019.xml EADID='ad_mc_019', Deleting file tamwag/tam_143.xml EADID='tam_143' (HEAD -> main) [jgpawletko]
# e4bfc536020c4477044633ac7a57242bb6f67cee 2025-03-24 19:53:30 -0400 | Updating nyuad/ad_mc_019.xml, Updating tamwag/tam_143.xml [jgpawletko]
# 2fee15ffc217a86d19756a6c816f59ca86e23893 2025-03-24 19:53:30 -0400 | Deleting file fales/mss_460.xml EADID='mss_460' [jgpawletko]
# fdd7ce5e54b88894460b52dd0dd27055ffb3bbdd 2025-03-24 19:53:30 -0400 | Updating fales/mss_460.xml [jgpawletko]
# e4fe6008decb5f26382fae903de40a4f3470d509 2025-03-24 19:53:30 -0400 | Deleting file akkasah/ad_mc_030.xml EADID='ad_mc_030', Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Deleting file edip/mos_2024.xml EADID='mos_2024', Deleting file fales/mss_420.xml EADID='mss_420', Deleting file fales/mss_460.xml EADID='mss_460', Deleting file nyhs/ms256_harmon_hendricks_goldstone.xml EADID='ms256_harmon_hendricks_goldstone', Deleting file nyhs/ms347_foundling_hospital.xml EADID='ms347_foundling_hospital', Deleting file nyuad/ad_mc_019.xml EADID='ad_mc_019', Deleting file tamwag/tam_143.xml EADID='tam_143' [jgpawletko]
# 5546ffda27581c4933aeb4102f6a0107c3e522ff 2025-03-24 19:53:30 -0400 | Updating akkasah/ad_mc_030.xml, Updating cbh/arc_212_plymouth_beecher.xml, Updating edip/mos_2024.xml, Updating fales/mss_420.xml, Updating fales/mss_460.xml, Updating nyhs/ms256_harmon_hendricks_goldstone.xml, Updating nyhs/ms347_foundling_hospital.xml, Updating nyuad/ad_mc_019.xml, Updating tamwag/tam_143.xml [jgpawletko]

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
# copy     2
#   add    2
# copy     3
#   add 1, delete 1, add 1, delete 1, add 1
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
ADD_AND_DELETE_ONE_EAD='fales/mss_460.xml'
ADD_AND_DELETE_TWO_EADS_NUM_01='nyuad/ad_mc_019.xml'
ADD_AND_DELETE_TWO_EADS_NUM_02='tamwag/tam_143.xml'
ADD_THREE_EADS_NUM_01='akkasah/ad_mc_030.xml'
ADD_THREE_EADS_NUM_02='cbh/arc_212_plymouth_beecher.xml'
ADD_THREE_EADS_NUM_03='edip/mos_2024.xml'

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

#------------------------------------------------------------------------------
# MAIN
#------------------------------------------------------------------------------
if [[ -d "$REPO_ROOT" ]]; then
    err_exit "'$REPO_ROOT' directory already exists. Please remove it before running this script."
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


echo "------------------------------------------------------------------------------"
echo "setting up 'delete all' commit"
echo "------------------------------------------------------------------------------"
commit_str=""
for f in $EAD_FILES; do
    rm_file "$f"
done
strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"


echo "------------------------------------------------------------------------------"
echo "setting up 'add one' and 'delete one' commits"
echo "------------------------------------------------------------------------------"
commit_str=""
cp_ead "$ADD_AND_DELETE_ONE_EAD"
add_file "$ADD_AND_DELETE_ONE_EAD"
strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"

commit_str=""
rm_file "$ADD_AND_DELETE_ONE_EAD"
strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"


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
echo "1. move ${REPO_NAME} to pkg/index/testdata/fixtures"
echo "2. Update the git pkg test scenarios with the new commit hash values"
echo "3. Run the git pkg tests"
echo "------------------------------------------------------------------------------"

exit 0
