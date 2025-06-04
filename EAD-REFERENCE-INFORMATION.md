# EAD Reference Information

Jira Ticket: https://jira.nyu.edu/projects/DLFA/issues/DLFA-241

## References:
[ArchivesSpace Manual for Local Usage at NYU](https://docs.google.com/document/d/11kWxbFTazB6q5fDNBWDHJxMf3wdVsp8cd7HzjEhE-ao/edit?tab=t.0)  
[EAD Validation Criteria for Publishing Finding Aids](https://github.com/nyudlts/findingaids_docs/blob/main/user/EAD_Validation_Criteria_for_Publishing.md)  
[Finding Aids Publishing Service documentation](https://github.com/nyudlts/findingaids_docs)  


## Q&A:
**Question:** `Are HTML tags allowed in free text fields in ArchivesSpace, according to the ACM style rules?`  
No, only EAD tags should be used. Occasionally in the past an ArchivesSpace user who knew HTML would put in HTML tags to italicize something, which ArchivesSpace would then render, making it look like it would work when it went through the indexing process. 

**Question:** `What does the first token in EAD IDs (also used for filenames) denote? E.g. The "photos" prefix in the tamwag/photos_223.xml file.`  
It is a legacy practice of using a prefix label for a collecting area. In the past there were prefixes based on format, such as "buttons". This naming practice became unwieldy after a while. There is still some adherence to the prefix system, but in general, when there is discussion it's the repositories that tend to be referred to more. The prefixes are treated as something incidental. 

(For official EADID construction rules, see FADESIGN-20, and EAD Validation Criteria for Publishing Finding Aids in References below)

**Question:** `Can the level attribute of <archdesc> ever have a value other than "collection"?`  
It's not supposed to. Technically it's possible to set it to something else via the ArchivesSpace dropdown, but it would be an error to select anything but "collection" for a record whose EAD file is used in Special Collections. 
