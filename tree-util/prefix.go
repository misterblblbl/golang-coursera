package main

import "fmt"

func printSize(size int64) string {
	if size == 0 {
		return " (empty)"
	}

	return fmt.Sprintf(" (%vb)", size)
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
