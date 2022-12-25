package main

import (
	"encoding/json"
	"fmt"

	"github.com/chakravarthy712/antdb/driver"
	"github.com/chakravarthy712/antdb/seed"
)

func main() {
	dir := "./"

	db, err := driver.New(dir, nil)

	if err != nil {
		fmt.Println("Error", err)
	}

	for _, value := range seed.Employees {
		//ranging over the values given and inserting into the database.
		db.Write("users", value.Name, seed.User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Address: value.Address,
		})
	}

	users, err := db.ReadAll("users")

	if err != nil {
		fmt.Println("Error", err)
	}

	fmt.Println(users)

	allUsers := []seed.User{}

	for _, user := range users {
		employeeFound := seed.User{}

		//unmarshiling the json data to golang understandlable language.
		if err := json.Unmarshal([]byte(user), &employeeFound); err != nil {
			fmt.Println("Error", err)
		}

		allUsers = append(allUsers, employeeFound)
	}

	fmt.Println("All users", allUsers)

}
