package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// type Contact struct {
// 	Id      string
// 	Name    string
// 	Email   string
// 	Subject string
// 	Message string
// }
type Contact struct {
	gorm.Model
	// ID      string `json:"Id"`
	Name    string `json:"Name"`
	Email   string `json:"Email"`
	Subject string `json:"Subject"`
	Message string `json:"Message"`
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// func dbConn() *sql.DB {
// 	db, err := sql.Open("mysql", "root:@/contact")
// 	checkError(err)
// 	return db
// }
func dbConn() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/contact?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	checkError(err)
	return db
}
func autoMigrate() {
	db := dbConn()
	db.AutoMigrate(&Contact{})
}
func createContact(w http.ResponseWriter, r *http.Request) {
	// data, err := ioutil.ReadAll(r.Body)
	DB := dbConn()
	var contact Contact
	json.NewDecoder(r.Body).Decode(&contact)
	DB.Create(&contact)
	// checkError(err)
	// var contact Contact
	// err = json.Unmarshal(data, &contact)
	// checkError(err)
	// db := dbConn()
	// defer db.Close()
	// qry := "INSERT INTO contact(name,email,subject,message) VALUES($1,$2,$3,$4)"
	// _, err = db.Exec(qry, contact.Name, contact.Email, contact.Subject, contact.Message)
	// checkError(err)
	// w.Write(data)
	json.NewEncoder(w).Encode(contact)
}
func getContacts(w http.ResponseWriter, r *http.Request) {
	var contacts []Contact
	DB := dbConn()
	DB.Find(&contacts)
	json.NewEncoder(w).Encode(contacts)
	// db := dbConn()
	// defer db.Close()
	// qry := "SELECT *FROM contact"
	// rows, err := db.Query(qry)
	// checkError(err)
	// var contacts []Contact
	// for rows.Next() {
	// 	var contact Contact
	// 	err = rows.Scan(&contact.Id, &contact.Name, &contact.Email, &contact.Subject, &contact.Message)
	// 	checkError(err)
	// 	contacts = append(contacts, contact)
	// }
	// pBytes, err := json.MarshalIndent(contacts, "", " ")
	// checkError(err)
	// w.Write(pBytes)
}
func getContact(w http.ResponseWriter, r *http.Request) {
	var contact Contact
	id := mux.Vars(r)
	DB := dbConn()
	DB.First(&contact, id["id"])
	json.NewEncoder(w).Encode(contact)
	// id := mux.Vars(r)
	// db := dbConn()
	// defer db.Close()
	// qry := "SELECT *FROM contact WHERE id=" + id["id"]
	// row := db.QueryRow(qry)
	// var contact Contact
	// row.Scan(&contact.Id, &contact.Name, &contact.Email, &contact.Subject, &contact.Message)
	// pBytes, err := json.MarshalIndent(contact, "", " ")
	// checkError(err)
	// w.Write(pBytes)
}

func updateContact(w http.ResponseWriter, r *http.Request) {
	var contact Contact
	id := mux.Vars(r)
	DB := dbConn()
	DB.First(&contact, id["id"]) //FIRST LAST & TAKE
	json.NewDecoder(r.Body).Decode(&contact)
	DB.Save(&contact)
	json.NewEncoder(w).Encode(contact)
	// id := mux.Vars(r)
	// var contact Contact
	// data, err := ioutil.ReadAll(r.Body)
	// err = json.Unmarshal(data, &contact)
	// db := dbConn()
	// defer db.Close()
	// qry := "UPDATE contact SET name = ? , email = ? , subject = ? , message = ? WHERE id = ?"
	// _, err = db.Exec(qry, contact.Name, contact.Email, contact.Subject, contact.Message, id["id"])
	// checkError(err)
	// w.Write(data)
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	var contact Contact
	id := mux.Vars(r)
	DB := dbConn()
	DB.Delete(&contact, id["id"])
	fmt.Println(contact)
	json.NewEncoder(w).Encode("The User is deleted successfully")
	// id := mux.Vars(r)
	// db := dbConn()
	// defer db.Close()
	// qry := "DELETE  FROM contact WHERE id=?"
	// db.Exec(qry, id["id"])
	// w.Write([]byte("SUCCESSFULLY DELETE RECORD WITH ID" + id["id"]))
}
func main() {
	autoMigrate()
	r := mux.NewRouter()
	r.HandleFunc("/api/contact", getContacts).Methods("GET")
	r.HandleFunc("/api/contact/{id}", getContact).Methods("GET")
	r.HandleFunc("/api/contact", createContact).Methods("POST")
	r.HandleFunc("/api/contact/{id}", updateContact).Methods("PUT")
	r.HandleFunc("/api/contact/{id}", deleteContact).Methods("DELETE")
	http.ListenAndServe(":8000", r)
}
