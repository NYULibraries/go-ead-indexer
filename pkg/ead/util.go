package ead

import (
	"github.com/lestrrat-go/libxml2/types"
)

func getValuesForXPathQuery(query string, node types.Node) ([]string, error) {
	var values []string

	xpathResult, err := node.Find(query)
	if err != nil {
		return nil, err
	}

	for _, node = range xpathResult.NodeList() {
		values = append(values, node.NodeValue())
	}

	return values, nil
}

func replaceMARCSubfieldDemarcatorsInSlice(stringSlice []string) []string {
	newSlice := []string{}
	for _, element := range stringSlice {
		newSlice = append(newSlice, replaceMARCSubfieldDemarcators(element))
	}

	return newSlice
}

// TODO: fix the bug we've intentionally preserved here -- for details, see:
// * https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10154897&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10154897
// * https://jira.nyu.edu/browse/DLFA-229?focusedCommentId=10153922&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10153922
func replaceMARCSubfieldDemarcators(str string) string {
	return marcSubfieldDemarcator.ReplaceAllString(str, "--")
}
