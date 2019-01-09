package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

type UserRow struct {
	Id        int `xml:"id"`
	Name      string
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
}

type Root struct {
	Row []UserRow `xml:"row"`
}

func decode(data []byte) []UserRow {
	input := bytes.NewReader(data)
	decoder := xml.NewDecoder(input)

	var users []UserRow
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
				var user UserRow
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

func SearchServer(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("AccessToken") == "bad_token" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	query := r.FormValue("query")
	if query == "timeout" {
		time.Sleep(time.Second * 2)
	}

	if query == "internal_error" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if query == "invalid_json" {
		w.Write([]byte("invalid_json"))
		return
	}

	if query == "bad_request" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if query == "bad_request_unknown" {
		resp, _ := json.Marshal(SearchErrorResponse{"UnknownError"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	orderField := r.FormValue("order_field")
	if orderField == "" {
		orderField = "Name"
	}

	if orderField != "Id" && orderField != "Age" && orderField != "Name" {
		resp, _ := json.Marshal(SearchErrorResponse{"ErrorBadOrderField"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
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

	if query == "less_data" {
		resp, _ := json.Marshal(users[offset:(limit + offset - 5)])
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(users[offset:(limit + offset)])
	w.Write(resp)
	return
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
		URL:         testServer.URL,
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

func TestFindUsersInternalError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer testServer.Close()

	client := &SearchClient{
		URL:         testServer.URL,
		AccessToken: "good_token",
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 1, Query: "internal_error"})
	if err == nil {
		t.Errorf("Response should return error, [%v]", resp)
	}

	if !strings.Contains(err.Error(), "SearchServer fatal error") {
		t.Errorf("Wrong error type received, [%v]", err)
	}
}

func TestFindUsersBadRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer testServer.Close()

	client := &SearchClient{
		URL:         testServer.URL,
		AccessToken: "good_token",
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 1, Query: "bad_request"})
	if err == nil {
		t.Errorf("Response should return error, [%v]", resp)
	}

	if !strings.Contains(err.Error(), "cant unpack error json") {
		t.Errorf("Wrong error type received, [%v]", err)
	}
}

func TestFindUsersBadRequestUnknown(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer testServer.Close()

	client := &SearchClient{
		URL:         testServer.URL,
		AccessToken: "good_token",
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 1, Query: "bad_request_unknown"})
	if err == nil {
		t.Errorf("Response should return error, [%v]", resp)
	}

	if !strings.Contains(err.Error(), "unknown bad request error") {
		t.Errorf("Wrong error type received, [%v]", err)
	}
}

func TestFindUsersInvalidJson(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer testServer.Close()

	client := &SearchClient{
		URL:         testServer.URL,
		AccessToken: "good_token",
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 1, Query: "invalid_json"})
	if err == nil {
		t.Errorf("Response should return error, [%v]", resp)
	}

	if !strings.Contains(err.Error(), "cant unpack result json") {
		t.Errorf("Wrong error type received, [%v]", err)
	}
}

func TestFindUsersOrderField(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer testServer.Close()

	client := &SearchClient{
		URL:         testServer.URL,
		AccessToken: "good_token",
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 1, OrderField: "order_field"})
	if err == nil {
		t.Errorf("Response should return error, [%v]", resp)
	}

	if !strings.Contains(err.Error(), "OrderFeld order_field invalid") {
		t.Errorf("Wrong error type received, [%v]", err)
	}
}

func TestFindUsersUnknownErr(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer testServer.Close()

	client := &SearchClient{
		URL:         "http://",
		AccessToken: "good_token",
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 1})
	if err == nil {
		t.Errorf("Response should return error, [%v]", resp)
	}

	if !strings.Contains(err.Error(), "unknown error") {
		t.Errorf("Wrong error type received, [%v]", err)
	}
}

func TestFindUsersAuth(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer testServer.Close()

	client := &SearchClient{
		URL:         testServer.URL,
		AccessToken: "bad_token",
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 1})
	if err == nil {
		t.Errorf("Response should return error, [%v]", resp)
	}

	if !strings.Contains(err.Error(), "Bad AccessToken") {
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
		TestCase{
			Description: "Offset: 5, Limit: 10",
			Request: SearchRequest{
				Limit: 14,
				Query: "less_data",
			},
			ExpectedTotal: 10,
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
