package cmd

import (
	"fmt"
	"os"

	"github.com/paulbuckley/mdmu/internal/output"
	"github.com/paulbuckley/mdmu/internal/store"
	"github.com/spf13/cobra"
)

var commentsCmd = &cobra.Command{
	Use:   "comments <file>",
	Short: "Print comments for a file to stdout",
	Long:  "Loads stored comments for a markdown file and prints them in a structured format suitable for LLM consumption.",
	Args:  cobra.ExactArgs(1),
	RunE:  runComments,
}

func init() {
	rootCmd.AddCommand(commentsCmd)
}

func runComments(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	if !isAbsPath(filePath) {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}
		filePath = wd + "/" + filePath
	}

	source, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	cf, err := store.Load(filePath)
	if err != nil {
		return fmt.Errorf("loading comments: %w", err)
	}

	if len(cf.Comments) == 0 {
		fmt.Fprintln(os.Stderr, "No comments found for this file.")
		return nil
	}

	out := output.Format(cf, source)
	fmt.Print(out)
	return nil
}
