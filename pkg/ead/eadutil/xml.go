// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Package xml implements a simple XML 1.0 parser that
// understands XML name spaces.

// This file contains an edited version of `xml.EscapeText()` from the standard
// library `xml` package:
// https://github.com/golang/go/blob/69234ded30614a471c35cef5d87b0e0d3c136cd9/src/encoding/xml/xml.go
//
// Changes:
//   - Only "&", "<", and ">" characters are escaped.  All other characters are passed.

// TODO: DLFA-238
// Delete this entire file and switch back to using standard library `xml.EscapeText()`.
// v1 indexer does not escape a lot of the characters that Go's `xml.EscapedText()`
// does, and these characters frequently appear in the <unittitle> data.
package eadutil

import (
	"io"
	"unicode/utf8"
)

var (
	escAmp = []byte("&amp;")
	escLT  = []byte("&lt;")
	escGT  = []byte("&gt;")
)

// EscapeText writes to w the properly escaped XML equivalent
// of the plain text data s.
func EscapeText(w io.Writer, s []byte) error {
	return escapeText(w, s, true)
}

// escapeText writes to w the properly escaped XML equivalent
// of the plain text data s. If escapeNewline is true, newline
// characters will be escaped.
func escapeText(w io.Writer, s []byte, escapeNewline bool) error {
	var esc []byte
	last := 0
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRune(s[i:])
		i += width
		switch r {
		case '&':
			esc = escAmp
		case '<':
			esc = escLT
		case '>':
			esc = escGT
		default:
			continue
		}
		if _, err := w.Write(s[last : i-width]); err != nil {
			return err
		}
		if _, err := w.Write(esc); err != nil {
			return err
		}
		last = i
	}
	_, err := w.Write(s[last:])
	return err
}
