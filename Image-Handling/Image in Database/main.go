package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Store struct {
	ID   string
	Name string
	File []byte
	Img  string
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func dbConn() *sql.DB {
	db, err := sql.Open("postgres", "user=postgres password=1234 dbname=contact sslmode=disable")
	checkError(err)
	return db
}
func getTemplate() *template.Template {
	tmpl := template.Must(template.ParseGlob("template/*"))
	return tmpl
}
func index(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate()
	if r.Method == "POST" {

		file, header, err := r.FormFile("image")
		checkError(err)

		fileBytes, err := ioutil.ReadAll(file)
		checkError(err)

		db := dbConn()
		defer db.Close()
		qry := "INSERT INTO store (name, image) VALUES($1,$2)"
		_, err = db.Exec(qry, header.Filename, fileBytes)
		checkError(err)
		http.Redirect(w, r, "/display", 301)
	}
	tmpl.ExecuteTemplate(w, "index.html", nil)
}
func show(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate()
	db := dbConn()
	defer db.Close()
	qry := "SELECT *FROM store"
	rows, err := db.Query(qry)
	checkError(err)
	var store []Store
	for rows.Next() {
		var img Store
		err := rows.Scan(&img.ID, &img.Name, &img.File)

		base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(img.File)))
		base64.StdEncoding.Encode(base64Text, img.File)

		img.Img = string(base64Text)

		checkError(err)
		store = append(store, img)
	}

	tmpl.ExecuteTemplate(w, "show.html", store)
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/", index).Methods("POST")
	r.HandleFunc("/show", show).Methods("GET")
	fmt.Println("Server started at 8000")
	log.Fatal(http.ListenAndServe(":8588", r))
}

//Iterations i tried to get this work done
// imgFile, _, _ := image.Decode(bytes.NewReader(img.File))
// imgFile := base64.StdEncoding.EncodeToString(img.File)
// images = append(images, imagestring)
// l, _ := base64.StdEncoding.Decode(base64Text, img.File)
