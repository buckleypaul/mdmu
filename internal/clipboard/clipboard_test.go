package clipboard

import (
	"os/exec"
	"runtime"
	"testing"
)

func TestCopyDoesNotPanic(t *testing.T) {
	// Skip if no clipboard command is available
	switch runtime.GOOS {
	case "darwin":
		if _, err := exec.LookPath("pbcopy"); err != nil {
			t.Skip("pbcopy not available")
		}
	case "linux":
		if _, err := exec.LookPath("xclip"); err != nil {
			if _, err := exec.LookPath("xsel"); err != nil {
				t.Skip("no clipboard command available")
			}
		}
	case "windows":
		if _, err := exec.LookPath("clip.exe"); err != nil {
			t.Skip("clip.exe not available")
		}
	default:
		t.Skip("unsupported platform")
	}

	err := Copy("test clipboard content")
	if err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	// On macOS, verify clipboard contents
	if runtime.GOOS == "darwin" {
		out, err := exec.Command("pbpaste").Output()
		if err != nil {
			t.Fatalf("pbpaste failed: %v", err)
		}
		if string(out) != "test clipboard content" {
			t.Errorf("clipboard contents = %q, want %q", string(out), "test clipboard content")
		}
	}
}
