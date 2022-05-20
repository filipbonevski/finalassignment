package main

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite", filepath)

	if err != nil {
		panic(err)
	}

	if db == nil {
		panic("db nil")
	}
	return db
}

func migrate(db *sql.DB) {
	sql := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		username VARCHAR NOT NULL UNIQUE,
		password VARCHAR NOT NULL
	);
    CREATE TABLE IF NOT EXISTS lists(
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name VARCHAR NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
    );
	CREATE TABLE IF NOT EXISTS tasks(
		id INTEGER NOT NULL,
		name VARCHAR NOT NULL,
		list_id INTEGER NOT NULL,
		completed INTEGER,
		PRIMARY KEY (id),
		FOREIGN KEY(list_id) REFERENCES lists(id)
	);
    `
	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
}

func CreateUser(db *sql.DB, username, password string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}

	db.Exec("INSERT or IGNORE INTO users(username, password) VALUES(?,?)", username, string(hashedPassword))
}
