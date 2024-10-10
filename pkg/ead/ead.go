package ead

type EAD struct {
	Collection Collection
	Components *[]Component
}

func New(eadXML string) (EAD, error) {
	ead := EAD{}

	// TODO: Remove this fake data
	ead.Collection.SolrAddMessage = ""
	ead.Components = &[]Component{
		{
			ID:             "mos_2021additional-daos",
			SolrAddMessage: "",
		},
		{
			ID:             "mos_2021dao1",
			SolrAddMessage: "",
		},
		{
			ID:             "mos_2021dao1",
			SolrAddMessage: "",
		},
	}

	return ead, nil
}
