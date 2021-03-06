package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	c := make(chan string)
	go count("Hello", c)
	for msg := range c {
		fmt.Println(msg)
	}

}
func count(thing string, c chan string) {
	for i := 0; i < 5; i++ {
		// fmt.Println(thing)
		c <- thing
		time.Sleep(time.Millisecond * 500)
	}
	close(c)
}
