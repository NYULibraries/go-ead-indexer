package sanitize

import (
	"fmt"
	"regexp"
)

// https://github.com/rgrove/sanitize/blob/v6.0.1/lib/sanitize.rb#L27
// REGEX_HTML_CONTROL_CHARACTERS = /[\u0001-\u0008\u000b\u000e-\u001f\u007f-\u009f]+/u
const regexpHTMLControlCharacters = `[\x{0001}-\x{0008}\x{000b}\x{000e}-\x{001f}\x{007f}-\x{009f}]`

// https://github.com/rgrove/sanitize/blob/v6.0.1/lib/sanitize.rb#L34
// REGEX_HTML_NON_CHARACTERS = /[\ufdd0-\ufdef\ufffe\uffff\u{1fffe}\u{1ffff}\u{2fffe}\u{2ffff}\u{3fffe}\u{3ffff}\u{4fffe}\u{4ffff}\u{5fffe}\u{5ffff}\u{6fffe}\u{6ffff}\u{7fffe}\u{7ffff}\u{8fffe}\u{8ffff}\u{9fffe}\u{9ffff}\u{afffe}\u{affff}\u{bfffe}\u{bffff}\u{cfffe}\u{cffff}\u{dfffe}\u{dffff}\u{efffe}\u{effff}\u{ffffe}\u{fffff}\u{10fffe}\u{10ffff}]+/u
const regexpHTMLNonCharacters = `[\x{fdd0}-\x{fdef}\x{fffe}\x{ffff}\x{1fffe}\x{1ffff}\x{2fffe}\x{2ffff}\x{3fffe}\x{3ffff}\x{4fffe}\x{4ffff}\x{5fffe}\x{5ffff}\x{6fffe}\x{6ffff}\x{7fffe}\x{7ffff}\x{8fffe}\x{8ffff}\x{9fffe}\x{9ffff}\x{afffe}\x{affff}\x{bfffe}\x{bffff}\x{cfffe}\x{cffff}\x{dfffe}\x{dffff}\x{efffe}\x{effff}\x{ffffe}\x{fffff}\x{10fffe}\x{10ffff}]`

// https://github.com/rgrove/sanitize/blob/v6.0.1/lib/sanitize.rb#L48
// REGEX_UNSUITABLE_CHARS = /(?:#{REGEX_HTML_CONTROL_CHARACTERS}|#{REGEX_HTML_NON_CHARACTERS})/u
var sanitizeRegexpString = fmt.Sprintf("(?:%s+|%s+)", regexpHTMLControlCharacters,
	regexpHTMLNonCharacters)

var sanitizeRegexp = regexp.MustCompile(sanitizeRegexpString)

func Clean(text string) string {
	return sanitizeRegexp.ReplaceAllString(text, "")
}
