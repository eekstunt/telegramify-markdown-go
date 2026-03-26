package tgmd

import (
	"strings"
	"unicode/utf8"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

// renderTable walks a GFM table node and renders it as a monospace pre block.
func (w *walker) renderTable(table *east.Table) {
	// Collect all rows of cell text.
	var rows [][]string
	for row := table.FirstChild(); row != nil; row = row.NextSibling() {
		var cells []string
		for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
			cells = append(cells, w.extractText(cell))
		}
		rows = append(rows, cells)
	}

	if len(rows) == 0 {
		return
	}

	// Compute column widths.
	numCols := 0
	for _, row := range rows {
		if len(row) > numCols {
			numCols = len(row)
		}
	}
	colWidths := make([]int, numCols)
	for _, row := range rows {
		for i, cell := range row {
			cw := utf8.RuneCountInString(cell)
			if cw > colWidths[i] {
				colWidths[i] = cw
			}
		}
	}

	// Build formatted table text.
	var sb strings.Builder
	for i, row := range rows {
		sb.WriteString("| ")
		for j := 0; j < numCols; j++ {
			cell := ""
			if j < len(row) {
				cell = row[j]
			}
			padded := cell + strings.Repeat(" ", colWidths[j]-utf8.RuneCountInString(cell))
			sb.WriteString(padded)
			if j < numCols-1 {
				sb.WriteString(" | ")
			}
		}
		sb.WriteString(" |")
		sb.WriteString("\n")

		// After header row, add separator.
		if i == 0 {
			sb.WriteString("|-")
			for j := 0; j < numCols; j++ {
				sb.WriteString(strings.Repeat("-", colWidths[j]))
				if j < numCols-1 {
					sb.WriteString("-+-")
				}
			}
			sb.WriteString("-|")
			sb.WriteString("\n")
		}
	}

	tableText := strings.TrimRight(sb.String(), "\n")

	// Write as pre entity.
	w.pushEntity(Pre, "", "")
	w.writeText(tableText)
	w.popEntity()
	w.ensureNewlines(2)
}

// extractText recursively extracts plain text from an AST node and its children.
func (w *walker) extractText(node ast.Node) string {
	var sb strings.Builder
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch c := child.(type) {
		case *ast.Text:
			sb.Write(c.Segment.Value(w.source))
		case *ast.String:
			sb.Write(c.Value)
		case *ast.CodeSpan:
			// Extract text from code span children.
			for gc := c.FirstChild(); gc != nil; gc = gc.NextSibling() {
				if t, ok := gc.(*ast.Text); ok {
					sb.Write(t.Segment.Value(w.source))
				}
			}
		default:
			// Recurse for other inline nodes.
			sb.WriteString(w.extractText(child))
		}
	}
	return sb.String()
}

