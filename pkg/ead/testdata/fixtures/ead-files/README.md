## The EDIP/indexer Omega file

The indexer Omega file: _edip/mos_2024.xml_

* Originally based on [dlts\-finding\-aids\-ead\-go\-packages/ead/testdata/omega/v0\.1\.5 /Omega\-EAD\.xml](https://raw.githubusercontent.com/NYULibraries/dlts-finding-aids-ead-go-packages/7baee7dfde24a01422ec8e6470fdc8a76d84b3fb/ead/testdata/omega/v0.1.5/Omega-EAD.xml),
  the DLTS or Finding Aids or FADESIGN Omega file.
* Changes:
  * File was renamed _mos_2024.xml_ to match the `<eadid>` value used in the file.
  * Tags have been added to test indexer code that was not getting exercised by
    the Finding Aids (FADESIGN-originated) Omega file.
* See also Jira DLFA-221: [v1 indexer: get Solr HTTP requests and responses for Omega files](https://jira.nyu.edu/browse/DLFA-221)

## Curated "list of 10"

"The list of 10 (subset of the list of 11)" from
[DLFA-220: Lists of EAD files sample files for testing \(especially golden file testing\)](https://nyu.atlassian.net/browse/DLFA-220):

* akkasah/ad_mc_030
* akkasah/ad_mc_066
* cbh/arc_212_plymouth_beecher
* fales/mss_420
* fales/mss_460
* nyhs/ms256_harmon_hendricks_goldstone
* nyhs/ms347_foundling_hospital
* nyuad/ad_mc_019
* tamwag/aia_003
* tamwag/tam_143

## Random sample

A random selection of 15 files from each repository.

In statistics, 15 is considered a good minimum starting point for a sample size
that will yield results of statistical significance.  Here we treated each repository
as a population from which we take one hopefully statistically significant sample.
In reality, this view of the corpus isn't actually that accurate.  Based on past
development experience in both the Finding Aids and Special Collections projects
and feedback from domain experts (archivists), it's more the case that within each
repository there are subpopulations that formed over time due to changes in
archivist practice and adapting to exigencies that often arose out of the limitations
of various technical systems the archivists worked with.  Unfortunately, it would
take a significant amount of effort to identify these subpopulations, as there
isn't any one individual who has the required historical knowledge.

Even disregarding the variation introduced by the formation of subpopulations
(sticking with statistical terms here), the extreme flexibility of the EAD2002
specification and the very broad variations that exist inherently in archives
and their metadata means that there are often so many important edge cases that
need to be covered that it can sometimes feel like the entire corpus is made of
edge cases.

Here are the file counts of the EAD sample sets that were used for testing and QA
in the FADESIGN project (Finding Aids Redesign), in chronological order of creation
over the several years of the development of the new Finding Aids:

* [dlts\-finding\-aids\-ead\-sample\-set\-1](https://github.com/NYULibraries/dlts-finding-aids-ead-sample-set-1): 19 EAD files
* [dlts\-finding\-aids\-ead\-sample\-set\-2](https://github.com/NYULibraries/dlts-finding-aids-ead-sample-set-2): 2,205 EAD files
* [dlts\-finding\-aids\-ead\-sample\-set\-3](https://github.com/nyudlts/dlts-finding-aids-ead-sample-set-3): 2,249 EAD files
* [dlts\-finding\-aids\-ead\-sample\-set\-4](https://github.com/nyudlts/dlts-finding-aids-ead-sample-set-4): 2,257 EAD files

Initially QA was done by the dev team and selected stakeholders was done against
19 EAD files.  As time went on, the developers kept encountering EAD files and
patterns that violated their previous assumptions, requiring rework of their
data model and code.  The QA sample set had to be expanded, with the final QA
set including 2,257 EAD files.  At the time, this would likely have been more
than half of the EAD files in the entire corpus.

Thus, the sampling we have done here should just be considered a bare minimum
effort toward avoiding regressions.  Later we can consider making larger, more
comprehensive QA tests that perhaps we would only run on certain occasions and/or
only on the AppDev Lab host.  We can let experience be our guide.

Notes on the sampling:

* The only restriction placed on sampling was to exclude any EAD file whose size
  exceeded 313,331 bytes.  This was the size of _tamwag/tam_143.xml_, the largest
  of the "list of 10" (see above).  For the unit tests we seek to strike a balance
  between coverage and speed.
  * An exception was made for _vlp/_, which only had one EAD file that was 4.2M.
    It was manually included in order to get coverage of that repository.
* There were only 9 EAD files in _arabartarchive/_ at the time of the sampling,
  and all of these were included as fixtures.  So even though there are only 9 files
  and not 15, the coverage level for _arabartarchive/_ is 100% (as with _vlp/_).

Actual commands used to create the random sample:

```shell
go-ead-indexer> # Disallow existing fixture files from being selected in the random sampling.
go-ead-indexer> gfind $GO_EAD_INDEXER/pkg/ead/testdata/fixtures/ead-files/ -type f -name *.xml -printf '%f\n' | sort > /tmp/exclude-existing.txt
go-ead-indexer> cat /tmp/exclude-existing.txt
ad_mc_019.xml
ad_mc_030.xml
ad_mc_066.xml
aia_003.xml
arc_212_plymouth_beecher.xml
mos_2024.xml
ms256_harmon_hendricks_goldstone.xml
ms347_foundling_hospital.xml
mss_420.xml
mss_460.xml
tam_143.xml
go-ead-indexer> 
```

```shell
findingaids_eads_v2> git log -1
commit 96191c044470980f3ebb4f29b5dc2d5bcb5ac1b8 (HEAD -> master, origin/master, origin/HEAD)
Author: nyudl-ead-committer <eadcommitter@nyu.edu>
Date:   Thu Jun 25 15:08:04 2026 -0400

    Updating file archives/mc_334.xml
findingaids_eads_v2> for f in $( for repository in $( ls | grep -v README ); do \
>                            gfind $repository -type f -size -323331c -printf "%p\n" | grep -v -f /tmp/exclude-existing.txt | shuf -n 15; \
>                            done ); do \
>                          destdir=$GO_EAD_INDEXER/pkg/ead/testdata/fixtures/ead-files/$( dirname $f ); \
>                          echo "mkdir -p $destdir/ && cp -p $f $destdir/"; \
>                        done | bash
findingaids_eads_v2> mkdir -p $GO_EAD_INDEXER/pkg/ead/testdata/fixtures/ead-files/vlp
findingaids_eads_v2> cp -p vlp/mss_lapietra_001.xml $GO_EAD_INDEXER/pkg/ead/testdata/fixtures/ead-files/vlp/
findingaids_eads_v2> 
```


