package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {

	c := make(chan string)
	// c := make(chan string,2)
	c <- "Hello"
	// c <- "World"
	msg := <-c
	fmt.Println(msg)
	// msg = <-c
	// fmt.Println(msg)
	// defer close(c)
}
