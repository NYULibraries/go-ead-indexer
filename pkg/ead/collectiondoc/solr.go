package collectiondoc

type SolrAddMessage struct {
	Add AddElement `xml:"add"`
}

type AddElement struct {
	Doc DocElement `xml:"doc"`
}

type DocElement struct {
	EAD_SSI        string `xml:"ead_ssi"`
	ID             string `xml:"id"`
	Repository_SIM string `xml:"repository_sim"`
	Repository_SSI string `xml:"repository_ssi"`
	Repository_SSM string `xml:"repository_ssm"`
}

func (collectionDoc *CollectionDoc) setSolrAddMessage() {
	docElement := &collectionDoc.SolrAddMessage.Add.Doc

	docElement.EAD_SSI = collectionDoc.Parts.EADID.Values[0]
	docElement.ID = collectionDoc.Parts.EADID.Values[0]
	docElement.Repository_SIM = collectionDoc.Parts.RepositoryCode.Values[0]
	docElement.Repository_SSI = collectionDoc.Parts.RepositoryCode.Values[0]
	docElement.Repository_SSM = collectionDoc.Parts.RepositoryCode.Values[0]
}

func (solrAddMessage *SolrAddMessage) String() string {
	return "test"
}
