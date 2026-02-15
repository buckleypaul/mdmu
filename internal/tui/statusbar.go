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
		hints = fmt.Sprintf(" %s save  %s cancel",
			statusKeyStyle.Render("Ctrl+S"),
			statusKeyStyle.Render("Esc"))

	case m.mode == modeSelecting:
		hints = fmt.Sprintf(" %s extend  %s comment  %s cancel",
			statusKeyStyle.Render("Shift+↑↓"),
			statusKeyStyle.Render("C"),
			statusKeyStyle.Render("Esc"))

	case m.focusPane == paneComments:
		hints = fmt.Sprintf(" %s navigate  %s delete  %s markdown  %s quit",
			statusKeyStyle.Render("↑↓"),
			statusKeyStyle.Render("d"),
			statusKeyStyle.Render("Tab"),
			statusKeyStyle.Render("q"))

	default:
		hints = fmt.Sprintf(" %s navigate  %s select  %s comment  %s comments  %s quit",
			statusKeyStyle.Render("↑↓"),
			statusKeyStyle.Render("Shift+↑↓"),
			statusKeyStyle.Render("C"),
			statusKeyStyle.Render("Tab"),
			statusKeyStyle.Render("q"))
	}

	return statusBarStyle.Width(width).Render(hints)
}
