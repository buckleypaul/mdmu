package markdown

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

// ANSI escape codes for styling
const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	italic    = "\033[3m"
	dim       = "\033[2m"
	underline = "\033[4m"

	fgCyan    = "\033[36m"
	fgYellow  = "\033[33m"
	fgGreen   = "\033[32m"
	fgMagenta = "\033[35m"
	fgWhite   = "\033[37m"
	fgGray    = "\033[90m"

	bgDarkGray = "\033[48;5;236m"
)

type ansiRenderer struct {
	source      []byte
	width       int
	lines       []string
	mappings    []LineMapping
	lineOffsets []int // byte offsets where each source line starts

	// State for inline rendering
	inlineStyles []string
}

func newANSIRenderer(source []byte, width int) *ansiRenderer {
	// Precompute line offset table for O(log n) byte-to-line lookups
	offsets := []int{0}
	for i, b := range source {
		if b == '\n' {
			offsets = append(offsets, i+1)
		}
	}

	return &ansiRenderer{
		source:      source,
		width:       width,
		lineOffsets: offsets,
	}
}

func (r *ansiRenderer) addLine(text string, sourceStart, sourceEnd int) {
	idx := len(r.lines)
	r.lines = append(r.lines, text)
	r.mappings = append(r.mappings, LineMapping{
		RenderedLine: idx,
		SourceStart:  sourceStart,
		SourceEnd:    sourceEnd,
	})
}

func (r *ansiRenderer) addBlankLine(sourceStart, sourceEnd int) {
	r.addLine("", sourceStart, sourceEnd)
}

// sourceLineRange returns the 1-indexed start and end source lines for a node.
func (r *ansiRenderer) sourceLineRange(node ast.Node) (int, int) {
	start := 0
	end := 0

	if node.Lines() != nil && node.Lines().Len() > 0 {
		firstSeg := node.Lines().At(0)
		lastSeg := node.Lines().At(node.Lines().Len() - 1)
		start = r.byteOffsetToLine(firstSeg.Start)
		end = r.byteOffsetToLine(lastSeg.Stop - 1)
	}

	if start == 0 {
		// Fall back to child nodes
		start, end = r.sourceLineRangeFromChildren(node)
	}

	if start == 0 {
		start = 1
		end = 1
	}

	return start, end
}

func (r *ansiRenderer) sourceLineRangeFromChildren(node ast.Node) (int, int) {
	start := 0
	end := 0

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		cs, ce := r.sourceLineRange(child)
		if cs > 0 {
			if start == 0 || cs < start {
				start = cs
			}
			if ce > end {
				end = ce
			}
		}
	}

	return start, end
}

// byteOffsetToLine converts a byte offset to a 1-indexed source line number.
func (r *ansiRenderer) byteOffsetToLine(offset int) int {
	// Binary search: find the number of line starts at or before offset
	return sort.Search(len(r.lineOffsets), func(i int) bool {
		return r.lineOffsets[i] > offset
	})
}

func (r *ansiRenderer) render(node ast.Node) {
	r.renderNode(node, 0)
}

func (r *ansiRenderer) renderNode(node ast.Node, depth int) {
	switch n := node.(type) {
	case *ast.Document:
		r.renderChildren(n, depth)

	case *ast.Heading:
		r.renderHeading(n)

	case *ast.Paragraph:
		r.renderParagraph(n, depth)

	case *ast.TextBlock:
		r.renderParagraph(n, depth)

	case *ast.FencedCodeBlock:
		r.renderFencedCodeBlock(n)

	case *ast.CodeBlock:
		r.renderCodeBlock(n)

	case *ast.List:
		r.renderList(n, depth)

	case *ast.ThematicBreak:
		start, end := r.sourceLineRange(n)
		r.addLine(fgGray+strings.Repeat("─", min(r.width, 40))+reset, start, end)
		r.addBlankLine(start, end)

	case *ast.Blockquote:
		r.renderBlockquote(n, depth)

	case *ast.HTMLBlock:
		r.renderHTMLBlock(n)

	default:
		// For unknown block nodes, try rendering children
		if node.HasChildren() {
			r.renderChildren(node, depth)
		}
	}
}

func (r *ansiRenderer) renderChildren(node ast.Node, depth int) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		r.renderNode(child, depth)
	}
}

func (r *ansiRenderer) renderHeading(node *ast.Heading) {
	start, end := r.sourceLineRange(node)
	text := r.renderInlineChildren(node)

	var prefix string
	var color string
	switch node.Level {
	case 1:
		color = fgCyan + bold
		prefix = "# "
	case 2:
		color = fgGreen + bold
		prefix = "## "
	case 3:
		color = fgYellow + bold
		prefix = "### "
	case 4:
		color = fgMagenta + bold
		prefix = "#### "
	default:
		color = bold
		prefix = strings.Repeat("#", node.Level) + " "
	}

	r.addLine(color+prefix+text+reset, start, end)
	r.addBlankLine(start, end)
}

func (r *ansiRenderer) renderParagraph(node ast.Node, depth int) {
	start, end := r.sourceLineRange(node)
	text := r.renderInlineChildren(node)

	// Word wrap
	wrapped := r.wordWrap(text, r.width-depth*2)
	for _, line := range wrapped {
		r.addLine(line, start, end)
	}
	r.addBlankLine(start, end)
}

func (r *ansiRenderer) renderFencedCodeBlock(node *ast.FencedCodeBlock) {
	start, end := r.sourceLineRange(node)
	lang := ""
	if node.Language(r.source) != nil {
		lang = string(node.Language(r.source))
	}

	// Header line
	if lang != "" {
		r.addLine(bgDarkGray+fgGray+" "+lang+" "+reset, start, start)
	}

	// Render each line of the code block
	for i := 0; i < node.Lines().Len(); i++ {
		seg := node.Lines().At(i)
		line := string(seg.Value(r.source))
		line = strings.TrimRight(line, "\n")
		r.addLine(bgDarkGray+fgWhite+" "+line+" "+reset, r.byteOffsetToLine(seg.Start), r.byteOffsetToLine(seg.Start))
	}

	r.addBlankLine(end, end)
}

func (r *ansiRenderer) renderCodeBlock(node *ast.CodeBlock) {
	_, end := r.sourceLineRange(node)

	for i := 0; i < node.Lines().Len(); i++ {
		seg := node.Lines().At(i)
		line := string(seg.Value(r.source))
		line = strings.TrimRight(line, "\n")
		r.addLine(bgDarkGray+fgWhite+" "+line+" "+reset, r.byteOffsetToLine(seg.Start), r.byteOffsetToLine(seg.Start))
	}

	r.addBlankLine(end, end)
}

func (r *ansiRenderer) renderList(node *ast.List, depth int) {
	itemNum := node.Start
	if itemNum == 0 {
		itemNum = 1
	}

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		listItem, ok := child.(*ast.ListItem)
		if !ok {
			continue
		}

		start, end := r.sourceLineRange(listItem)

		var prefix string
		if node.IsOrdered() {
			prefix = fmt.Sprintf("%d. ", itemNum)
			itemNum++
		} else {
			prefix = "  • "
		}

		// Render the list item's inline content
		firstBlock := true
		for blockChild := listItem.FirstChild(); blockChild != nil; blockChild = blockChild.NextSibling() {
			if _, ok := blockChild.(*ast.List); ok {
				// Nested list — render with increased depth
				r.renderList(blockChild.(*ast.List), depth+1)
				continue
			}

			text := r.renderInlineChildren(blockChild)
			cs, ce := r.sourceLineRange(blockChild)
			if cs > 0 {
				start = cs
				end = ce
			}

			indent := strings.Repeat("  ", depth)
			if firstBlock {
				wrapped := r.wordWrap(text, r.width-len(indent)-len(prefix))
				for i, line := range wrapped {
					if i == 0 {
						r.addLine(indent+prefix+line, start, end)
					} else {
						r.addLine(indent+strings.Repeat(" ", len(prefix))+line, start, end)
					}
				}
				firstBlock = false
			} else {
				wrapped := r.wordWrap(text, r.width-len(indent)-len(prefix))
				for _, line := range wrapped {
					r.addLine(indent+strings.Repeat(" ", len(prefix))+line, start, end)
				}
			}
		}
	}

	r.addBlankLine(r.sourceLineRange(node))
}

func (r *ansiRenderer) renderBlockquote(node *ast.Blockquote, depth int) {
	// Render children into a temporary renderer, then prefix each line
	subRenderer := newANSIRenderer(r.source, r.width-4)
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		subRenderer.renderNode(child, depth)
	}

	prefix := fgGray + "│ " + reset
	for i, line := range subRenderer.lines {
		mapping := subRenderer.mappings[i]
		r.addLine(prefix+line, mapping.SourceStart, mapping.SourceEnd)
	}
}

func (r *ansiRenderer) renderHTMLBlock(node *ast.HTMLBlock) {
	_, end := r.sourceLineRange(node)
	for i := 0; i < node.Lines().Len(); i++ {
		seg := node.Lines().At(i)
		line := string(seg.Value(r.source))
		line = strings.TrimRight(line, "\n")
		r.addLine(fgGray+line+reset, r.byteOffsetToLine(seg.Start), r.byteOffsetToLine(seg.Start))
	}
	r.addBlankLine(end, end)
}

// renderInlineChildren renders all inline children of a node to a string with ANSI codes.
func (r *ansiRenderer) renderInlineChildren(node ast.Node) string {
	var buf bytes.Buffer
	r.renderInline(&buf, node)
	return buf.String()
}

func (r *ansiRenderer) renderInline(buf *bytes.Buffer, node ast.Node) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch n := child.(type) {
		case *ast.Text:
			buf.Write(n.Segment.Value(r.source))
			if n.SoftLineBreak() {
				buf.WriteByte(' ')
			}
			if n.HardLineBreak() {
				buf.WriteByte('\n')
			}

		case *ast.String:
			buf.Write(n.Value)

		case *ast.CodeSpan:
			buf.WriteString(fgYellow)
			buf.WriteString("`")
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				if t, ok := c.(*ast.Text); ok {
					buf.Write(t.Segment.Value(r.source))
				}
			}
			buf.WriteString("`")
			buf.WriteString(reset)

		case *ast.Emphasis:
			if n.Level == 2 {
				buf.WriteString(bold)
				r.renderInline(buf, n)
				buf.WriteString(reset)
			} else {
				buf.WriteString(italic)
				r.renderInline(buf, n)
				buf.WriteString(reset)
			}

		case *ast.Link:
			buf.WriteString(underline + fgCyan)
			r.renderInline(buf, n)
			buf.WriteString(reset)
			buf.WriteString(fgGray + " (" + string(n.Destination) + ")" + reset)

		case *ast.Image:
			buf.WriteString(fgGray + "[img: ")
			r.renderInline(buf, n)
			buf.WriteString("]" + reset)

		case *ast.AutoLink:
			url := string(n.URL(r.source))
			buf.WriteString(underline + fgCyan + url + reset)

		case *ast.RawHTML:
			for i := 0; i < n.Segments.Len(); i++ {
				seg := n.Segments.At(i)
				buf.Write(seg.Value(r.source))
			}

		case *east.Strikethrough:
			buf.WriteString(dim)
			r.renderInline(buf, n)
			buf.WriteString(reset)

		default:
			// Unknown inline node — try rendering children
			if child.HasChildren() {
				r.renderInline(buf, child)
			}
		}
	}
}

// wordWrap wraps text to the given width, respecting ANSI escape sequences.
func (r *ansiRenderer) wordWrap(s string, width int) []string {
	if width <= 0 {
		width = 40
	}

	var lines []string
	for _, paragraph := range strings.Split(s, "\n") {
		if paragraph == "" {
			lines = append(lines, "")
			continue
		}
		wrapped := wrapLine(paragraph, width)
		lines = append(lines, wrapped...)
	}

	if len(lines) == 0 {
		lines = []string{""}
	}
	return lines
}

// wrapLine wraps a single line of text, accounting for ANSI escape codes.
func wrapLine(s string, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	currentLine := ""
	currentWidth := 0

	for _, word := range words {
		wordWidth := VisibleLen(word)

		if currentWidth == 0 {
			currentLine = word
			currentWidth = wordWidth
		} else if currentWidth+1+wordWidth <= width {
			currentLine += " " + word
			currentWidth += 1 + wordWidth
		} else {
			lines = append(lines, currentLine)
			currentLine = word
			currentWidth = wordWidth
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// VisibleLen returns the visible width of a string, ignoring ANSI escape sequences
// and accounting for wide characters (CJK, emoji).
func VisibleLen(s string) int {
	width := 0
	inEscape := false
	for _, r := range s {
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		if r == '\033' {
			inEscape = true
			continue
		}
		width += runewidth.RuneWidth(r)
	}
	return width
}
