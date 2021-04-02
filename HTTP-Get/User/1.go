package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type AutoGenerated struct {
	Code int `json:"code"`
	Meta struct {
		Pagination struct {
			Total int `json:"total"`
			Pages int `json:"pages"`
			Page  int `json:"page"`
			Limit int `json:"limit"`
		} `json:"pagination"`
	} `json:"meta"`
	Data []struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Gender    string    `json:"gender"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
}

func getData(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://gorest.co.in/public-api/users") //resp contains header and body
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body) //We will be needing body only
	if err != nil {
		panic(err)
	}

	var userData AutoGenerated
	json.Unmarshal(data, &userData)
	pBytes, _ := json.MarshalIndent(userData, " ", "  ")
	err1 := ioutil.WriteFile("data.json", pBytes, 0644)
	if err1 != nil {
		fmt.Println(err1)
	}

}
func main() {
	http.HandleFunc("/", getData)
	fmt.Println("Server started")
	http.ListenAndServe(":8080", nil)
}
