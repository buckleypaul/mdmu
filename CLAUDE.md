# mdmu — Markdown Markup

Terminal UI for annotating markdown files with line-level comments, designed for LLM-consumable feedback.

## Commands

```bash
go build -o mdmu .          # Build
go test ./...               # Run all tests (markdown: 6, store: 3, output: 3)
go test -v -race ./...      # Tests with race detection (matches CI)
go vet ./...                # Lint
./mdmu testfile.md          # Smoke test the TUI
./mdmu comments testfile.md # Export comments as markdown
```

**Go version**: 1.25+ (see `go.mod`)

## CI & Release

- **CI** (`.github/workflows/ci.yml`): Runs on push/PR to `main`. Steps: `go vet`, `go test -race`, `go build`.
- **Release** (`.github/workflows/release.yml`): Triggers on `v*` tags. Uses GoReleaser v2.
- **GoReleaser** (`.goreleaser.yml`): Builds for linux/darwin/windows (amd64/arm64). Publishes to Homebrew tap `buckleypaul/homebrew-tap`.

**Note**: Module path in `go.mod` is `github.com/paulbuckley/mdmu` but the GitHub repo is `github.com/buckleypaul/mdmu`. Be aware of this mismatch.

## Architecture

### Core Components

1. **Markdown Renderer** (`internal/markdown/`) — Custom goldmark-based ANSI renderer. Maintains bidirectional mapping between source lines and rendered output. Word-wrapping aware of ANSI escape codes.

2. **TUI** (`internal/tui/`) — Bubble Tea framework. Split-pane layout (65% markdown, 35% comments). Modes: normal, selecting, commenting. Auto-resizes and re-renders on terminal resize.

3. **Storage** (`internal/store/`) — JSON persistence in `/tmp/mdmu/<sha256-of-filepath>.json`. File hash tracking for change detection.

4. **Output** (`internal/output/`) — Formats comments as structured markdown for LLM consumption.

### Key Design Decisions

- **Line Mapping**: Renderer tracks which rendered lines correspond to which source lines (critical because word-wrapping changes line counts).
- **Selection Model**: Line-level selection (not character-level).
- **Storage**: `/tmp/mdmu/` for ephemeral sessions. Cleared on reboot.
- **Re-rendering**: On terminal resize, markdown is re-parsed at new width.

## Code Style

- Standard Go conventions
- `lipgloss` for all TUI styling (defined in `internal/tui/styles.go`)
- Keep renderer logic separate from TUI logic

## Adding Features

**New markdown elements:**
1. Add rendering logic in `internal/markdown/renderer.go` → `renderNode()`
2. Track source line mapping via `addLine()`
3. Add test in `renderer_test.go`

**New keybindings:**
1. Add handler in `internal/tui/model.go` (`handleMarkdownKeys` or `handleCommentKeys`)
2. Update status bar in `internal/tui/statusbar.go`
3. Update README.md keybindings section

**New TUI modes:**
1. Add mode constant in `internal/tui/model.go`
2. Add mode logic to `Update()` and `View()`
3. Add status bar context in `statusbar.go`

## Gotchas

### Line Mapping Bugs

If comments appear on wrong lines:
1. Check `renderer.go` — ensure `addLine()` is called with correct source line numbers
2. Verify `byteOffsetToLine()` computes line numbers correctly
3. Test with `TestParseAndRender_SourceLineMappings`

### Debugging TUI State

Temporary debug output in `model.go::View()`:
```go
debug := fmt.Sprintf("cursor=%d scroll=%d sel=%d mode=%d", m.cursor, m.scrollOffset, m.selectionStart, m.mode)
return panels + "\n" + debug + "\n" + statusBar
```

### Storage

Inspect comments: `cat /tmp/mdmu/<hash>.json`

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
├── .goreleaser.yml              # Release config
├── .github/workflows/
│   ├── ci.yml                   # CI pipeline
│   └── release.yml              # Tag-triggered release
└── testfile.md                  # Sample file for testing
```
