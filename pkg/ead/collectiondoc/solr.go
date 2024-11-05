package collectiondoc

type SolrAddMessage struct {
	Add AddElement `xml:"add"`
}

type AddElement struct {
	Doc DocElement `xml:"doc"`
}

type DocElement struct {
	EAD_ssi        string   `xml:"ead_ssi"`
	ID             string   `xml:"id"`
	Repository_sim string   `xml:"repository_sim"`
	Repository_ssi string   `xml:"repository_ssi"`
	Repository_ssm string   `xml:"repository_ssm"`
}

func (collectionDoc *CollectionDoc) setSolrAddMessage() {
	docElement := &collectionDoc.SolrAddMessage.Add.Doc

	docElement.EAD_ssi = collectionDoc.Parts.EADID.Values[0]
	docElement.ID = collectionDoc.Parts.EADID.Values[0]
	docElement.Repository_sim = collectionDoc.Parts.RepositoryCode.Values[0]
	docElement.Repository_ssi = collectionDoc.Parts.RepositoryCode.Values[0]
	docElement.Repository_ssm = collectionDoc.Parts.RepositoryCode.Values[0]
}

func (solrAddMessage *SolrAddMessage) String() string {
	return "test"
}
