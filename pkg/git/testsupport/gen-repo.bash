#!/bin/bash

# This script generates a git repo test fixture for use in the git package tests.
# The script creates a directory named 'simple-repo', then generates and commits various files.
# Finally, the script creates a tarball named 'simple-repo.tar.gz'.
# The tarball can be copied into the pkg/git/testdata directory for use in git pkg tests.

# Commit history replicated in repo (NOTE: commit hashes may differ)
# * 95f2f904ad261e7d31632021fa10768d2b4096c9 2025-01-24 17:10:44 -0500 | Updating file fales/mss_001.xml (HEAD -> main) [jgpawletko]
# * aa58b2314e11ae5af61129ebfe1ceb07b49c2d33 2025-01-24 17:10:44 -0500 | Updating file archives/mc_1.xml, Deleting file fales/mss_002.xml EADID='mss_002', Updating file fales/mss_005.xml, Updating file tamwag/aia_002.xml [jgpawletko]
# * 3dc6fabe0fcd990e95cdd3f88cff821196fccdbd 2025-01-24 17:10:44 -0500 | Updating file archives/cap_1.xml, Updating file fales/mss_004.xml, Updating file tamwag/aia_001.xml [jgpawletko]
# * 7fe6de7c56d30149889f8d24eaf2fa66ed9f2e2d 2025-01-24 17:10:44 -0500 | Updating file fales/mss_002.xml, Updating file fales/mss_003.xml [jgpawletko]
# * 155309f674b5acffd7473c1648f3647a2a3d242b 2025-01-24 17:10:44 -0500 | Updating file fales/mss_001.xml [jgpawletko]

err_exit() {
    echo "$@" 1>&2
    exit 1
}

if [[ -d simple-repo ]]; then
    err_exit "simple-repo directory already exists. Please remove it before running this script."
fi

mkdir -p simple-repo/archives simple-repo/fales simple-repo/tamwag
pushd simple-repo/archives || err_exit "Failed to change directory to simple-repo/archives"
for e in 'mc_1' 'cap_1' ; do
    echo "$e" > "${e}.xml"
done
popd || err_exit "Failed to popd after creating archives files"

pushd simple-repo/fales || err_exit "Failed to change directory to simple-repo/fales"
for i in {1..5}; do
    echo "mss_00${i}" > "mss_00${i}.xml"
done
popd || err_exit "Failed to popd after creating fales files"

pushd simple-repo/tamwag || err_exit "Failed to change directory to simple-repo/tamwag"
for i in {1..2}; do
    echo "aia_00${i}" > "aia_00${i}.xml"
done
popd || err_exit "Failed to popd after creating tamwag files"

pushd simple-repo || err_exit "Failed to change directory to simple-repo"
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

popd || err_exit "Failed to popd after git operations"

# NOTE: you MUST include the trailing /. or the .git directory will not be included in the tarball
tar cvfz simple-repo.tar.gz simple-repo/. 
rm -rf simple-repo

exit 0
