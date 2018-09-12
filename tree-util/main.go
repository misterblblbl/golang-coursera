package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

var ignoredFiles = []string{
	".DS_Store",
	".gitignore",
}

func checkIgnoredFile(file string) bool {
	var isIgnored bool
	for _, ignored := range ignoredFiles {
		if ignored == file {
			isIgnored = true
		}
	}

	return isIgnored
}

func printDir(output io.Writer, path string, openedDirs map[int]bool, depth int, printFiles bool) {
	stats, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	fileInfo, _ := file.Readdir(-1)
	if !printFiles {
		fileInfo = filterDirs(fileInfo)
	}
	openedDirs[depth] = true
	sort.Sort(byFilename(fileInfo))

	for i, file := range fileInfo {
		fileName := file.Name()
		if checkIgnoredFile(fileName) {
			continue
		}

		isLast := i+1 == len(fileInfo)
		openedDirs[depth] = !isLast
		prefix := getPrefix(depth, isLast, openedDirs)

		if file.IsDir() {
			fmt.Fprintf(output, "%s%s\n", prefix, fileName)
			nextPath := fmt.Sprintf("%s/%s", path, fileName)

			printDir(output, nextPath, openedDirs, depth+1, printFiles)
		} else if printFiles {
			fileSize := printSize(stats.Size())

			fmt.Fprintf(output, "%s%s%s\n", prefix, fileName, fileSize)
		}
	}
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	openedDirs := map[int]bool{}

	stats, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	if mode := stats.Mode(); mode.IsDir() {
		printDir(output, path, openedDirs, 0, printFiles)
	}

	return nil
}

func main() {
	out := os.Stdout

	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
