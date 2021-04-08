package main

import (
	"Sessions/Controller"

	"fmt"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/login", Controller.Login)
	http.HandleFunc("/customer", Controller.Customer)
	http.HandleFunc("/logout", Controller.Logout)
	fmt.Println("Server Started at : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
