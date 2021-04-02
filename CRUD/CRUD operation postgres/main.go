package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

type Contact struct {
	Id      string
	Name    string
	Email   string
	Subject string
	Message string
}

func dbConn() (db *sql.DB) {
	var err error

	db, err = sql.Open("postgres", "postgres://postgres:1234@localhost/contact?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("data base connection -", db)
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)

	}
	return db
}
func getTemplate() *template.Template {
	tmpl := template.Must(template.ParseGlob("templates//*"))
	return tmpl
}
func contact(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate()
	tmpl.ExecuteTemplate(w, "index.html", nil)
}
func success(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db := dbConn()
		defer db.Close()
		contact := Contact{
			Name:    r.FormValue("name"),
			Email:   r.FormValue("email"),
			Subject: r.FormValue("subject"),
			Message: r.FormValue("message"),
		}
		qry := "INSERT INTO contact(name,email,subject,message) VALUES($1,$2,$3,$4)"
		_, e := db.Exec(qry, contact.Name, contact.Email, contact.Subject, contact.Message)
		if e != nil {
			log.Fatal(e)
		}
		fmt.Println("Values entered successfully")
		http.Redirect(w, r, "/show", 301)
	}
	http.Redirect(w, r, "/", 301)
}
func show(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate()
	db := dbConn()
	defer db.Close()
	rows, _ := db.Query("SELECT *FROM contact")
	var contact Contact
	var store []Contact
	for rows.Next() {
		err := rows.Scan(&contact.Id, &contact.Name, &contact.Email, &contact.Subject, &contact.Message)
		if err != nil {
			log.Fatal(err)
		}
		store = append(store, contact)
	}
	tmpl.ExecuteTemplate(w, "show.html", store)
}
func edit(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate()
	db := dbConn()
	defer db.Close()
	id := r.URL.Query().Get("id")
	qry := "SELECT *FROM contact WHERE id=" + id
	row := db.QueryRow(qry)
	var contact Contact
	err := row.Scan(&contact.Id, &contact.Name, &contact.Email, &contact.Subject, &contact.Message)
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "edit.html", contact)
}
func update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db := dbConn()
		defer db.Close()
		contact := Contact{
			Id:      r.FormValue("id"),
			Name:    r.FormValue("name"),
			Email:   r.FormValue("email"),
			Subject: r.FormValue("subject"),
			Message: r.FormValue("message"),
		}
		qry := "UPDATE contact SET name=$1,email=$2,subject=$3,message=$4 WHERE id=$5"
		_, err := db.Exec(qry, contact.Name, contact.Email, contact.Subject, contact.Message, contact.Id)
		if err != nil {
			log.Fatal(err)
		}
	}
	http.Redirect(w, r, "/show", 301)
}
func delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()
	id := r.URL.Query().Get("id")
	qry := "DELETE FROM contact WHERE id=$1"
	_, err := db.Exec(qry, id)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/show", 301)
}
func main() {
	http.HandleFunc("/", contact)
	http.HandleFunc("/success", success)
	http.HandleFunc("/show", show)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)
	fmt.Println("server started at 8000")
	http.ListenAndServe(":8000", nil)

}
