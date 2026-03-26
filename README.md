# telegramify-markdown-go

[![Go Reference](https://pkg.go.dev/badge/github.com/eekstunt/telegramify-markdown-go.svg)](https://pkg.go.dev/github.com/eekstunt/telegramify-markdown-go)
[![Test](https://github.com/eekstunt/telegramify-markdown-go/actions/workflows/test.yml/badge.svg)](https://github.com/eekstunt/telegramify-markdown-go/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/eekstunt/telegramify-markdown-go)](https://goreportcard.com/report/github.com/eekstunt/telegramify-markdown-go)
[![Coverage](https://raw.githubusercontent.com/eekstunt/telegramify-markdown-go/coverage/coverage.svg)](https://raw.githubusercontent.com/eekstunt/telegramify-markdown-go/coverage/coverage.html)

A Go library that converts Markdown to Telegram-compatible plain text with [MessageEntity](https://core.telegram.org/bots/api#messageentity) objects.

Uses [goldmark](https://github.com/yuin/goldmark) for Markdown parsing.

## Why entities instead of MarkdownV2?

Telegram Bot API offers two ways to format messages:

**1. ParseMode (MarkdownV2 or HTML)** — you embed formatting markers in the text:
```
Text: "*bold* and _italic_"
ParseMode: "MarkdownV2"
```
This requires escaping 18 special characters (`_ * [ ] ( ) ~ ` > # + - = | { } . !`) with context-dependent rules. Different rules apply inside code blocks, URLs, and normal text. A single unescaped `.` or `!` will make Telegram reject the entire message. This is the #1 source of bugs in every Telegram formatting library.

**2. Entities** — you send plain text plus an array of entity objects with numeric offsets:
```
Text: "bold and italic"
Entities: [
  {"type": "bold", "offset": 0, "length": 4},
  {"type": "italic", "offset": 9, "length": 6}
]
```
No markers in the text. No escaping. No parse failures. Telegram applies formatting based on UTF-16 offsets. The complexity shifts to computing offsets correctly — which is deterministic math, not fragile string manipulation.

This library uses the entities approach exclusively.

## Install

```bash
go get github.com/eekstunt/telegramify-markdown-go
```

## Usage

```go
import tgmd "github.com/eekstunt/telegramify-markdown-go"

// Convert markdown to plain text + entities
msg := tgmd.Convert("**bold** and *italic*")
// msg.Text     = "bold and italic"
// msg.Entities = [{Bold, 0, 4}, {Italic, 9, 6}]

// Convert and split into <=4096 UTF-16 unit messages
msgs := tgmd.ConvertAndSplit(longMarkdown)
```

### Example with [go-telegram/bot](https://github.com/go-telegram/bot)

```go
import (
    "github.com/go-telegram/bot"
    "github.com/go-telegram/bot/models"
    tgmd "github.com/eekstunt/telegramify-markdown-go"
)

// Convert and send
msgs := tgmd.ConvertAndSplit(markdown)
for _, msg := range msgs {
    b.SendMessage(ctx, &bot.SendMessageParams{
        ChatID:   chatID,
        Text:     msg.Text,
        Entities: toEntities(msg.Entities),
    })
}

// tgmd.Entity fields map 1:1 to Telegram's MessageEntity,
// so the conversion is a straightforward struct copy:
func toEntities(ents []tgmd.Entity) []models.MessageEntity {
    out := make([]models.MessageEntity, len(ents))
    for i, e := range ents {
        out[i] = models.MessageEntity{
            Type:     models.MessageEntityType(e.Type),
            Offset:   e.Offset,
            Length:   e.Length,
            URL:      e.URL,
            Language: e.Language,
        }
    }
    return out
}
```

## Supported Markdown elements

- **Bold**, *italic*, ~~strikethrough~~
- `inline code` and fenced code blocks (with language)
- [Links](https://example.com) and images (rendered as links)
- Headings (bold/underline/italic depending on level, with configurable emoji prefix)
- Ordered and unordered lists (with nesting)
- Blockquotes
- GFM tables (rendered as monospace pre blocks)
- Task lists (with unicode checkmarks)
- Horizontal rules

## Configuration

```go
msg := tgmd.Convert(markdown,
    tgmd.WithHeadingSymbols([6]string{"📌", "✏", "📚", "🔖", "", ""}),
    tgmd.WithTaskMarkers("✅", "☐"),
    tgmd.WithBulletMarker("⦁"),
    tgmd.WithMaxMessageLen(4096),
)
```

## License

MIT

## Acknowledgments

This library is inspired by the excellent [telegramify-markdown](https://github.com/sudoskys/telegramify-markdown) Python library by [@sudoskys](https://github.com/sudoskys). Their pioneering work on the entity-based approach (moving away from MarkdownV2 strings in v1.0.0) directly influenced the architecture of this Go implementation. Thank you for building such a well-thought-out solution and sharing it with the community.
