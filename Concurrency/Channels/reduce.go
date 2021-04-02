package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup
var mt sync.Mutex

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func apiCall(i int, ch chan string) {
	defer wg.Done()

	ch <- "API call for" + strconv.Itoa(i) + "started"
	time.Sleep(100 * time.Millisecond)

}

func main() {
	ch := make(chan string, 1000)
	numArray := makeRange(0, 1000)

	for i, _ := range numArray {
		wg.Add(1)
		go apiCall(i, ch)

	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for msg := range ch {

		fmt.Println(msg)
	}

}
