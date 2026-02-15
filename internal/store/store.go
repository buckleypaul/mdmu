package store

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const storeDir = "/tmp/mdmu"

// StorePath returns the deterministic JSON file path for a given source file.
func StorePath(filePath string) string {
	abs, _ := filepath.Abs(filePath)
	hash := sha256.Sum256([]byte(abs))
	return filepath.Join(storeDir, fmt.Sprintf("%x.json", hash))
}

// Load reads the comment file for the given source file path.
// Returns an empty CommentFile if none exists.
func Load(filePath string) (*CommentFile, error) {
	abs, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("resolving path: %w", err)
	}

	sp := StorePath(abs)
	data, err := os.ReadFile(sp)
	if err != nil {
		if os.IsNotExist(err) {
			hash, hashErr := ComputeFileHash(abs)
			if hashErr != nil {
				return nil, hashErr
			}
			return &CommentFile{
				File:     abs,
				FileHash: hash,
				Comments: []Comment{},
			}, nil
		}
		return nil, fmt.Errorf("reading store file: %w", err)
	}

	var cf CommentFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("parsing store file: %w", err)
	}

	return &cf, nil
}

// Save writes the CommentFile to its deterministic store path.
func Save(cf *CommentFile) error {
	if err := os.MkdirAll(storeDir, 0o755); err != nil {
		return fmt.Errorf("creating store dir: %w", err)
	}

	data, err := json.MarshalIndent(cf, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling comments: %w", err)
	}

	sp := StorePath(cf.File)
	if err := os.WriteFile(sp, data, 0o644); err != nil {
		return fmt.Errorf("writing store file: %w", err)
	}

	return nil
}

// ComputeFileHash returns the SHA256 hex digest of a file's contents.
func ComputeFileHash(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("reading file for hash: %w", err)
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}
