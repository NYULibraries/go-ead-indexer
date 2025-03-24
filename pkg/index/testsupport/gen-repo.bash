#!/bin/bash

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

SCRIPT_ROOT=$(dirname "$0")
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

REPO_ROOT="${SCRIPT_ROOT}/git-repo"

err_exit() {
    echo "$@" 1>&2
    exit 1
}

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

strip_commit_str_trailing_comma_space() {
    commit_str=$(echo "$commit_str" | sed -e 's/, $//')
}

if [[ -d "$REPO_ROOT" ]]; then
    err_exit "'$REPO_ROOT' directory already exists. Please remove it before running this script."
fi

if [[ ! -d "$EAD_FILE_ROOT" ]]; then
    err_exit "EAD directory '$EAD_FILE_ROOT' does not exist. Please locate or set up the EAD directory hierarchy."
fi

# set EAD fixture root
# set directories
# cp -rpv 
# commit all
# delete all
# copy one
# commit 1
# delete 1
# copy 2
# commit 2
# copy 3
# delete 3, commit 2
# 

# create directory hierarchy
echo "------------------------------------------------------------------------------"
echo "creating directory hierarchy"
echo "------------------------------------------------------------------------------"
for f in $EAD_FILES; do
    dir=$(dirname "$f")
    tgt_dir="${REPO_ROOT}/${dir}"
    mkdir -p "${tgt_dir}" || err_exit "problem creating '$tgt_dir'"
done

# copy files
echo "------------------------------------------------------------------------------"
echo "copying EAD files"
echo "------------------------------------------------------------------------------"
for f in $EAD_FILES; do
    cp "${EAD_FILE_ROOT}/${f}"  "${REPO_ROOT}/${f}" || err_exit "problem copying EAD '$f'"
done


echo "------------------------------------------------------------------------------"
echo "setting up git repository"
echo "------------------------------------------------------------------------------"
pushd git-repo &>/dev/null || err_exit "Failed to change directory to git-repo"

git init . || err_exit "Failed to init git repo"


# echo "------------------------------------------------------------------------------"
# echo "setting up 'add all' commit"
# echo "------------------------------------------------------------------------------"
# commit_str=""
# for f in $EAD_FILES; do
#     git add "$f"  || err_exit "Failed to add '$f' to git repo"
#     commit_str+="Updating $f, "
# done

# # strip of trailing ', '
# commit_str=$(echo "$commit_str" | sed -e 's/, $//')

# # commit
# git commit -m "$commit_str" || err_exit "problem committing: $commit_str"


# echo "------------------------------------------------------------------------------"
# echo "setting up 'delete all' commit"
# echo "------------------------------------------------------------------------------"
# commit_str=""
# for f in $EAD_FILES; do
#     git rm "$f"  || err_exit "Failed to rm '$f' to git repo"
#     eadid=$(echo "$f" | cut -d/ -f2 | cut -d\. -f1)
#     commit_str+="Deleting file ${f} EADID='${eadid}', "
# done

# # strip of trailing ', '
# commit_str=$(echo "$commit_str" | sed -e 's/, $//')

# # commit
# git commit -m "$commit_str" || err_exit "problem committing: $commit_str"


echo "setting up 'add all' commit"
echo "------------------------------------------------------------------------------"
commit_str=""
for f in $EAD_FILES; do
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

exit

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
# NOTE: you MUST include the trailing /. or the .git directory will not be included in the tarball
mv -nv git-repo/.git git-repo/dot-git &>/dev/null || err_exit "Failed to rename git-repo/.git to git-repo/dot-git"

echo "------------------------------------------------------------------------------"
echo "NEXT STEPS:"
echo "1. move git-repo to pkg/git/testdata/fixtures"
echo "2. Update the git pkg test scenarios with the new commit hash values"
echo "3. Run the git pkg tests"
echo "------------------------------------------------------------------------------"

exit 0


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

