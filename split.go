package tgmd

import "strings"

// Split divides a Message into multiple Messages, each within maxLen UTF-16 code units.
// Entities spanning a split boundary are clipped into both chunks.
func Split(msg Message, maxLen int) []Message {
	if UTF16Len(msg.Text) <= maxLen {
		return []Message{msg}
	}
	return splitRecursive(msg.Text, msg.Entities, maxLen)
}

func splitRecursive(text string, entities []Entity, maxLen int) []Message {
	textLen := UTF16Len(text)
	if textLen <= maxLen {
		return []Message{{Text: text, Entities: entities}}
	}

	splitAt := findSplitPoint(text, maxLen)
	prefix, suffix := SplitAtUTF16(text, splitAt)

	var chunkEntities, restEntities []Entity
	for _, e := range entities {
		eEnd := e.Offset + e.Length

		if eEnd <= splitAt {
			// Fully in current chunk.
			chunkEntities = append(chunkEntities, e)
		} else if e.Offset >= splitAt {
			// Fully in next chunk — shift offset.
			restEntities = append(restEntities, Entity{
				Type:     e.Type,
				Offset:   e.Offset - splitAt,
				Length:   e.Length,
				URL:      e.URL,
				Language: e.Language,
			})
		} else {
			// Spans boundary — clip into both.
			leftLen := splitAt - e.Offset
			rightLen := e.Length - leftLen

			if leftLen > 0 {
				chunkEntities = append(chunkEntities, Entity{
					Type:     e.Type,
					Offset:   e.Offset,
					Length:   leftLen,
					URL:      e.URL,
					Language: e.Language,
				})
			}
			if rightLen > 0 {
				restEntities = append(restEntities, Entity{
					Type:     e.Type,
					Offset:   0,
					Length:   rightLen,
					URL:      e.URL,
					Language: e.Language,
				})
			}
		}
	}

	result := []Message{{Text: prefix, Entities: chunkEntities}}
	result = append(result, splitRecursive(suffix, restEntities, maxLen)...)
	return result
}

// findSplitPoint finds the best UTF-16 offset to split text at, preferring newline boundaries.
func findSplitPoint(text string, maxLen int) int {
	prefix, _ := SplitAtUTF16(text, maxLen)

	// Prefer splitting at a newline.
	if idx := strings.LastIndex(prefix, "\n"); idx > 0 {
		return UTF16Len(prefix[:idx+1])
	}

	// Fall back to space.
	if idx := strings.LastIndex(prefix, " "); idx > 0 {
		return UTF16Len(prefix[:idx+1])
	}

	// Hard split.
	return maxLen
}
