package store

import "time"

type Comment struct {
	ID           string    `json:"id"`
	SourceStart  int       `json:"source_start"` // 1-indexed
	SourceEnd    int       `json:"source_end"`
	SelectedText string    `json:"selected_text"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`
}

type CommentFile struct {
	File     string    `json:"file"`      // absolute path
	FileHash string    `json:"file_hash"` // sha256 of contents at comment time
	Comments []Comment `json:"comments"`
}
