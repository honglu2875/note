package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/honglu2875/note/note"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

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
				raiseIfIdOOB(id, len(nodes))
				path = nodes[id].Path
			}
			editor := note.GetEditor()
			err := OpenNote(path, editor, output)
			raiseIfError(err)
		},
	}

	newCmd := &cobra.Command{
		Use:   "new [title]",
		Short: "Create and open a quick note",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var name string
			if len(args) > 0 {
				name = args[0]
			} else {
				fmt.Print("Please provide a name (without .md extension): ")
				fmt.Scan(&name)
				if name == "" {
					fmt.Println("Note: empty input leads to a random hash.")
				}
			}
			path, err := CreateNewNote(basePath, name, output)
			raiseIfError(err)
			err = OpenNote(path, editor, output)
			raiseIfError(err)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all notes",
		Run: func(cmd *cobra.Command, args []string) {
			err := ListNotes(basePath, output)
			raiseIfError(err)
		},
	}

	renameCmd := &cobra.Command{
		Use:   "rename [id] [newName]",
		Short: "Rename the note",
		Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			idStr := args[0]
			newName := maybeAddSuffix(args[1])
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "The first argument needs to be an id according to the result of `note list`: %v\n", err)
			}
			nodes, err := note.BuildTree(basePath)
			raiseIfError(err)
			raiseIfIdOOB(id, len(nodes))
			path := nodes[id].Path
			dir, _ := filepath.Split(path)
			newPath := filepath.Join(dir, newName)
			err = RenameNote(nodes[id].Path, newPath, output)
			raiseIfError(err)
		},
	}

	removeCmd := &cobra.Command{
		Use:   "rm [id]",
		Short: "Remove the note",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			idStr := args[0]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "The first argument needs to be an id according to the result of `note list`: %v\n", err)
			}
			nodes, err := note.BuildTree(basePath)
			raiseIfError(err)
			raiseIfIdOOB(id, len(nodes))
			err = RemoveNote(nodes[id].Path, output)
			raiseIfError(err)
		},
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the base folder",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			isGit, err := note.CheckGitRepo(basePath)
			raiseIfError(err)
			if isGit {
				fmt.Printf("%s\n", color("The base folder is already a git repo.", 9))
				os.Exit(0)
			}

			var ret string
			fmt.Printf("It will initialize a git repo on the folder %s. Okay to continue? (y/N): ", color(basePath, 75))
			fmt.Scanf("%s", &ret)
			if ret == "y" || ret == "Y" {
				err = note.InitGitRepo(basePath)
				raiseIfError(err)
				err = note.CommitChanges(basePath, "initial commit")
				raiseIfError(err)
			}
		},
	}

	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(renameCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
