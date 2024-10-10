package ead

type SolrAddMessages struct {
	Collection string
	Components *[]Component
}

type ComponentSolrAddMessage struct {
	ID      string
	Message string
}

type SolrAddMessages struct {
	Collection string
	Components *[]ComponentSolrAddMessage
}

func ParseSolrAddMessages(eadXML string) (SolrAddMessages, error) {
	// TODO: Remove this fake data
	return SolrAddMessages{
		Collection: "",
		Components: &[]ComponentSolrAddMessage{
			{
				ID:      "mos_2021additional-daos",
				Message: "",
			},
			{
				ID:      "mos_2021dao1",
				Message: "",
			},
			{
				ID:      "mos_2021non-existent",
				Message: "",
			},
		},
	}, nil
}
