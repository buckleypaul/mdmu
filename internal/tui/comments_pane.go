package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/paulbuckley/mdmu/internal/store"
)

func (m Model) renderCommentsPane() string {
	width := m.rightWidth()
	height := m.contentHeight()

	if len(m.commentFile.Comments) == 0 {
		content := emptyStateStyle.Render("No comments yet\nSelect lines and press C")
		return m.commentsPaneBorder(width, height, content)
	}

	// Sort comments by source line
	sorted := m.sortedComments()

	var visibleLines []string
	lineIdx := 0

	for i, c := range sorted {
		// Header: line range
		var header string
		if c.SourceStart == c.SourceEnd {
			header = commentHeaderStyle.Render(fmt.Sprintf("L%d", c.SourceStart))
		} else {
			header = commentHeaderStyle.Render(fmt.Sprintf("L%d-%d", c.SourceStart, c.SourceEnd))
		}

		// Truncate comment text for display
		text := c.Comment
		maxTextWidth := width - 6
		if maxTextWidth < 10 {
			maxTextWidth = 10
		}
		if runewidth.StringWidth(text) > maxTextWidth {
			text = runewidth.Truncate(text, maxTextWidth, "...")
		}

		commentLine := header + " " + commentTextStyle.Render(text)

		// Highlight if this comment is focused
		if m.focusPane == paneComments && i == m.commentCursor {
			commentLine = commentHighlightStyle.Width(width - 4).Render(commentLine)
		}

		if lineIdx >= m.commentScrollOffset && lineIdx < m.commentScrollOffset+height {
			visibleLines = append(visibleLines, commentLine)
		}
		lineIdx++

		// Add separator between comments
		if i < len(sorted)-1 {
			sep := commentLineRefStyle.Render(strings.Repeat("─", width-6))
			if lineIdx >= m.commentScrollOffset && lineIdx < m.commentScrollOffset+height {
				visibleLines = append(visibleLines, sep)
			}
			lineIdx++
		}
	}

	// Pad remaining height
	for len(visibleLines) < height {
		visibleLines = append(visibleLines, "")
	}

	content := strings.Join(visibleLines, "\n")
	return m.commentsPaneBorder(width, height, content)
}

func (m Model) commentsPaneBorder(width, height int, content string) string {
	title := paneTitle.Render(fmt.Sprintf("Comments (%d)", len(m.commentFile.Comments)))

	style := inactiveBorderStyle
	if m.focusPane == paneComments {
		style = activeBorderStyle
	}

	return style.Width(width - 2).Render(title + "\n" + content)
}

func (m Model) sortedComments() []store.Comment {
	sorted := make([]store.Comment, len(m.commentFile.Comments))
	copy(sorted, m.commentFile.Comments)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].SourceStart < sorted[j].SourceStart
	})
	return sorted
}
