package cmd

import (
	"fmt"
	"os"

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

func Execute() error {
	return rootCmd.Execute()
}

func runTUI(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Resolve to absolute path
	if !isAbsPath(filePath) {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}
		filePath = wd + "/" + filePath
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

	// Load existing comments
	cf, err := store.Load(filePath)
	if err != nil {
		return fmt.Errorf("loading comments: %w", err)
	}

	// Initialize the TUI model
	model := tui.NewModel(doc, cf, source)

	// Run Bubble Tea
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("running TUI: %w", err)
	}

	return nil
}

func isAbsPath(p string) bool {
	return len(p) > 0 && p[0] == '/'
}
