package tgmd

import "testing"

func TestUTF16Len_ASCII(t *testing.T) {
	if got := UTF16Len("hello"); got != 5 {
		t.Errorf("UTF16Len(\"hello\") = %d, want 5", got)
	}
}

func TestUTF16Len_Empty(t *testing.T) {
	if got := UTF16Len(""); got != 0 {
		t.Errorf("UTF16Len(\"\") = %d, want 0", got)
	}
}

func TestUTF16Len_Emoji(t *testing.T) {
	// 🌍 is U+1F30D — above U+FFFF, takes 2 UTF-16 code units (surrogate pair).
	if got := UTF16Len("🌍"); got != 2 {
		t.Errorf("UTF16Len(\"🌍\") = %d, want 2", got)
	}
}

func TestUTF16Len_Mixed(t *testing.T) {
	// "Hello 🌍" = 6 ASCII (6 units) + 1 emoji (2 units) = 8
	if got := UTF16Len("Hello 🌍"); got != 8 {
		t.Errorf("UTF16Len(\"Hello 🌍\") = %d, want 8", got)
	}
}

func TestUTF16Len_Chinese(t *testing.T) {
	// Chinese characters are in BMP — 1 UTF-16 unit each.
	if got := UTF16Len("寓言"); got != 2 {
		t.Errorf("UTF16Len(\"寓言\") = %d, want 2", got)
	}
}

func TestUTF16Len_MultipleEmoji(t *testing.T) {
	// "👋🌍" = 2 + 2 = 4 UTF-16 units
	if got := UTF16Len("👋🌍"); got != 4 {
		t.Errorf("UTF16Len(\"👋🌍\") = %d, want 4", got)
	}
}

func TestUTF16RuneLen_ASCII(t *testing.T) {
	if got := UTF16RuneLen('a'); got != 1 {
		t.Errorf("UTF16RuneLen('a') = %d, want 1", got)
	}
}

func TestUTF16RuneLen_Emoji(t *testing.T) {
	if got := UTF16RuneLen('🌍'); got != 2 {
		t.Errorf("UTF16RuneLen('🌍') = %d, want 2", got)
	}
}

func TestSplitAtUTF16_ASCII(t *testing.T) {
	prefix, suffix := SplitAtUTF16("hello world", 5)
	if prefix != "hello" || suffix != " world" {
		t.Errorf("SplitAtUTF16(\"hello world\", 5) = (%q, %q), want (\"hello\", \" world\")", prefix, suffix)
	}
}

func TestSplitAtUTF16_AtStart(t *testing.T) {
	prefix, suffix := SplitAtUTF16("hello", 0)
	if prefix != "" || suffix != "hello" {
		t.Errorf("SplitAtUTF16(\"hello\", 0) = (%q, %q), want (\"\", \"hello\")", prefix, suffix)
	}
}

func TestSplitAtUTF16_BeyondEnd(t *testing.T) {
	prefix, suffix := SplitAtUTF16("hi", 10)
	if prefix != "hi" || suffix != "" {
		t.Errorf("SplitAtUTF16(\"hi\", 10) = (%q, %q), want (\"hi\", \"\")", prefix, suffix)
	}
}

func TestSplitAtUTF16_BeforeEmoji(t *testing.T) {
	// "Hi 🌍" — split at offset 3 should be before the emoji
	prefix, suffix := SplitAtUTF16("Hi 🌍", 3)
	if prefix != "Hi " || suffix != "🌍" {
		t.Errorf("SplitAtUTF16(\"Hi 🌍\", 3) = (%q, %q), want (\"Hi \", \"🌍\")", prefix, suffix)
	}
}

func TestSplitAtUTF16_AfterEmoji(t *testing.T) {
	// "🌍Hi" — emoji is 2 units, split at 2 should be after the emoji
	prefix, suffix := SplitAtUTF16("🌍Hi", 2)
	if prefix != "🌍" || suffix != "Hi" {
		t.Errorf("SplitAtUTF16(\"🌍Hi\", 2) = (%q, %q), want (\"🌍\", \"Hi\")", prefix, suffix)
	}
}

func TestSplitAtUTF16_Empty(t *testing.T) {
	prefix, suffix := SplitAtUTF16("", 0)
	if prefix != "" || suffix != "" {
		t.Errorf("SplitAtUTF16(\"\", 0) = (%q, %q), want (\"\", \"\")", prefix, suffix)
	}
}

func TestUTF16Encode(t *testing.T) {
	units := utf16Encode("A🌍B")
	// A=0x41, 🌍=0xD83C 0xDF0D (surrogate pair), B=0x42
	if len(units) != 4 {
		t.Errorf("utf16Encode(\"A🌍B\") length = %d, want 4", len(units))
	}
}
