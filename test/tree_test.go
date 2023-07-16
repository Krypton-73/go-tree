package test

import (
	"go-tree/constant"
	"go-tree/internal"
	"os"
	"testing"
	// Replace with your package import path
)

func TestTree(t *testing.T) {
	// check if current directory is valid
	_, err := internal.IsValid(".")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Test Case 1: Empty Directory
	t.Run("Empty Directory", func(t *testing.T) {
		// tree with default flags
		tree := newTree()
		// Set up the empty directory
		dir := createEmptyDirectory()
		defer os.RemoveAll(dir)

		// Call the BuildTree function
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		// Add assertions for the expected output
		expected := internal.NewTreeSummary(2, 3)
		// Check if the output matches the expected summary
		if tree.Summary != expected {
			t.Errorf("BuildTree() for empty directory: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// has only directories tag
		val := true
		tree.Flags[constant.Dir] = &val
		tree.Summary = internal.NewTreeSummary(1, 0)
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for empty directory with dir tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// sorted by date modified
		m := true
		tree.Flags[constant.Time] = &m
		tree.Summary = internal.NewTreeSummary(1, 0)
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for empty directory with date modified tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// has level tag
		l := 1
		tree.Flags[constant.Level] = &l
		tree.Summary = internal.NewTreeSummary(1, 0)
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		expected = internal.NewTreeSummary(2, 3)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for empty directory with level tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}
	})

	// Test Case 2: Nested Empty Directories
	t.Run("Nested Empty Directories", func(t *testing.T) {
		// tree with default flags
		tree := newTree()
		// Set up the nested empty directories
		dir := createNestedEmptyDirectories()
		defer os.RemoveAll(dir)

		// Call the BuildTree function
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		// Add assertions for the expected output
		expected := internal.NewTreeSummary(4, 3)
		// Check if the output matches the expected summary
		if tree.Summary != expected {
			t.Errorf("BuildTree() for nested empty directories: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// only directories tag
		val := true
		tree.Flags[constant.Dir] = &val
		tree.Summary = internal.NewTreeSummary(1, 0)
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for nested empty directories with dir tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// sorted by date modified
		m := true
		tree.Flags[constant.Time] = &m
		tree.Summary = internal.NewTreeSummary(1, 0)
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for nested empty directories with date modified tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// if has level tag
		l := 2
		tree.Flags[constant.Level] = &l
		tree.Summary = internal.NewTreeSummary(1, 0)
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		expected = internal.NewTreeSummary(3, 3)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for nested empty directories with level tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}
	})

	// Test Case 3: Directory with Multiple Files
	t.Run("Directory with Multiple Files", func(t *testing.T) {
		// tree with default flags
		tree := newTree()
		// Set up the directory with multiple files
		noOfFiles := 10
		dir := createDirectoryWithFiles(noOfFiles)
		defer os.RemoveAll(dir)

		// Call the BuildTree function
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		// Add assertions for the expected output
		expected := internal.NewTreeSummary(2, 13)
		// Check if the output matches the expected summary
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with multiple files: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// only directories tag
		val := true
		tree.Flags[constant.Dir] = &val
		tree.Summary = internal.NewTreeSummary(1, 0)
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with multiple files and dir tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// sorted by date modified
		m := true
		tree.Flags[constant.Time] = &m
		tree.Summary = internal.NewTreeSummary(1, 0)
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with multiple files and date modified tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// if has level tag
		l := 1
		tree.Flags[constant.Level] = &l
		tree.Summary = internal.NewTreeSummary(1, 0)
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		expected = internal.NewTreeSummary(2, 3)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with multiple files and level tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}
	})

	// Test Case 4: Directory with Permission Issue
	t.Run("Directory with Permission Issue", func(t *testing.T) {
		// tree with default flags
		tree := newTree()
		// Set up the directory with permission issue
		dir := createDirectoryWithPermissionIssue()
		defer os.RemoveAll(dir)

		// Call the BuildTree function
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		// Add assertions for the expected output
		expected := internal.NewTreeSummary(3, 3)
		// Check if the output matches the expected summary
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with permission issue: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// only directories tag
		val := true
		tree.Flags[constant.Dir] = &val
		tree.Summary = internal.NewTreeSummary(1, 0)
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with permission issue and directory tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// sorted by date modified
		m := true
		tree.Flags[constant.Time] = &m
		tree.Summary = internal.NewTreeSummary(1, 0)
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with permission issue and date modified tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}

		// if has level tag
		l := 1
		tree.Flags[constant.Level] = &l
		tree.Summary = internal.NewTreeSummary(1, 0)
		tree.Root.BuildTree(tree.Flags, &tree.Summary)
		expected = internal.NewTreeSummary(2, 3)
		if tree.Summary != expected {
			t.Errorf("BuildTree() for directory with permission issue and level tag: \n output = %#v\n expected = %#v\n", tree.Summary, expected)
		}
	})
}
