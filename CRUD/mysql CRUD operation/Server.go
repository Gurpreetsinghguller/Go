package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Contact struct {
	Id      string
	Name    string
	Email   string
	Subject string
	Message string
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:@/contact")
	if err != nil {
		panic(err)
	}
	return db
}
func getTemplate() *template.Template {
	tmpl := template.Must(template.ParseGlob("templates/*"))
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
			Message: r.FormValue("message")}

		insForm, err := db.Prepare("INSERT INTO contact (name,email,subject,messsage) VALUES(?,?,?,?)")
		if err != nil {
			panic(err)
		}

		_, err1 := insForm.Exec(contact.Name, contact.Email, contact.Subject, contact.Message)
		if err1 != nil {
			fmt.Println(err1)

		}

		http.Redirect(w, r, "/show", 301)
	}
	http.Redirect(w, r, "/", 301)
}
func show(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate()
	db := dbConn()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM contact")
	if err != nil {
		panic(err)
	}

	var store []Contact
	for rows.Next() {
		var contact Contact
		err = rows.Scan(&contact.Id, &contact.Name, &contact.Email, &contact.Subject, &contact.Message)
		if err != nil {
			panic(err.Error())
		}
		store = append(store, contact)
	}
	tmpl.ExecuteTemplate(w, "show.html", store)
}
func edit(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate()

	id := r.URL.Query().Get("id")
	db := dbConn()
	fmt.Println(id)
	row := db.QueryRow("SELECT *FROM contact WHERE id=" + id)
	var contact Contact
	err := row.Scan(&contact.Id, &contact.Name, &contact.Email, &contact.Subject, &contact.Message)
	fmt.Println(contact.Id, contact.Name, contact.Email)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	tmpl.ExecuteTemplate(w, "edit.html", contact)
}
func update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update called")
	if r.Method == "POST" {
		db := dbConn()
		contact := Contact{
			Id:      r.FormValue("id"),
			Name:    r.FormValue("name"),
			Email:   r.FormValue("email"),
			Subject: r.FormValue("subject"),
			Message: r.FormValue("message")}
		updateRow, err := db.Prepare("UPDATE contact SET name=?,email=?,subject=?,messsage=? WHERE id=?")
		if err != nil {
			panic(err)
		}

		_, err1 := updateRow.Exec(contact.Name, contact.Email, contact.Subject, contact.Message, contact.Id)
		if err1 != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/show", 301)
	}
	http.Redirect(w, r, "/", 301)
}
func delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	db := dbConn()
	defer db.Close()
	qry := "DELETE FROM contact WHERE id=?"
	// del, err := db.Prepare(qry)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	res, err := db.Exec(qry, id)
	if err != nil {
		log.Fatal(err)
	}
	count, err := res.RowsAffected()
	if count == 0 {
		http.Redirect(w, r, "/show", 301)
	} else {
		http.Redirect(w, r, "/show", 301)
	}
	if err != nil {
		log.Fatal(err)

	}
	http.Redirect(w, r, "/show", 301)
}
func main() {
	http.HandleFunc("/", contact)
	http.HandleFunc("/show", show)
	http.HandleFunc("/success", success)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)
	fmt.Println("Server started at 8525")
	http.ListenAndServe(":8525", nil)
}
