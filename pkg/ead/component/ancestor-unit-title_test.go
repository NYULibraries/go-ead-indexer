package component

import (
	"fmt"
	"github.com/lestrrat-go/libxml2/parser"
	"reflect"
	"testing"
)

const sampleXML = `<ead>
	<dsc>
		<c id="component_1_level_1">
			<did><unittitle>Component #1, level 1</unittitle></did>
			<c id="component_1_level_2">
				<did><unittitle>Component #1, level 2</unittitle></did>
				<c id="component_1_level_3">
					<did><unittitle>Component #1, level 3</unittitle></did>
				</c>
			</c>
		</c>
		<c id="component_2_level_1">
			<did><unitdate>1901</unitdate></did>
			<c id="component_2_level_2">
				<did><unitdate>1902</unitdate></did>
				<c id="component_2_level_3">
					<did><unitdate>1903</unitdate></did>
				</c>
			</c>
		</c>
		<c id="component_3_level_1">
			<c id="component_3_level_2">
			</c>
		</c>
		<c id="component_4_level_1">
			<did><unittitle>Component #4, level 1</unittitle></did>
			<c id="component_4_level_2_subcomponent_a_level_1">
				<did><unittitle>Component #4, level 2, Subcomponent A, level 1</unittitle></did>
				<c id="component_4_level_2_subcomponent_a_subcomponent_a_level1">
					<did><unittitle>Component #4, level 2, Subcomponent A, Subcomponent A, level 1</unittitle></did>
				</c>
				<c id="component_4_level_2_subcomponent_a_subcomponent_b_level1">
					<did><unittitle>Component #4, level 2, Subcomponent A, Subcomponent B, level 1</unittitle></did>
					<c id="component_4_level_2_subcomponent_a_subcomponent_b_level2">
						<did><unittitle>Component #4, level 2, Subcomponent A, Subcomponent B, level 2</unittitle></did>
					</c>
				</c>
			</c>
			<c id="component_4_level_2_subcomponent_b_level_1">
				<did><unittitle>Component #4, level 2, Subcomponent B, level 1</unittitle></did>
			</c>
			<c id="component_4_level_2_subcomponent_c_level_1">
				<did><unittitle>Component #4, level 2, Subcomponent C, level 1</unittitle></did>
				<c id="component_4_level_2_subcomponent_c_subcomponent_a_level1">
					<did><unittitle>Component #4, level 2, Subcomponent C, Subcomponent A, level 1</unittitle></did>
				</c>
				<c id="component_4_level_2_subcomponent_c_subcomponent_b_level1">
					<did><unittitle>Component #4, level 2, Subcomponent C, Subcomponent B, level 1</unittitle></did>
					<c id="component_4_level_2_subcomponent_c_subcomponent_b_level2">
						<did><unittitle>Component #4, level 2, Subcomponent C, Subcomponent B, level 2</unittitle></did>
					</c>
				</c>
			</c>
		</c>
	</dsc>
</ead>
`

var expectedAncestorUnitTitleListMap = map[string][]string{
	// Component #1: <unittitle>
	"component_1_level_1": {},
	"component_1_level_2": {
		"Component #1, level 1",
	},
	"component_1_level_3": {
		"Component #1, level 1",
		"Component #1, level 2",
	},
	// Component #2: <unitdate>
	"component_2_level_1": {},
	"component_2_level_2": {
		"1901",
	},
	"component_2_level_3": {
		"1901",
		"1902",
	},
	// Component #3: no <unititle> or <unitdate>
	"component_3_level_1": {},
	"component_3_level_2": {
		noTitleAvailable,
	},
	// Component #4: hierarchy with sibling subcomponents
	// Regression test case for https://jira.nyu.edu/browse/DLFA-282.
	"component_4_level_1": {},
	"component_4_level_2_subcomponent_a_level_1": {
		"Component #4, level 1",
	},
	"component_4_level_2_subcomponent_a_subcomponent_a_level1": {
		"Component #4, level 1",
		"Component #4, level 2, Subcomponent A, level 1",
	},
	"component_4_level_2_subcomponent_a_subcomponent_b_level1": {
		"Component #4, level 1",
		"Component #4, level 2, Subcomponent A, level 1",
	},
	"component_4_level_2_subcomponent_a_subcomponent_b_level2": {
		"Component #4, level 1",
		"Component #4, level 2, Subcomponent A, level 1",
		"Component #4, level 2, Subcomponent A, Subcomponent B, level 1",
	},
	"component_4_level_2_subcomponent_b_level_1": {
		"Component #4, level 1",
	},
	"component_4_level_2_subcomponent_c_level_1": {
		"Component #4, level 1",
	},
	"component_4_level_2_subcomponent_c_subcomponent_a_level1": {
		"Component #4, level 1",
		"Component #4, level 2, Subcomponent C, level 1",
	},
	"component_4_level_2_subcomponent_c_subcomponent_b_level1": {
		"Component #4, level 1",
		"Component #4, level 2, Subcomponent C, level 1",
	},
	"component_4_level_2_subcomponent_c_subcomponent_b_level2": {
		"Component #4, level 1",
		"Component #4, level 2, Subcomponent C, level 1",
		"Component #4, level 2, Subcomponent C, Subcomponent B, level 1",
	},
}

// This unit test is mainly for illustrative purposes.  It should not be considered
// an adequate regression test.  The golden files tests in ead_test.go provide
// much stronger coverage and should be considered the actual regression test
// suite for this function.
func TestMakeAncestorUnitTitleListMap(t *testing.T) {
	xmlParser := parser.New()
	xmlDoc, err := xmlParser.ParseString(sampleXML)
	defer xmlDoc.Free()
	if err != nil {
		t.Fatalf("`xmlParser.ParseString(sampleXML)` failed with error: %s", err.Error())
	}

	rootNode, err := xmlDoc.DocumentElement()
	if err != nil {
		t.Fatalf("`xmlDoc.DocumentElement()` failed with error: %s", err.Error())
	}

	ancestorUnitTitleListMap, err := makeAncestorUnitTitleListMap(rootNode)
	if err != nil {
		t.Errorf("`makeAncestorUnitTitleListMap(rootNode) failed with error: %s`", err.Error())
	}

	if !reflect.DeepEqual(ancestorUnitTitleListMap, expectedAncestorUnitTitleListMap) {
		actual := fmt.Sprintf("%v", ancestorUnitTitleListMap)
		expected := fmt.Sprintf("%v", expectedAncestorUnitTitleListMap)
		t.Errorf("\nExpected: %s\n\nActual:   %s", expected, actual)
	}
}
