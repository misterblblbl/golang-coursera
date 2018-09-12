package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
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

func getPrefix(depth int, isLast bool, openedDirs map[int]bool) string {
	prefix := "├───"
	tabulation := ""

	if isLast {
		prefix = "└───"
	}

	for level := 0; level < depth; level++ {
		dirOpened, ok := openedDirs[level]
		if ok {
			nextTab := "\t"
			if dirOpened {
				nextTab = "│\t"
			}

			tabulation = tabulation + nextTab
		}
	}

	return tabulation + prefix
}

func printDir(output io.Writer, path string, openedDirs map[int]bool, depth int) {
	stats, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	if mode := stats.Mode(); mode.IsDir() {
		fileInfo, _ := file.Readdir(-1)
		openedDirs[depth] = true
		sort.Sort(byFilename(fileInfo))

		for i, file := range fileInfo {
			isLast := i+1 == len(fileInfo)
			prefix := getPrefix(depth, isLast, openedDirs)
			openedDirs[depth] = !isLast

			if file.IsDir() {
				nextPath := fmt.Sprintf("%s/%s", path, file.Name())
				fmt.Fprintf(output, "%s%s\n", prefix, file.Name())

				printDir(output, nextPath, openedDirs, depth+1)
			} else {
				if file.Name() != ".DS_Store" {
					fmt.Fprintf(output, "%s%s\n", prefix, file.Name())
				}
			}
		}
	} else {
		prefix := getPrefix(depth, false, openedDirs)

		if file.Name() != ".DS_Store" {
			fmt.Fprintf(output, "%s%s\n", prefix, file.Name())
		}
	}
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	openedDirs := map[int]bool{}
	printDir(output, path, openedDirs, 0)

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
