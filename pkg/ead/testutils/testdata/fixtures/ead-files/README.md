* _mos_2021.xml_: the AppDev or indexer Omega file
  * Orginally based on [dlts\-finding\-aids\-ead\-go\-packages/ead/testdata/omega/v0\.1\.5 /Omega\-EAD\.xml](https://raw.githubusercontent.com/NYULibraries/dlts-finding-aids-ead-go-packages/7baee7dfde24a01422ec8e6470fdc8a76d84b3fb/ead/testdata/omega/v0.1.5/Omega-EAD.xml), the DLTS or Finding Aids or FADESIGN Omega file
  * Changes:
    * File was renamed _mos_2021.xml_ to match the `<eadid>` value used in the file.
    * Tags have been added to test indexer code that was not getting exercised by the Finding Aids (FADESIGN-originated) Omega file.
  * See also Jira DLFA-221: [v1 indexer: get Solr HTTP requests and responses for Omega files](https://jira.nyu.edu/browse/DLFA-221)