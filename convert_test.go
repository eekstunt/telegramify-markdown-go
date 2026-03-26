package tgmd

import (
	"strings"
	"testing"
)

func TestConvert_Bold(t *testing.T) {
	msg := Convert("**bold**")
	if msg.Text != "bold" {
		t.Errorf("text = %q, want %q", msg.Text, "bold")
	}
	assertEntity(t, msg, 0, Entity{Type: Bold, Offset: 0, Length: 4})
}

func TestConvert_Italic(t *testing.T) {
	msg := Convert("*italic*")
	if msg.Text != "italic" {
		t.Errorf("text = %q, want %q", msg.Text, "italic")
	}
	assertEntity(t, msg, 0, Entity{Type: Italic, Offset: 0, Length: 6})
}

func TestConvert_BoldItalic(t *testing.T) {
	msg := Convert("***both***")
	if msg.Text != "both" {
		t.Errorf("text = %q, want %q", msg.Text, "both")
	}
	// Should have both bold and italic entities covering the same range.
	hasBold := false
	hasItalic := false
	for _, e := range msg.Entities {
		if e.Type == Bold && e.Offset == 0 && e.Length == 4 {
			hasBold = true
		}
		if e.Type == Italic && e.Offset == 0 && e.Length == 4 {
			hasItalic = true
		}
	}
	if !hasBold {
		t.Error("missing Bold entity")
	}
	if !hasItalic {
		t.Error("missing Italic entity")
	}
}

func TestConvert_Strikethrough(t *testing.T) {
	msg := Convert("~~struck~~")
	if msg.Text != "struck" {
		t.Errorf("text = %q, want %q", msg.Text, "struck")
	}
	assertEntity(t, msg, 0, Entity{Type: Strikethrough, Offset: 0, Length: 6})
}

func TestConvert_InlineCode(t *testing.T) {
	msg := Convert("`code`")
	if msg.Text != "code" {
		t.Errorf("text = %q, want %q", msg.Text, "code")
	}
	assertEntity(t, msg, 0, Entity{Type: Code, Offset: 0, Length: 4})
}

func TestConvert_FencedCodeBlock(t *testing.T) {
	msg := Convert("```go\nfmt.Println()\n```")
	if msg.Text != "fmt.Println()" {
		t.Errorf("text = %q, want %q", msg.Text, "fmt.Println()")
	}
	assertEntity(t, msg, 0, Entity{Type: Pre, Offset: 0, Length: UTF16Len("fmt.Println()"), Language: "go"})
}

func TestConvert_FencedCodeBlockNoLang(t *testing.T) {
	msg := Convert("```\nhello\n```")
	if msg.Text != "hello" {
		t.Errorf("text = %q, want %q", msg.Text, "hello")
	}
	assertEntity(t, msg, 0, Entity{Type: Pre, Offset: 0, Length: 5})
}

func TestConvert_Link(t *testing.T) {
	msg := Convert("[click here](https://example.com)")
	if msg.Text != "click here" {
		t.Errorf("text = %q, want %q", msg.Text, "click here")
	}
	assertEntity(t, msg, 0, Entity{Type: TextLink, Offset: 0, Length: 10, URL: "https://example.com"})
}

func TestConvert_Image(t *testing.T) {
	msg := Convert("![alt text](https://img.png)")
	if !strings.Contains(msg.Text, "alt text") {
		t.Errorf("text = %q, should contain 'alt text'", msg.Text)
	}
	// Should have a TextLink entity.
	hasLink := false
	for _, e := range msg.Entities {
		if e.Type == TextLink && e.URL == "https://img.png" {
			hasLink = true
		}
	}
	if !hasLink {
		t.Error("missing TextLink entity for image")
	}
}

func TestConvert_Heading(t *testing.T) {
	msg := Convert("# Title")
	// Default h1 symbol is 📌
	if !strings.Contains(msg.Text, "📌") {
		t.Errorf("text = %q, should contain heading symbol", msg.Text)
	}
	if !strings.Contains(msg.Text, "Title") {
		t.Errorf("text = %q, should contain 'Title'", msg.Text)
	}
	// H1 should have Bold and Underline entities.
	hasBold := false
	hasUnderline := false
	for _, e := range msg.Entities {
		if e.Type == Bold {
			hasBold = true
		}
		if e.Type == Underline {
			hasUnderline = true
		}
	}
	if !hasBold {
		t.Error("missing Bold entity for heading")
	}
	if !hasUnderline {
		t.Error("missing Underline entity for H1 heading")
	}
}

func TestConvert_HeadingLevels(t *testing.T) {
	symbols := [6]string{"📌", "✏", "📚", "🔖", "", ""}
	for i := 1; i <= 6; i++ {
		md := strings.Repeat("#", i) + " Level"
		msg := Convert(md)
		if symbols[i-1] != "" && !strings.Contains(msg.Text, symbols[i-1]) {
			t.Errorf("h%d: text = %q, should contain %q", i, msg.Text, symbols[i-1])
		}
	}
}

func TestConvert_Blockquote(t *testing.T) {
	msg := Convert("> quoted text")
	if !strings.Contains(msg.Text, "quoted text") {
		t.Errorf("text = %q, should contain 'quoted text'", msg.Text)
	}
	hasBlockquote := false
	for _, e := range msg.Entities {
		if e.Type == Blockquote {
			hasBlockquote = true
		}
	}
	if !hasBlockquote {
		t.Error("missing Blockquote entity")
	}
}

func TestConvert_UnorderedList(t *testing.T) {
	msg := Convert("- one\n- two\n- three")
	if !strings.Contains(msg.Text, "⦁ one") {
		t.Errorf("text = %q, should contain '• one'", msg.Text)
	}
	if !strings.Contains(msg.Text, "⦁ two") {
		t.Errorf("text = %q, should contain '• two'", msg.Text)
	}
}

func TestConvert_OrderedList(t *testing.T) {
	msg := Convert("1. first\n2. second\n3. third")
	if !strings.Contains(msg.Text, "1. first") {
		t.Errorf("text = %q, should contain '1. first'", msg.Text)
	}
	if !strings.Contains(msg.Text, "2. second") {
		t.Errorf("text = %q, should contain '2. second'", msg.Text)
	}
}

func TestConvert_NestedList(t *testing.T) {
	msg := Convert("- outer\n  - inner")
	if !strings.Contains(msg.Text, "⦁ outer") {
		t.Errorf("text = %q, should contain '• outer'", msg.Text)
	}
	// Inner item should be indented.
	if !strings.Contains(msg.Text, "  ⦁ inner") {
		t.Errorf("text = %q, should contain '  • inner'", msg.Text)
	}
}

func TestConvert_TaskList(t *testing.T) {
	msg := Convert("- [x] done\n- [ ] todo")
	if !strings.Contains(msg.Text, "✅ done") {
		t.Errorf("text = %q, should contain '✅ done'", msg.Text)
	}
	if !strings.Contains(msg.Text, "☐ todo") {
		t.Errorf("text = %q, should contain '☐ todo'", msg.Text)
	}
}

func TestConvert_Table(t *testing.T) {
	md := "| Name | Age |\n|------|-----|\n| Alice | 30 |"
	msg := Convert(md)
	// Should render as monospace (Pre entity).
	hasPre := false
	for _, e := range msg.Entities {
		if e.Type == Pre {
			hasPre = true
		}
	}
	if !hasPre {
		t.Error("missing Pre entity for table")
	}
	if !strings.Contains(msg.Text, "Alice") {
		t.Errorf("text = %q, should contain 'Alice'", msg.Text)
	}
}

func TestConvert_ThematicBreak(t *testing.T) {
	msg := Convert("above\n\n---\n\nbelow")
	if !strings.Contains(msg.Text, "————————") {
		t.Errorf("text = %q, should contain em-dash separator", msg.Text)
	}
}

func TestConvert_PlainText(t *testing.T) {
	msg := Convert("just plain text")
	if msg.Text != "just plain text" {
		t.Errorf("text = %q, want %q", msg.Text, "just plain text")
	}
	if len(msg.Entities) != 0 {
		t.Errorf("entities = %d, want 0", len(msg.Entities))
	}
}

func TestConvert_Empty(t *testing.T) {
	msg := Convert("")
	if msg.Text != "" {
		t.Errorf("text = %q, want empty", msg.Text)
	}
}

func TestConvert_SpecialChars(t *testing.T) {
	// Unlike MarkdownV2, entities don't need escaping. Dots, exclamation marks etc. pass through.
	msg := Convert("file.go says Hello!")
	if msg.Text != "file.go says Hello!" {
		t.Errorf("text = %q, want %q", msg.Text, "file.go says Hello!")
	}
}

func TestConvert_EmojiUTF16(t *testing.T) {
	// Verify UTF-16 offsets are correct with emoji.
	msg := Convert("🌍 **bold**")
	// "🌍 " = 3 UTF-16 units (2 for emoji + 1 for space), then "bold" = 4 units
	hasBold := false
	for _, e := range msg.Entities {
		if e.Type == Bold {
			hasBold = true
			if e.Offset != 3 {
				t.Errorf("bold offset = %d, want 3 (after emoji + space)", e.Offset)
			}
			if e.Length != 4 {
				t.Errorf("bold length = %d, want 4", e.Length)
			}
		}
	}
	if !hasBold {
		t.Error("missing Bold entity")
	}
}

func TestConvert_NestedFormatting(t *testing.T) {
	msg := Convert("**bold and *italic* inside**")
	hasBold := false
	hasItalic := false
	for _, e := range msg.Entities {
		if e.Type == Bold {
			hasBold = true
		}
		if e.Type == Italic {
			hasItalic = true
		}
	}
	if !hasBold {
		t.Error("missing Bold entity")
	}
	if !hasItalic {
		t.Error("missing Italic entity inside bold")
	}
}

func TestConvert_MixedContent(t *testing.T) {
	md := `# Summary

Here is **bold** text and *italic* text.

` + "```go\nfmt.Println(\"hello\")\n```" + `

- Item one
- Item two

> A blockquote

[A link](https://example.com)
`
	msg := Convert(md)

	// Should have text content.
	if msg.Text == "" {
		t.Error("text is empty")
	}

	// Should have entities of various types.
	types := make(map[EntityType]bool)
	for _, e := range msg.Entities {
		types[e.Type] = true
	}
	for _, want := range []EntityType{Bold, Italic, Pre, Blockquote, TextLink} {
		if !types[want] {
			t.Errorf("missing entity type %q", want)
		}
	}
}

func TestConvert_CustomOptions(t *testing.T) {
	msg := Convert("# Title",
		WithHeadingSymbols([6]string{"H1", "H2", "H3", "H4", "H5", "H6"}),
	)
	if !strings.Contains(msg.Text, "H1") {
		t.Errorf("text = %q, should contain custom heading symbol 'H1'", msg.Text)
	}
}

func TestConvert_CustomBullet(t *testing.T) {
	msg := Convert("- item", WithBulletMarker("•"))
	if !strings.Contains(msg.Text, "• item") {
		t.Errorf("text = %q, should use custom bullet '•'", msg.Text)
	}
}

func TestConvert_CodeBlockWithBackticks(t *testing.T) {
	// Code containing backticks should pass through as-is (no escaping needed with entities).
	md := "```\nfoo `bar` baz\n```"
	msg := Convert(md)
	if !strings.Contains(msg.Text, "`bar`") {
		t.Errorf("text = %q, should contain backticks in code block", msg.Text)
	}
}

func TestConvert_MultiParagraph(t *testing.T) {
	msg := Convert("First paragraph.\n\nSecond paragraph.")
	if !strings.Contains(msg.Text, "First paragraph.") {
		t.Errorf("text = %q, missing first paragraph", msg.Text)
	}
	if !strings.Contains(msg.Text, "Second paragraph.") {
		t.Errorf("text = %q, missing second paragraph", msg.Text)
	}
	// Should have double newline between paragraphs.
	if !strings.Contains(msg.Text, "\n\n") {
		t.Errorf("text = %q, should have double newline between paragraphs", msg.Text)
	}
}

func TestConvert_BoldInListItem(t *testing.T) {
	msg := Convert("- **bold item**")
	if !strings.Contains(msg.Text, "bold item") {
		t.Errorf("text = %q, should contain 'bold item'", msg.Text)
	}
	hasBold := false
	for _, e := range msg.Entities {
		if e.Type == Bold {
			hasBold = true
		}
	}
	if !hasBold {
		t.Error("missing Bold entity in list item")
	}
}

// assertEntity checks that msg has an entity at index i matching the expected values.
func assertEntity(t *testing.T, msg Message, i int, want Entity) {
	t.Helper()
	if len(msg.Entities) <= i {
		t.Fatalf("expected at least %d entities, got %d: %+v", i+1, len(msg.Entities), msg.Entities)
	}
	got := msg.Entities[i]
	if got.Type != want.Type {
		t.Errorf("entity[%d].Type = %q, want %q", i, got.Type, want.Type)
	}
	if got.Offset != want.Offset {
		t.Errorf("entity[%d].Offset = %d, want %d", i, got.Offset, want.Offset)
	}
	if got.Length != want.Length {
		t.Errorf("entity[%d].Length = %d, want %d", i, got.Length, want.Length)
	}
	if want.URL != "" && got.URL != want.URL {
		t.Errorf("entity[%d].URL = %q, want %q", i, got.URL, want.URL)
	}
	if want.Language != "" && got.Language != want.Language {
		t.Errorf("entity[%d].Language = %q, want %q", i, got.Language, want.Language)
	}
}
