package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paulbuckley/mdmu/internal/markdown"
	"github.com/paulbuckley/mdmu/internal/store"
)

// reRender re-parses and re-renders the markdown at the given width.
func (m *Model) reRender() {
	renderWidth := m.leftWidth() - 4 // account for borders and padding
	if renderWidth < 20 {
		renderWidth = 20
	}
	doc, err := markdown.ParseAndRender(m.source, renderWidth)
	if err != nil {
		m.statusMessage = "Render error: " + err.Error()
		return
	}
	m.doc = doc
	if m.cursor >= len(m.doc.Lines) {
		m.cursor = len(m.doc.Lines) - 1
		if m.cursor < 0 {
			m.cursor = 0
		}
	}
	m.ensureCursorVisible()
}

type mode int

const (
	modeNormal     mode = iota
	modeSelecting
	modeCommenting
	modePreview
)

type pane int

const (
	paneMarkdown pane = iota
	paneComments
)

type Model struct {
	doc         *markdown.RenderedDocument
	commentFile *store.CommentFile
	source      []byte
	filename    string

	// Window dimensions
	width  int
	height int

	// Markdown pane state
	cursor       int // current rendered line (0-indexed)
	scrollOffset int // first visible line

	// Selection state
	selectionStart int // -1 means no selection
	mode           mode

	// Comments pane state
	commentCursor       int
	commentScrollOffset int

	// Focus
	focusPane pane

	// Comment input
	textarea textarea.Model

	// Preview state
	previewContent string
	previewScroll  int
	copiedMessage  bool

	// Status
	statusMessage string
}

func NewModel(doc *markdown.RenderedDocument, cf *store.CommentFile, source []byte, filename string) Model {
	return Model{
		doc:            doc,
		commentFile:    cf,
		source:         source,
		filename:       filename,
		selectionStart: -1,
		focusPane:      paneMarkdown,
		textarea:       newCommentTextarea(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		oldWidth := m.width
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(m.width - 6)
		if m.width != oldWidth {
			m.reRender()
		}
		return m, nil

	case tea.KeyMsg:
		// Global keys
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Route to comment input handler if in comment mode
		if m.mode == modeCommenting {
			return m.handleCommentInput(msg)
		}

		// Route to preview handler if in preview mode
		if m.mode == modePreview {
			return m.handlePreviewKeys(msg)
		}

		return m.handleKeypress(msg)
	}

	// Route non-key messages to textarea when in comment mode (cursor blink, etc.)
	if m.mode == modeCommenting {
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKeypress(msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	switch {
	// Quit
	case key == "q":
		return m, tea.Quit

	// Enter preview mode
	case key == "enter" && len(m.commentFile.Comments) > 0:
		return m.enterPreviewMode(), nil

	// Tab: switch focus
	case key == "tab":
		if m.focusPane == paneMarkdown {
			m.focusPane = paneComments
			if len(m.commentFile.Comments) > 0 {
				m.scrollToCommentTarget()
			}
		} else {
			m.focusPane = paneMarkdown
		}
		return m, nil

	// Navigation when in markdown pane
	case m.focusPane == paneMarkdown:
		return m.handleMarkdownKeys(key)

	// Navigation when in comments pane
	case m.focusPane == paneComments:
		return m.handleCommentKeys(key)
	}

	return m, nil
}

func (m Model) handleMarkdownKeys(key string) (Model, tea.Cmd) {
	maxLine := len(m.doc.Lines) - 1
	if maxLine < 0 {
		maxLine = 0
	}

	switch key {
	case "up":
		m.selectionStart = -1
		m.mode = modeNormal
		if m.cursor > 0 {
			m.cursor--
		}
		m.ensureCursorVisible()

	case "down":
		m.selectionStart = -1
		m.mode = modeNormal
		if m.cursor < maxLine {
			m.cursor++
		}
		m.ensureCursorVisible()

	case "shift+up":
		if m.selectionStart < 0 {
			m.selectionStart = m.cursor
			m.mode = modeSelecting
		}
		if m.cursor > 0 {
			m.cursor--
		}
		m.ensureCursorVisible()

	case "shift+down":
		if m.selectionStart < 0 {
			m.selectionStart = m.cursor
			m.mode = modeSelecting
		}
		if m.cursor < maxLine {
			m.cursor++
		}
		m.ensureCursorVisible()

	case "pgup":
		m.selectionStart = -1
		m.mode = modeNormal
		m.cursor -= m.contentHeight()
		if m.cursor < 0 {
			m.cursor = 0
		}
		m.ensureCursorVisible()

	case "pgdown":
		m.selectionStart = -1
		m.mode = modeNormal
		m.cursor += m.contentHeight()
		if m.cursor > maxLine {
			m.cursor = maxLine
		}
		m.ensureCursorVisible()

	case "home":
		m.selectionStart = -1
		m.mode = modeNormal
		m.cursor = 0
		m.ensureCursorVisible()

	case "end":
		m.selectionStart = -1
		m.mode = modeNormal
		m.cursor = maxLine
		m.ensureCursorVisible()

	case "esc":
		m.selectionStart = -1
		m.mode = modeNormal

	case "c", "C":
		// Enter comment mode
		m.mode = modeCommenting
		m.textarea = newCommentTextarea()
		m.textarea.SetWidth(m.width - 6)
		if m.selectionStart < 0 {
			// Comment on current line only
			m.selectionStart = m.cursor
		}
		return m, m.textarea.Focus()
	}

	return m, nil
}

func (m Model) handleCommentKeys(key string) (Model, tea.Cmd) {
	sorted := m.sortedComments()
	maxIdx := len(sorted) - 1
	if maxIdx < 0 {
		maxIdx = 0
	}

	switch key {
	case "up":
		if m.commentCursor > 0 {
			m.commentCursor--
			m.scrollToCommentTarget()
		}

	case "down":
		if m.commentCursor < maxIdx {
			m.commentCursor++
			m.scrollToCommentTarget()
		}

	case "d":
		if len(sorted) > 0 && m.commentCursor < len(sorted) {
			// Delete the comment
			target := sorted[m.commentCursor]
			for i, c := range m.commentFile.Comments {
				if c.ID == target.ID {
					m.commentFile.Comments = append(m.commentFile.Comments[:i], m.commentFile.Comments[i+1:]...)
					break
				}
			}
			if m.commentCursor >= len(m.commentFile.Comments) && m.commentCursor > 0 {
				m.commentCursor--
			}
		}
	}

	return m, nil
}

// scrollToCommentTarget scrolls the markdown pane to show the lines referenced by the focused comment.
func (m *Model) scrollToCommentTarget() {
	sorted := m.sortedComments()
	if m.commentCursor >= len(sorted) {
		return
	}

	c := sorted[m.commentCursor]

	// Find the first rendered line that maps to this source line
	for i, mapping := range m.doc.Mappings {
		if mapping.SourceStart >= c.SourceStart && mapping.SourceStart <= c.SourceEnd {
			m.cursor = i
			m.ensureCursorVisible()
			break
		}
	}
}

func (m *Model) ensureCursorVisible() {
	height := m.contentHeight()
	if height <= 0 {
		return
	}

	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}
	if m.cursor >= m.scrollOffset+height {
		m.scrollOffset = m.cursor - height + 1
	}
}

func (m Model) leftWidth() int {
	if m.width <= 0 {
		return 40
	}
	return int(float64(m.width) * 0.65)
}

func (m Model) rightWidth() int {
	if m.width <= 0 {
		return 20
	}
	return m.width - m.leftWidth()
}

func (m Model) contentHeight() int {
	h := m.height - 4 // borders + status bar + title
	if h < 1 {
		h = 1
	}
	return h
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Preview mode: full-screen formatted output
	if m.mode == modePreview {
		return m.renderPreview()
	}

	// Render panes side by side
	left := m.renderMarkdownPane()
	right := m.renderCommentsPane()
	panels := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	// Status bar
	statusBar := m.renderStatusBar()

	// Comment input modal overlay
	if m.mode == modeCommenting {
		modal := m.renderCommentInput()
		// Place modal at the bottom, above the status bar
		return panels + "\n" + modal + "\n" + statusBar
	}

	return panels + "\n" + statusBar
}
