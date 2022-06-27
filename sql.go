package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() {
	var err error
	db, err = sql.Open("mysql", "root:root@123@tcp(localhost:3306)/library")

	if err != nil {
		fmt.Println("Failed during Connection")
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Ping Giving error ")
		log.Fatal(err)
	}
	fmt.Println("Connection Established..!!")
}
func CloseDB() {
	fmt.Println("Connection Closed..!!")
	db.Close()
}

func CheckAuthor(auth *Author) bool {
	result := db.QueryRow("select * from Author where id=?", auth.AuthId)
	if result != nil {
		return true
	}
	return false
}
