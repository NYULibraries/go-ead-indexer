#!/usr/bin/bash -e

# check that our arg is a directory
ead_repo_path=$1
if [ ! -d "$ead_repo_path" ]
then
    echo >&2 "Usage: $0 <full path to EAD repository>"
    exit 1
fi

# check that go indexer has its env var
if [ -z "$SOLR_ORIGIN_WITH_PORT" ]
then
    echo >&2 "Must set SOLR_ORIGIN_WITH_PORT; aborting!"
    exit 1
fi

# assume all files ending in .xml in our dir are EADs: index them in descendinding file size order
ead_list=`find $ead_repo_path -name '*.xml' -printf '%s %p\n' | sort -n -r | awk -F ' ' '{print $2}' | tr '\n' ' '`
total_eads=$(echo $ead_list | wc -w | tr -d ' ')
echo "Indexing all EADs found in $ead_repo_path: $total_eads"
count=0
for ead in $ead_list
do
    count=$((count + 1))
    echo "Indexing $ead ($count/$total_eads)"
    ./eadindexer index --file $ead
done
