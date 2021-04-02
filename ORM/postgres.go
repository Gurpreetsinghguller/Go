package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Contact struct {
	Id      string
	Name    string
	Email   string
	Subject string
	Message string
}

func dbConn() *sql.DB {
	db, _ := sql.Open("postgres", "postgres://postgres:1234@localhost/contact?sslmode=disable")
	return db
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func saveData(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()
	body, err := ioutil.ReadAll(r.Body)
	checkError(err)
	var contact Contact
	err1 := json.Unmarshal(body, &contact)
	checkError(err1)
	qry := "INSERT INTO contact(name,email,subject,message) VALUES($1,$2,$3,$4)"
	_, err3 := db.Exec(qry, contact.Name, contact.Email, contact.Subject, contact.Message)
	checkError(err3)
	w.Write(body)
}
func getContacts(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()
	qry := "SELECT *FROM contact"
	rows, err := db.Query(qry)
	checkError(err)
	var contact Contact
	var contacts []Contact
	for rows.Next() {
		err := rows.Scan(&contact.Id, &contact.Name, &contact.Email, &contact.Subject, &contact.Message)
		checkError(err)
		contacts = append(contacts, contact)

	}
	pBytes, err := json.MarshalIndent(contacts, "", " ")
	w.Write(pBytes)
}
func getContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	db := dbConn()
	defer db.Close()
	qry := "SELECT *FROM contact WHERE id=" + vars["id"]
	row := db.QueryRow(qry)
	var contact Contact
	err := row.Scan(&contact.Id, &contact.Name, &contact.Email, &contact.Subject, &contact.Message)
	flag := true
	if sql.ErrNoRows == err {
		flag = false
	} else {
		checkError(err)
	}
	pBytes, err := json.MarshalIndent(contact, "", " ")
	checkError(err)
	if flag {
		w.Write([]byte(pBytes))
	} else {
		w.Write([]byte("No rows found"))
	}

}

func updateContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	db := dbConn()
	defer db.Close()
	body, err := ioutil.ReadAll(r.Body)
	checkError(err)
	var contact Contact
	json.Unmarshal(body, &contact)
	qry := "UPDATE contact SET name = $1 , email = $2 , subject = $3 , message = $4 WHERE id = $5"

	_, err1 := db.Exec(qry, contact.Name, contact.Email, contact.Subject, contact.Message, vars["id"])
	checkError(err1)
	w.Write(body)
}
func deleteContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	db := dbConn()
	defer db.Close()
	qry := "DELETE FROM contact WHERE id=$1"
	_, err := db.Exec(qry, vars["id"])
	checkError(err)
	w.Write([]byte("DELETE REQUEST SUCCESS FOR " + vars["id"]))
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/Contact", saveData).Methods("POST")
	r.HandleFunc("/api/Contact", getContacts).Methods("GET")
	r.HandleFunc("/api/Contact/{id}", getContact).Methods("GET")
	r.HandleFunc("/api/Contact/{id}", updateContact).Methods("PUT")
	r.HandleFunc("/api/Contact/{id}", deleteContact).Methods("DELETE")
	http.ListenAndServe(":8000", r)
}
