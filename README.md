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

# Image publishing and deployment

CircleCI builds, tests, and pushes three Quay image tags for every green
build on every branch:

```text
quay.io/nyulibraries/go-ead-indexer:<branch>
quay.io/nyulibraries/go-ead-indexer:<branch>-<sha>
quay.io/nyulibraries/go-ead-indexer:dev
```

Deployment is image-only; this repository never triggers indexing runs.
Hermes is deliberately not used for image updates: its `createJob -tag`
endpoint always starts an indexing run as a side effect, and `setImage`
only targets Deployments, not CronJobs. Indexing runs in both dev and prod
are triggered exclusively by
[`findingaids_eads_v2`](https://github.com/NYULibraries/findingaids_eads_v2)
data commits, through Hermes. The CronJobs (defined in
[`nyulibraries_kubernetes`](https://github.com/NYULibraries/nyulibraries_kubernetes))
use `imagePullPolicy: Always`, so each run pulls its mutable tag fresh:

* prod runs the `main` tag: merging to `main` updates the prod image, and
  the next data-triggered run uses it.
* dev runs the `dev` tag: every green build on every branch (including
  `main`) pushes its image as `dev`, so the next dev indexing run uses the
  last completed green build push. Merging to `main` re-publishes `dev`
  from `main`, which resets dev.

Only one branch can be active in dev at a time, so a green branch build is
a dev deployment, not just CI. `dev` is a single mutable tag, and the last
*completed* `:dev` push wins: with overlapping builds, an older build that
finishes later overwrites a newer one. The Job spec stores the tag string,
not a digest — each Pod pulls whatever `dev` points to when its container
starts, so dev keeps running that image — including for data-triggered
runs — until another build overwrites it. To
see what `dev` currently points at, compare its digest on
[quay.io](https://quay.io/repository/nyulibraries/go-ead-indexer?tab=tags)
with the branch tags. To reset dev to main without waiting for a merge,
re-run the latest `main` workflow in CircleCI (or retag manually:
`docker pull ...:main && docker tag ...:main ...:dev && docker push ...:dev`).
