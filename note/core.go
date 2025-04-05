package note

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DEFAULT_BASE_PATH = "notes"
	DEFAULT_EDITOR    = "nvim"
	DEFAULT_NOTE_NAME = "notes.md"
)

func GetBasePath() string {
	basePath := os.Getenv("NOTE_PATH")
	if basePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get home directory: %v\n", err)
			os.Exit(1)
		}
		basePath = filepath.Join(homeDir, DEFAULT_BASE_PATH)
	}
	return basePath
}

func GetEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DEFAULT_EDITOR
	}
	return editor
}

func GetPath(elem ...string) string {
	return filepath.Join(elem...)
}
