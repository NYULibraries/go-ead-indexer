package component

import (
	"errors"
	"fmt"
	"github.com/lestrrat-go/libxml2/types"
	"github.com/nyulibraries/go-ead-indexer/pkg/ead/eadutil"
	"github.com/nyulibraries/go-ead-indexer/pkg/sanitize"
)

const noTitleAvailable = "[No title available]"

func getAncestorUnitTitle(node types.Node) (string, error) {
	var ancestorUnitTitle string

	// Try `<unittitle>` first, then try <unitdate>, and if neither worked, return
	// no title available.
	xpathResult, err := node.Find("did/unittitle")
	if err != nil {
		return ancestorUnitTitle, err
	}
	defer xpathResult.Free()

	unitTitleNodes := xpathResult.NodeList()
	if len(unitTitleNodes) > 0 {
		unitTitleContents, err := eadutil.ParseEscapedNodeTextContent(unitTitleNodes[0])
		if err != nil {
			return ancestorUnitTitle, errors.New(
				fmt.Sprintf(`eadutil.ParseEscapedNodeTextContent(unitTitleNodes[0]) error: %s`, err.Error()))
		}

		// TODO: DLFA-238
		// Replace this with `util.IsNonEmptyString(unitDateContents)`
		if unitTitleContents != "" {
			// TODO: Find out if `sanitize.Clean()` is necessary.
			// We are doing this because v1 indexer seems to suggest it might be
			// necessary.  See `get_title()` in this comment:
			// https://jira.nyu.edu/browse/DLFA-212?focusedCommentId=8495151&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-8495151
			ancestorUnitTitle = sanitize.Clean(unitTitleContents)

			// TODO: DLFA-238
			// Remove this left- and right- padding for matching v1 indexer bug
			// behavior described here:
			// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10849506&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10849506
			ancestorUnitTitle = eadutil.PadUnitTitleIfNeeded(
				eadutil.StripOpenAndCloseTags(unitTitleNodes[0].String()),
				ancestorUnitTitle)
		}
	}

	// <unittitle> didn't work.  Try <unitdate>.
	if ancestorUnitTitle == "" {
		xpathResult, err := node.Find("did/unitdate")
		if err != nil {
			return ancestorUnitTitle, err
		}
		defer xpathResult.Free()

		unitDateNodes := xpathResult.NodeList()
		if len(unitDateNodes) > 0 {
			unitDateContents := unitDateNodes[0].TextContent()
			// TODO: DLFA-238
			// Replace this with `util.IsNonEmptyString(unitDateContents)`
			if unitDateContents != "" {
				// TODO: Find out if `sanitize.Clean()` is necessary.
				// We are doing this because v1 indexer seems to suggest it might be
				// necessary.  See `get_title()` in this comment:
				// https://jira.nyu.edu/browse/DLFA-212?focusedCommentId=8495151&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-8495151
				ancestorUnitTitle = sanitize.Clean(unitDateContents)

				// TODO: DLFA-238
				// Remove this left- and right- padding for matching v1 indexer bug
				// behavior described here:
				// https://jira.nyu.edu/browse/DLFA-211?focusedCommentId=10849506&page=com.atlassian.jira.plugin.system.issuetabpanels:comment-tabpanel#comment-10849506
				ancestorUnitTitle = eadutil.PadUnitTitleIfNeeded(
					eadutil.StripOpenAndCloseTags(unitDateNodes[0].String()),
					ancestorUnitTitle)
			}
		}
	}

	// Can't create a title.
	if ancestorUnitTitle == "" {
		return noTitleAvailable, nil
	}

	return ancestorUnitTitle, nil
}

func makeAncestorUnitTitleListMap(node types.Node) (map[string][]string, error) {
	ancestorUnitTitleListMap := map[string][]string{}

	dscNode, err := eadutil.GetFirstNode("//dsc", node)
	if err != nil {
		return ancestorUnitTitleListMap, err
	}

	if dscNode == nil {
		return ancestorUnitTitleListMap,
			errors.New("makeAncestorUnitTitleListMap() error: no <dsc> element found")
	}

	childCNodes, err := eadutil.GetNodeList("child::c", dscNode)
	if err != nil {
		return ancestorUnitTitleListMap, err
	}

	if len(childCNodes) == 0 {
		return ancestorUnitTitleListMap, nil
	}

	for _, cNode := range childCNodes {
		ancestorUnitTitleList := []string{}
		err = makeAncestorUnitTitleListMap_add(ancestorUnitTitleListMap,
			ancestorUnitTitleList, cNode)
		if err != nil {
			return ancestorUnitTitleListMap, err
		}
	}

	return ancestorUnitTitleListMap, nil
}

func makeAncestorUnitTitleListMap_add(ancestorUnitTitleListMap map[string][]string,
	ancestorUnitTitleList []string, node types.Node) error {

	idNode, err := node.(types.Element).GetAttribute("id")
	if err != nil {
		return errors.New(
			fmt.Sprintf("makeAncestorUnitTitleListMap_add() error: can't get `id` attribute of node: %s",
				node.String()))
	}

	ancestorUnitTitleListMap[idNode.NodeValue()] = ancestorUnitTitleList

	ancestorUnitTitle, err := getAncestorUnitTitle(node)
	if err != nil {
		return err
	}
	ancestorUnitTitleList = append(ancestorUnitTitleList, ancestorUnitTitle)

	childCNodes, err := eadutil.GetNodeList("child::c", node)
	if err != nil {
		return err
	}

	if len(childCNodes) == 0 {
		return nil
	}

	for _, childCNode := range childCNodes {
		err = makeAncestorUnitTitleListMap_add(ancestorUnitTitleListMap,
			ancestorUnitTitleList, childCNode)
		if err != nil {
			return err
		}
	}

	return nil
}
