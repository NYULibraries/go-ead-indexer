package component

import "regexp"

const ARCHIVAL_OBJECT_FORMAT = "Archival Object"
const ARCHIVAL_SERIES_FORMAT = "Archival Series"

var archivalSeriesRegExp = regexp.MustCompile(`/\Aseries|subseries/`)

// TODO: Do we need to have anything in `CollectionDoc.Part.Source` for these?
func (component *Component) setComplexParts() []error {
	errs := []error{}

	component.setFormat()

	return errs
}

func (component *Component) setFormat() {
	parts := &component.Parts

	var formatLevel string
	level := parts.Level.Values[0]
	if archivalSeriesRegExp.Match([]byte(level)) {
		formatLevel = ARCHIVAL_SERIES_FORMAT
	} else {
		formatLevel = ARCHIVAL_OBJECT_FORMAT
	}

	parts.Format.Values = []string{formatLevel}
}

func (component *Component) setLocation() {
	parts := &component.Parts

	locationValues := []string{}

	parts.Location.Values = locationValues
}
