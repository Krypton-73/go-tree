package internal

import (
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"
)

func IsValid(rootPath string) (fs.FileInfo, error) {
	fileInfo, err := os.Stat(rootPath)
	var dir, file int
	invalid := false

	if err != nil {
		invalid = true
	} else if !fileInfo.IsDir() {
		invalid = true
		file = 1
	} else if fileInfo.Mode().Perm()&0400 == 0 {
		invalid = true
		dir = 1
	}

	if invalid {
		fmt.Printf("%s%s[error opening dir]\n", rootPath, strings.Repeat(" ", 4))
		fmt.Printf("\n%v directories, %v files\n", dir, file)
		return nil, fs.ErrExist
	}
	return fileInfo, nil
}

func exceptHiddens(files []fs.DirEntry) []fs.DirEntry {
	result := []fs.DirEntry{}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		result = append(result, file)
	}
	return result
}

func justDirs(files []fs.DirEntry) []fs.DirEntry {
	dirs := []fs.DirEntry{}
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
	}
	return dirs
}

func sortByModifiedTime(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		a := getFileInfo(files[i])
		b := getFileInfo(files[j])
		if a.ModTime().Before(b.ModTime()) {
			return false
		}
		return true
	})
}

func sortByName(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
}

func getFileInfo(file fs.DirEntry) fs.FileInfo {
	fileInfo, err := file.Info()
	if err != nil {
		fmt.Println(err)
	}
	return fileInfo
}

func hasReadPermission(file fs.DirEntry) bool {
	fileInfo := getFileInfo(file)
	if fileInfo.Mode().Perm()&0400 != 0 {
		return true
	}
	return false
}

func getFileType(f fs.FileInfo) string {
	filetype := "directory"
	if !f.IsDir() {
		filetype = "file"
	}
	return filetype
}
