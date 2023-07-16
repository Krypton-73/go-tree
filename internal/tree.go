package internal

import (
	"bytes"
	"fmt"
	"go-tree/constant"
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

// Draws a tree map
func DrawTree(flags map[string]interface{}) {
	rootPath := *(flags[constant.Root].(*string))

	// Check if path is a existing directory with read permission
	info, err := IsValid(rootPath)
	if err != nil {
		return
	}

	node := TreeNode{
		Root:     nil,
		Children: nil,
		Depth:    0,
		IsLast:   false,
		Path:     rootPath,
		Info:     info,
	}
	var out bytes.Buffer
	summary := TreeSummary{Directories: 1, Files: 0}
	tree := Tree{node, flags, summary, &out}
	tree.draw()
}

func (t *Tree) draw() {
	// build directory tree map
	if err := t.Root.BuildTree(t.Flags, &t.Summary); err != nil {
		fmt.Println(err)
	}

	indent := ""
	// draw in xml format
	if xml := *(t.Flags[constant.XML].(*bool)); xml { // draw in xml format
		t.Root.drawXmlTree(indent+strings.Repeat(" ", 2), t.Flags, t.Out)
		t.printXmlTree()
		return
	}
	// draw in json format
	if json := *(t.Flags[constant.JSON].(*bool)); json {
		t.Root.drawJsonTree(indent+strings.Repeat(" ", 2), t.Flags, t.Out)
		t.printJsonTree()
		return
	}
	// draw directory tree map
	t.Root.draw(indent, t.Flags, t.Out)
	t.printTreeMap()
	return
}

func (t *Tree) printTreeMap() {
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
	hasIndent := *(t.Flags[constant.Indent].(*bool))
	if hasIndent {
		newline = ""
	}
	fmt.Printf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>%s", newline)
	fmt.Printf("<tree>%s", newline)
	fmt.Printf("%s", t.Out)
	// print summary report
	indent := strings.Repeat(" ", 2)
	if hasIndent {
		indent = ""
	}
	fmt.Printf("%s<report>%s", indent, newline)
	fmt.Printf("%s<directories>%v</directories>%s", indent+strings.Repeat(" ", 2), t.Summary.Directories, newline)
	if justDirs := *(t.Flags[constant.Dir].(*bool)); !justDirs {
		fmt.Printf("%s<files>%v</files>%s", indent+strings.Repeat(" ", 2), t.Summary.Files, newline)
	}
	fmt.Printf("%s</report>%s", indent, newline)
	fmt.Printf("</tree>%s", newline)
}

func (t *Tree) printJsonTree() {
	newline := fmt.Sprintf("\n")
	// print without indentation
	hasIndent := *(t.Flags[constant.Indent].(*bool))
	if hasIndent {
		newline = ""
	}
	fmt.Printf("[%s", newline)
	fmt.Printf("%s%s", t.Out, newline)
	fmt.Printf(",%s", newline)
	// print summary report
	indent := strings.Repeat(" ", 2)
	if hasIndent {
		indent = ""
	}
	fmt.Printf("%s{\"type\":\"report\",\"directories\":%v", indent, t.Summary.Directories)
	if justDirs := *(t.Flags[constant.Dir].(*bool)); !justDirs {
		fmt.Printf(",\"files\":%v", t.Summary.Directories)
	}
	fmt.Printf("}%s", newline)
	fmt.Printf("]\n")
}
