package store

import "time"

type Comment struct {
	ID           string
	SourceStart  int // 1-indexed
	SourceEnd    int
	SelectedText string
	Comment      string
	CreatedAt    time.Time
}

type CommentFile struct {
	Comments []Comment
}
