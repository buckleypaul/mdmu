package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/paulbuckley/mdmu/internal/store"
)

func newCommentTextarea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = ""
	ta.ShowLineNumbers = false
	ta.SetHeight(3)
	return ta
}

func (m Model) handleCommentInput(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mode = modeNormal
			return m, nil

		case "alt+enter":
			// Insert a newline into the textarea
			enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(enterMsg)
			return m, cmd

		case "enter":
			// Save the comment
			comment := m.textarea.Value()
			if comment == "" {
				m.mode = modeNormal
				return m, nil
			}

			// Map rendered line selection to source lines
			selStart, selEnd := m.selectionRange()
			sourceStart, sourceEnd := m.renderedToSourceRange(selStart, selEnd)

			// Extract selected text
			selectedText := m.extractSourceText(sourceStart, sourceEnd)

			c := store.Comment{
				ID:           uuid.New().String(),
				SourceStart:  sourceStart,
				SourceEnd:    sourceEnd,
				SelectedText: selectedText,
				Comment:      comment,
				CreatedAt:    time.Now(),
			}

			m.commentFile.Comments = append(m.commentFile.Comments, c)

			m.mode = modeNormal
			m.selectionStart = -1
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m Model) renderCommentInput() string {
	selStart, selEnd := m.selectionRange()
	sourceStart, sourceEnd := m.renderedToSourceRange(selStart, selEnd)

	title := modalTitleStyle.Render(fmt.Sprintf("Comment on lines %d-%d", sourceStart, sourceEnd))
	ta := m.textarea.View()

	content := title + "\n" + ta

	inputWidth := m.width - 4
	if inputWidth < 20 {
		inputWidth = 20
	}

	return modalStyle.Width(inputWidth).Render(content)
}

func (m Model) renderedToSourceRange(renderedStart, renderedEnd int) (int, int) {
	sourceStart := 0
	sourceEnd := 0

	for i := renderedStart; i <= renderedEnd && i < len(m.doc.Mappings); i++ {
		mapping := m.doc.Mappings[i]
		if mapping.SourceStart > 0 {
			if sourceStart == 0 || mapping.SourceStart < sourceStart {
				sourceStart = mapping.SourceStart
			}
			if mapping.SourceEnd > sourceEnd {
				sourceEnd = mapping.SourceEnd
			}
		}
	}

	if sourceStart == 0 {
		sourceStart = 1
		sourceEnd = 1
	}

	return sourceStart, sourceEnd
}

func (m Model) extractSourceText(sourceStart, sourceEnd int) string {
	lines := splitLines(string(m.source))
	start := sourceStart - 1
	end := sourceEnd
	if start < 0 {
		start = 0
	}
	if end > len(lines) {
		end = len(lines)
	}

	result := ""
	for i := start; i < end; i++ {
		if i > start {
			result += "\n"
		}
		result += lines[i]
	}
	return result
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start <= len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
