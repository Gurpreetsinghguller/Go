package main

import (
	"fmt"
	"os"
)

//Using os package to get the environment variables which is already set
func envvr(key string) string {

	//set env variable using os package
	os.Setenv(key, "gopher")

	//return the env variable  using os package
	return os.Getenv(key)
}

func main() {
	//os package
	value := envvr("name")
	fmt.Printf("os package: %s = %s \n", "name", value)
}
