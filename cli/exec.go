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
)

type FileNode = note.FileNode

func OpenNote(path string, editor string, output *termenv.Output) {
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open editor: %v\n", err)
		os.Exit(1)
	}
}

func CreateFastNote(basePath, name string, output *termenv.Output) string {
	timestamp := time.Now().Format("2006-01-02")
	dir := filepath.Join(basePath, timestamp)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create directory: %v\n", err)
		os.Exit(1)
	}

	if name == "" {
		name = note.GenerateRandomHash(5)
	}
	fileName := name + ".md"
	filePath := filepath.Join(dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Printf("Created new note: %s\n", filePath)

	return filePath
}

func ListNotes(basePath string, output *termenv.Output) {

	nodes, err := note.BuildTree(basePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list notes: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(output, termenv.String("Notes:").Bold().Foreground(termenv.ANSIBlue))
	printNode(output, nodes[0], "", true, 0)

}

func getLabeledName(id int64, name string) string {
	return fmt.Sprintf("(%d) %s", id, name)
}

func printNode(output *termenv.Output, node *FileNode, prefix string, isLast bool, level int) {
	// Skip the root node display
	if level > 0 {
		// Create the branch prefix
		branch := "├── "
		if isLast {
			branch = "└── "
		}

		// Format the name based on type
		var nameDisplay termenv.Style
		if node.IsDir {
			// Directories in bold cyan
			nameDisplay = termenv.String(node.Name).Bold().Foreground(termenv.ANSICyan)
		} else {
			// Files in green
			nameDisplay = termenv.String(getLabeledName(node.Id, node.Name)).Foreground(termenv.ANSIGreen)
		}

		// Print the node
		fmt.Fprintf(output, "%s%s%s\n", prefix, branch, nameDisplay)
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

		printNode(
			output,
			childNode,
			newPrefix,
			i == len(keys)-1, // Is this the last child?
			level+1,
		)
	}
}
