package tgmd

import "unicode/utf16"

// UTF16Len returns the number of UTF-16 code units needed to encode s.
// BMP characters (< U+10000) take 1 unit; supplementary characters (emoji, etc.) take 2.
func UTF16Len(s string) int {
	n := 0
	for _, r := range s {
		if r >= 0x10000 {
			n += 2
		} else {
			n++
		}
	}
	return n
}

// UTF16RuneLen returns the number of UTF-16 code units for a single rune.
func UTF16RuneLen(r rune) int {
	if r >= 0x10000 {
		return 2
	}
	return 1
}

// SplitAtUTF16 splits s at the given UTF-16 offset.
// Returns the prefix (up to the offset) and the suffix (from the offset onward).
// If offset is beyond the string length, returns (s, "").
func SplitAtUTF16(s string, offset int) (prefix, suffix string) {
	pos := 0
	for i, r := range s {
		if pos >= offset {
			return s[:i], s[i:]
		}
		pos += UTF16RuneLen(r)
	}
	return s, ""
}

// utf16Encode is a helper that encodes a string to UTF-16 code units.
// Exposed for internal use and testing.
func utf16Encode(s string) []uint16 {
	runes := []rune(s)
	return utf16.Encode(runes)
}
