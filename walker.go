package tgmd

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

// walker walks a goldmark AST and produces plain text + entities.
type walker struct {
	source   []byte // original markdown source bytes
	cfg      config

	buf      strings.Builder // accumulated plain text
	utf16Pos int             // current position in UTF-16 code units

	entities []Entity     // completed entities
	stack    []stackEntry // open entity scopes

	listStack    []listState // nested list tracking
	lastNewlines int         // how many trailing newlines we've written
	inCodeBlock  bool        // suppress entity creation inside code blocks
}

type stackEntry struct {
	entityType EntityType
	startPos   int
	url        string
	language   string
}

type listState struct {
	ordered bool
	counter int
	depth   int
}

// writeText appends text to the buffer and advances the UTF-16 position.
func (w *walker) writeText(s string) {
	if s == "" {
		return
	}
	w.buf.WriteString(s)
	w.utf16Pos += UTF16Len(s)
	// Track trailing newlines for spacing control.
	nlCount := 0
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '\n' {
			nlCount++
		} else {
			break
		}
	}
	if nlCount == len(s) {
		// Entire string is newlines — accumulate with existing count.
		w.lastNewlines += nlCount
	} else if nlCount > 0 {
		w.lastNewlines = nlCount
	} else {
		w.lastNewlines = 0
	}
}

// ensureNewlines ensures at least n trailing newlines exist.
// Does nothing if the buffer is empty (no leading newlines needed).
func (w *walker) ensureNewlines(n int) {
	if w.buf.Len() == 0 {
		return
	}
	for w.lastNewlines < n {
		w.writeText("\n")
	}
}

// pushEntity opens a new entity scope at the current position.
func (w *walker) pushEntity(t EntityType, url, language string) {
	w.stack = append(w.stack, stackEntry{
		entityType: t,
		startPos:   w.utf16Pos,
		url:        url,
		language:   language,
	})
}

// popEntity closes the most recent entity scope and creates the entity.
func (w *walker) popEntity() {
	if len(w.stack) == 0 {
		return
	}
	entry := w.stack[len(w.stack)-1]
	w.stack = w.stack[:len(w.stack)-1]

	length := w.utf16Pos - entry.startPos
	if length <= 0 {
		return
	}
	w.entities = append(w.entities, Entity{
		Type:     entry.entityType,
		Offset:   entry.startPos,
		Length:   length,
		URL:      entry.url,
		Language: entry.language,
	})
}

// walk is the goldmark AST visitor function.
func (w *walker) walk(n ast.Node, entering bool) (ast.WalkStatus, error) {
	switch node := n.(type) {
	case *ast.Document:
		// No-op — just walk children.

	case *ast.Paragraph:
		if !entering {
			// Don't add paragraph spacing inside list items or blockquotes.
			if n.Parent() != nil && n.Parent().Kind() == ast.KindListItem {
				// Only add newline if there are siblings after this paragraph.
				if n.NextSibling() != nil {
					w.ensureNewlines(1)
				}
			} else {
				w.ensureNewlines(2)
			}
		}

	case *ast.Heading:
		if entering {
			w.ensureNewlines(2)
			// Write heading prefix emoji.
			level := node.Level
			if level >= 1 && level <= 6 {
				symbol := w.cfg.headingSymbols[level-1]
				if symbol != "" {
					w.writeText(symbol + " ")
				}
			}
			// H1-H2: bold+underline, H3-H4: bold, H5-H6: italic (matching telegramify-markdown).
			if level <= 2 {
				w.pushEntity(Bold, "", "")
				w.pushEntity(Underline, "", "")
			} else if level <= 4 {
				w.pushEntity(Bold, "", "")
			} else {
				w.pushEntity(Italic, "", "")
			}
		} else {
			w.popEntity()
			if node.Level <= 2 {
				w.popEntity() // pop the second entity (bold, underline was popped first)
			}
			w.ensureNewlines(2)
		}

	case *ast.ThematicBreak:
		if entering {
			w.ensureNewlines(1)
			w.writeText("————————")
			w.ensureNewlines(2)
		}

	case *ast.CodeBlock:
		if entering {
			w.ensureNewlines(1)
			// Push entity AFTER newlines so offset starts at code content.
			w.pushEntity(Pre, "", "")
			w.inCodeBlock = true
			w.writeCodeBlockLines(node)
			return ast.WalkSkipChildren, nil
		}

	case *ast.FencedCodeBlock:
		if entering {
			w.ensureNewlines(1)
			lang := ""
			if node.Language(w.source) != nil {
				lang = string(node.Language(w.source))
			}
			// Push entity AFTER newlines so offset starts at code content.
			w.pushEntity(Pre, "", lang)
			w.inCodeBlock = true
			w.writeCodeBlockLines(node)
			return ast.WalkSkipChildren, nil
		}

	case *ast.Blockquote:
		if entering {
			w.pushEntity(Blockquote, "", "")
		} else {
			w.popEntity()
			w.ensureNewlines(2)
		}

	case *ast.List:
		if entering {
			w.ensureNewlines(1)
			depth := len(w.listStack)
			w.listStack = append(w.listStack, listState{
				ordered: node.IsOrdered(),
				counter: node.Start,
				depth:   depth,
			})
		} else {
			if len(w.listStack) > 0 {
				w.listStack = w.listStack[:len(w.listStack)-1]
			}
			// Only add extra spacing at top-level list exit.
			if len(w.listStack) == 0 {
				w.ensureNewlines(2)
			}
		}

	case *ast.ListItem:
		if entering {
			if len(w.listStack) > 0 {
				ls := &w.listStack[len(w.listStack)-1]
				indent := strings.Repeat("  ", ls.depth)
				if ls.ordered {
					w.writeText(fmt.Sprintf("%s%d. ", indent, ls.counter))
					ls.counter++
				} else {
					// Skip bullet marker if this list item starts with a task checkbox.
					hasCheckbox := false
					for child := node.FirstChild(); child != nil; child = child.NextSibling() {
						// Check the first grandchild of the first child (paragraph).
						if fc := child.FirstChild(); fc != nil {
							if _, ok := fc.(*east.TaskCheckBox); ok {
								hasCheckbox = true
							}
						}
						break
					}
					if !hasCheckbox {
						w.writeText(indent + w.cfg.bulletMarker + " ")
					} else {
						w.writeText(indent)
					}
				}
			}
		} else {
			w.ensureNewlines(1)
		}

	case *ast.Text:
		if entering {
			w.writeText(string(node.Segment.Value(w.source)))
			if node.SoftLineBreak() {
				w.writeText("\n")
			}
			if node.HardLineBreak() {
				w.writeText("\n")
			}
		}

	case *ast.String:
		if entering {
			w.writeText(string(node.Value))
		}

	case *ast.CodeSpan:
		if entering {
			w.pushEntity(Code, "", "")
			w.inCodeBlock = true
			// Walk children manually to get raw text.
			for child := node.FirstChild(); child != nil; child = child.NextSibling() {
				if t, ok := child.(*ast.Text); ok {
					w.writeText(string(t.Segment.Value(w.source)))
				}
			}
			w.inCodeBlock = false
			w.popEntity()
			return ast.WalkSkipChildren, nil
		}

	case *ast.Emphasis:
		if entering {
			if w.inCodeBlock {
				return ast.WalkContinue, nil
			}
			if node.Level == 2 {
				w.pushEntity(Bold, "", "")
			} else {
				w.pushEntity(Italic, "", "")
			}
		} else {
			if !w.inCodeBlock {
				w.popEntity()
			}
		}

	case *ast.Link:
		if entering {
			dest := string(node.Destination)
			w.pushEntity(TextLink, dest, "")
		} else {
			w.popEntity()
		}

	case *ast.AutoLink:
		if entering {
			url := string(node.URL(w.source))
			w.pushEntity(TextLink, url, "")
			w.writeText(url)
			w.popEntity()
			return ast.WalkSkipChildren, nil
		}

	case *ast.Image:
		if entering {
			dest := string(node.Destination)
			w.writeText("🖼 ")
			w.pushEntity(TextLink, dest, "")
			// Walk children for alt text.
			for child := node.FirstChild(); child != nil; child = child.NextSibling() {
				if t, ok := child.(*ast.Text); ok {
					w.writeText(string(t.Segment.Value(w.source)))
				}
			}
			w.popEntity()
			return ast.WalkSkipChildren, nil
		}

	case *ast.RawHTML:
		// Skip raw HTML tags.
		if entering {
			return ast.WalkSkipChildren, nil
		}

	case *ast.HTMLBlock:
		// Skip HTML blocks.
		if entering {
			return ast.WalkSkipChildren, nil
		}

	// GFM extensions
	case *east.Strikethrough:
		if entering {
			if !w.inCodeBlock {
				w.pushEntity(Strikethrough, "", "")
			}
		} else {
			if !w.inCodeBlock {
				w.popEntity()
			}
		}

	case *east.TaskCheckBox:
		if entering {
			if node.IsChecked {
				w.writeText(w.cfg.checkedMarker + " ")
			} else {
				w.writeText(w.cfg.uncheckedMarker + " ")
			}
		}

	case *east.Table:
		if entering {
			w.ensureNewlines(1)
			w.renderTable(node)
			return ast.WalkSkipChildren, nil
		}
	}

	return ast.WalkContinue, nil
}

// writeCodeBlockLines reads lines from a code block node and writes them.
func (w *walker) writeCodeBlockLines(node ast.Node) {
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		text := string(line.Value(w.source))
		w.writeText(text)
	}
	// Strip trailing newline inside the code block.
	s := w.buf.String()
	if strings.HasSuffix(s, "\n") {
		w.buf.Reset()
		trimmed := strings.TrimRight(s, "\n")
		w.buf.WriteString(trimmed)
		w.utf16Pos -= UTF16Len(s) - UTF16Len(trimmed)
		w.lastNewlines = 0
	}
	w.inCodeBlock = false
	w.popEntity()
	w.ensureNewlines(2)
}

// result returns the final message.
func (w *walker) result() Message {
	raw := w.buf.String()

	// Trim leading newlines and adjust entity offsets.
	trimmed := strings.TrimLeft(raw, "\n")
	leadingRemoved := UTF16Len(raw) - UTF16Len(trimmed)

	// Trim trailing newlines.
	text := strings.TrimRight(trimmed, "\n")
	textLen := UTF16Len(text)

	var entities []Entity
	for _, e := range w.entities {
		e.Offset -= leadingRemoved
		if e.Offset < 0 {
			// Entity started in the trimmed leading region.
			e.Length += e.Offset // shrink length
			e.Offset = 0
		}
		if e.Length <= 0 {
			continue
		}
		if e.Offset >= textLen {
			continue
		}
		if e.Offset+e.Length > textLen {
			e.Length = textLen - e.Offset
		}
		if e.Length > 0 {
			entities = append(entities, e)
		}
	}
	return Message{Text: text, Entities: entities}
}
