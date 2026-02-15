package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Pane border styles
	activeBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62"))

	inactiveBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240"))

	// Line highlighting
	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236"))

	selectedLineStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("24"))

	// Comment pane
	commentHighlightStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("236"))

	commentHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("62")).
				Bold(true)

	commentTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	commentLineRefStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("243"))

	// Status bar
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	statusKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("236")).
			Bold(true)

	// Title styles
	paneTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Bold(true).
			Padding(0, 1)

	// Comment input modal
	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	modalTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Bold(true)

	// Empty state
	emptyStateStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true)

	// Warning style
	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)
)
