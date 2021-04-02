package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Image struct {
	ID       string
	Name     string
	Path     string
	FullPath string
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func getTemplate() *template.Template {
	tmpl := template.Must(template.ParseGlob("template/*"))
	return tmpl
}
func dbConn() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:1234@localhost/contact?sslmode=disable")
	checkError(err)
	return db
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate()
	if r.Method == "POST" {
		//Making db connection
		db := dbConn()
		defer db.Close()
		//Uploading file in range 10MB to 20Mb
		err := r.ParseMultipartForm(10 << 20)
		checkError(err)
		//Getting file from frontEnd
		file, header, err := r.FormFile("image")
		checkError(err)
		defer file.Close()
		//Storing temporary file in static folder
		dst, err := os.Create(filepath.Join("C:/Users/Gurpreet/Desktop/Golang -web/Image handling/static", filepath.Base(header.Filename)))
		// dst, err := os.Create("/static")
		checkError(err)
		defer dst.Close()
		//Copying filedata into that temporary file we created
		_, err = io.Copy(dst, file)
		checkError(err)
		//Inserting into Database

		qry := "INSERT INTO image (name,path) VALUES($1,$2)"
		filepath := "./static/" + header.Filename
		_, err = db.Exec(qry, header.Filename, filepath)
		checkError(err)
		http.Redirect(w, r, "/display", 301)
	}
	tmpl.ExecuteTemplate(w, "index.html", nil)
}
func save(w http.ResponseWriter, r *http.Request) {

}
func display(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate()
	db := dbConn()
	defer db.Close()
	qry := "SELECT *FROM image"
	rows, err := db.Query(qry)
	checkError(err)
	var img []Image
	// var getPath []string
	for rows.Next() {
		var image Image
		err = rows.Scan(&image.ID, &image.Name, &image.Path)
		fullPath := image.Path
		image.FullPath = fullPath
		// getPath = append(getPath, get)
		checkError(err)
		img = append(img, image)
	}
	// fmt.Println(img[0].FullPath)
	tmpl.ExecuteTemplate(w, "display.html", img)
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/", index).Methods("POST")
	r.HandleFunc("/display", display).Methods("GET")
	r.HandleFunc("/", index).Methods("GET")
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fs))
	http.Handle("/", r)
	fmt.Println("Server started at 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

//Some iterations i tried to get this work done

// tempfile, err := ioutil.TempFile("static", header.Filename+"*.png")
// defer tempfile.Close()

// fileBytes, err := ioutil.ReadAll(file)
// checkError(err)
// tempfile.Write(fileBytes)
