package internal

import (
	"bytes"
	"fmt"
	"go-tree/constant"
	"os"
	"strings"
)

type TreeSummary struct {
	Directories int
	Files       int
}

type Tree struct {
	Root    TreeNode
	Flags   map[string]interface{}
	Summary TreeSummary
	Out     *bytes.Buffer
}

func NewTreeSummary(noOfDirectories int, noOfFiles int) TreeSummary {
	return TreeSummary{
		Directories: noOfDirectories,
		Files:       noOfFiles,
	}
}

func NewTree(root TreeNode, flags map[string]interface{}, summary TreeSummary, out *bytes.Buffer) Tree {
	return Tree{
		Root:    root,
		Flags:   flags,
		Summary: summary,
		Out:     out,
	}
}

// Draws a tree map
func DrawTree(flags map[string]interface{}) {
	rootPath := *(flags[constant.Root].(*string))
	// Check if path is a existing directory with read permission
	info, err := IsValid(rootPath)
	if err != nil {
		return
	}

	var out bytes.Buffer
	rootNode := NewTreeNode(nil, nil, 0, false, rootPath, info)
	summary := NewTreeSummary(1, 0)
	tree := NewTree(rootNode, flags, summary, &out)
	tree.draw()
}

func (t *Tree) draw() {
	// build directory tree map
	if err := t.Root.BuildTree(t.Flags, &t.Summary); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	indent := ""
	// draw in xml format
	if xml := *(t.Flags[constant.XML].(*bool)); xml { // draw in xml format
		t.Root.drawxml(indent+strings.Repeat(" ", 2), t.Flags, t.Out)
		t.printXmlTree()
		return
	}
	// draw in json format
	if json := *(t.Flags[constant.JSON].(*bool)); json {
		t.Root.drawjson(indent+strings.Repeat(" ", 2), t.Flags, t.Out)
		t.printJsonTree()
		return
	}
	// draw tree map
	t.Root.draw(indent, t.Flags, t.Out)
	t.printTree()
	return
}

func (t *Tree) printTree() {
	// print tree
	fmt.Println(t.Out)
	// print tree summary
	if justDirs := *(t.Flags[constant.Dir].(*bool)); justDirs {
		fmt.Printf("%v directories\n", t.Summary.Directories)
	} else {
		fmt.Printf("%v directories, %v files\n", t.Summary.Directories, t.Summary.Files)
	}
}

func (t *Tree) printXmlTree() {
	newline := fmt.Sprintf("\n")
	// print without indentation
	noIndent := *(t.Flags[constant.Indent].(*bool))
	if noIndent {
		newline = ""
	}
	fmt.Printf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>%s", newline)
	fmt.Printf("<tree>%s", newline)
	fmt.Printf("%s", t.Out)
	// print summary report
	indent := strings.Repeat(" ", 2)
	if noIndent {
		indent = ""
	}
	fmt.Printf("%s<report>%s", indent, newline)
	fmt.Printf("%s<directories>%v</directories>%s", strings.Repeat(indent, 2), t.Summary.Directories, newline)
	if justDirs := *(t.Flags[constant.Dir].(*bool)); !justDirs {
		fmt.Printf("%s<files>%v</files>%s", strings.Repeat(indent, 2), t.Summary.Files, newline)
	}
	fmt.Printf("%s</report>%s", indent, newline)
	fmt.Printf("</tree>%s", newline)
}

func (t *Tree) printJsonTree() {
	newline := fmt.Sprintf("\n")
	// print without indentation
	noIndent := *(t.Flags[constant.Indent].(*bool))
	if noIndent {
		newline = ""
	}
	fmt.Printf("[%s", newline)
	fmt.Printf("%s%s", t.Out, newline)
	fmt.Printf(",%s", newline)
	// print summary report
	indent := strings.Repeat(" ", 2)
	if noIndent {
		indent = ""
	}
	fmt.Printf("%s{\"type\":\"report\",\"directories\":%v", indent, t.Summary.Directories)
	if justDirs := *(t.Flags[constant.Dir].(*bool)); !justDirs {
		fmt.Printf(",\"files\":%v", t.Summary.Files)
	}
	fmt.Printf("}%s", newline)
	fmt.Printf("]\n")
}
