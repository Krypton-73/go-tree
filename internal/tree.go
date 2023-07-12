package internal

import (
	"bytes"
	"fmt"
	"go-tree/constant"
)

const (
	totalDirs  = "dir"
	totalFiles = "file"
)

type tree struct {
	root  treeNode
	flags map[string]interface{}
	out   *bytes.Buffer
}

// Draws a tree map
func DrawTree(flags map[string]interface{}) {
	rootPath := *(flags[constant.Root].(*string))

	// Check if path is a existing directory with read permission
	info, err := isValid(rootPath)
	if err != nil {
		return
	}

	var out bytes.Buffer
	node := treeNode{
		root:     nil,
		children: nil,
		depth:    0,
		isLast:   false,
		path:     rootPath,
		summary: map[string]int{
			totalDirs:  1,
			totalFiles: 0,
		},
		info: info,
	}

	tree := tree{node, flags, &out}
	tree.draw()
}

func (t tree) draw() {
	if err := t.root.buildTree(t.flags); err != nil {
		fmt.Println(err)
	}

	t.root.draw(t.out, t.flags)

	// Print tree
	fmt.Println(t.out)
	// Print Summary
	if justDirs := *(t.flags[constant.Dir].(*bool)); justDirs {
		fmt.Printf("%v directories\n", t.root.summary[totalDirs])
	} else {
		fmt.Printf("%v directories, %v files\n", t.root.summary[totalDirs], t.root.summary[totalFiles])
	}
}
