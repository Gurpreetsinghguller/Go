package Controller

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func GetTemplate() *template.Template {
	tmpl := template.Must(template.ParseGlob("template/*"))
	return tmpl
}

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl := GetTemplate()
	if r.Method == "POST" {
		email := r.FormValue("email")
		pass := r.FormValue("pass")
		fmt.Println(email, pass)
		if email == "gurpreet@gmail.com" && pass == "78788866" {
			session, _ := store.Get(r, "carbazar")
			session.Values["authenticated"] = true
			session.Values["username"] = "Gurpreet"
			//email[0:10]
			session.Save(r, w)
			http.Redirect(w, r, "/customer", http.StatusSeeOther)
		} else {
			data := map[string]interface{}{
				"err": "Invalid",
			}
			tmpl.ExecuteTemplate(w, "signin.html", data)
			return
		}
	}
	fmt.Println("Get Method called")
	tmpl.ExecuteTemplate(w, "signin.html", nil)
}

func Customer(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "carbazar")
	username := session.Values["username"]
	data := map[string]interface{}{
		"username": username,
	}
	tmpl := GetTemplate()
	tmpl.ExecuteTemplate(w, "index.html", data)
}
func Logout(w http.ResponseWriter, r *http.Request) {
	tmpl := GetTemplate()
	session, _ := store.Get(r, "carbazar")
	session.Values["authenticated"] = false
	session.Save(r, w)
	tmpl.ExecuteTemplate(w, "signin.html", nil)
}
