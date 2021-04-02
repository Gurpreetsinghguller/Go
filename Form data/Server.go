package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Contact struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// var tmpl =template.Must(template.ParseGlob("templates/*"))
func getTemplate() *template.Template {
	tmpl := template.Must(template.ParseGlob("templates/*"))
	return tmpl
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate()
	tmpl.ExecuteTemplate(w, "index.html", nil)

}

func contact(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate()
	tmpl.ExecuteTemplate(w, "contact.html", nil)

}

func success(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		read, _ := ioutil.ReadFile("Contact.Json")
		var contacts []Contact
		json.Unmarshal(read, &contacts)

		contact := Contact{
			Name:    r.FormValue("name"),
			Email:   r.FormValue("email"),
			Subject: r.FormValue("subject"),
			Message: r.FormValue("message")}

		contacts = append(contacts, contact)
		pBytes, _ := json.MarshalIndent(contacts, "", "  ")
		err1 := ioutil.WriteFile("Contact.Json", pBytes, 0644)
		if err1 != nil {
			fmt.Println(err1)
		}
		tmpl := getTemplate()
		tmpl.ExecuteTemplate(w, "success.html", nil)
	}
	http.Redirect(w, r, "/contact", 301)
}
func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/contact", contact)
	http.HandleFunc("/success", success)

	fmt.Println("Server has started")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

/*
fileServer := http.FileServer(http.Dir("./static")) //adds static folder
	http.Handle("/", fileServer)                        // I dont know why we added this line?
*/
