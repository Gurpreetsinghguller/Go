package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func godot(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	return os.Getenv(key)
}
func main() {
	value := godot("STRONGEST_AVENGER")
	fmt.Println(value)
}
