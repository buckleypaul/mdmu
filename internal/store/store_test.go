package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStorePathDeterministic(t *testing.T) {
	p1 := StorePath("/tmp/test/file.md")
	p2 := StorePath("/tmp/test/file.md")
	if p1 != p2 {
		t.Errorf("StorePath not deterministic: %q != %q", p1, p2)
	}

	p3 := StorePath("/tmp/test/other.md")
	if p1 == p3 {
		t.Error("different files should have different store paths")
	}
}

func TestLoadSaveRoundTrip(t *testing.T) {
	// Create a temp file to comment on
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("# Test\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Load (should return empty)
	cf, err := Load(testFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(cf.Comments) != 0 {
		t.Errorf("expected 0 comments, got %d", len(cf.Comments))
	}

	// Add a comment and save
	cf.Comments = append(cf.Comments, Comment{
		ID:          "test-1",
		SourceStart: 1,
		SourceEnd:   1,
		Comment:     "test comment",
	})
	if err := Save(cf); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Reload
	cf2, err := Load(testFile)
	if err != nil {
		t.Fatalf("Load after save failed: %v", err)
	}
	if len(cf2.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(cf2.Comments))
	}
	if cf2.Comments[0].Comment != "test comment" {
		t.Errorf("comment text = %q, want %q", cf2.Comments[0].Comment, "test comment")
	}

	// Cleanup
	os.Remove(StorePath(testFile))
}

func TestComputeFileHash(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	h1, err := ComputeFileHash(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if h1 == "" {
		t.Error("hash should not be empty")
	}

	// Same content = same hash
	h2, err := ComputeFileHash(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Error("same file should produce same hash")
	}

	// Different content = different hash
	if err := os.WriteFile(testFile, []byte("world\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	h3, err := ComputeFileHash(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if h1 == h3 {
		t.Error("different content should produce different hash")
	}
}
