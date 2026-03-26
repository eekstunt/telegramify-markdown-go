package tgmd

import "testing"

func TestSplit_ShortMessage(t *testing.T) {
	msg := Message{Text: "short", Entities: []Entity{{Type: Bold, Offset: 0, Length: 5}}}
	msgs := Split(msg, 4096)
	if len(msgs) != 1 {
		t.Fatalf("got %d messages, want 1", len(msgs))
	}
	if msgs[0].Text != "short" {
		t.Errorf("text = %q, want %q", msgs[0].Text, "short")
	}
}

func TestSplit_SplitsAtNewline(t *testing.T) {
	text := "line1\nline2\nline3"
	msg := Message{Text: text}
	// maxLen=7 should split after "line1\n" (6 UTF-16 units).
	msgs := Split(msg, 7)
	if len(msgs) < 2 {
		t.Fatalf("got %d messages, want >= 2", len(msgs))
	}
	if msgs[0].Text != "line1\n" {
		t.Errorf("chunk[0] = %q, want %q", msgs[0].Text, "line1\n")
	}
}

func TestSplit_EntityFullyInFirstChunk(t *testing.T) {
	msg := Message{
		Text:     "bold plain text here",
		Entities: []Entity{{Type: Bold, Offset: 0, Length: 4}},
	}
	msgs := Split(msg, 10)
	if len(msgs) < 2 {
		t.Fatalf("got %d messages, want >= 2", len(msgs))
	}
	// Bold entity (0-4) should be in first chunk.
	found := false
	for _, e := range msgs[0].Entities {
		if e.Type == Bold {
			found = true
		}
	}
	if !found {
		t.Error("Bold entity not found in first chunk")
	}
}

func TestSplit_EntityFullyInSecondChunk(t *testing.T) {
	msg := Message{
		Text:     "plain text bold end",
		Entities: []Entity{{Type: Bold, Offset: 11, Length: 4}}, // "bold"
	}
	msgs := Split(msg, 12)
	if len(msgs) < 2 {
		t.Fatalf("got %d messages, want >= 2", len(msgs))
	}
	// Bold entity should be in second chunk with adjusted offset.
	found := false
	for _, e := range msgs[1].Entities {
		if e.Type == Bold {
			found = true
			if e.Offset < 0 {
				t.Errorf("entity offset = %d, should be >= 0", e.Offset)
			}
		}
	}
	if !found {
		t.Error("Bold entity not found in second chunk")
	}
}

func TestSplit_EntitySpansBoundary(t *testing.T) {
	// Bold covers "longbold" (8 chars), split at 5.
	msg := Message{
		Text:     "longboldtext",
		Entities: []Entity{{Type: Bold, Offset: 0, Length: 8}},
	}
	msgs := Split(msg, 5)
	if len(msgs) < 2 {
		t.Fatalf("got %d messages, want >= 2", len(msgs))
	}
	// First chunk should have clipped Bold entity.
	if len(msgs[0].Entities) == 0 {
		t.Fatal("first chunk has no entities")
	}
	if msgs[0].Entities[0].Length != 5 {
		t.Errorf("first chunk entity length = %d, want 5", msgs[0].Entities[0].Length)
	}
	// Second chunk should have continuation Bold entity starting at 0.
	if len(msgs[1].Entities) == 0 {
		t.Fatal("second chunk has no entities")
	}
	if msgs[1].Entities[0].Offset != 0 {
		t.Errorf("second chunk entity offset = %d, want 0", msgs[1].Entities[0].Offset)
	}
	if msgs[1].Entities[0].Length != 3 {
		t.Errorf("second chunk entity length = %d, want 3", msgs[1].Entities[0].Length)
	}
}

func TestSplit_EmptyMessage(t *testing.T) {
	msgs := Split(Message{}, 4096)
	if len(msgs) != 1 {
		t.Fatalf("got %d messages, want 1", len(msgs))
	}
}

func TestSplit_EmojiPreservesUTF16(t *testing.T) {
	// "🌍🌍🌍" = 6 UTF-16 units. Split at 4 should give "🌍🌍" + "🌍".
	msg := Message{Text: "🌍🌍🌍"}
	msgs := Split(msg, 4)
	if len(msgs) != 2 {
		t.Fatalf("got %d messages, want 2", len(msgs))
	}
	if msgs[0].Text != "🌍🌍" {
		t.Errorf("chunk[0] = %q, want %q", msgs[0].Text, "🌍🌍")
	}
	if msgs[1].Text != "🌍" {
		t.Errorf("chunk[1] = %q, want %q", msgs[1].Text, "🌍")
	}
}

func TestSplit_PreservesMultipleEntities(t *testing.T) {
	msg := Message{
		Text: "bold and italic text",
		Entities: []Entity{
			{Type: Bold, Offset: 0, Length: 4},
			{Type: Italic, Offset: 9, Length: 6},
		},
	}
	msgs := Split(msg, 10)
	if len(msgs) < 2 {
		t.Fatalf("got %d messages, want >= 2", len(msgs))
	}
	// Bold should be in chunk 0.
	hasBold := false
	for _, e := range msgs[0].Entities {
		if e.Type == Bold {
			hasBold = true
		}
	}
	if !hasBold {
		t.Error("Bold entity not in first chunk")
	}
}
