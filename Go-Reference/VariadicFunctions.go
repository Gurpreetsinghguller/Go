package main

import (
	"fmt"
)

func main() {
	sum("Hello", 2)
	sum("Hie", 2, 2)
	summation(2, 2, 4)
}
func sum(s string, nums ...int) {
	fmt.Println(s)
	total := 0
	for _, value := range nums {
		total += value
	}
	fmt.Println("sum :", total)
}
func summation(nums ...int) {

	total := 0
	for _, value := range nums {
		total += value
	}
	fmt.Println("Summation: ", total)
	fmt.Println("Hello")
}
