package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/honglu2875/note/note"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

func raiseIfError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	basePath := note.GetBasePath()
	editor := note.GetEditor()
	output := termenv.NewOutput(os.Stdout)

	rootCmd := &cobra.Command{
		Use:   "note",
		Short: "A note-taking CLI",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var path string

			if len(args) == 0 {
				path = note.GetPath(basePath, note.DEFAULT_NOTE_NAME)
			} else {
				nodes, err := note.BuildTree(basePath)
				raiseIfError(err)
				id, err := strconv.Atoi(args[0])
				raiseIfError(err)
				if id >= len(nodes) {
					fmt.Fprintf(os.Stderr, "id supplied is more than the notes collected. Need to be less than %d.\n", len(nodes))
					os.Exit(1)
				}
				path = nodes[id].Path
			}
			editor := note.GetEditor()
			OpenNote(path, editor, output)
		},
	}

	fastCmd := &cobra.Command{
		Use:   "fast [title]",
		Short: "Create and open a quick note",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var path string
			if len(args) > 2 {
				path = CreateFastNote(basePath, args[2], output)
			} else {
				path = CreateFastNote(basePath, "", output)
			}
			OpenNote(path, editor, output)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all notes",
		Run: func(cmd *cobra.Command, args []string) {
			ListNotes(basePath, output)
		},
	}

	rootCmd.AddCommand(fastCmd)
	rootCmd.AddCommand(listCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
