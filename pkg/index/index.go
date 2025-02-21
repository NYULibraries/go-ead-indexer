package index

import (
	"github.com/nyulibraries/go-ead-indexer/pkg/ead"
)

const MessageKey = "message"

func hello() {
	_, err := ead.New("hello", "world")
	if err != nil {
		panic(err)
	}
}
