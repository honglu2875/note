package note

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CheckBasePath verifies if the basePath exists and is a directory
func CheckBasePath(basePath string) (bool, error) {
	info, err := os.Stat(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error checking base path: %w", err)
	}
	return info.IsDir(), nil
}

// CheckGitRepo checks if the basePath is initialized as a git repository
func CheckGitRepo(basePath string) (bool, error) {
	gitPath := filepath.Join(basePath, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error checking git repository: %w", err)
	}
	return info.IsDir(), nil
}

// execGitCommand executes a git command in the specified directory
func execGitCommand(basePath string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = basePath
	return cmd.CombinedOutput()
}

// ensureGitRepo validates that basePath is a git repository
func ensureGitRepo(basePath string) error {
	isRepo, err := CheckGitRepo(basePath)
	if err != nil {
		return err
	}
	if !isRepo {
		return fmt.Errorf("not a git repository: %s", basePath)
	}
	return nil
}

// hasChanges checks if there are any staged or unstaged changes in the repository
func hasChanges(basePath string) (bool, error) {
	output, err := execGitCommand(basePath, "status", "--porcelain")
	if err != nil {
		return false, fmt.Errorf("git status failed: %w, output: %s", err, output)
	}
	return len(output) > 0, nil
}

// InitGitRepo initializes a git repository in the specified path
func InitGitRepo(basePath string) error {
	// Check if directory exists first
	exists, err := CheckBasePath(basePath)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	if !exists {
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Check if it's already a git repo
	isRepo, err := CheckGitRepo(basePath)
	if err != nil {
		return err
	}
	if isRepo {
		return nil // Already a git repo, nothing to do
	}

	// Initialize git repository
	output, err := execGitCommand(basePath, "init")
	if err != nil {
		return fmt.Errorf("git init failed: %w, output: %s", err, output)
	}

	return nil
}

// CommitChanges adds all changes and commits them with the given message
func CommitChanges(basePath, message string) error {
	// Ensure it's a git repo
	if err := ensureGitRepo(basePath); err != nil {
		return err
	}

	// Git add all
	output, err := execGitCommand(basePath, "add", ".")
	if err != nil {
		return fmt.Errorf("git add failed: %w, output: %s", err, output)
	}

	// Check if there are changes to commit
	hasChanges, err := hasChanges(basePath)
	if err != nil {
		return err
	}

	// If there are no changes, return early
	if !hasChanges {
		return nil
	}

	// Git commit
	output, err = execGitCommand(basePath, "commit", "-m", message)
	if err != nil {
		return fmt.Errorf("git commit failed: %w, output: %s", err, output)
	}

	return nil
}

// StashChanges stashes the current changes in the git repository
func StashChanges(basePath string, message string) error {
	// Ensure it's a git repo
	if err := ensureGitRepo(basePath); err != nil {
		return err
	}

	// Check if there are changes to stash
	hasChanges, err := hasChanges(basePath)
	if err != nil {
		return err
	}

	// If there are no changes, return early
	if !hasChanges {
		return nil
	}

	// Execute the stash command
	var output []byte
	if message == "" {
		output, err = execGitCommand(basePath, "stash")
	} else {
		output, err = execGitCommand(basePath, "stash", "push", "-m", message)
	}

	if err != nil {
		return fmt.Errorf("git stash failed: %w, output: %s", err, output)
	}

	return nil
}
