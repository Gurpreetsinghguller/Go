package main

import "fmt"

func main() {
	m := map[string]string{"Name": "Gurpreet"}
	m1 := map[string]string{"Surame": "Guller"}
	m2 := map[string]string{"Fathername": "Surjit"}
	m3 := make(map[string][]map[string]string)
	m3["Name"] = []map[string]string{m, m1, m2}
	for _, value := range m3 {
		for _, indexValue := range value {
			for Key, value1 := range indexValue {
				fmt.Println(Key, value1)
			}
		}
	}

}
