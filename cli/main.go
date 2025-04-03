package main

import (
	"fmt"
	"os"

	"github.com/honglu2875/note/note"
)

func main() {
	basePath := note.GetBasePath()
	editor := note.GetEditor()

	if len(os.Args) == 1 {
		// `note` along will open the default note.
		note.OpenNote(note.GetPath(basePath, note.DEFAULT_NOTE_NAME), editor)
		return
	}

	switch os.Args[1] {
	case "fast":
		var path string
		if len(os.Args) > 2 {
			path = note.CreateFastNote(basePath, os.Args[2])
		} else {
			path = note.CreateFastNote(basePath, "")
		}
		note.OpenNote(path, editor)
	case "list":
		note.ListNotes(basePath)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
