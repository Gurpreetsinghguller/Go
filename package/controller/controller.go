package controller

import (
	"fmt"
	"net/http"
	"package/app"
	"package/data/models"
	"package/data/views"
)

func Contact(w http.ResponseWriter, r *http.Request) {
	tmpl := app.GetTemplate()
	tmpl.ExecuteTemplate(w, "index.html", nil)
}
func Success(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		views.Success(w, r)
	}
	http.Redirect(w, r, "/", 301)
}
func Show(w http.ResponseWriter, r *http.Request) {
	tmpl := app.GetTemplate()
	var store []models.Contact
	store = views.Show(w, r)
	tmpl.ExecuteTemplate(w, "show.html", store)
}
func Edit(w http.ResponseWriter, r *http.Request) {
	tmpl := app.GetTemplate()
	// var contact models.Contact
	contact := views.Edit(w, r)
	fmt.Println(contact)
	tmpl.ExecuteTemplate(w, "edit.html", contact)
}
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		views.Update(w, r)
	}
	http.Redirect(w, r, "/show", 301)
}
func Delete(w http.ResponseWriter, r *http.Request) {
	views.Delete(w, r)
}
