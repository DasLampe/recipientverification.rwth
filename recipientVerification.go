package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var tmpl = template.Must(template.ParseGlob("templates/*"))

type User struct {
	Username string
}

func List(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load() //Load .env file

	if err != nil {
		log.Fatal(err.Error())
	}

	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/mail_production", username, password, host))
	if err != nil {
		log.Fatal(err.Error())
	}

	rows, err := db.Query("SELECT username FROM people WHERE local = 1")
	if err != nil {
		log.Panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()

	user := User{}
	var res []User

	for rows.Next() {
		var username string

		err = rows.Scan(&username)
		if err != nil {
			panic(err.Error())
		}
		user.Username = username

		res = append(res, user)
	}

	err = tmpl.ExecuteTemplate(w, "List", res)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/adressenliste", List)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}