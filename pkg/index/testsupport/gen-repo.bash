#!/bin/bash
set -uo pipefail
# This script generates a git repo test fixture for use in the git package tests.
# The script creates a directory named 'git-repo', then generates and commits various files.
# Finally, the script renames the 'git-repo/.git' directory to 'git-repo/dot-git'.
# The git-repo directory can be moved into the pkg/git/testdata/fixtures directory 
# for use in git pkg tests.

# Commit history replicated in repo (NOTE: commit hashes WILL differ)
# * 95f2f904ad261e7d31632021fa10768d2b4096c9 2025-01-24 17:10:44 -0500 | Updating file fales/mss_001.xml (HEAD -> main) [jgpawletko]
# * aa58b2314e11ae5af61129ebfe1ceb07b49c2d33 2025-01-24 17:10:44 -0500 | Updating file archives/mc_1.xml, Deleting file fales/mss_002.xml EADID='mss_002', Updating file fales/mss_005.xml, Updating file tamwag/aia_002.xml [jgpawletko]
# * 3dc6fabe0fcd990e95cdd3f88cff821196fccdbd 2025-01-24 17:10:44 -0500 | Updating file archives/cap_1.xml, Updating file fales/mss_004.xml, Updating file tamwag/aia_001.xml [jgpawletko]
# * 7fe6de7c56d30149889f8d24eaf2fa66ed9f2e2d 2025-01-24 17:10:44 -0500 | Updating file fales/mss_002.xml, Updating file fales/mss_003.xml [jgpawletko]
# * 155309f674b5acffd7473c1648f3647a2a3d242b 2025-01-24 17:10:44 -0500 | Updating file fales/mss_001.xml [jgpawletko]

err_exit() {
    echo "$@" 1>&2
    exit 1
}

#------------------------------------------------------------------------------
# VARIABLES
#------------------------------------------------------------------------------
SCRIPT_ROOT=$(dirname "$(realpath "$0")") || err_exit "Failed to get script root"
REPO_ROOT="${SCRIPT_ROOT}/git-repo"
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
ADD_AND_DELETE_TWO_EADS='nyuad/ad_mc_019.xml tamwag/tam_143.xml'
ADD_THREE_EADS='akkasah/ad_mc_030.xml cbh/arc_212_plymouth_beecher.xml edip/mos_2024.xml'

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
    pwd
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

# set EAD fixture root
# set directories
# copy     all
#   add    all
#   delete all
# copy     1
#   add    1
#   delete 1
# copy     2
#   add    2
# copy     3
#   add 3, delete 2
# 

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
for f in $ADD_AND_DELETE_TWO_EADS; do
    cp_ead "$f"
    add_file "$f"
done

strip_commit_str_trailing_comma_space
git commit -m "$commit_str" || err_exit "problem committing: $commit_str"


echo "------------------------------------------------------------------------------"
echo "setting up 'add three and delete two' commit"
echo "------------------------------------------------------------------------------"
commit_str=""
# process "add three" files
for f in $ADD_THREE_EADS; do
    cp_ead "$f"
    add_file "$f"
done
# process "delete two" files
for f in $ADD_AND_DELETE_TWO_EADS; do
    rm_file "$f"
done
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
echo "1. move git-repo to pkg/git/testdata/fixtures"
echo "2. Update the git pkg test scenarios with the new commit hash values"
echo "3. Run the git pkg tests"
echo "------------------------------------------------------------------------------"

exit 0

# git add fales/mss_001.xml
# git commit -m "Updating file fales/mss_001.xml"

# git add fales/mss_002.xml fales/mss_003.xml
# git commit -m "Updating file fales/mss_002.xml, Updating file fales/mss_003.xml"

# git add archives/cap_1.xml fales/mss_004.xml tamwag/aia_001.xml
# git commit -m "Updating file archives/cap_1.xml, Updating file fales/mss_004.xml, Updating file tamwag/aia_001.xml"

# git add archives/mc_1.xml
# git rm fales/mss_002.xml
# git add fales/mss_005.xml
# git add tamwag/aia_002.xml
# git commit -m "Updating file archives/mc_1.xml, Deleting file fales/mss_002.xml EADID='mss_002', Updating file fales/mss_005.xml, Updating file tamwag/aia_002.xml"

# echo "mss_001 update" > fales/mss_001.xml
# git add fales/mss_001.xml
# git commit -m 'Updating file fales/mss_001.xml'


# echo "------------------------------------------------------------------------------"
# echo "creating directory hierarchy and test files"
# echo "------------------------------------------------------------------------------"
# mkdir -p git-repo/archives git-repo/fales git-repo/tamwag
# pushd git-repo/archives &>/dev/null || err_exit "Failed to change directory to git-repo/archives"
# for e in 'mc_1' 'cap_1' ; do
#     echo "$e" > "${e}.xml"
# done
# popd &>/dev/null || err_exit "Failed to popd after creating archives files"

# pushd git-repo/fales &>/dev/null || err_exit "Failed to change directory to git-repo/fales"
# for i in {1..5}; do
#     echo "mss_00${i}" > "mss_00${i}.xml"
# done
# popd &>/dev/null || err_exit "Failed to popd after creating fales files"

# pushd git-repo/tamwag &>/dev/null || err_exit "Failed to change directory to git-repo/tamwag"
# for i in {1..2}; do
#     echo "aia_00${i}" > "aia_00${i}.xml"
# done
# popd &>/dev/null || err_exit "Failed to popd after creating tamwag files"

# create directory hierarchy
# echo "------------------------------------------------------------------------------"
# echo "creating directory hierarchy"
# echo "------------------------------------------------------------------------------"
# for f in $EAD_FILES; do
#     dir=$(dirname "$f")
#     tgt_dir="${REPO_ROOT}/${dir}"
#     mkdir -p "${tgt_dir}" || err_exit "problem creating '$tgt_dir'"
# done

