package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var slugInvalidChars = regexp.MustCompile(`[^a-z0-9]+`)

// SlugMaxLen leaves room, within a 150-char `slug` column, for a "-NNN"
// uniqueness suffix a caller's insert-and-retry loop may need to append.
const SlugMaxLen = 140

// Slugify strips diacritics (e.g. "Café" -> "Cafe"), lowercases, replaces
// runs of remaining non-alphanumeric characters (including any script that
// isn't Latin/digits, e.g. Arabic or CJK, or symbols/emoji) with a single
// hyphen, trims leading/trailing hyphens, and caps the length.
//
// Inputs that carry no representable ASCII content at all (e.g. purely
// Arabic, purely emoji) collapse to an empty string here — Slugify returns
// fallback in that case, so a caller's insert-and-retry loop still produces
// a valid, unique slug by appending "-2", "-3", etc.
func Slugify(s, fallback string) string {
	s = stripDiacritics(s)
	s = strings.ToLower(s)
	s = slugInvalidChars.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")

	if len(s) > SlugMaxLen {
		s = strings.Trim(s[:SlugMaxLen], "-")
	}

	if s == "" {
		return fallback
	}
	return s
}

// stripDiacritics transliterates accented Latin characters to their plain
// ASCII base (NFD-decompose, drop combining marks, NFC-recompose) so e.g.
// "café" slugifies to "cafe" instead of being silently dropped.
func stripDiacritics(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		return s
	}
	return result
}
