#!/usr/bin/env bash

# USAGE EXAMPLES
# ==============
#
# Return all docs for Omega file from prod Solr (default):
#
#     $ ./solr-query.sh mos_2021
#
# Return first 10 docs (sorted by `id`) for photos_223 from dev Solr:
#
#     $ SPECIALCOLLECTIONS_SOLR_ORIGIN=http://44.216.225.190:8080 ./solr-query.sh photos_223 10

# Default to prod Solr, for which you must be on NYU-NET.
DEFAULT_SPECIALCOLLECTIONS_SOLR_ORIGIN=http://44.218.37.122:8080
# SPECIALCOLLECTIONS_SOLR_ORIGIN environment var must be a valid https://developer.mozilla.org/en-US/docs/Glossary/Origin
specialCollectionsSolrOrigin=${SPECIALCOLLECTIONS_SOLR_ORIGIN:-$DEFAULT_SPECIALCOLLECTIONS_SOLR_ORIGIN}

# Get the only require script arg, the EAD ID.
eadid=$1
if [ -z "$eadid" ]
then
    echo >&2 "Usage: $0 eadid [rows]"
    exit 1
fi

# Default to returning all Solr docs. There's no "give me all rows" option in Solr.
# You have to set `rows` to a number higher than the total number of docs in the
# index.
DEFAULT_ROWS=999999999
# Allow the user the option to set a row limit.
rows="${2:-$DEFAULT_ROWS}"

# We build the query used by the FAB landing page, which is a search that
# returns everything.
QUERY_PARAMS[0]='bq=format_sim%3A%22Archival+Collection%22%5E250'
QUERY_PARAMS[1]='bq=level_sim%3Afile%5E20'
QUERY_PARAMS[2]='bq=level_sim%3Aitem'
QUERY_PARAMS[3]='bq=level_sim%3Aseries%5E150'
QUERY_PARAMS[4]='bq=level_sim%3Asubseries%5E50'
QUERY_PARAMS[5]='defType=edismax'
QUERY_PARAMS[6]='f.collection_sim.facet.limit=21'
QUERY_PARAMS[7]='f.creator_sim.facet.limit=21'
QUERY_PARAMS[8]='f.dao_sim.facet.limit=21'
QUERY_PARAMS[9]='f.date_range_sim.facet.limit=21'
QUERY_PARAMS[10]='f.format_sim.facet.limit=21'
QUERY_PARAMS[11]='f.language_sim.facet.limit=21'
QUERY_PARAMS[12]='f.name_sim.facet.limit=21'
QUERY_PARAMS[13]='f.place_sim.facet.limit=21'
QUERY_PARAMS[14]='f.repository_sim.facet.limit=21'
QUERY_PARAMS[15]='f.subject_sim.facet.limit=21'
QUERY_PARAMS[16]='facet.field=collection_sim'
QUERY_PARAMS[17]='facet.field=creator_sim'
QUERY_PARAMS[18]='facet.field=dao_sim'
QUERY_PARAMS[19]='facet.field=date_range_sim'
QUERY_PARAMS[20]='facet.field=format_sim'
QUERY_PARAMS[21]='facet.field=language_sim'
QUERY_PARAMS[22]='facet.field=name_sim'
QUERY_PARAMS[23]='facet.field=place_sim'
QUERY_PARAMS[24]='facet.field=repository_sim'
QUERY_PARAMS[25]='facet.field=subject_sim'
QUERY_PARAMS[26]='facet.mincount=1'
QUERY_PARAMS[27]='facet=true'
QUERY_PARAMS[28]='fl=*'
QUERY_PARAMS[29]='indent=true'
QUERY_PARAMS[30]='pf=unittitle_teim%5E145.0+parent_unittitles_teim+collection_teim+unitid_teim%5E60+collection_unitid_teim%5E40+lang\'
QUERY_PARAMS[31]='uage_ssm+unitdate_start_teim+unitdate_end_teim+unitdate_teim+name_teim+subject_teim%5E60.0+abstract_teim%5E55.0+cr\'
QUERY_PARAMS[32]='eator_teim%5E60.0+scopecontent_teim%5E60.0+bioghist_teim%5E55.0+title_teim+material_type_teim+place_teim+dao_teim+\'
QUERY_PARAMS[33]='chronlist_teim+appraisal_teim+custodhist_teim%5E15+acqinfo_teim%5E20.0+address_teim+note_teim%5E30.0+phystech_teim\'
QUERY_PARAMS[34]='%5E30.0+author_teim%5E10.0'
QUERY_PARAMS[35]='ps=50'
QUERY_PARAMS[36]='qf=unittitle_teim%5E145.0+parent_unittitles_teim+collection_teim+unitid_teim%5E60+collection_unitid_teim%5E40+lang\'
QUERY_PARAMS[37]='uage_ssm+unitdate_start_teim+unitdate_end_teim+unitdate_teim+name_teim+subject_teim%5E60.0+abstract_teim%5E55.0+cr\'
QUERY_PARAMS[38]='eator_teim%5E60.0+scopecontent_teim%5E60.0+bioghist_teim%5E55.0+title_teim+material_type_teim+place_teim+dao_teim+\'
QUERY_PARAMS[39]='chronlist_teim+appraisal_teim+custodhist_teim%5E15+acqinfo_teim%5E20.0+address_teim+note_teim%5E30.0+phystech_teim\'
QUERY_PARAMS[40]='%5E30.0+author_teim%5E10.0'
QUERY_PARAMS[41]='sort=score+desc'
QUERY_PARAMS[42]='timeAllowed=-1'
QUERY_PARAMS[43]='wt=json'

# Join
# Source: https://stackoverflow.com/questions/1527049/how-can-i-join-elements-of-a-bash-array-into-a-delimited-string
queryString=$(IFS='&' ; echo "${QUERY_PARAMS[*]}")

# This is the basic Solr query to get all docs.  We are using the FAB query instead.
# QUERY_URL="${SPECIALCOLLECTIONS_SOLR_ORIGIN}/solr/findingaids/select?q=ead_ssi:${eadid}&wt=json&indent=true&rows=${rows}&sort=id+asc"

# FAB request
# See: "Sample Solr requests/queries used by the FAB"
# https://jira.nyu.edu/browse/DLFA-182?jql=text%20~%20%22solr%20queries%22
QUERY_URL="${specialCollectionsSolrOrigin}/solr/findingaids/select?q=ead_ssi:${eadid}&rows=${rows}&${queryString}"

curl --silent $QUERY_URL

