package output

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/paulbuckley/mdmu/internal/store"
)

// Format renders comments as structured markdown for LLM consumption.
func Format(cf *store.CommentFile, source []byte) string {
	if len(cf.Comments) == 0 {
		return ""
	}

	sourceLines := strings.Split(string(source), "\n")

	// Sort comments by source line position
	sorted := make([]store.Comment, len(cf.Comments))
	copy(sorted, cf.Comments)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].SourceStart < sorted[j].SourceStart
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## Comments on %s\n\n", filepath.Base(cf.File)))

	for i, c := range sorted {
		if c.SourceStart == c.SourceEnd {
			sb.WriteString(fmt.Sprintf("### Line %d:\n", c.SourceStart))
		} else {
			sb.WriteString(fmt.Sprintf("### Lines %d-%d:\n", c.SourceStart, c.SourceEnd))
		}

		// Quote the selected source text
		start := c.SourceStart - 1
		end := c.SourceEnd
		if start < 0 {
			start = 0
		}
		if end > len(sourceLines) {
			end = len(sourceLines)
		}
		for _, line := range sourceLines[start:end] {
			sb.WriteString("> " + line + "\n")
		}
		sb.WriteString("\n")

		sb.WriteString("**Comment:** " + c.Comment + "\n")

		if i < len(sorted)-1 {
			sb.WriteString("\n---\n\n")
		} else {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
