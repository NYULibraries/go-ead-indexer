package collectiondoc

type SolrAddMessage struct {
	Add AddElement `xml:"add"`
}

type AddElement struct {
	Doc DocElement `xml:"doc"`
}

type DocElement struct {
	ID string `xml:"id"`
}

func (collectionDoc *CollectionDoc) setSolrAddMessage() error {
	solrAddMessage := SolrAddMessage{
		AddElement{
			DocElement{
				ID: "testid",
			},
		},
	}

	collectionDoc.SolrAddMessage = solrAddMessage

	return nil
}

func (solrAddMessage *SolrAddMessage) String() string {
	return "test"
}
