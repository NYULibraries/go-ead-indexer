#!/bin/bash

# This script generates a git repo test fixture for use in the git package tests.
# The script creates a directory named 'git-repo', then generates and commits various files.
# Finally, the script renames the 'git-repo/.git' directory to 'git-repo/dot-git'.
# The git-repo directory can be moved into the pkg/git/testdata/fixtures directory
# for use in git pkg tests.

# Commit history replicated in repo (NOTE: commit hashes WILL differ)
# 039021182a1d291b7d0638b19e1a3cd91a0eace9 2025-09-25 15:19:56 -0400 | Rename archives/mc_1.xml -> archives/mc_0.xml and add note about it to README.md (HEAD -> master) [David]
# 49e0b399e657d626c7d5e6a6cf937e0aa8481863 2025-09-25 15:19:56 -0400 | Updating .circleci/config.yml with [whatever] [David]
# 51d8153c9fee66173b7dd926f460d1bf34f647cb 2025-09-25 15:19:56 -0400 | Updating README.md with [whatever] [David]
# 60d7bc6b81bd415154b5cb8a4ef7f132dbb10733 2025-09-25 15:19:56 -0400 | Updating file fales/mss_001.xml [David]
# 161c00ae8a330d03b0e7c374a81ef05e12ea7084 2025-09-25 15:19:56 -0400 | Updating file archives/mc_1.xml, Deleting file fales/mss_002.xml EADID='mss_002', Updating file fales/mss_005.xml, Updating file tamwag/aia_002.xml [David]
# 8667d870f483e13d08abcaccd5908d6103e430a8 2025-09-25 15:19:55 -0400 | Updating file archives/cap_1.xml, Updating file fales/mss_004.xml, Updating file tamwag/aia_001.xml [David]
# 0b3d154283e14443874c2c743ac93666b84ce4ab 2025-09-25 15:19:55 -0400 | Updating file fales/mss_002.xml, Updating file fales/mss_003.xml [David]
# c4ef6e204f4e7da3f72c1c0798759d0d7780da83 2025-09-25 15:19:55 -0400 | Updating file fales/mss_001.xml [David]
# e0c4b58339f58250e0ec86cc19a8a41aa1910fbb 2025-09-25 15:19:55 -0400 | Initial commit of README.md and .circle/config.yml [David]

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
mkdir -p git-repo/.circleci git-repo/archives git-repo/fales git-repo/tamwag

pushd git-repo &>/dev/null || err_exit "Failed to change directory to git-repo/archives"
echo 'README.md' > README.md
popd &>/dev/null || err_exit "Failed to popd after creating README.md file"

pushd git-repo/.circleci &>/dev/null || err_exit "Failed to change directory to git-repo/archives"
echo 'config.yml' > config.yml
popd &>/dev/null || err_exit "Failed to popd after creating .circleci/config.yml file"

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

git add README.md .circleci/config.yml
git commit -m "Initial commit of README.md and .circle/config.yml"

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
git add README.md
git commit -m 'Updating README.md with [whatever]'

echo 'config.yml update' > .circleci/config.yml
git add .circleci/config.yml
git commit -m 'Updating .circleci/config.yml with [whatever]'

echo 'README.md [had to do something special to archives/mc_1.xml' > README.md
git add README.md
echo 'Do something special to mc_1.xml' > archives/mc_1.xml
git add archives/mc_1.xml
git commit -m '[Do something special to] archives/mc_1.xml and add note about it to README.md'

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
