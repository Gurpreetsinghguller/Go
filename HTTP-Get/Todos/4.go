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
		UserID    int       `json:"user_id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
}

func getData(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://gorest.co.in/public-api/todos") //resp contains header and body
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body) //We will be needing body only
	if err != nil {
		panic(err)
	}

	var todos AutoGenerated
	json.Unmarshal(data, &todos)
	pBytes, _ := json.MarshalIndent(todos, " ", "  ")
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