package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
)

type Root struct {
	Row []User `xml:"row"`
}

func decode(data []byte) []User {
	input := bytes.NewReader(data)
	decoder := xml.NewDecoder(input)

	var users []User
	for {
		token, tokenErr := decoder.Token()
		if tokenErr == io.EOF {
			break
		}

		if tokenErr != nil {
			fmt.Println("error happend", tokenErr)
			break
		}

		if token == nil {
			fmt.Println("t is nil break")
		}

		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == "row" {
				var user User
				err := decoder.DecodeElement(&user, &token)
				if err != nil {
					fmt.Println("error happend", err)
				}
				user.Name = user.FirstName + " " + user.LastName

				users = append(users, user)
			}
		}
	}

	return users
}

func main() {
	data, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		fmt.Println("Failed to read file")
	}

	decode(data)
}
