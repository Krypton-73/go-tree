package test

import (
	"fmt"
	"os"
	"path/filepath"
)

// Helper function to create an empty directory for testing
func createEmptyDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, "emptyDir")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}
	return dir
}

// Helper function to create nested empty directories for testing
func createNestedEmptyDirectories() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, "nested_empty_directories")
	err = os.MkdirAll(filepath.Join(dir, "subdir1", "subdir2"), 0755)
	if err != nil {
		panic(err)
	}
	return dir
}

// Helper function to create a directory with multiple files for testing
func createDirectoryWithFiles(N int) string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, "directory_with_files")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	// Create multiple files in the directory
	for i := 1; i <= N; i++ {
		filePath := filepath.Join(dir, fmt.Sprintf("file%v.txt", i))
		file, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	return dir
}

// Helper function to create a directory with permission issue for testing
func createDirectoryWithPermissionIssue() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, "directory_with_permission_issue")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	// Create a nested directory without read/execute permission
	subdir := filepath.Join(dir, "subdir")
	err = os.Mkdir(subdir, 0333)
	if err != nil {
		panic(err)
	}

	return dir
}
