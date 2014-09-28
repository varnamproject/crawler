package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func initDb() *sql.DB {
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		panic("Error while creating db")
	}
	createTable(db)
	return db

}

func createTable(db *sql.DB) {
	sqlStmt := `
	create table word_count (name text not null primary key, value integer);
	`
	db.Exec(sqlStmt)
}

func wordCollector(db *sql.DB) (chan string, chan struct{}) {
	wordsChannel := make(chan string, 100)
	done := make(chan struct{})
	go storeWords(db, wordsChannel, done)
	return wordsChannel, done
}

func storeWords(db *sql.DB, wordsChannel chan string, done chan struct{}) {
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

func generateVarnamFiles(db *sql.DB) {
	rows, err := db.Query("select * from word_count order by value desc")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fo, err := os.Create("output.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer fo.Close()
	w := bufio.NewWriter(fo)
	for rows.Next() {
		var value uint64
		var word string
		rows.Scan(&word, &value)
		w.WriteString(fmt.Sprintf("%s %d\n", word, value))
	}
	w.Flush()
}
