package markdown

import (
	"strings"
	"testing"
)

func TestParseAndRender_BasicHeading(t *testing.T) {
	source := []byte("# Hello World\n\nSome text here.\n")
	doc, err := ParseAndRender(source, 80)
	if err != nil {
		t.Fatalf("ParseAndRender failed: %v", err)
	}

	if len(doc.Lines) == 0 {
		t.Fatal("expected rendered lines, got none")
	}

	// First line should contain "Hello World"
	if !containsVisible(doc.Lines[0], "Hello World") {
		t.Errorf("first line should contain 'Hello World', got: %q", doc.Lines[0])
	}

	// Check that mappings exist for each line
	if len(doc.Mappings) != len(doc.Lines) {
		t.Errorf("mappings count (%d) != lines count (%d)", len(doc.Mappings), len(doc.Lines))
	}

	// First mapping should point to source line 1
	if doc.Mappings[0].SourceStart != 1 {
		t.Errorf("first mapping source start = %d, want 1", doc.Mappings[0].SourceStart)
	}
}

func TestParseAndRender_CodeBlock(t *testing.T) {
	source := []byte("```go\nfmt.Println(\"hello\")\n```\n")
	doc, err := ParseAndRender(source, 80)
	if err != nil {
		t.Fatalf("ParseAndRender failed: %v", err)
	}

	if len(doc.Lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d", len(doc.Lines))
	}

	// Should contain the code content
	found := false
	for _, line := range doc.Lines {
		if containsVisible(line, "fmt.Println") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find 'fmt.Println' in rendered output")
	}
}

func TestParseAndRender_List(t *testing.T) {
	source := []byte("- item one\n- item two\n- item three\n")
	doc, err := ParseAndRender(source, 80)
	if err != nil {
		t.Fatalf("ParseAndRender failed: %v", err)
	}

	found := 0
	for _, line := range doc.Lines {
		if containsVisible(line, "item") {
			found++
		}
	}
	if found < 3 {
		t.Errorf("expected 3 list items rendered, found %d", found)
	}
}

func TestParseAndRender_Blockquote(t *testing.T) {
	source := []byte("> This is a quote\n")
	doc, err := ParseAndRender(source, 80)
	if err != nil {
		t.Fatalf("ParseAndRender failed: %v", err)
	}

	found := false
	for _, line := range doc.Lines {
		if containsVisible(line, "This is a quote") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected blockquote text in output")
	}
}

func TestParseAndRender_SourceLineMappings(t *testing.T) {
	source := []byte("# Title\n\nParagraph one.\n\nParagraph two.\n")
	doc, err := ParseAndRender(source, 80)
	if err != nil {
		t.Fatalf("ParseAndRender failed: %v", err)
	}

	// Verify all mappings have positive source lines
	for i, m := range doc.Mappings {
		if m.SourceStart < 1 {
			t.Errorf("mapping %d has SourceStart %d, want >= 1", i, m.SourceStart)
		}
		if m.SourceEnd < m.SourceStart {
			t.Errorf("mapping %d has SourceEnd %d < SourceStart %d", i, m.SourceEnd, m.SourceStart)
		}
	}
}

func TestVisibleLen(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"hello", 5},
		{"\033[1mhello\033[0m", 5},
		{"\033[36m# Title\033[0m", 7},
		{"", 0},
	}

	for _, tt := range tests {
		got := visibleLen(tt.input)
		if got != tt.want {
			t.Errorf("visibleLen(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

// containsVisible strips ANSI codes and checks if the string contains the substring.
func containsVisible(s, substr string) bool {
	stripped := stripANSI(s)
	return strings.Contains(stripped, substr)
}

func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false
	for _, r := range s {
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		if r == '\033' {
			inEscape = true
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}
