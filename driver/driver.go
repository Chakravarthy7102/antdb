package driver

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

func (d *Driver) Write(collection string, filename string, v interface{}) error {

	if collection == "" {
		return fmt.Errorf("missing the collection name")
	}

	if filename == "" {
		return fmt.Errorf("missing the filename")
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
		fmt.Errorf("missing collection name : Please enter the collection name")
	}

	if resource == "" {
		fmt.Errorf("missing resource")
	}

	record := filepath.Join(d.directory, collection, resource)

	if _, err := utils.stat(record); err != nil {
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

	if _, err := utils.stat(dir); err != nil {
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

	switch fi, err := utils.stat(dir); {
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
