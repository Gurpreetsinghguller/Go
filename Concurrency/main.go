package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// var wg sync.WaitGroup
// var mt sync.Mutex
func sendRequest(url string) {
	defer wg.Done()
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	mt.Lock()
	defer mt.Unlock()
	fmt.Printf("[%d]  %s \n :-", res.StatusCode, url)

}
func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: ro run main.go <url> <url2>..<urln>")
	}
	for _, url := range os.Args[1:] {
		go sendRequest("https://" + url)
		wg.Add(1)
	}
	wg.Wait()
}

//5.795sec
