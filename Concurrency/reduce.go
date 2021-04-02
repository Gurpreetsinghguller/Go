package main

import (
	"fmt"
	"log"
	"time"
)

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func apiCall(i int, ch chan string) {
	// defer wg.Done()
	ch <- "API call for" + i + "started"
	time.Sleep(100 * time.Millisecond)
	close(ch)
}

func main() {
	ch := make(chan string)
	numArray := makeRange(0, 1000)

	start := time.Now()

	for i, _ := range numArray {

		go apiCall(i, ch)
	}
	for msg := range ch {
		fmt.Println(msg)
	}
	elapsed := time.Since(start)
	log.Printf("Time taken %s", elapsed)
}
