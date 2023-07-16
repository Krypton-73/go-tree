package test

import (
	"bytes"
	"go-tree/constant"
	"go-tree/internal"
	"os"
	"testing"

	"github.com/spf13/cobra"
	// Replace with your package import path
)

var flags map[string]interface{}

var goTree = &cobra.Command{
	Use:   "./main",
	Short: "unix command \"tree\" implementation in go",
	Long:  "go-tree is a cli tool which draws a tree of directory structure",
	Run: func(cmd *cobra.Command, args []string) {
		internal.DrawTree(flags)
	},
}

func init() {
	flags = map[string]interface{}{}
	flags[constant.Root] = goTree.PersistentFlags().StringP(constant.Root, "r", ".", "Root path of the tree")
	flags[constant.Path] = goTree.PersistentFlags().BoolP(constant.Path, "f", false, "Flag to show fullpaths")
	flags[constant.Dir] = goTree.PersistentFlags().BoolP(constant.Dir, "d", false, "Flag to only list directories")
	flags[constant.Level] = goTree.PersistentFlags().IntP(constant.Level, "L", 0, "Max level of tree depth")
	flags[constant.Permission] = goTree.PersistentFlags().BoolP(constant.Permission, "p", false, "Flag to show permission modes")
	flags[constant.Time] = goTree.PersistentFlags().BoolP(constant.Time, "t", false, "Flag to sort output by modified time")
	flags[constant.JSON] = goTree.PersistentFlags().BoolP(constant.JSON, "J", false, "Prints tree in JSON format")
	flags[constant.XML] = goTree.PersistentFlags().BoolP(constant.XML, "X", false, "Prints tree in XML format")
	//Do not print indentation lines
}

func TestTree(t *testing.T) {
	info, err := internal.IsValid(".")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	node := internal.TreeNode{
		Root:     nil,
		Children: nil,
		Depth:    0,
		IsLast:   false,
		Path:     ".",
		Info:     info,
	}
	var out bytes.Buffer
	summary := internal.TreeSummary{Directories: 1, Files: 0}
	tree := internal.Tree{Root: node, Flags: flags, Summary: summary, Out: &out}

	// Test Case 1: Empty Directory
	t.Run("Empty Directory", func(t *testing.T) {
		// Set up the empty directory
		dir := createEmptyDirectory()
		defer os.RemoveAll(dir)

		// Call the BuildTree function
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		// Add assertions for the expected output
		expected := internal.TreeSummary{Directories: 2, Files: 2}
		// Check if the output matches the expected summary
		if tree.Summary != expected {
			t.Errorf("BuildTree() for empty directory: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// has only directories tag
		val := true
		flags[constant.Dir] = &val
		tree.Summary = internal.TreeSummary{Directories: 1, Files: 0}
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for empty directory with dir tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// sorted by date modified
		m := true
		flags[constant.Time] = &m
		tree.Summary = internal.TreeSummary{Directories: 1, Files: 0}
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for empty directory with date modified tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// has level tag
		l := 1
		flags[constant.Level] = &l
		tree.Summary = internal.TreeSummary{Directories: 1, Files: 0}
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		expected = internal.TreeSummary{Directories: 2, Files: 2}
		if tree.Summary != expected {
			t.Errorf("BuildTree() for empty directory with level tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}
	})

	// Test Case 2: Nested Empty Directories
	t.Run("Nested Empty Directories", func(t *testing.T) {

		// Set up the nested empty directories
		dir := createNestedEmptyDirectories()
		defer os.RemoveAll(dir)

		// Call the BuildTree function
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		// Add assertions for the expected output
		expected := internal.TreeSummary{Directories: 4, Files: 2}
		// Check if the output matches the expected summary
		if tree.Summary != expected {
			t.Errorf("BuildTree() for nested empty directories: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// only directories tag
		val := true
		flags[constant.Dir] = &val
		tree.Summary = internal.TreeSummary{Directories: 1, Files: 0}
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for nested empty directories with dir tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// sorted by date modified
		m := true
		flags[constant.Time] = &m
		tree.Summary = internal.TreeSummary{Directories: 1, Files: 0}
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for nested empty directories with date modified tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// if has level tag
		l := 2
		flags[constant.Level] = &l
		tree.Summary = internal.TreeSummary{Directories: 1, Files: 0}
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		expected = internal.TreeSummary{Directories: 3, Files: 2}
		if tree.Summary != expected {
			t.Errorf("BuildTree() for nested empty directories with level tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}
	})

	// Test Case 3: Directory with Multiple Files
	t.Run("Directory with Multiple Files", func(t *testing.T) {
		// Set up the directory with multiple files
		noOfFiles := 10
		dir := createDirectoryWithFiles(noOfFiles)
		defer os.RemoveAll(dir)

		// Call the BuildTree function
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		// Add assertions for the expected output
		expected := internal.TreeSummary{Directories: 2, Files: 12}
		// Check if the output matches the expected summary
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with multiple files: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// only directories tag
		val := true
		flags[constant.Dir] = &val
		tree.Summary = internal.TreeSummary{Directories: 1, Files: 0}
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with multiple files and dir tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// sorted by date modified
		m := true
		flags[constant.Time] = &m
		tree.Summary = internal.TreeSummary{Directories: 1, Files: 0}
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with multiple files and date modified tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// if has level tag
		l := 1
		flags[constant.Level] = &l
		tree.Summary = internal.TreeSummary{Directories: 1, Files: 0}
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		expected = internal.TreeSummary{Directories: 2, Files: 2}
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with multiple files and level tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}
	})

	// Test Case 4: Directory with Permission Issue
	t.Run("Directory with Permission Issue", func(t *testing.T) {
		// Set up the directory with permission issue
		dir := createDirectoryWithPermissionIssue()
		defer os.RemoveAll(dir)

		// Call the BuildTree function
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		// Add assertions for the expected output
		expected := internal.TreeSummary{Directories: 3, Files: 2}
		// Check if the output matches the expected summary
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with permission issue: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// only directories tag
		val := true
		flags[constant.Dir] = &val
		tree.Summary = internal.TreeSummary{Directories: 1, Files: 0}
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with permission issue and directory tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// sorted by date modified
		m := true
		flags[constant.Time] = &m
		tree.Summary = internal.TreeSummary{Directories: 1, Files: 0}
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with permission issue and date modified tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// if has level tag
		l := 1
		flags[constant.Level] = &l
		tree.Summary = internal.TreeSummary{Directories: 1, Files: 0}
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		expected = internal.TreeSummary{Directories: 2, Files: 2}
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with permission issue and level tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}
	})
}
