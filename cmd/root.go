package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulbuckley/mdmu/internal/markdown"
	"github.com/paulbuckley/mdmu/internal/store"
	"github.com/paulbuckley/mdmu/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mdmu <file>",
	Short: "Annotate markdown files with comments",
	Long:  "A terminal UI for navigating rendered markdown files and adding line-level comments.",
	Args:  cobra.ExactArgs(1),
	RunE:  runTUI,
}

func SetVersion(v string) {
	rootCmd.Version = v
}

func Execute() error {
	return rootCmd.Execute()
}

func runTUI(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Resolve to absolute path
	if !filepath.IsAbs(filePath) {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}
		filePath = filepath.Join(wd, filePath)
	}

	// Read the source file
	source, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	// Parse and render the markdown
	doc, err := markdown.ParseAndRender(source, 80)
	if err != nil {
		return fmt.Errorf("parsing markdown: %w", err)
	}

	// Create empty in-memory comment store
	cf := &store.CommentFile{}

	// Initialize the TUI model
	model := tui.NewModel(doc, cf, source, filepath.Base(filePath))

	// Run Bubble Tea
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("running TUI: %w", err)
	}

	return nil
}
