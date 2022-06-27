package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestGetbooks(t *testing.T) {
	test := []struct {
		desc     string
		endpoint string
		output   *Book
	}{
		{"success:valid details ", "books", &Book{1, Author{1, "Chetan", "Bhagat", "14/05/1989", "Chetan"}, "400 Days", "Penguin", "04/11/2006"}},
		{"invalid id", "books", &Book{-11, Author{3, "Rao", "kumar", "11/12/1999", "Rao"}, "earth", "East", "12/01/1888"}},
		{"missing first name", "books", &Book{10, Author{4, "", "last", "17/12/2000", ""}, "", "east", "14/12/2001"}},
		{"missing last name", "books", &Book{10, Author{4, "anurag", "", "17/12/2000", ""}, "", "east", "14/12/2001"}},
		{"missing dob", "books", &Book{10, Author{4, "nitish", "kashyap", "", ""}, "", "east", "14/12/2001"}},
		{"invalid Publication", "books", &Book{8, Author{8, "james", "bond", "02/09/1992", "bond"}, "village fields", "village", "21/03/2012"}},
	}

	for i, v := range test {

		req := httptest.NewRequest(http.MethodGet, "localhost:8000/"+v.endpoint, nil)

		w := httptest.NewRecorder()

		GetBooks(w, req)

		body, err := io.ReadAll(w.Body)
		if err != nil {
			log.Printf("%v", err)
		}
		var output *Book

		err = json.Unmarshal(body, &output)
		if err != nil {
			log.Printf("%v", err)
		}

		if reflect.DeepEqual(output, v.output) {
			t.Errorf("Desc : %v,Testcase[%v], expected output : %v, actual output : %v", v.desc, i, output, v.output)
		}

	}
}

func TestGetId(t *testing.T) {
	test := []struct {
		desc    string
		inputid string
		output  Book
	}{
		{"valid", "1", Book{1, Author{1, "Chetan", "Bhagat", "14/05/1989", "Chetan"}, "400 Days", "Penguin", "04/11/2006"}},
		{"invalid id", "-1", Book{2, Author{3, "Rao", "kumar", "11/12/1999", "Rao"}, "Talking", "Penguin", "12/03/2012"}},
		{"missing fields", "", Book{10, Author{}, "", "east", "14/12/2001"}},
	}
	for i, v := range test {
		params := url.Values{}
		params.Add("bookid", v.inputid)

		req := httptest.NewRequest(http.MethodGet, "localhost:8000/book?"+params.Encode(), nil)

		w := httptest.NewRecorder()

		GetBookId(w, req)

		res := w.Body

		body, err := io.ReadAll(res)
		if err != nil {
			log.Printf("%v", err)
		}
		var output Book

		err = json.Unmarshal(body, &output)
		if err != nil {
			log.Printf("%v", err)
		}

		if reflect.DeepEqual(output, v.output) {
			t.Errorf("Desc : %v,Testcase[%v], expected output : %v, actual output : %v", v.desc, i, output, v.output)
		}
	}
}

func TestPostAuthor(t *testing.T) {
	test := []struct {
		Desc       string `json:"desc"`
		Auth       Author `json:"author"`
		StatusCode int    `json:"status_code"`
	}{
		{"valid details", Author{1, "Chetan", "Bhagat", "06/04/2001", "Chetan"}, http.StatusOK},
		{"empty first name", Author{4, "", "Bhagat", "06/04/2001", "Chetan"}, http.StatusBadRequest},
		{"empty last name", Author{3, "Chetan", "", "06/04/2001", "Chetan"}, http.StatusBadRequest},
		{"empty dob", Author{6, "Chetan", "Bhagat", "", "Chetan"}, http.StatusBadRequest},
		{"empty pen name", Author{7, "Chetan", "Bhagat", "06/04/2001", ""}, http.StatusBadRequest},
		{"invalid id", Author{-9, "Soni", "Raj", "01/10/1999", "sk"}, http.StatusBadRequest},
	}

	ConnectDB()

	for i, v := range test {

		// encoding
		jsonfile, err := json.Marshal(v.Auth)
		if err != nil {
			log.Printf("%v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "localhost:8000/author", bytes.NewReader(jsonfile))

		w := httptest.NewRecorder()

		PostAuth(w, req)

		status := w.Result().StatusCode

		if status != v.StatusCode {
			t.Errorf("TestCase[%v] = Desc : %v,Expected : %v, Actual : %v", i, v.Desc, v.StatusCode, status)
		}
	}

	CloseDB()

}

func TestPostBook(t *testing.T) {
	test := []struct {
		desc       string
		input      *Book
		statuscode int
	}{
		{"valid details", &Book{1, Author{1, "Chetan", "Bhagat", "06/04/2001", "Chetan"},
			"2 States", "Scholastic", "16/03/2016"}, http.StatusOK},
		{"author not exist", &Book{2, Author{3, "Mukesh", "Seth", "06/08/1990", "Seth"},
			"Journey", "Penguin", "11/12/2001"}, http.StatusBadRequest},
		{"invalid published date", &Book{3, Author{9, "Chetan", "Bhagat", "06/04/2001", "Chetan"},
			"Beauty of time", "Arihant", "12/12/1986"}, http.StatusBadRequest},
		{"invalid publication", &Book{4, Author{10, "Chetan", "Bhagat", "06/04/2001", "Chetan"},
			"Golden boy", "lenin", "19/02/2012"}, http.StatusBadRequest},
		{"empty Title", &Book{1, Author{1, "Chetan", "Bhagat", "06/04/2001", "Chetan"},
			"", "Scholastic", "16/03/2016"}, http.StatusBadRequest},
		{"empty Publication", &Book{1, Author{1, "Chetan", "Bhagat", "06/04/2001", "Chetan"},
			"2 States", "", "16/03/2016"}, http.StatusBadRequest},
		{"empty date", &Book{1, Author{1, "Chetan", "Bhagat", "06/04/2001", "Chetan"},
			"2 States", "Penguin", ""}, http.StatusBadRequest},
		{"invalid book id", &Book{-21, Author{1, "Chetan", "Bhagat", "06/04/2001", "Chetan"},
			"2 States", "Scholastic", ""}, http.StatusBadRequest},
	}

	ConnectDB()

	for _, v := range test {

		jsonfile, err := json.Marshal(v.input)
		if err != nil {
			log.Printf("err: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "localhost:8000/books", bytes.NewReader(jsonfile))

		w := httptest.NewRecorder()

		PostBook(w, req)

		status := w.Result().StatusCode

		if status != v.statuscode {
			t.Errorf("Desc : %v,Expected : %v, Actual : %v", v.desc, v.statuscode, status)
		}
	}
	CloseDB()

}

/*
func TestPutBook(t *testing.T) {
	test := []struct {
		desc       string
		inputid    string
		updatedata *Book
		statuscode int
	}{
		{"valid", "1", &Book{1, Author{1, "Chetan", "Bhagat", "14/05/1989", "Chetan"}, "400 Days", "Penguin", "04/11/2006"}, http.StatusOK},
		{"Invalid id", "-11", &Book{1, Author{1, "Chetan", "Bhagat", "14/05/1989", "Chetan"}, "400 Days", "Penguin", "04/11/2006"}, http.StatusBadRequest},
	}
}
*/
