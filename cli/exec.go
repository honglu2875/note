package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"github.com/honglu2875/note/note"
	"github.com/muesli/termenv"
	"github.com/pkg/errors"
)

type FileNode = note.FileNode

func OpenNote(path string, editor string, output *termenv.Output) error {
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to open %s with %s.", path, editor))
	}
	return nil
}

func RenameNote(path string, newName string, output *termenv.Output) error {
	if err := os.Rename(path, newName); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to rename %s to %s.", path, newName))
	}
	return nil
}

func RemoveNote(path string, output *termenv.Output) error {
	var ret string
	fmt.Printf("Are you sure to delete %s (y/N):", path)
	fmt.Scanf("%s", &ret)

	if ret == "y" || ret == "Y" {
		if err := os.Remove(path); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to remove %s.", path))
		}

	}
	return nil
}

func CreateNewNote(basePath, name string, output *termenv.Output) (string, error) {
	timestamp := time.Now().Format("2006-01-02")
	dir := filepath.Join(basePath, timestamp)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to create directory %s.", dir))
	}

	if name == "" {
		name = note.GenerateRandomHash(5)
	}
	fileName := name + ".md"
	filePath := filepath.Join(dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to create file %s.", filePath))
	}
	defer file.Close()

	fmt.Printf("Created new note: %s\n", filePath)

	return filePath, nil
}

func ListNotes(basePath string, output *termenv.Output) error {

	nodes, err := note.BuildTree(basePath)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to scan for markdown files under %s.", basePath))
	}

	fmt.Fprintln(output, termenv.String("Notes:").Bold().Foreground(termenv.ANSIBlue))
	err = printNode(output, nodes[0], "", true, 0)
	return err
}

func printNode(output *termenv.Output, node *FileNode, prefix string, isLast bool, level int) error {
	// Skip the root node display
	if level > 0 {
		// Create the branch prefix
		branch := "├── "
		if isLast {
			branch = "└── "
		}

		// Format the name based on type
		var idDisplay termenv.Style
		var nameDisplay termenv.Style
		if node.IsDir {
			// Directories in bold cyan
			idDisplay = termenv.String("")
			nameDisplay = termenv.String(node.Name).Bold().Foreground(termenv.ANSICyan)
		} else {
			// Files in green
			idDisplay = termenv.String(fmt.Sprintf("(%d) ", node.Id)).Foreground(termenv.ANSI256Color(155))
			nameDisplay = termenv.String(node.Name).Foreground(termenv.ANSIGreen)
		}

		// Print the node
		fmt.Fprintf(output, "%s%s%s%s\n", prefix, branch, idDisplay, nameDisplay)
	}

	// Sort children for consistent display
	var keys []string
	for k := range node.Children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Print children
	for i, key := range keys {
		childNode := node.Children[key]
		newPrefix := prefix

		if level > 0 {
			if isLast {
				newPrefix += "    " // Add space where the branch was
			} else {
				newPrefix += "│   " // Continue the line down
			}
		}

		err := printNode(
			output,
			childNode,
			newPrefix,
			i == len(keys)-1, // Is this the last child?
			level+1,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
