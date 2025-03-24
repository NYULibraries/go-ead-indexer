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

err_exit() {
    echo "$@" 1>&2
    exit 1
}

if [[ -d git-repo ]]; then
    err_exit "'git-repo' directory already exists. Please remove it before running this script."
fi

echo "------------------------------------------------------------------------------"
echo "creating directory hierarchy and test files"
echo "------------------------------------------------------------------------------"
mkdir -p git-repo/archives git-repo/fales git-repo/tamwag
pushd git-repo/archives &>/dev/null || err_exit "Failed to change directory to git-repo/archives"
for e in 'mc_1' 'cap_1' ; do
    echo "$e" > "${e}.xml"
done
popd &>/dev/null || err_exit "Failed to popd after creating archives files"

pushd git-repo/fales &>/dev/null || err_exit "Failed to change directory to git-repo/fales"
for i in {1..5}; do
    echo "mss_00${i}" > "mss_00${i}.xml"
done
popd &>/dev/null || err_exit "Failed to popd after creating fales files"

pushd git-repo/tamwag &>/dev/null || err_exit "Failed to change directory to git-repo/tamwag"
for i in {1..2}; do
    echo "aia_00${i}" > "aia_00${i}.xml"
done
popd &>/dev/null || err_exit "Failed to popd after creating tamwag files"

pushd git-repo &>/dev/null || err_exit "Failed to change directory to git-repo"

echo "------------------------------------------------------------------------------"
echo "setting up git repository"
echo "------------------------------------------------------------------------------"
git init .

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
