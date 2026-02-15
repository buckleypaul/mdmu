# mdmu Test File

This is a test file for the **mdmu** markdown annotation tool.

## Features

- Navigate rendered markdown with arrow keys
- Select line ranges with **Shift+Up/Down**
- Add comments to selected lines with `C`
- View and manage comments in the right panel

## Code Example

```go
func main() {
    fmt.Println("Hello, mdmu!")
}
```

## Blockquote

> This is a blockquote that demonstrates
> how blockquotes are rendered in the TUI.

## Links and Formatting

Visit [GitHub](https://github.com) for more information.

Here is some *italic text* and some **bold text** and some `inline code`.

---

### Nested Lists

1. First item
   - Sub-item A
   - Sub-item B
2. Second item
3. Third item

## Final Section

This is the last section of the test file. It contains enough content
to test scrolling behavior when the terminal window is small enough
that not all lines fit on screen at once.

More text here to ensure we have plenty of lines for testing the
viewport scrolling and cursor movement functionality of mdmu.

And yet another paragraph to really make sure we can test all the
navigation features properly.
