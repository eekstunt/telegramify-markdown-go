package tgmd

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// md is the goldmark parser configured with GFM extensions.
var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)

// Convert parses Markdown text and returns plain text with Telegram entities.
func Convert(markdown string, opts ...Option) Message {
	cfg := applyOptions(opts)
	source := []byte(markdown)

	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader, parser.WithContext(parser.NewContext()))

	w := &walker{
		source: source,
		cfg:    cfg,
	}
	ast.Walk(doc, w.walk)
	return w.result()
}

// ConvertAndSplit parses Markdown and splits the result into messages
// that each fit within the configured max message length (UTF-16 units).
func ConvertAndSplit(markdown string, opts ...Option) []Message {
	cfg := applyOptions(opts)
	msg := Convert(markdown, opts...)
	return Split(msg, cfg.maxMessageLen)
}
