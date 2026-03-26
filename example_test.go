package tgmd_test

import (
	"fmt"

	tgmd "github.com/eekstunt/telegramify-markdown-go"
)

func ExampleConvert() {
	msg := tgmd.Convert("**bold** and *italic*")
	fmt.Println(msg.Text)
	for _, e := range msg.Entities {
		fmt.Printf("%s at %d len %d\n", e.Type, e.Offset, e.Length)
	}
	// Output:
	// bold and italic
	// bold at 0 len 4
	// italic at 9 len 6
}

func ExampleConvert_codeBlock() {
	msg := tgmd.Convert("```go\nfmt.Println(\"hello\")\n```")
	fmt.Println(msg.Text)
	for _, e := range msg.Entities {
		fmt.Printf("%s lang=%q at %d len %d\n", e.Type, e.Language, e.Offset, e.Length)
	}
	// Output:
	// fmt.Println("hello")
	// pre lang="go" at 0 len 20
}

func ExampleConvert_options() {
	msg := tgmd.Convert("# Hello",
		tgmd.WithHeadingSymbols([6]string{">>", "", "", "", "", ""}),
	)
	fmt.Println(msg.Text)
	// Output:
	// >> Hello
}

func ExampleConvertAndSplit() {
	msgs := tgmd.ConvertAndSplit("short message")
	fmt.Println(len(msgs))
	fmt.Println(msgs[0].Text)
	// Output:
	// 1
	// short message
}
