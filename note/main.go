package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	DefaultBasePath = "notes"
	DefaultEditor   = "nvim"
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

func generateRandomHash(length int) string {
	// Implement a simple random hash generator
	// For simplicity, using a fixed string here
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range length {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate random hash: %v\n", err)
			os.Exit(1)
		}
		result[i] = charset[index.Int64()]
	}
	return string(result)
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

