package markdown

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// ParseAndRender parses markdown source and renders it with ANSI styling,
// tracking the mapping from rendered lines to source lines.
func ParseAndRender(source []byte, width int) (*RenderedDocument, error) {
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := newANSIRenderer(source, width)
	renderer.render(doc)

	return &RenderedDocument{
		Lines:    renderer.lines,
		Mappings: renderer.mappings,
	}, nil
}
