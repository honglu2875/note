package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/muesli/termenv"
)

func maybeAddSuffix(name string) string {
	if name[len(name)-3:] == ".md" {
		return name
	}
	return name + ".md"
}

func color(msg string, num int) termenv.Style {
	return termenv.String(msg).Foreground(termenv.ANSI256Color(num))
}

func removeEmptyFirstLevelFolders(rootPath string) ([]string, error) {
	// Check if rootPath exists and is a directory
	rootInfo, err := os.Stat(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to access path %s: %w", rootPath, err)
	}
	if !rootInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", rootPath)
	}

	// Get all first-level entries
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", rootPath, err)
	}

	var removedFolders []string

	// Iterate through entries
	for _, entry := range entries {
		// Skip if not a directory
		if !entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(rootPath, entry.Name())

		// Check if directory is empty
		isEmpty, err := isDirEmpty(fullPath)
		if err != nil {
			fmt.Printf("Warning: couldn't check if %s is empty: %v\n", fullPath, err)
			continue
		}

		// Remove empty directory
		if isEmpty {
			if err := os.Remove(fullPath); err != nil {
				fmt.Printf("Warning: couldn't remove empty folder %s: %v\n", fullPath, err)
				continue
			}
			removedFolders = append(removedFolders, entry.Name())
		}
	}

	return removedFolders, nil
}

// isDirEmpty checks if a directory is empty
func isDirEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Read just one entry
	_, err = f.Readdirnames(1)
	if err == fs.ErrNotExist {
		return false, err
	}

	// If we got EOF, the directory is empty
	return err == io.EOF, nil
}
