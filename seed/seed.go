package seed

import "encoding/json"

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Address Address
}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type Uinterface interface {
}

var Employees = []User{
	{
		Name:    "Chakravarthy",
		Age:     "23",
		Contact: "9475757473",
		Address: Address{
			City:    "Hyderabad",
			State:   "Telangana",
			Country: "India",
			Pincode: "500038",
		},
	},
	{
		Name:    "Chakravarthy",
		Age:     "23",
		Contact: "9475757473",
		Address: Address{
			City:    "Hyderabad",
			State:   "Telangana",
			Country: "India",
			Pincode: "500038",
		},
	},
	{
		Name:    "Chakravarthy",
		Age:     "23",
		Contact: "9475757473",
		Address: Address{
			City:    "Hyderabad",
			State:   "Telangana",
			Country: "India",
			Pincode: "500038",
		},
	},
	{
		Name:    "Chakravarthy",
		Age:     "23",
		Contact: "9475757473",
		Address: Address{
			City:    "Hyderabad",
			State:   "Telangana",
			Country: "India",
			Pincode: "500038",
		},
	},
	{
		Name:    "Chakravarthy",
		Age:     "23",
		Contact: "9475757473",
		Address: Address{
			City:    "Hyderabad",
			State:   "Telangana",
			Country: "India",
			Pincode: "500038",
		},
	},
	{
		Name:    "Chakravarthy",
		Age:     "23",
		Contact: "9475757473",
		Address: Address{
			City:    "Hyderabad",
			State:   "Telangana",
			Country: "India",
			Pincode: "500038",
		},
	},
	{
		Name:    "Chakravarthy",
		Age:     "23",
		Contact: "9475757473",
		Address: Address{
			City:    "Hyderabad",
			State:   "Telangana",
			Country: "India",
			Pincode: "500038",
		},
	},
}
