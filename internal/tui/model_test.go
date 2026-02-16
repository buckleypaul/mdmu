package tui

import (
	"testing"

	"github.com/paulbuckley/mdmu/internal/markdown"
	"github.com/paulbuckley/mdmu/internal/store"
)

func TestSelectionRange(t *testing.T) {
	tests := []struct {
		name           string
		cursor         int
		selectionStart int
		wantStart      int
		wantEnd        int
	}{
		{"no selection", 5, -1, 5, 5},
		{"forward selection", 7, 3, 3, 7},
		{"backward selection", 2, 8, 2, 8},
		{"single line selection", 4, 4, 4, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{cursor: tt.cursor, selectionStart: tt.selectionStart}
			start, end := m.selectionRange()
			if start != tt.wantStart || end != tt.wantEnd {
				t.Errorf("selectionRange() = (%d, %d), want (%d, %d)", start, end, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

func TestRenderedToSourceRange(t *testing.T) {
	m := Model{
		doc: &markdown.RenderedDocument{
			Mappings: []markdown.LineMapping{
				{RenderedLine: 0, SourceStart: 1, SourceEnd: 1},
				{RenderedLine: 1, SourceStart: 1, SourceEnd: 1},
				{RenderedLine: 2, SourceStart: 3, SourceEnd: 4},
				{RenderedLine: 3, SourceStart: 5, SourceEnd: 5},
				{RenderedLine: 4, SourceStart: 7, SourceEnd: 9},
			},
		},
	}

	tests := []struct {
		name      string
		start     int
		end       int
		wantStart int
		wantEnd   int
	}{
		{"single rendered line", 0, 0, 1, 1},
		{"span with gap", 2, 4, 3, 9},
		{"out of range defaults to 1,1", 10, 15, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := m.renderedToSourceRange(tt.start, tt.end)
			if start != tt.wantStart || end != tt.wantEnd {
				t.Errorf("renderedToSourceRange(%d, %d) = (%d, %d), want (%d, %d)",
					tt.start, tt.end, start, end, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

func TestSortedComments(t *testing.T) {
	m := Model{
		commentFile: &store.CommentFile{
			Comments: []store.Comment{
				{ID: "c", SourceStart: 10},
				{ID: "a", SourceStart: 1},
				{ID: "b", SourceStart: 5},
			},
		},
	}

	sorted := m.sortedComments()
	if len(sorted) != 3 {
		t.Fatalf("expected 3 comments, got %d", len(sorted))
	}
	if sorted[0].ID != "a" || sorted[1].ID != "b" || sorted[2].ID != "c" {
		t.Errorf("comments not sorted: %v, %v, %v", sorted[0].ID, sorted[1].ID, sorted[2].ID)
	}

	// Verify original slice is not modified
	if m.commentFile.Comments[0].ID != "c" {
		t.Error("sortedComments modified the original slice")
	}
}

func TestExtractSourceText(t *testing.T) {
	m := Model{source: []byte("line one\nline two\nline three\n")}

	tests := []struct {
		name      string
		start     int
		end       int
		want      string
	}{
		{"single line", 1, 1, "line one"},
		{"multiple lines", 1, 3, "line one\nline two\nline three"},
		{"middle line", 2, 2, "line two"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.extractSourceText(tt.start, tt.end)
			if got != tt.want {
				t.Errorf("extractSourceText(%d, %d) = %q, want %q", tt.start, tt.end, got, tt.want)
			}
		})
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"empty", "", []string{""}},
		{"no newline", "hello", []string{"hello"}},
		{"trailing newline", "hello\n", []string{"hello", ""}},
		{"multiple lines", "a\nb\nc", []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitLines(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("splitLines(%q) returned %d lines, want %d", tt.input, len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("splitLines(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}
