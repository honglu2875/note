package note

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func openNote(path string, editor string) {
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

func createFastNote(basePath, name string) string {
	timestamp := time.Now().Format("2006-01-02")
	dir := filepath.Join(basePath, timestamp)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create directory: %v\n", err)
		os.Exit(1)
	}

	if name == "" {
		name = generateRandomHash(5)
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

func listNotes(basePath string) {
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			fmt.Println(path)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list notes: %v\n", err)
		os.Exit(1)
	}
}
