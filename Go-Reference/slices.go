package main

func main() {
	// Creating a slice in different ways

	// Slice Literal syntax
	// var slice []int = []int{1, 2, 3, 4, 5} //Length will be 5 and capacity will be 5
	// fmt.Println(slice, len(slice), cap(slice))
	// slice = append(slice, 6)
	// fmt.Println(slice, len(slice), cap(slice)) //Length will become 6 and capacity will become 10

	//Using make function
	// s := make([]int, 5)            //length will be 5 and capacity will be 5
	// fmt.Println(s, len(s), cap(s)) // will print initial values that is [0,0,0,0,0]
	// s = append(s, 6)               // length will become 6 and capacity will become 10
	// fmt.Println(s, len(s), cap(s)) //Will print [0,0,0,0,0,6]

	//Another literal syntax
	// s := []int{1, 2, 3, 4, 5}

}
