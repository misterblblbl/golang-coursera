package main

import (
	"os"
)

type byFilename []os.FileInfo

func (fi byFilename) Len() int {
	return len(fi)
}

func (fi byFilename) Swap(i, j int) {
	fi[i], fi[j] = fi[j], fi[i]
}

func (fi byFilename) Less(i, j int) bool {
	return fi[i].Name() < fi[j].Name()
}

func filterDirs(fileList []os.FileInfo) []os.FileInfo {
	var filteredList []os.FileInfo

	for _, file := range fileList {
		if file.IsDir() {
			filteredList = append(filteredList, file)
		}
	}

	return filteredList
}
