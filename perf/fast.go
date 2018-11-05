package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type User struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Job      string   `json:"job"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
}

func hasItem(array []string, item string) bool {
	isFound := false
	for _, arrItem := range array {
		if arrItem == item {
			isFound = true
			return isFound
		}
	}

	return isFound
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var seenBrowsers []string
	var foundUsers string

	lines := strings.Split(string(fileContents), "\n")
	linesNumber := len(lines)
	users := make([]User, linesNumber)

	for i, line := range lines {
		var user User

		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}
		users[i] = user

		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				if !hasItem(seenBrowsers, browser) {
					seenBrowsers = append(seenBrowsers, browser)
				}
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				if !hasItem(seenBrowsers, browser) {
					seenBrowsers = append(seenBrowsers, browser)
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		email := strings.Replace(user.Email, "@", " [at] ", 1)
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))

	// fmt.Println("found users:\n" + foundUsers)
	// fmt.Println("Total unique browsers", len(seenBrowsers))

	// for _, b := range seenBrowsers {
	// 	fmt.Printf("B %s \n", b)
	// }
}

func main() {
	FastSearch(ioutil.Discard)

	fmt.Println("*********************************************************")

	SlowSearch(ioutil.Discard)
}
