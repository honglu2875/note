package note

import (
	"os"
	"path/filepath"
	"strings"
)

type FileNode struct {
	Name     string
	Path     string
	IsDir    bool
	Children map[string]*FileNode
	Id       int64
}

func BuildTree(basePath string) ([]*FileNode, error) {
	var counter int64 = 1 // root is No. 0
	var all_nodes = []*FileNode{}

	root := &FileNode{
		Name:     filepath.Base(basePath),
		Path:     basePath,
		IsDir:    true,
		Children: make(map[string]*FileNode),
		Id:       0,
	}

	all_nodes = append(all_nodes, root)

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == basePath {
			return nil
		}

		// Only include markdown files and directories
		if !info.IsDir() && !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Get relative path from the base
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}

		// Split path into components
		parts := strings.Split(relPath, string(os.PathSeparator))

		// Add to the tree (max 2 levels deep)
		currentNode := root
		for i, part := range parts {
			if i >= 2 { // Limit to 2 levels deep (excluding root)
				break
			}

			if _, exists := currentNode.Children[part]; !exists {
				new_node := &FileNode{
					Name:     part,
					Path:     filepath.Join(currentNode.Path, part),
					IsDir:    i < len(parts)-1 || info.IsDir(),
					Children: make(map[string]*FileNode),
					Id:       counter,
				}
				currentNode.Children[part] = new_node
				if !info.IsDir() {
					counter++
					all_nodes = append(all_nodes, new_node)
				}
			}
			currentNode = currentNode.Children[part]
		}

		return nil
	})
	return all_nodes, err
}
