package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: ro run main.go <url> <url2>..<urln>")
	}
	for _, url := range os.Args[1:] {
		wg.Add(1)
		go Get("https://" + url)
	}
	wg.Wait()
}
func Get(url string) {
	defer wg.Done()
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d %s \n ", res.StatusCode, url)
}
