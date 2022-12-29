package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

func main() {
	db, err := sql.Open("mysql", "root:root123@(127.0.0.1:3306)/go_web_examples?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	{
		query := `
			CREATE TABLE IF NOT EXISTS users (
				id INT AUTO_INCREMENT,
				username TEXT NOT NULL,
				password TEXT NOT NULL,
				created_at DATETIME,
				PRIMARY KEY (id)
			);
		`
		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
		}
	}

	{
		username := "johndoe"
		password := "secret"
		createdAt := time.Now()

		result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
		if err != nil {
			log.Fatal(err)
		}

		id, err := result.LastInsertId()

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id)
	}

	{
		var (
			id        int
			username  string
			password  string
			createdAt time.Time
		)
		query := "SELECT id, username, password, created_at FROM users WHERE id = ?"
		if err := db.QueryRow(query, 1).Scan(&id, &username, &password, &createdAt); err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, username, password, createdAt)
	}

	{
		type user struct {
			id        int
			username  string
			password  string
			createdAt time.Time
		}

		rows, err := db.Query(`SELECT id, username, password, created_at FROM users`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var users []user
		for rows.Next() {
			var u user
			err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, u)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%#v", users)
	}
	{
		_, err := db.Exec(`DELETE FROM users WHERE id = ?`, 1)
		if err != nil {
			log.Fatal(err)
		}
	}

	r := mux.NewRouter()

	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]
		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	http.ListenAndServe(":80", r)
}
