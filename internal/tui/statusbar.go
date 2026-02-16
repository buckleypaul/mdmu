package tui

import "fmt"

func (m Model) renderStatusBar() string {
	width := m.width
	if width <= 0 {
		width = 80
	}

	var hints string
	switch {
	case m.mode == modeCommenting:
		hints = fmt.Sprintf(" %s save  %s newline  %s cancel",
			statusKeyStyle.Render("Enter"),
			statusKeyStyle.Render("Alt+Enter"),
			statusKeyStyle.Render("Esc"))

	case m.mode == modeSelecting:
		hints = fmt.Sprintf(" %s extend  %s comment  %s cancel",
			statusKeyStyle.Render("Shift+↑↓"),
			statusKeyStyle.Render("Enter"),
			statusKeyStyle.Render("Esc"))

	case m.focusPane == paneComments:
		hints = fmt.Sprintf(" %s navigate  %s delete  %s markdown  %s quit",
			statusKeyStyle.Render("↑↓"),
			statusKeyStyle.Render("d"),
			statusKeyStyle.Render("Tab"),
			statusKeyStyle.Render("q"))

	default:
		hints = fmt.Sprintf(" %s navigate  %s select  %s comment  %s comments  %s preview  %s copy  %s quit",
			statusKeyStyle.Render("↑↓"),
			statusKeyStyle.Render("Shift+↑↓"),
			statusKeyStyle.Render("Enter"),
			statusKeyStyle.Render("Tab"),
			statusKeyStyle.Render("P"),
			statusKeyStyle.Render("C"),
			statusKeyStyle.Render("q"))
	}

	// Prepend status message if present
	if m.statusMessage != "" {
		hints = " " + m.statusMessage + "  |" + hints
	}

	return statusBarStyle.Width(width).Render(hints)
}
