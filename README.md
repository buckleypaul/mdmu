# mdmu (Markdown Markup)

A terminal UI for annotating markdown files with line-level comments.

## Overview

When working with AI-generated markdown files (like plan files from Claude Code), you often need to provide feedback on specific sections. `mdmu` provides a TUI where you can:

- Navigate rendered markdown with syntax highlighting
- Select line ranges and add comments
- View all comments in a sidebar
- Preview formatted output in full-screen mode
- Copy comments as structured markdown to clipboard for LLM consumption

## Installation

### Homebrew (macOS/Linux)

```bash
brew install buckleypaul/tap/mdmu
```

### Go Install

```bash
go install github.com/buckleypaul/mdmu@latest
```

### Build from Source

```bash
git clone https://github.com/buckleypaul/mdmu.git
cd mdmu
go build -o mdmu .
```

## Usage

```bash
mdmu <file.md>
```

**Keybindings:**

**Normal mode:**
- `↑↓` - Navigate lines
- `PgUp/PgDn` - Jump by page
- `Home/End` - Jump to start/end of document
- `Shift+↑↓` - Select line ranges
- `C` - Add comment to current line or selection
- `Tab` - Switch between markdown and comments pane
- `Enter` - Preview formatted output (when comments exist)
- `Esc` - Clear selection
- `q` - Quit

**Comment input mode:**
- `Enter` - Save comment
- `Alt+Enter` - Insert newline in comment
- `Esc` - Cancel comment input

**Comments pane:**
- `↑↓` - Navigate comments
- `d` - Delete focused comment
- `Tab` - Switch back to markdown pane

**Preview mode:**
- `y` - Copy formatted output to clipboard
- `↑↓` or `PgUp/PgDn` - Scroll preview
- `Esc` - Return to normal mode
- `q` - Quit

## Claude Code Integration

`mdmu` is designed to work seamlessly with Claude Code for reviewing AI-generated plans and documents.

**Workflow:**

1. **Claude Code generates a file** (e.g., `plan.md`)
2. **Open mdmu in a separate terminal:**
   ```bash
   mdmu plan.md
   ```
3. **Add your comments** interactively using the TUI
   - Navigate with `↑↓`
   - Select ranges with `Shift+↑↓`
   - Press `C` to add comments
4. **Preview your comments:**
   - Press `Enter` to view formatted output
   - Review the structured markdown
5. **Copy to clipboard:**
   - Press `y` in preview mode
   - Paste directly into your Claude Code session

The formatted output includes a ready-to-use prompt:

```markdown
Please address my comments on plan.md:

## Comments on plan.md

### Lines 5-12:
> This section covers the main
> components of the system...

**Comment:** Need more detail on the auth flow here

---
```

**Note:** Comments are ephemeral and exist only during your mdmu session. This encourages a focused review workflow without persistent file clutter.

## Features

- **Rich markdown rendering** - Headings, code blocks, lists, blockquotes, emphasis, links
- **Source line mapping** - Accurate tracking from rendered output to source lines (handles word-wrapping)
- **Preview mode** - Full-screen formatted output view before copying
- **Clipboard integration** - Cross-platform clipboard copy (macOS, Linux, Windows)
- **Ephemeral comments** - Session-only storage encourages focused review workflow
- **Responsive resize** - Automatically re-renders markdown when terminal is resized

## Architecture

- **Language**: Go
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Markdown Parser**: [goldmark](https://github.com/yuin/goldmark)
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)

See `CLAUDE.md` for detailed architecture decisions and development guidelines.

## License

MIT
