package internal

import (
	"bytes"
	"fmt"
	"go-tree/constant"
	"os"
	"path/filepath"
	"strings"
)

type treeNode struct {
	root     *treeNode
	children []treeNode
	depth    int
	isLast   bool
	path     string
	summary  map[string]int
	info     os.FileInfo
}

func (node *treeNode) buildTree(flags map[string]interface{}) error {
	dir, err := os.Open(node.path)
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.ReadDir(-1)
	if err != nil {
		return err
	}

	// Skip hidden files and directories
	files = exceptHiddens(files)

	dirs := justDirs(files)
	node.summary[totalDirs] = node.summary[totalDirs] + len(dirs)
	node.summary[totalFiles] = node.summary[totalFiles] + len(files) - len(dirs)

	// only list directories
	if justDir := *(flags[constant.Dir].(*bool)); justDir {
		files = dirs
	}

	// Sort files by time modified
	if sortByTime := *(flags[constant.Time].(*bool)); sortByTime {
		sortFilesByTimeModified(files)
	} else { // Sort files by name
		sortFilesByName(files)
	}

	for i, file := range files {
		isLast := false
		if i+1 == len(files) {
			isLast = true
		}

		childNode := treeNode{
			root:     node,
			children: nil,
			depth:    node.depth + 1,
			isLast:   isLast,
			path:     filepath.Join(node.path, file.Name()),
			summary:  node.summary,
			info:     getFileInfo(file),
		}

		maxDepth := *(flags[constant.Level].(*int))

		// Build tree upto max level
		if childNode.info.IsDir() && (maxDepth == 0 || childNode.depth < maxDepth) {
			// Build child node if directory has read permission
			if childNode.info.Mode().Perm()&0400 != 0 {
				if err := childNode.buildTree(flags); err != nil {
					return err
				}
			}
		}

		node.children = append(node.children, childNode)
	}
	return nil
}

func (node *treeNode) draw(out *bytes.Buffer, flags map[string]interface{}) {
	node.print(out, flags)
	for _, child := range node.children {
		if child.children != nil {
			child.draw(out, flags)
		} else {
			child.print(out, flags)
		}
	}
}

func (node *treeNode) print(out *bytes.Buffer, flags map[string]interface{}) {
	line := ""

	if node.root != nil {
		// indentation prefix for the line
		for p := node.root; p.root != nil; p = p.root {
			prefix := fmt.Sprintf("%s%s", "│", strings.Repeat(" ", 3))
			if p.isLast {
				prefix = strings.Repeat(" ", 4)
			}
			line = fmt.Sprintf("%s%s", prefix, line)
		}

		// adding suffix for the line
		suffix := "├── "
		if node.isLast {
			suffix = "└── "
		}
		line = fmt.Sprintf("%s%s", line, suffix)
	}

	name := filepath.Base(node.path)

	// print full path
	if hasPath := *(flags[constant.Path].(*bool)); hasPath || node.root == nil {
		name = node.path
	}

	// print file permissions
	if hasMode := *(flags[constant.Mode].(*bool)); hasMode && node.root != nil {
		name = fmt.Sprintf("[%v] %v", node.info.Mode(), name)
	}

	// Print msg if no read permission on directory
	errMsg := ""
	if node.info.Mode().Perm()&0400 == 0 {
		errMsg = fmt.Sprintf("%s[error opening dir]", strings.Repeat(" ", 4))
	}

	fmt.Fprintf(out, "%s%s%s\n", line, name, errMsg)
}
