package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulbuckley/mdmu/internal/clipboard"
	"github.com/paulbuckley/mdmu/internal/output"
)

func (m Model) enterPreviewMode() Model {
	m.previewContent = output.Format(m.commentFile, m.source, m.filename)
	m.previewScroll = 0
	m.copiedMessage = false
	m.mode = modePreview
	return m
}

func (m Model) handlePreviewKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "esc":
		m.mode = modeNormal
		m.copiedMessage = false
		m.statusMessage = ""
		return m, nil

	case "c", "C":
		if err := clipboard.Copy(m.previewContent); err != nil {
			m.statusMessage = "✗ Failed to copy: " + err.Error()
		} else {
			m.statusMessage = "✓ Copied to clipboard"
		}
		m.mode = modeNormal
		m.copiedMessage = false
		return m, nil

	case "up":
		if m.previewScroll > 0 {
			m.previewScroll--
		}
		return m, nil

	case "down":
		lines := strings.Split(m.previewContent, "\n")
		maxScroll := len(lines) - m.previewHeight()
		if maxScroll < 0 {
			maxScroll = 0
		}
		if m.previewScroll < maxScroll {
			m.previewScroll++
		}
		return m, nil

	case "pgup":
		m.previewScroll -= m.previewHeight()
		if m.previewScroll < 0 {
			m.previewScroll = 0
		}
		return m, nil

	case "pgdown":
		lines := strings.Split(m.previewContent, "\n")
		maxScroll := len(lines) - m.previewHeight()
		if maxScroll < 0 {
			maxScroll = 0
		}
		m.previewScroll += m.previewHeight()
		if m.previewScroll > maxScroll {
			m.previewScroll = maxScroll
		}
		return m, nil
	}

	return m, nil
}

func (m Model) previewHeight() int {
	h := m.height - 4 // title + status bar + borders
	if h < 1 {
		h = 1
	}
	return h
}

func (m Model) renderPreview() string {
	width := m.width
	if width <= 0 {
		width = 80
	}

	// Title
	title := previewTitleStyle.Render("Comment Preview")

	// Content
	lines := strings.Split(m.previewContent, "\n")
	height := m.previewHeight()

	var visibleLines []string
	for i := m.previewScroll; i < m.previewScroll+height && i < len(lines); i++ {
		visibleLines = append(visibleLines, lines[i])
	}

	// Pad remaining height
	for len(visibleLines) < height {
		visibleLines = append(visibleLines, "")
	}

	content := strings.Join(visibleLines, "\n")

	// Border around content
	bordered := activeBorderStyle.Width(width - 2).Render(title + "\n" + content)

	// Status bar
	statusBar := m.renderPreviewStatusBar()

	return bordered + "\n" + statusBar
}

func (m Model) renderPreviewStatusBar() string {
	width := m.width
	if width <= 0 {
		width = 80
	}

	hints := " " +
		statusKeyStyle.Render("C") + " copy  " +
		statusKeyStyle.Render("↑↓") + " or " + statusKeyStyle.Render("PgUp/PgDn") + " scroll  " +
		statusKeyStyle.Render("Esc") + " return  " +
		statusKeyStyle.Render("q") + " quit"

	return statusBarStyle.Width(width).Render(hints)
}
