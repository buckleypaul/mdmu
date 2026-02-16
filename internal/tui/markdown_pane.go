package tui

import (
	"strings"

	"github.com/paulbuckley/mdmu/internal/markdown"
)

func (m Model) renderMarkdownPane() string {
	width := m.leftWidth()
	height := m.contentHeight()

	if len(m.doc.Lines) == 0 {
		content := emptyStateStyle.Render("No content to display")
		return m.markdownPaneBorder(width, height, content)
	}

	// Determine visible lines
	var visibleLines []string
	for i := m.scrollOffset; i < m.scrollOffset+height && i < len(m.doc.Lines); i++ {
		line := m.doc.Lines[i]

		// Pad or truncate to width
		visibleWidth := markdown.VisibleLen(line)
		padding := width - 2 - visibleWidth // -2 for border padding
		if padding < 0 {
			padding = 0
		}

		styledLine := line + strings.Repeat(" ", padding)

		// Apply highlight styles
		if m.isLineSelected(i) {
			styledLine = selectedLineStyle.Render(styledLine)
		} else if i == m.cursor && m.focusPane == paneMarkdown {
			styledLine = cursorLineStyle.Render(styledLine)
		}

		visibleLines = append(visibleLines, styledLine)
	}

	// Pad remaining height with empty lines
	for len(visibleLines) < height {
		visibleLines = append(visibleLines, strings.Repeat(" ", width-2))
	}

	content := strings.Join(visibleLines, "\n")
	return m.markdownPaneBorder(width, height, content)
}

func (m Model) markdownPaneBorder(width, height int, content string) string {
	title := paneTitle.Render("Markdown")

	style := inactiveBorderStyle
	if m.focusPane == paneMarkdown {
		style = activeBorderStyle
	}

	return style.Width(width - 2).Render(title + "\n" + content)
}

func (m Model) isLineSelected(renderedLine int) bool {
	if m.selectionStart < 0 {
		return false
	}

	start, end := m.selectionRange()
	return renderedLine >= start && renderedLine <= end
}

func (m Model) selectionRange() (int, int) {
	if m.selectionStart < 0 {
		return m.cursor, m.cursor
	}

	start := m.selectionStart
	end := m.cursor

	if start > end {
		start, end = end, start
	}

	return start, end
}

