package component

type SolrAddMessage struct {
	Add AddElement `xml:"add"`
}

type AddElement struct {
	Doc DocElement `xml:"doc"`
}

type DocElement struct {
	ID string `xml:"id"`
}

func (component *Component) setSolrAddMessage() error {
	solrAddMessage := SolrAddMessage{
		AddElement{
			DocElement{
				ID: "testid",
			},
		},
	}

	component.SolrAddMessage = solrAddMessage

	return nil
}

func (solrAddMessage *SolrAddMessage) String() string {
	return "test"
}
