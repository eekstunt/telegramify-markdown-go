// Package tgmd converts Markdown to Telegram-compatible plain text with MessageEntity objects.
//
// Instead of producing MarkdownV2 strings (which require escaping 18 special characters
// with context-dependent rules), tgmd outputs plain text paired with entity objects that
// use UTF-16 offsets — matching what the Telegram Bot API expects natively.
//
// All functions are safe for concurrent use.
package tgmd

// EntityType represents the type of a Telegram message entity.
// Values match Telegram Bot API MessageEntity type strings.
type EntityType string

const (
	Bold          EntityType = "bold"
	Italic        EntityType = "italic"
	Underline     EntityType = "underline"
	Strikethrough EntityType = "strikethrough"
	Code          EntityType = "code"
	Pre           EntityType = "pre"
	TextLink      EntityType = "text_link"
	Blockquote    EntityType = "blockquote"
)

// Entity represents a Telegram MessageEntity.
// Offset and Length are in UTF-16 code units (Telegram's native encoding).
type Entity struct {
	Type     EntityType
	Offset   int    // UTF-16 offset from start of text
	Length   int    // length in UTF-16 code units
	URL      string // only for TextLink
	Language string // only for Pre (fenced code block language)
}

// Message is a single Telegram message: plain text plus formatting entities.
type Message struct {
	Text     string
	Entities []Entity
}
