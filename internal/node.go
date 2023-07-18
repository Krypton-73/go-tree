package internal

import (
	"bytes"
	"fmt"
	"go-tree/constant"
	"os"
	"path/filepath"
	"strings"
)

type TreeNode struct {
	Root     *TreeNode
	Children []TreeNode
	Depth    int
	IsLast   bool
	Path     string
	Info     os.FileInfo
}

func NewTreeNode(root *TreeNode, children []TreeNode, depth int, isLast bool, path string, info os.FileInfo) TreeNode {
	return TreeNode{
		Root:     root,
		Children: children,
		Depth:    depth,
		IsLast:   isLast,
		Path:     path,
		Info:     info,
	}
}

func (node *TreeNode) BuildTree(flags map[string]interface{}, summary *TreeSummary) error {
	dir, err := os.Open(node.Path)
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
	// list of directories
	dirs := justDirs(files)
	// Add to tree summary
	summary.Directories += len(dirs)
	summary.Files += len(files) - len(dirs)
	// only list directories
	if justDir := *(flags[constant.Dir].(*bool)); justDir {
		files = dirs
	}
	// Sort files by time modified
	if sortByTime := *(flags[constant.Time].(*bool)); sortByTime {
		sortByModifiedTime(files)
	} else { // Sort files by name
		sortByName(files)
	}

	for i, file := range files {
		isLast := false
		if i+1 == len(files) {
			isLast = true
		}
		path := filepath.Join(node.Path, file.Name())
		info := getFileInfo(file)
		childNode := NewTreeNode(node, nil, node.Depth+1, isLast, path, info)

		// Build tree upto max level
		maxDepth := *(flags[constant.Level].(*int))
		if childNode.Info.IsDir() && (maxDepth == 0 || childNode.Depth < maxDepth) {
			// Build child node if directory has read permission
			if childNode.Info.Mode().Perm()&0400 != 0 {
				if err := childNode.BuildTree(flags, summary); err != nil {
					return err
				}
			}
		}
		node.Children = append(node.Children, childNode)
	}
	return nil
}

func (node *TreeNode) draw(indent string, flags map[string]interface{}, out *bytes.Buffer) {
	node.print(node.addSuffix(indent), flags, out)

	subIndent := node.addIndentation(indent)
	for _, child := range node.Children {
		if child.Children != nil {
			child.draw(subIndent, flags, out)
		} else {
			child.print(child.addSuffix(subIndent), flags, out)
		}
	}
}

// Print line
func (node *TreeNode) print(indent string, flags map[string]interface{}, out *bytes.Buffer) {
	// print without indentation
	if noIndent := *(flags[constant.Indent].(*bool)); noIndent {
		indent = ""
	}
	name := filepath.Base(node.Path)
	// print full path
	if hasPath := *(flags[constant.Path].(*bool)); hasPath || node.Root == nil {
		name = node.Path
	}
	// print file permissions
	if hasMode := *(flags[constant.Permission].(*bool)); hasMode && node.Root != nil {
		name = fmt.Sprintf("[%v] %v", node.Info.Mode(), name)
	}
	// print msg if no read permission on directory
	msg := ""
	if node.Info.Mode().Perm()&0400 == 0 {
		msg = fmt.Sprintf("%s[error opening dir]", strings.Repeat(" ", 4))
	}
	fmt.Fprintf(out, "%s%s%s\n", indent, name, msg)
}

// Indentation prefix
func (node *TreeNode) addIndentation(indent string) string {
	subIndent := ""
	if node.Root != nil {
		subIndent = fmt.Sprintf("%s%s%s", indent, "│", strings.Repeat(" ", 3))
		if node.IsLast {
			subIndent = fmt.Sprintf("%s%s", indent, strings.Repeat(" ", 4))
		}
	}
	return subIndent
}

// Add suffix for the line
func (node *TreeNode) addSuffix(prefix string) string {
	line := ""
	if node.Root != nil {
		suffix := "├── "
		if node.IsLast {
			suffix = "└── "
		}
		line = fmt.Sprintf("%s%s", prefix, suffix)
	}
	return line
}

// Prints the directory tree in JSON format
func (node *TreeNode) drawjson(indent string, flags map[string]interface{}, out *bytes.Buffer) {
	newline := fmt.Sprintf("\n")
	// print without indentation
	noIndent := *(flags[constant.Indent].(*bool))
	if noIndent {
		indent = ""
		newline = ""
	}
	filetype := getFileType(node.Info)
	name := node.Info.Name()
	if hasPath := *(flags[constant.Path].(*bool)); hasPath || node.Root == nil {
		name = node.Path
	}
	line := fmt.Sprintf("%s{\"type\":\"%s\",\"name\":\"%s\"", indent, filetype, name)
	if hasMode := *(flags[constant.Permission].(*bool)); hasMode {
		line = fmt.Sprintf("%s,\"mode\":\"%04o\",\"prot\":\"%v\"", line, node.Info.Mode().Perm(), node.Info.Mode())
	}

	if len(node.Children) > 0 {
		line = fmt.Sprintf("%s,\"contents\":[", line)
		fmt.Fprintf(out, "%s%s", line, newline)
		for i, child := range node.Children {
			if i > 0 {
				fmt.Fprintf(out, ",%s", newline)
			}
			child.drawjson(indent+strings.Repeat(" ", 2), flags, out)
		}
		fmt.Fprintf(out, "%s%s]", newline, indent)
	} else {
		fmt.Fprintf(out, line)
	}
	fmt.Fprintf(out, "}")
}

// Prints the directory tree in XML format.
func (node *TreeNode) drawxml(indent string, flags map[string]interface{}, out *bytes.Buffer) {
	newline := fmt.Sprintf("\n")
	noIndent := *(flags[constant.Indent].(*bool))
	if noIndent {
		indent = ""
		newline = ""
	}
	filetype := getFileType(node.Info)
	name := node.Info.Name()
	if hasPath := *(flags[constant.Path].(*bool)); hasPath || node.Root == nil {
		name = node.Path
	}
	line := fmt.Sprintf("%s<%s name=\"%s\"", indent, filetype, name)
	if hasMode := *(flags[constant.Permission].(*bool)); hasMode {
		line = fmt.Sprintf("%s mode=\"%04o\" prot=\"%v\"", line, node.Info.Mode().Perm(), node.Info.Mode())
	}

	if len(node.Children) > 0 {
		fmt.Fprintf(out, "%s>\n", line)
		for _, child := range node.Children {
			child.drawxml(indent+strings.Repeat(" ", 2), flags, out)
		}
		fmt.Fprintf(out, "%s</%s>%s", indent, filetype, newline)
	} else {
		fmt.Fprintf(out, "%s></%s>%s", line, filetype, newline)
	}
}
