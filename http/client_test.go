package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func SearchServer(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	if query == "timeout" {
		time.Sleep(time.Second * 2)
	}

	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad offset param"))
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad limit param"))
	}

	rawData, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		fmt.Println("Failed to read file")
	}

	users := decode(rawData)
	resp, _ := json.Marshal(users[offset:(limit + offset)])

	w.Write(resp)
}

type TestCase struct {
	Description   string
	Request       SearchRequest
	ExpectedTotal int
}

func TestFindUsersTimeout(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer testServer.Close()

	client := &SearchClient{
		URL:         testServer.URL + "?timeout=true",
		AccessToken: "good_token",
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 1, Query: "timeout"})
	if err == nil {
		t.Errorf("Response should return error, [%v]", resp)
	}

	if !strings.Contains(err.Error(), "timeout for") {
		t.Errorf("Wrong error type received, [%v]", err)
	}
}

func TestFindUsersBadParams(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer testServer.Close()

	client := &SearchClient{
		URL:         testServer.URL,
		AccessToken: "good_token",
	}

	cases := []TestCase{
		TestCase{
			Description: "Bad limit param",
			Request: SearchRequest{
				Limit: -1,
			},
		},
		TestCase{
			Description: "Bad limit param",
			Request: SearchRequest{
				Offset: -30,
			},
		},
	}

	for _, testCase := range cases {
		resp, err := client.FindUsers(testCase.Request)
		if err == nil {
			t.Errorf("Response should return error, [%v]", resp)
		}
	}
}

func TestFindUsersGoodParam(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer testServer.Close()

	client := &SearchClient{
		URL:         testServer.URL,
		AccessToken: "good_token",
	}

	cases := []TestCase{
		TestCase{
			Description: "Offset: 0, Limit: 1",
			Request: SearchRequest{
				Limit: 1,
			},
			ExpectedTotal: 1,
		},
		TestCase{
			Description: "Offset: 5, Limit: 10",
			Request: SearchRequest{
				Offset: 5,
				Limit:  5,
			},
			ExpectedTotal: 5,
		},
		TestCase{
			Description: "Offset: 5, Limit: 10",
			Request: SearchRequest{
				Offset: 4,
				Limit:  31,
			},
			ExpectedTotal: 25,
		},
		TestCase{
			Description: "Offset: 5, Limit: 10",
			Request: SearchRequest{
				Limit:  26,
				Offset: 1,
			},
			ExpectedTotal: 25,
		},
	}

	for _, testCase := range cases {
		resp, err := client.FindUsers(testCase.Request)

		if resp == nil || err != nil {
			t.Errorf("Expected response, got error: [%s] [%v]", err, resp)
		}

		if len(resp.Users) != testCase.ExpectedTotal {
			t.Errorf("Expected response length doesn't match: recieved [%d] - expected [%d]", len(resp.Users), testCase.ExpectedTotal)
		}
	}
}
