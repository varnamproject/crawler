package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func initDb() (chan string, chan struct{}) {
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		panic("Error while creating db")
	}
	createTable(db)
	wordsChannel := make(chan string, 100)
	done := make(chan struct{})
	go storeWords(db, wordsChannel, done)
	return wordsChannel, done
}

func createTable(db *sql.DB) {
	sqlStmt := `
	create table word_count (name text not null primary key, value integer);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q\n", err)
	}
}

func storeWords(db *sql.DB, wordsChannel chan string, done chan struct{}) {
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert or replace into word_count(name, value) values(?, COALESCE((SELECT value FROM word_count WHERE name=?),0) + 1)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for word := range wordsChannel {
		_, err = stmt.Exec(word, word)
		if err != nil {
			fmt.Println(err)
		}
	}
	tx.Commit()
	fmt.Println("DONE")
	done <- struct{}{}
}
