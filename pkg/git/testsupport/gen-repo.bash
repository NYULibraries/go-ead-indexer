#!/bin/bash

# This script generates a git repo test fixture for use in the git package tests.
# The script creates a directory named 'simple-repo', then generates and commits various files
# finally, the script creates a tarball named 'simple-repo.tar.gz'
# This file can be copied into the pkg/git/testdata directory for use in git pkg tests

# Commit history replicated in repo (NOTE: commit hashes may differ)
# * ba7d9b00023a8d6ed962e46465d800265a6d06b9 2025-01-03 17:05:48 -0500 | Updating file archives/mc_1.xml, Deleting file fales/mss_002.xml EADID='mss_002', Updating file fales/mss_005.xml, Updating file tamwag/aia_002.xml (HEAD -> main) [jgpawletko]
# * ca2dd426b9f22b52e101e71ce8db83c80508df06 2025-01-03 17:00:32 -0500 | Updating file archives/cap_1.xml, Updating file fales/mss_004.xml, Updating file tamwag/aia_001.xml [jgpawletko]
# * ae25d50165da2befdfc21624ba52241ad36070de 2025-01-03 16:55:07 -0500 | Updating file fales/mss_002.xml, Updating file fales/mss_003.xml [jgpawletko]
# * d9fa76ef7c89994d8d3ed458e5c06b2c5bb9f414 2025-01-03 16:54:16 -0500 | Updating file fales/mss_001.xml [jgpawletko]

mkdir -p simple-repo/archives simple-repo/fales simple-repo/tamwag
pushd simple-repo/archives
for e in 'mc_1' 'cap_1' ; do
    echo "$e" > ${e}.xml
done
popd

pushd simple-repo/fales
for i in {1..5}; do
    echo "mss_00${i}" > mss_00${i}.xml
done
popd

pushd simple-repo/tamwag
for i in {1..2}; do
    echo "aia_00${i}" > aia_00${i}.xml
done
popd

pushd simple-repo
git init .

git add fales/mss_001.xml
git commit -m 'Updating file fales/mss_001.xml'

git add fales/mss_002.xml fales/mss_003.xml
git commit -m 'Updating file fales/mss_002.xml, Updating file fales/mss_003.xml'

git add archives/cap_1.xml fales/mss_004.xml tamwag/aia_001.xml
git commit -m 'Updating file archives/cap_1.xml, Updating file fales/mss_004.xml, Updating file tamwag/aia_001.xml'

git add archives/mc_1.xml
git rm fales/mss_002.xml
git add fales/mss_005.xml
git add tamwag/aia_002.xml
git commit -m 'Updating file archives/mc_1.xml, Deleting file fales/mss_002.xml EADID='mss_002', Updating file fales/mss_005.xml, Updating file tamwag/aia_002.xml'

popd

tar cvfz simple-repo.tar.gz simple-repo/.
rm -rf simple-repo
