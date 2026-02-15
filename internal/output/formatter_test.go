package output

import (
	"strings"
	"testing"

	"github.com/paulbuckley/mdmu/internal/store"
)

func TestFormatEmpty(t *testing.T) {
	cf := &store.CommentFile{
		File:     "/tmp/test.md",
		Comments: []store.Comment{},
	}
	result := Format(cf, []byte("# Test\n"))
	if result != "" {
		t.Errorf("expected empty string for no comments, got %q", result)
	}
}

func TestFormatSingleComment(t *testing.T) {
	source := []byte("# Title\n\nSome text here.\n")
	cf := &store.CommentFile{
		File: "/tmp/test.md",
		Comments: []store.Comment{
			{
				ID:          "1",
				SourceStart: 3,
				SourceEnd:   3,
				Comment:     "This needs more detail",
			},
		},
	}

	result := Format(cf, source)

	if !strings.Contains(result, "Comments on test.md") {
		t.Error("output should contain file name")
	}
	if !strings.Contains(result, "Line 3") {
		t.Error("output should contain line reference")
	}
	if !strings.Contains(result, "This needs more detail") {
		t.Error("output should contain comment text")
	}
	if !strings.Contains(result, "> Some text here.") {
		t.Error("output should contain quoted source text")
	}
}

func TestFormatMultipleComments_Sorted(t *testing.T) {
	source := []byte("line1\nline2\nline3\nline4\nline5\n")
	cf := &store.CommentFile{
		File: "/tmp/test.md",
		Comments: []store.Comment{
			{ID: "2", SourceStart: 4, SourceEnd: 5, Comment: "second"},
			{ID: "1", SourceStart: 1, SourceEnd: 2, Comment: "first"},
		},
	}

	result := Format(cf, source)

	// "first" should appear before "second" in output
	firstIdx := strings.Index(result, "first")
	secondIdx := strings.Index(result, "second")
	if firstIdx > secondIdx {
		t.Error("comments should be sorted by source line")
	}
}
