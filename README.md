# EAD Indexer (project EDIP)

Jira Epic: https://jira.nyu.edu/projects/DLFA/issues/DLFA-181


### `CLI`
#### Indexing EADs or Git Commits
```
Index EAD file or commit

Usage:
  go-ead-indexer index [flags]

Examples:
  go-ead-indexer index --file=[path to EAD file] --logging-level="debug"
  go-ead-indexer index --git-repo=[path] --commit=[hash] --logging-level="error"

Flags:
  -c, --commit string          hash of git commit
  -f, --file string            path to EAD file
  -g, --git-repo string        path to EAD files git repo
  -h, --help                   help for index
  -l, --logging-level string   Sets logging level: debug, info, error (default "info")
```

#### Deleting data for an EAD from the Solr index
```
Delete data from the index using the EADID

Usage:
  go-ead-indexer delete [flags]

Examples:
go-ead-indexer delete --eadid=[EADID] --logging-level="debug" --assume-yes

Flags:
  -y, --assume-yes             disable interactive mode
  -e, --eadid string           EADID value of EAD data to delete
  -h, --help                   help for delete
  -l, --logging-level string   Sets logging level: debug, info, error (default "info")
  ```

# Additional documentation

* [EAD Reference Information](EAD-REFERENCE-INFORMATION.md)
* Package README files
  * [git](pkg/git/README.md) 
  * [index](pkg/index/README.md) 
