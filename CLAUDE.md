# Claude Project Instructions for mdmu

## Project Overview

`mdmu` (Markdown Markup) is a terminal UI tool for annotating markdown files with line-level comments. It's designed to help users provide structured feedback on AI-generated markdown files (like plan files) that can be easily consumed by LLMs.

## Quick Start

**Installation:**
```bash
go install github.com/buckleypaul/mdmu@latest
```

**Usage:**
```bash
# Interactive TUI
mdmu <file.md>

# Export comments
mdmu comments <file.md>
```

**Dev workflow:**
```bash
go build -o mdmu .
./mdmu testfile.md
```

## Architecture

### Core Components

1. **Markdown Renderer** (`internal/markdown/`)
   - Custom goldmark-based ANSI renderer
   - Maintains bidirectional mapping between source lines and rendered output
   - Supports headings, code blocks, lists, blockquotes, emphasis, links
   - Word-wrapping aware of ANSI escape codes

2. **TUI** (`internal/tui/`)
   - Built with Bubble Tea framework
   - Split-pane layout (65% markdown, 35% comments)
   - Modes: normal, selecting, commenting
   - Supports cursor navigation, line selection, comment input
   - Auto-resizes and re-renders on terminal resize

3. **Storage** (`internal/store/`)
   - JSON persistence in `/tmp/mdmu/<sha256-of-filepath>.json`
   - File hash tracking for change detection
   - Comments store: source line range, selected text, comment, timestamp

4. **Output** (`internal/output/`)
   - Formats comments as structured markdown
   - Designed for LLM consumption
   - Includes quoted source text with each comment

### Key Design Decisions

- **Line Mapping**: The renderer tracks which rendered lines correspond to which source lines. This is critical because rendered output differs from source (word-wrapping, styling, blank lines).
- **Selection Model**: Line-level selection (not character-level) for simplicity.
- **Storage Location**: `/tmp/mdmu/` for ephemeral comment sessions. Comments are tied to file path via SHA256 hash.
- **Re-rendering**: On terminal resize, the markdown is re-parsed at the new width to maintain proper word-wrapping.

## Development Guidelines

### Code Style

- Follow standard Go conventions
- Use `lipgloss` for all TUI styling (defined in `internal/tui/styles.go`)
- Keep renderer logic separate from TUI logic
- Test all core logic (parser, store, formatter)

### Testing

Run tests with:
```bash
go test ./internal/...
```

Current coverage:
- `internal/markdown`: 6 tests (parser, renderer, line mapping)
- `internal/store`: 3 tests (load/save, hashing)
- `internal/output`: 3 tests (formatting, sorting)

### Building

```bash
go build -o mdmu .
```

### Linting

```bash
go vet ./...
```

### Adding Features

**New markdown elements:**
1. Add rendering logic to `internal/markdown/renderer.go` in `renderNode()`
2. Ensure source line mapping is tracked via `addLine()`
3. Add test case in `renderer_test.go`

**New keybindings:**
1. Add to appropriate handler in `internal/tui/model.go` (`handleMarkdownKeys` or `handleCommentKeys`)
2. Update status bar hints in `internal/tui/statusbar.go`
3. Update README.md keybindings section

**New TUI modes:**
1. Add mode constant to `internal/tui/model.go`
2. Add mode-specific logic to `Update()` and `View()`
3. Add status bar context in `statusbar.go`

## Common Tasks

### Debugging Line Mappings

If comments are appearing on wrong lines:
1. Check `renderer.go` — ensure `addLine()` is called with correct source line numbers
2. Verify `byteOffsetToLine()` is computing line numbers correctly
3. Test with `internal/markdown/renderer_test.go::TestParseAndRender_SourceLineMappings`

### Debugging TUI State

Add temporary debug output in `model.go::View()`:
```go
debug := fmt.Sprintf("cursor=%d scroll=%d sel=%d mode=%d", m.cursor, m.scrollOffset, m.selectionStart, m.mode)
return panels + "\n" + debug + "\n" + statusBar
```

### Storage Issues

Comments are stored at `StorePath(filePath)` which is `/tmp/mdmu/<hash>.json`. To inspect:
```bash
ls -la /tmp/mdmu/
cat /tmp/mdmu/<hash>.json
```

## Known Limitations

1. **No undo**: Deleted comments are permanently removed
2. **No comment editing**: Comments can only be deleted and re-added
3. **Single-file focus**: No cross-file comment management
4. **Ephemeral storage**: `/tmp/` is cleared on reboot
5. **No authentication**: Comments have no author tracking

## Future Enhancements

See [GitHub issues](https://github.com/buckleypaul/mdmu/issues) for planned improvements.

## Testing Checklist

Before releasing changes:

- [ ] `go build` succeeds
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes
- [ ] Manual smoke test: `./mdmu testfile.md`
  - [ ] Navigation works (arrows, PgUp/PgDn)
  - [ ] Selection works (Shift+Up/Down)
  - [ ] Comment input works (C, type, Ctrl+S)
  - [ ] Comment deletion works (Tab to comments, d)
  - [ ] Quit works (q)
- [ ] `./mdmu comments testfile.md` outputs clean markdown
- [ ] Terminal resize doesn't crash

## Dependencies

- `bubbletea` - TUI framework
- `bubbles` - TUI components (textarea)
- `lipgloss` - Terminal styling
- `goldmark` - Markdown parsing
- `cobra` - CLI framework
- `uuid` - Comment ID generation

All dependencies are vendored in `go.mod`/`go.sum`.

## Project Structure

```
mdmu/
├── main.go                      # Entry point
├── cmd/
│   ├── root.go                  # TUI command
│   └── comments.go              # Export command
├── internal/
│   ├── markdown/
│   │   ├── types.go             # RenderedDocument, LineMapping
│   │   ├── parser.go            # Parse entry point
│   │   ├── renderer.go          # ANSI renderer with line tracking
│   │   └── renderer_test.go
│   ├── tui/
│   │   ├── model.go             # Bubble Tea model
│   │   ├── markdown_pane.go     # Left pane rendering
│   │   ├── comments_pane.go     # Right pane rendering
│   │   ├── comment_input.go     # Modal input
│   │   ├── statusbar.go         # Keybinding hints
│   │   └── styles.go            # lipgloss styles
│   ├── store/
│   │   ├── types.go             # Comment, CommentFile
│   │   ├── store.go             # Load/Save/Hash
│   │   └── store_test.go
│   └── output/
│       ├── formatter.go         # Markdown export
│       └── formatter_test.go
├── testfile.md                  # Sample file for testing
└── README.md
```

## Support & Issues

Report bugs at: https://github.com/buckleypaul/mdmu/issues

When reporting:
1. Include Go version (`go version`)
2. Include OS and terminal type
3. Provide minimal reproduction steps
4. Attach sample markdown file if relevant
