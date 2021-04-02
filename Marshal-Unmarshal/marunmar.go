package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Address `json:"address"`
}

type Tech struct {
	Id   int           `json:"id"`
	Tech []TechDetails `json:"tech"`
}

type Contact struct {
	Id             int `json:"id"`
	ContactDetails `json:"contactdetails"`
}

type Address struct {
	Area    string `json:"area"`
	Country string `json:"country"`
}

type TechDetails struct {
	Tech string  `json:"tech"`
	Exp  float64 `json:"experience"`
}
type ContactDetails struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}
type MergeUser struct {
	Id   int
	Name string
	Address
	Techdets []TechDetails
	Email    string
	Phone    string
}

func merge(user []User, tech []Tech, contact []Contact) {
	countryCode := map[string]string{
		"IND": "+91",
		"UK":  "+41",
	}
	for i := 0; i < len(user); i++ {
		if user[i].Country == "IND" {
			getPhone := contact[i].Phone
			contact[i].Phone = countryCode["IND"] + getPhone
		} else {
			getPhone := contact[i].Phone
			contact[i].Phone = countryCode["UK"] + getPhone
		}

	}
	store := make([]MergeUser, len(user))
	for index, user := range user {

		mergeStructVariable := MergeUser{}

		//Entering values of user

		mergeStructVariable.Id = user.Id
		mergeStructVariable.Name = user.Name
		mergeStructVariable.Address = user.Address

		for _, value := range tech {
			if value.Id == user.Id {
				mergeStructVariable.Techdets = value.Tech
			}
		}

		for _, value := range contact {
			if value.Id == user.Id {
				mergeStructVariable.Email = value.Email
				mergeStructVariable.Phone = value.Phone
			}
		}
		store[index] = mergeStructVariable
	}
	// fmt.Println(store)
	// pBytes, _ := json.Marshal(store)
	pBytes, _ := json.MarshalIndent(store, "", "  ")
	err7 := ioutil.WriteFile("Merge.Json", pBytes, 0644)
	if err7 != nil {
		fmt.Println(err7)
	}
}
func main() {

	//Unmarshalling

	//User
	user1, err := ioutil.ReadFile("User.Json")
	if err != nil {
		fmt.Println(err.Error())
	}
	var users []User
	err1 := json.Unmarshal(user1, &users)
	if err1 != nil {
		fmt.Println(err1.Error())
	}

	//Technology
	tech1, err2 := ioutil.ReadFile("Tech.Json")
	if err2 != nil {
		fmt.Println(err2.Error())
	}
	var tech []Tech
	err3 := json.Unmarshal(tech1, &tech)
	if err3 != nil {
		fmt.Println(err3.Error())
	}

	//Contact
	contact1, err4 := ioutil.ReadFile("Contact.Json")
	if err4 != nil {
		fmt.Println(err4.Error())
	}
	var contact []Contact
	err5 := json.Unmarshal(contact1, &contact)
	if err1 != nil {
		fmt.Println(err5.Error())
	}

	//Calling Merge Function
	merge(users, tech, contact)

}

//Marshalling
/*
Add JSON tug to structure entries

p := Person{name Gurpreet:}

pBytes,err := Json.Marshal(p)
lo.print(err)
log.print(pBytes)
*/

//UnMarshalling

/*
err2 := json.Unmarshal([]byte(usercontent), &users)
*/

//Variour Tugs are

/*
`json:"customName"`
`json:"age,omitempty"`            omits entry if empty
`json:"-"`                        doesnot enters Secret entry like creditcard details
*/
