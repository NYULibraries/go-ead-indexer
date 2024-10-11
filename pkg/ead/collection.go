package ead

import (
	"github.com/lestrrat-go/libxml2/types"
)

type Collection struct {
	SolrAddMessage string
	XPathParts     XPathParts
}

func (collection *Collection) populateXPathParts(node types.Node) error {
	xp := collection.XPathParts

	xp["TEST"] = XPathPart{
		Query:  "QUERY",
		Values: []string{"VALUE"},
	}

	return nil
}
