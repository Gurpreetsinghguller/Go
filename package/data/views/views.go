package views

import (
	"net/http"
	"package/app"
	"package/data/models"
)

func Success(w http.ResponseWriter, r *http.Request) {

	db := app.DbConn()
	defer db.Close()
	contact := models.Contact{
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Subject: r.FormValue("subject"),
		Message: r.FormValue("message"),
	}
	qry := "INSERT INTO contact(name,email,subject,message) VALUES($1,$2,$3,$4)"
	_, e := db.Exec(qry, contact.Name, contact.Email, contact.Subject, contact.Message)
	app.CheckError(e)

	http.Redirect(w, r, "/show", 301)

}
func Show(w http.ResponseWriter, r *http.Request) []models.Contact {

	db := app.DbConn()
	defer db.Close()
	rows, _ := db.Query("SELECT *FROM contact")
	var contact models.Contact
	var store []models.Contact
	for rows.Next() {
		err := rows.Scan(&contact.Id, &contact.Name, &contact.Email, &contact.Subject, &contact.Message)
		app.CheckError(err)
		store = append(store, contact)
	}
	return store
}
func Edit(w http.ResponseWriter, r *http.Request) models.Contact {

	db := app.DbConn()
	defer db.Close()
	id := r.URL.Query().Get("id")
	qry := "SELECT *FROM contact WHERE id=" + id
	row := db.QueryRow(qry)
	var contact models.Contact
	err := row.Scan(&contact.Id, &contact.Name, &contact.Email, &contact.Subject, &contact.Message)
	app.CheckError(err)
	// fmt.Println(contact)
	return contact
}
func Update(w http.ResponseWriter, r *http.Request) {

	db := app.DbConn()
	defer db.Close()
	contact := models.Contact{
		Id:      r.FormValue("id"),
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Subject: r.FormValue("subject"),
		Message: r.FormValue("message"),
	}
	qry := "UPDATE contact SET name=$1,email=$2,subject=$3,message=$4 WHERE id=$5"
	_, err := db.Exec(qry, contact.Name, contact.Email, contact.Subject, contact.Message, contact.Id)
	app.CheckError(err)

}
func Delete(w http.ResponseWriter, r *http.Request) {
	db := app.DbConn()
	defer db.Close()
	id := r.URL.Query().Get("id")
	qry := "DELETE FROM contact WHERE id=$1"
	_, err := db.Exec(qry, id)
	app.CheckError(err)
	http.Redirect(w, r, "/show", 301)
}
