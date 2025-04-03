package note

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultBasePath   = "notes"
	DefaultEditor     = "nvim"
	DEFAULT_NOTE_NAME = "notes.md"
)

func getBasePath() string {
	basePath := os.Getenv("NOTE_PATH")
	if basePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get home directory: %v\n", err)
			os.Exit(1)
		}
		basePath = filepath.Join(homeDir, DefaultBasePath)
	}
	return basePath
}

func getEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}
	return editor
}

func getPath(elem ...string) string {
	return filepath.Join(elem...)
}

func main() {
	basePath := getBasePath()
	editor := getEditor()

	if len(os.Args) == 1 {
		// `note` along will open the default note.
		openNote(getPath(basePath, DEFAULT_NOTE_NAME), editor)
		return
	}

	switch os.Args[1] {
	case "fast":
		var path string
		if len(os.Args) > 2 {
			path = createFastNote(basePath, os.Args[2])
		} else {
			path = createFastNote(basePath, "")
		}
		openNote(path, editor)
	case "list":
		listNotes(basePath)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
