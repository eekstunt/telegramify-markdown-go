package tgmd

// Compatibility tests inspired by Python's telegramify-markdown test suite.
// See: https://github.com/sudoskys/telegramify-markdown/tree/main/tests

import (
	"strings"
	"testing"
)

// --- test_converter.py: Core Conversion Tests ---

func TestCompat_BoldSurrounded(t *testing.T) {
	msg := Convert("foo **bar** baz")
	if !strings.Contains(msg.Text, "bar") {
		t.Fatalf("text = %q, want 'bar'", msg.Text)
	}
	assertHasEntity(t, msg, Bold, "bar")
}

func TestCompat_NestedBoldItalic(t *testing.T) {
	msg := Convert("**bold *italic* bold**")
	assertHasEntityType(t, msg, Bold)
	assertHasEntityType(t, msg, Italic)
	// Italic should cover "italic".
	assertHasEntity(t, msg, Italic, "italic")
}

func TestCompat_InlineCodeInContext(t *testing.T) {
	msg := Convert("use `print()` here")
	assertHasEntity(t, msg, Code, "print()")
	if !strings.Contains(msg.Text, "use") {
		t.Errorf("text = %q, should contain 'use'", msg.Text)
	}
	if !strings.Contains(msg.Text, "here") {
		t.Errorf("text = %q, should contain 'here'", msg.Text)
	}
}

func TestCompat_CodeBlockPython(t *testing.T) {
	msg := Convert("```python\nprint('hello')\n```")
	assertHasEntity(t, msg, Pre, "print('hello')")
	// Check language.
	for _, e := range msg.Entities {
		if e.Type == Pre && e.Language != "python" {
			t.Errorf("Pre entity language = %q, want 'python'", e.Language)
		}
	}
}

func TestCompat_CodeBlockNoLang(t *testing.T) {
	msg := Convert("```\nsome code\n```")
	assertHasEntity(t, msg, Pre, "some code")
	for _, e := range msg.Entities {
		if e.Type == Pre && e.Language != "" {
			t.Errorf("Pre entity language = %q, want empty", e.Language)
		}
	}
}

func TestCompat_HeadingH1(t *testing.T) {
	msg := Convert("# Title")
	if !strings.Contains(msg.Text, "📌") {
		t.Errorf("text = %q, should contain 📌", msg.Text)
	}
	assertHasEntityType(t, msg, Bold)
}

func TestCompat_HeadingH1_Underline(t *testing.T) {
	msg := Convert("# Title")
	assertHasEntityType(t, msg, Bold)
	assertHasEntityType(t, msg, Underline)
}

func TestCompat_HeadingH2(t *testing.T) {
	msg := Convert("## Subtitle")
	if !strings.Contains(msg.Text, "✏") {
		t.Errorf("text = %q, should contain ✏", msg.Text)
	}
	assertHasEntityType(t, msg, Bold)
	assertHasEntityType(t, msg, Underline)
}

func TestCompat_HeadingH3(t *testing.T) {
	msg := Convert("### Section")
	if !strings.Contains(msg.Text, "📚") {
		t.Errorf("text = %q, should contain 📚", msg.Text)
	}
	assertHasEntityType(t, msg, Bold)
}

func TestCompat_HeadingH4(t *testing.T) {
	msg := Convert("#### Sub")
	if !strings.Contains(msg.Text, "🔖") {
		t.Errorf("text = %q, should contain 🔖", msg.Text)
	}
	assertHasEntityType(t, msg, Bold)
}

func TestCompat_HeadingH5_Italic(t *testing.T) {
	msg := Convert("##### H5")
	if !strings.Contains(msg.Text, "H5") {
		t.Fatalf("text = %q, should contain 'H5'", msg.Text)
	}
	assertHasEntityType(t, msg, Italic)
}

func TestCompat_HeadingH6_Italic(t *testing.T) {
	msg := Convert("###### H6")
	if !strings.Contains(msg.Text, "H6") {
		t.Fatalf("text = %q, should contain 'H6'", msg.Text)
	}
	assertHasEntityType(t, msg, Italic)
}

func TestCompat_Link(t *testing.T) {
	msg := Convert("[Google](https://google.com)")
	if !strings.Contains(msg.Text, "Google") {
		t.Fatalf("text = %q, should contain 'Google'", msg.Text)
	}
	found := false
	for _, e := range msg.Entities {
		if e.Type == TextLink && e.URL == "https://google.com" {
			found = true
		}
	}
	if !found {
		t.Error("missing TextLink entity with url 'https://google.com'")
	}
}

func TestCompat_Image(t *testing.T) {
	msg := Convert("![alt](https://example.com/img.png)")
	found := false
	for _, e := range msg.Entities {
		if e.Type == TextLink && e.URL == "https://example.com/img.png" {
			found = true
		}
	}
	if !found {
		t.Error("missing TextLink entity for image URL")
	}
}

func TestCompat_Blockquote(t *testing.T) {
	msg := Convert("> quoted text")
	if !strings.Contains(msg.Text, "quoted text") {
		t.Fatalf("text = %q, should contain 'quoted text'", msg.Text)
	}
	assertHasEntityType(t, msg, Blockquote)
}

func TestCompat_Table(t *testing.T) {
	msg := Convert("| a | b |\n| --- | --- |\n| 1 | 2 |")
	assertHasEntityType(t, msg, Pre)
	if !strings.Contains(msg.Text, "a") {
		t.Errorf("text = %q, should contain 'a'", msg.Text)
	}
	if !strings.Contains(msg.Text, "b") {
		t.Errorf("text = %q, should contain 'b'", msg.Text)
	}
	if !strings.Contains(msg.Text, "1") {
		t.Errorf("text = %q, should contain '1'", msg.Text)
	}
	if !strings.Contains(msg.Text, "2") {
		t.Errorf("text = %q, should contain '2'", msg.Text)
	}
}

func TestCompat_UnorderedList(t *testing.T) {
	msg := Convert("- item1\n- item2")
	if !strings.Contains(msg.Text, "item1") {
		t.Errorf("text = %q, should contain 'item1'", msg.Text)
	}
	if !strings.Contains(msg.Text, "item2") {
		t.Errorf("text = %q, should contain 'item2'", msg.Text)
	}
}

func TestCompat_OrderedList(t *testing.T) {
	msg := Convert("1. first\n2. second")
	if !strings.Contains(msg.Text, "1. first") {
		t.Errorf("text = %q, should contain '1. first'", msg.Text)
	}
	if !strings.Contains(msg.Text, "2. second") {
		t.Errorf("text = %q, should contain '2. second'", msg.Text)
	}
}

func TestCompat_TaskListChecked(t *testing.T) {
	msg := Convert("- [x] done")
	if !strings.Contains(msg.Text, "✅") {
		t.Errorf("text = %q, should contain ✅", msg.Text)
	}
	if !strings.Contains(msg.Text, "done") {
		t.Errorf("text = %q, should contain 'done'", msg.Text)
	}
}

func TestCompat_TaskListUnchecked(t *testing.T) {
	msg := Convert("- [ ] todo")
	if !strings.Contains(msg.Text, "☐") {
		t.Errorf("text = %q, should contain ☐", msg.Text)
	}
	if !strings.Contains(msg.Text, "todo") {
		t.Errorf("text = %q, should contain 'todo'", msg.Text)
	}
}

func TestCompat_TaskListNoBullet(t *testing.T) {
	// Task list items should NOT have a bullet marker — the checkbox replaces it.
	// telegramify-markdown: "- [x] done" → no "⦁" in text
	msg := Convert("- [x] done\n- [ ] todo")
	if strings.Contains(msg.Text, "•") {
		t.Errorf("text = %q, should not contain bullet marker alongside checkbox", msg.Text)
	}
}

func TestCompat_HorizontalRule(t *testing.T) {
	msg := Convert("above\n\n---\n\nbelow")
	if !strings.Contains(msg.Text, "————————") {
		t.Errorf("text = %q, should contain em-dash separator", msg.Text)
	}
}

func TestCompat_ParagraphSpacing(t *testing.T) {
	msg := Convert("para1\n\npara2")
	if !strings.Contains(msg.Text, "para1") {
		t.Errorf("text = %q, should contain 'para1'", msg.Text)
	}
	if !strings.Contains(msg.Text, "para2") {
		t.Errorf("text = %q, should contain 'para2'", msg.Text)
	}
	if !strings.Contains(msg.Text, "\n\n") {
		t.Errorf("text = %q, should have double newline between paragraphs", msg.Text)
	}
}

func TestCompat_HeadingThenContent(t *testing.T) {
	msg := Convert("# Title\n\nContent")
	if !strings.Contains(msg.Text, "Title") {
		t.Errorf("text = %q, should contain 'Title'", msg.Text)
	}
	if !strings.Contains(msg.Text, "Content") {
		t.Errorf("text = %q, should contain 'Content'", msg.Text)
	}
}

// --- test_converter.py: UTF-16 Offset Tests ---

func TestCompat_UTF16_EmojiOffset(t *testing.T) {
	// "📌 **bold**" — 📌 is 2 UTF-16 units, space is 1, so bold starts at offset 3.
	msg := Convert("📌 **bold**")
	for _, e := range msg.Entities {
		if e.Type == Bold {
			if e.Offset != 3 {
				t.Errorf("bold offset = %d, want 3", e.Offset)
			}
			if e.Length != 4 {
				t.Errorf("bold length = %d, want 4", e.Length)
			}
			return
		}
	}
	t.Error("missing Bold entity")
}

func TestCompat_UTF16_ChineseOffset(t *testing.T) {
	// "你好 **世界**" — 你好 are BMP (1 UTF-16 unit each), space is 1, so bold at offset 3.
	// 世界 are BMP, so length = 2.
	msg := Convert("你好 **世界**")
	for _, e := range msg.Entities {
		if e.Type == Bold {
			if e.Offset != 3 {
				t.Errorf("bold offset = %d, want 3", e.Offset)
			}
			if e.Length != 2 {
				t.Errorf("bold length = %d, want 2", e.Length)
			}
			return
		}
	}
	t.Error("missing Bold entity")
}

// --- test_converter.py: Complex Document ---

func TestCompat_ComplexDocument(t *testing.T) {
	md := `# Hello World

This is **bold** and *italic* text.

- item 1
- item 2

> A quote

` + "```python\nprint(\"hello\")\n```"

	msg := Convert(md)

	types := make(map[EntityType]bool)
	for _, e := range msg.Entities {
		types[e.Type] = true
	}
	for _, want := range []EntityType{Bold, Italic, Blockquote, Pre} {
		if !types[want] {
			t.Errorf("missing entity type %q", want)
		}
	}
	for _, substr := range []string{"Hello World", "item 1", "A quote", `print("hello")`} {
		if !strings.Contains(msg.Text, substr) {
			t.Errorf("text should contain %q", substr)
		}
	}
}

// --- test_converter.py: Nested Lists ---

func TestCompat_NestedList_ParentChild(t *testing.T) {
	msg := Convert("- parent\n    - child")
	if !strings.Contains(msg.Text, "parent") {
		t.Errorf("text = %q, should contain 'parent'", msg.Text)
	}
	if !strings.Contains(msg.Text, "child") {
		t.Errorf("text = %q, should contain 'child'", msg.Text)
	}
}

func TestCompat_NestedList_ThreeLevels(t *testing.T) {
	msg := Convert("- a\n    - b\n        - c")
	for _, item := range []string{"a", "b", "c"} {
		if !strings.Contains(msg.Text, item) {
			t.Errorf("text = %q, should contain %q", msg.Text, item)
		}
	}
	// Should produce at least 3 lines.
	lines := strings.Split(strings.TrimSpace(msg.Text), "\n")
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines, got %d: %q", len(lines), msg.Text)
	}
}

func TestCompat_MixedOrderedUnordered(t *testing.T) {
	msg := Convert("1. step\n    - detail")
	if !strings.Contains(msg.Text, "step") {
		t.Errorf("text = %q, should contain 'step'", msg.Text)
	}
	if !strings.Contains(msg.Text, "detail") {
		t.Errorf("text = %q, should contain 'detail'", msg.Text)
	}
}

// --- test_entity.py: UTF-16 Length Tests ---

func TestCompat_UTF16Len(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"hello", 5},
		{"你好", 2},
		{"📌", 2},
		{"A📌B", 4},
		{"你好世界", 4},
		{"📌✅🔗", 5}, // ✅ is U+2705 (BMP, 1 unit), 📌 and 🔗 are supplementary (2 units each)
		{"A📌B你好C", 7},
	}
	for _, tc := range tests {
		got := UTF16Len(tc.input)
		if got != tc.want {
			t.Errorf("UTF16Len(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

// --- test_entity.py: Split Tests ---

func TestCompat_SplitShortMessage(t *testing.T) {
	msg := Message{
		Text:     "hello",
		Entities: []Entity{{Type: Bold, Offset: 0, Length: 5}},
	}
	chunks := Split(msg, 100)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
	if chunks[0].Text != "hello" {
		t.Errorf("chunk text = %q, want %q", chunks[0].Text, "hello")
	}
	if len(chunks[0].Entities) != 1 {
		t.Errorf("chunk entities = %d, want 1", len(chunks[0].Entities))
	}
}

func TestCompat_SplitEmpty(t *testing.T) {
	msg := Message{Text: "", Entities: nil}
	chunks := Split(msg, 100)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk for empty, got %d", len(chunks))
	}
}

func TestCompat_SplitPreservesContent(t *testing.T) {
	msg := Message{Text: "aaa\nbbb\nccc", Entities: nil}
	chunks := Split(msg, 5)
	combined := ""
	for i, c := range chunks {
		if i > 0 {
			// Account for the newline at the split point being in one of the chunks.
		}
		combined += c.Text
	}
	// All original content should be present.
	if !strings.Contains(combined, "aaa") || !strings.Contains(combined, "bbb") || !strings.Contains(combined, "ccc") {
		t.Errorf("split lost content: combined = %q", combined)
	}
}

func TestCompat_SplitEntityInFirstChunk(t *testing.T) {
	msg := Message{
		Text:     "bold\nnormal",
		Entities: []Entity{{Type: Bold, Offset: 0, Length: 4}},
	}
	chunks := Split(msg, 5)
	if len(chunks) < 2 {
		t.Fatalf("expected >=2 chunks, got %d", len(chunks))
	}
	// First chunk should have the bold entity.
	hasBold := false
	for _, e := range chunks[0].Entities {
		if e.Type == Bold {
			hasBold = true
		}
	}
	if !hasBold {
		t.Error("first chunk should have Bold entity")
	}
}

func TestCompat_SplitEntitySpansBoundary(t *testing.T) {
	text := "aabbcc\nddee"
	msg := Message{
		Text:     text,
		Entities: []Entity{{Type: Bold, Offset: 0, Length: UTF16Len(text)}},
	}
	chunks := Split(msg, 7)
	if len(chunks) < 2 {
		t.Fatalf("expected >=2 chunks, got %d", len(chunks))
	}
	// Both chunks should have a bold entity (clipped).
	for i, chunk := range chunks {
		hasBold := false
		for _, e := range chunk.Entities {
			if e.Type == Bold {
				hasBold = true
			}
		}
		if !hasBold {
			t.Errorf("chunk %d should have Bold entity", i)
		}
	}
}

func TestCompat_SplitEmoji(t *testing.T) {
	msg := Message{Text: "📌\n📌\n📌", Entities: nil}
	chunks := Split(msg, 4)
	combined := ""
	for _, c := range chunks {
		combined += c.Text
	}
	if strings.Count(combined, "📌") != 3 {
		t.Errorf("split lost emoji: combined = %q", combined)
	}
}

// --- test_converter.py: Loose List Paragraphs ---

func TestCompat_LooseListParagraphs(t *testing.T) {
	msg := Convert("- para1\n\n  para2")
	if !strings.Contains(msg.Text, "para1") {
		t.Errorf("text = %q, should contain 'para1'", msg.Text)
	}
	if !strings.Contains(msg.Text, "para2") {
		t.Errorf("text = %q, should contain 'para2'", msg.Text)
	}
}

// --- Integration: Full document with many elements ---

func TestCompat_FullDocument(t *testing.T) {
	md := `# Main Title

## Introduction

This is a **comprehensive** test with *various* formatting elements.

### Features

- Bold: **text**
- Italic: *text*
- Code: ` + "`inline`" + `
- ~~Strikethrough~~

### Code Example

` + "```python\ndef hello():\n    print('world')\n```" + `

### Data Table

| Name | Value |
|------|-------|
| foo  | 42    |
| bar  | 99    |

### Task List

- [x] Completed task
- [ ] Pending task

### Quote

> This is a blockquote
> with multiple lines

---

[Visit Example](https://example.com)
`

	msg := Convert(md)

	// All content should be present.
	for _, substr := range []string{
		"Main Title", "Introduction", "comprehensive", "various",
		"Features", "Code Example", "Data Table", "Task List",
		"hello", "print", "foo", "42", "bar", "99",
		"✅", "Completed task", "☐", "Pending task",
		"blockquote", "Visit Example", "————————",
	} {
		if !strings.Contains(msg.Text, substr) {
			t.Errorf("text should contain %q, got:\n%s", substr, msg.Text)
		}
	}

	// All entity types should be present.
	types := make(map[EntityType]bool)
	for _, e := range msg.Entities {
		types[e.Type] = true
	}
	for _, want := range []EntityType{Bold, Italic, Code, Strikethrough, Pre, Blockquote, TextLink} {
		if !types[want] {
			t.Errorf("missing entity type %q", want)
		}
	}
}

// --- Helpers ---

// assertHasEntityType checks that msg contains at least one entity of the given type.
func assertHasEntityType(t *testing.T, msg Message, typ EntityType) {
	t.Helper()
	for _, e := range msg.Entities {
		if e.Type == typ {
			return
		}
	}
	t.Errorf("missing entity type %q in %+v", typ, msg.Entities)
}

// assertHasEntity checks that msg contains an entity of the given type covering the given substring.
func assertHasEntity(t *testing.T, msg Message, typ EntityType, substr string) {
	t.Helper()
	for _, e := range msg.Entities {
		if e.Type != typ {
			continue
		}
		// Extract the text covered by this entity using UTF-16 offsets.
		_, after := SplitAtUTF16(msg.Text, e.Offset)
		covered, _ := SplitAtUTF16(after, e.Length)
		if covered == substr {
			return
		}
	}
	t.Errorf("no %q entity covering %q in text=%q entities=%+v", typ, substr, msg.Text, msg.Entities)
}
