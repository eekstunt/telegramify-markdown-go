package tgmd

// Smoke tests using telegramify-markdown's fixture files (exp1.md, exp2.md).
// These verify our library doesn't crash on complex real-world markdown
// and produces valid output with correct entity offsets.

import (
	"os"
	"strings"
	"testing"
)

func TestFixture_Exp1_NoCrash(t *testing.T) {
	data, err := os.ReadFile("testdata/exp1.md")
	if err != nil {
		t.Fatalf("failed to read exp1.md: %v", err)
	}
	msg := Convert(string(data))
	if msg.Text == "" {
		t.Fatal("Convert produced empty text for exp1.md")
	}
}

func TestFixture_Exp2_NoCrash(t *testing.T) {
	data, err := os.ReadFile("testdata/exp2.md")
	if err != nil {
		t.Fatalf("failed to read exp2.md: %v", err)
	}
	msg := Convert(string(data))
	if msg.Text == "" {
		t.Fatal("Convert produced empty text for exp2.md")
	}
}

func TestFixture_Exp1_HasEntities(t *testing.T) {
	data, err := os.ReadFile("testdata/exp1.md")
	if err != nil {
		t.Fatalf("failed to read exp1.md: %v", err)
	}
	msg := Convert(string(data))

	// exp1.md contains many formatting elements — should produce entities.
	if len(msg.Entities) == 0 {
		t.Fatal("Convert produced no entities for exp1.md")
	}

	// Check for expected entity types.
	types := make(map[EntityType]bool)
	for _, e := range msg.Entities {
		types[e.Type] = true
	}
	for _, want := range []EntityType{Bold, Italic, Strikethrough, Code, Pre, TextLink, Blockquote} {
		if !types[want] {
			t.Errorf("exp1.md: missing entity type %q", want)
		}
	}
}

func TestFixture_Exp2_HasEntities(t *testing.T) {
	data, err := os.ReadFile("testdata/exp2.md")
	if err != nil {
		t.Fatalf("failed to read exp2.md: %v", err)
	}
	msg := Convert(string(data))

	if len(msg.Entities) == 0 {
		t.Fatal("Convert produced no entities for exp2.md")
	}

	types := make(map[EntityType]bool)
	for _, e := range msg.Entities {
		types[e.Type] = true
	}
	for _, want := range []EntityType{Bold, Code, Pre} {
		if !types[want] {
			t.Errorf("exp2.md: missing entity type %q", want)
		}
	}
}

func TestFixture_Exp1_ContentPresent(t *testing.T) {
	data, err := os.ReadFile("testdata/exp1.md")
	if err != nil {
		t.Fatalf("failed to read exp1.md: %v", err)
	}
	msg := Convert(string(data))

	// Key content from exp1.md that should survive conversion.
	for _, substr := range []string{
		"TEST1", "TEST2", "TEST3",
		"Bold text", "Italic text", "Strikethrough text",
		"Inline code", "Code block",
		"Hello, World!",
		"inline URL",
		"Blockquote text",
		"Horizontal Rule",
		"item", "nested item",
		"numbered item",
		"Uncompleted task list item",
		"Completed task list item",
	} {
		if !strings.Contains(msg.Text, substr) {
			t.Errorf("exp1.md: text should contain %q", substr)
		}
	}
}

func TestFixture_Exp2_ContentPresent(t *testing.T) {
	data, err := os.ReadFile("testdata/exp2.md")
	if err != nil {
		t.Fatalf("failed to read exp2.md: %v", err)
	}
	msg := Convert(string(data))

	for _, substr := range []string{
		"0-1 Knapsack Problem",
		"Problem Description",
		"Recursive + Memoization",
		"Dynamic Programming",
		"knapsack_recursive_memo",
		"knapsack_dp",
		"knapsack_dp_optimized",
		"Key Points",
		"Time Complexity",
		"Space Complexity",
		"React UI Example",
		"KnapsackProblem",
	} {
		if !strings.Contains(msg.Text, substr) {
			t.Errorf("exp2.md: text should contain %q", substr)
		}
	}
}

func TestFixture_Exp1_ValidEntityOffsets(t *testing.T) {
	data, err := os.ReadFile("testdata/exp1.md")
	if err != nil {
		t.Fatalf("failed to read exp1.md: %v", err)
	}
	msg := Convert(string(data))
	validateEntityOffsets(t, "exp1.md", msg)
}

func TestFixture_Exp2_ValidEntityOffsets(t *testing.T) {
	data, err := os.ReadFile("testdata/exp2.md")
	if err != nil {
		t.Fatalf("failed to read exp2.md: %v", err)
	}
	msg := Convert(string(data))
	validateEntityOffsets(t, "exp2.md", msg)
}

func TestFixture_Exp2_SplitDoesNotCrash(t *testing.T) {
	data, err := os.ReadFile("testdata/exp2.md")
	if err != nil {
		t.Fatalf("failed to read exp2.md: %v", err)
	}
	msgs := ConvertAndSplit(string(data))
	if len(msgs) == 0 {
		t.Fatal("ConvertAndSplit produced no messages for exp2.md")
	}
	// exp2.md is long enough to require splitting.
	if len(msgs) < 2 {
		t.Logf("exp2.md produced %d message(s) — may not need splitting at 4096", len(msgs))
	}
	// Verify all chunks have valid offsets.
	for i, msg := range msgs {
		validateEntityOffsets(t, "exp2.md chunk "+string(rune('0'+i)), msg)
	}
}

// validateEntityOffsets checks that all entity offsets and lengths are within text bounds.
func validateEntityOffsets(t *testing.T, name string, msg Message) {
	t.Helper()
	textLen := UTF16Len(msg.Text)
	for i, e := range msg.Entities {
		if e.Offset < 0 {
			t.Errorf("%s: entity[%d] offset %d is negative", name, i, e.Offset)
		}
		if e.Length <= 0 {
			t.Errorf("%s: entity[%d] length %d is non-positive", name, i, e.Length)
		}
		if e.Offset+e.Length > textLen {
			t.Errorf("%s: entity[%d] offset=%d + length=%d = %d exceeds text UTF-16 length %d",
				name, i, e.Offset, e.Length, e.Offset+e.Length, textLen)
		}
	}
}
