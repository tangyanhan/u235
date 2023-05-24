package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		log.Fatalf("Failed to connect sqlite database: %s", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := db.ExecContext(ctx, `CREATE TABLE users(name TEXT, password TEXT);`); err != nil {
		log.Fatalf("Failed to create table: %s", err)
	}
	ret, err := db.ExecContext(ctx, `INSERT INTO users(name, password) VALUES('ethan', 'changeme')`)
	if err != nil {
		log.Fatalf("Failed to insert records: %s", err)
	}
	nRows, err := ret.RowsAffected()
	if err != nil {
		log.Fatalf("Failed to get rows affected:%s", err)
	}
	log.Println("Affacted rows=", nRows)
	rows, err := db.QueryContext(ctx, `SELECT password from users WHERE name=?`, "ethan")
	if err != nil {
		log.Fatalf("Failed to query password: %s", err)
	}
	var password string
	if rows.Next() {
		if err := rows.Scan(&password); err != nil {
			log.Fatalf("Failed to scan password: %s", err)
		}
		log.Println("Password=", password)
	}
}
