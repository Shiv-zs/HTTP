package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var db *sql.DB

type Author struct {
	AuthId    int    `json:"auth_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Dob       string `json:"dob"`
	PenName   string `json:"pen_name"`
}

type Book struct {
	BookId        int    `json:"book_id"`
	Auth          Author `json:"auth"`
	Title         string `json:"title"`
	Publication   string `json:"publication"`
	PublishedDate string `json:"published_date"`
}

func isValidPublishedDate(dob string) bool {

	split := strings.Split(dob, "/")
	yearInstr := split[2]

	yearInint, err := strconv.Atoi(yearInstr)

	if err != nil {
		log.Printf("Cannot convert dob in integer : %v", yearInint)
	}
	if yearInint < 2022 && yearInint > 1880 {
		return true
	}
	return false
}

func isValidPublication(pub string) bool {
	if pub == "Scholastic" || pub == "Arihant" || pub == "Penguin" {
		return true
	}
	return false
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("books endpoint hits")

}

func GetBookId(w http.ResponseWriter, r *http.Request) {
}

func PostAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post Author endpoint hits")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("body has nothing %v\n", err)
	}

	// To store author data
	var auth *Author

	// Decoding
	json.Unmarshal(body, &auth)

	//Checking All conditions for Valid Author
	if auth.AuthId <= 0 || auth.FirstName == "" || auth.LastName == "" || auth.Dob == "" || auth.PenName == "" {
		fmt.Printf("Fields are empty or author id is invalid of %v %v\n", auth.FirstName, auth.LastName)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Now Checking if author already present in DB(Author)
	if !CheckAuthor(auth) {
		fmt.Printf("Author is already exist wiht id : %v\n", auth.AuthId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Now insert the valid data into to DB
	_, err = db.Exec("insert into Author(authorId,firstName,lastName,dob,penName) values (?,?,?,?,?)", auth.AuthId, auth.FirstName, auth.LastName, auth.Dob, auth.PenName)
	if err != nil {
		log.Printf("Data is not able to insert in Author %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)

}

func PostBook(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Post Book endpoint hits")

	//reading body of request in json format
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("body has nothing %v\n", err)
	}

	// To store book data
	var book *Book

	// Decoding
	json.Unmarshal(body, &book)

	//Checking all condition for missing fields in book
	if book.BookId <= 0 || book.Title == "" || book.Publication == "" || book.PublishedDate == "" {
		fmt.Println("Fields are empty or author id is invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Checking Valid Published Date { 1880 < date < 2022 }
	if !isValidPublishedDate(book.PublishedDate) {
		fmt.Printf("Invalid Publication Date of %v %v\n", book.Auth.FirstName, book.Auth.LastName)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Checking for Valid Publications { Scholastic , Penguin , Arihant }
	if !isValidPublication(book.Publication) {
		fmt.Printf("Invalid Publication of book id : %v\n", book.BookId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Author Validation
	if !CheckAuthor(&book.Auth) {
		fmt.Printf("Author is Not Present in DB with Author id : %v\n", book.Auth.AuthId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Now insert the valid data into to DB
	_, err = db.Exec("insert into Book(bookId,title,Publication,PublishedDate) values (?,?,?,?)", book.BookId, book.Title, book.Publication, book.PublishedDate)
	if err != nil {
		log.Printf("Data is not able to insert in Book %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)

}

func main() {

	//http.HandleFunc("/author", PostAuth)
	//http.ListenAndServe(":8000", nil)
}
