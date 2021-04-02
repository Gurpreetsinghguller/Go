package main

import "fmt"

func main() {
	fmt.Println("Hello")
	nextInt := Intsequence()
	fmt.Println(nextInt)
}
func Intsequence() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}
