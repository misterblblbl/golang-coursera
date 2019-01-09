package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		fmt.Println("Failed to read file")
	}

	decode(data)
}
