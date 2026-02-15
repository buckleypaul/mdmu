# mdmu (Markdown Markup)

A terminal UI for annotating markdown files with line-level comments.

## Overview

When working with AI-generated markdown files (like plan files from Claude Code), you often need to provide feedback on specific sections. `mdmu` provides a TUI where you can:

- Navigate rendered markdown with syntax highlighting
- Select line ranges and add comments
- View all comments in a sidebar
- Export comments as structured markdown for LLM consumption

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

### Interactive TUI

```bash
mdmu <file.md>
```

**Keybindings:**
- `↑↓` - Navigate lines
- `Shift+↑↓` - Select line ranges
- `C` - Add comment to current line or selection
- `Ctrl+S` - Save comment (when in comment input mode)
- `Esc` - Cancel comment input or clear selection
- `Tab` - Switch between markdown and comments pane
- `d` - Delete focused comment (when in comments pane)
- `q` - Quit

### Export Comments

```bash
mdmu comments <file.md>
```

Prints all comments in a structured markdown format suitable for copying to Claude Code:

```markdown
Please address my comments on plan.md:

## Comments on plan.md

### Lines 5-12:
> This section covers the main
> components of the system...

**Comment:** Need more detail on the auth flow here

---
```

## Claude Code Integration

`mdmu` is designed to work seamlessly with Claude Code for reviewing AI-generated plans and documents.

**Workflow:**

1. **Claude Code generates a file** (e.g., `plan.md`)
2. **Open mdmu in a separate terminal:**
   ```bash
   mdmu plan.md
   ```
3. **Add your comments** interactively using the TUI
4. **Quit mdmu** (press `q`) - comments are auto-saved
5. **Export comments for Claude Code:**
   ```bash
   mdmu comments plan.md
   ```
6. **Copy the entire output and paste into your Claude Code session**

The output is formatted with a ready-to-use prompt that tells Claude Code to address your comments:

```markdown
Please address my comments on plan.md:

## Comments on plan.md

### Lines 5-12:
> This section covers the main
> components of the system...

**Comment:** Need more detail on the auth flow here

---
```

**Alternative:** You can also run `mdmu comments <file>` directly from within Claude Code using bash commands, and Claude will automatically read and process the comments.

## Features

- **Rich markdown rendering** - Headings, code blocks, lists, blockquotes, emphasis, links
- **Source line mapping** - Accurate tracking from rendered output to source lines
- **Persistent storage** - Comments stored as JSON in `/tmp/mdmu/`
- **File change detection** - Warns when file has been modified since comments were added
- **Responsive resize** - Automatically re-renders markdown when terminal is resized

## Architecture

- **Language**: Go
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Markdown Parser**: [goldmark](https://github.com/yuin/goldmark)
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)

See the [implementation plan](plan.md) for detailed architecture decisions.

## License

MIT
