package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const VERSION = "1.0.1"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Debug(string, ...interface{})
		Info(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex     sync.Mutex
		mutexes   map[string]*sync.Mutex
		directory string
		log       Logger
	}

	Options struct {
		Logger
	}
)

func New(directory string, options *Options) (*Driver, error) {

	directory = filepath.Clean(directory)

	opts := Options{}

	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := Driver{
		directory: directory,
		mutexes:   make(map[string]*sync.Mutex),
		log:       opts.Logger,
	}

	if _, err := os.Stat(directory); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists) \n", directory)
		return &driver, nil
	}

	opts.Logger.Debug("Creating the database at '%s' ...\n", directory)

	return &driver, os.MkdirAll(directory, 0755)

}

func stat(path string) (file os.FileInfo, err error) {
	if file, err = os.Stat(path); os.IsNotExist(err) {
		file, err = os.Stat(path + ".json")
	}
	return
}

func (d *Driver) Write(collection string, filename string, v interface{}) error {

	if collection == "" {
		return fmt.Errorf("Missing the collection name")
	}

	if filename == "" {
		return fmt.Errorf("Missing the filename")
	}

	mutex := d.GetOrCreateMutex(collection)

	mutex.Lock()

	defer mutex.Unlock()

	dir := filepath.Join(d.directory, collection)
	fnlPath := filepath.Join(dir, filename+".json")
	temporayPath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")

	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(temporayPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(temporayPath, fnlPath)

}

func (d *Driver) Read(collection string, resource string, v interface{}) error {

	if collection == "" {
		fmt.Errorf("Missing Collection name : Please enter the collection name")
	}

	if resource == "" {
		fmt.Errorf("Missing resource")
	}

	record := filepath.Join(d.directory, collection, resource)

	if _, err := stat(record); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(record + ".json")

	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)
}

func (d *Driver) ReadAll(collection string) ([]string, error) {

	if collection == "" {
		err := fmt.Errorf("Invalid Collection:")
		return nil, err
	}

	dir := filepath.Join(d.directory, collection+".json")

	if _, err := stat(dir); err != nil {
		return nil, err
	}

	files, _ := ioutil.ReadDir(dir)

	var records []string
	for _, file := range files {
		buffer, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))

		if err != nil {
			return nil, err
		}
		records = append(records, string(buffer))
	}

	return records, nil

}

func (d *Driver) Delete(collection string, filename string, v interface{}) error {

	pathToTheCollection := filepath.Join(collection, filename)
	mutex := d.GetOrCreateMutex(collection)

	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.directory, pathToTheCollection)

	switch fi, err := stat(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("Unable to the file or Directory '%v'\n .", pathToTheCollection)

	case fi.Mode().IsDir():
		return os.RemoveAll(dir)

	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	default:
		return nil
	}
}

func (d *Driver) GetOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

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

func main() {
	dir := "./"

	db, err := New(dir, nil)

	if err != nil {
		fmt.Println("Error", err)
	}

	employees := []User{
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

	for _, value := range employees {
		//ranging over the values given and inserting into the database.
		db.Write("users", value.Name, User{
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

	allUsers := []User{}

	for _, user := range users {
		employeeFound := User{}

		//unmarshiling the json data to golang understandlable language.
		if err := json.Unmarshal([]byte(user), &employeeFound); err != nil {
			fmt.Println("Error", err)
		}

		allUsers = append(allUsers, employeeFound)
	}

	fmt.Println("All users", allUsers)

	// if err := db.Delete("user", "john"); err != nil {
	// 	fmt.Printf("Error", err)
	// }

	// if err := db.Delete("user", ""); err != nil {
	// 	fmt.Printf("Error", err)
	// }

}
