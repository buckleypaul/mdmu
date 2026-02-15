package markdown

// LineMapping maps a rendered line to its source file line range.
type LineMapping struct {
	RenderedLine int // 0-indexed line in rendered output
	SourceStart  int // 1-indexed line in source file
	SourceEnd    int // 1-indexed line in source file
}

// RenderedDocument holds the rendered output and its line mappings.
type RenderedDocument struct {
	Lines    []string      // rendered lines (with ANSI codes)
	Mappings []LineMapping // one per rendered line
}
