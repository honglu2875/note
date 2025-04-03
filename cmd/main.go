package main

import (
	"fmt"
	"os"

	"github.com/honglu2875/note/note"
)

func main() {
	basePath := note.getBasePath()
	editor := note.getEditor()

	if len(os.Args) == 1 {
		// `note` along will open the default note.
		note.openNote(note.getPath(basePath, note.DEFAULT_NOTE_NAME), editor)
		return
	}

	switch os.Args[1] {
	case "fast":
		var path string
		if len(os.Args) > 2 {
			path = note.createFastNote(basePath, os.Args[2])
		} else {
			path = note.createFastNote(basePath, "")
		}
		note.openNote(path, editor)
	case "list":
		note.listNotes(basePath)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
