package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"text/template"

	_ "embed"
	_ "github.com/go-sql-driver/mysql"
)

//go:embed templates/List.tmpl
var tmpl string

type User struct {
	Username string
}

func getUsers(dbHost, dbUser, dbPass string) ([]User, error) {
	var users []User

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/mail_production", dbUser, dbPass, dbHost))
	if err != nil {
		return users, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT username FROM people WHERE local = 1 UNION SELECT name FROM role_accounts")
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var username string

		err = rows.Scan(&username)
		if err != nil {
			return users, err
		}
		users = append(users, User{Username: username})
	}

	return users, nil
}

func handleAdressenliste(dbHost, dbUsername, dbPassword string, getUserFunc func(dbHost, dbUsername, dbPassword string) ([]User, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request to verify adressen liste")

		users, err := getUserFunc(dbHost, dbUsername, dbPassword)
		if err != nil {
			slog.Error("Error getting users:", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t, err := template.New("").Parse(tmpl)
		if err != nil {
			slog.Error("Error parsing template:", "err", err, "users", users)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := t.Execute(w, users); err != nil {
			slog.Error("Error executing template:", "err", err, "users", users)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	slog.Info("Starting recipient verification server")
	http.HandleFunc("/adressenliste", handleAdressenliste(host, username, password, getUsers))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("Error starting http server", "err", err)
	}
	return
}
