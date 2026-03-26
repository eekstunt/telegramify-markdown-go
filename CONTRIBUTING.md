# Contributing

Contributions are welcome! Here's how to get started.

## Reporting issues

Open a [GitHub issue](https://github.com/eekstunt/telegramify-markdown-go/issues) with a clear description and, if applicable, a minimal Markdown input that reproduces the problem.

## Submitting changes

1. Fork the repository
2. Create a feature branch (`git checkout -b my-feature`)
3. Make your changes
4. Run tests: `go test -race ./...`
5. Run vet: `go vet ./...`
6. Commit and push
7. Open a pull request

## Code style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Add tests for new functionality
- Keep the single-dependency philosophy — avoid adding dependencies unless absolutely necessary
