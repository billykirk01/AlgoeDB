package AlgoeDB

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"sync"
)

type Database struct {
	config    DatabaseConfig
	mu        sync.Mutex
	documents []map[string]interface{}
}

type DatabaseConfig struct {
	Path            string
	Pretty          *bool
	Autoload        *bool
	Immutable       *bool
	OnlyInMemory    *bool
	SchemaValidator *SchemaValidator
}

type SchemaValidator func(document interface{}) bool

type QueryFunc func(value interface{}) bool

func NewDatabase(config *DatabaseConfig) (*Database, error) {

	pretty := true
	autoload := true
	immutable := true
	onlyInMemory := false

	if config.Pretty == nil {
		config.Pretty = &pretty
	}

	if config.Autoload == nil {
		config.Autoload = &autoload
	}

	if config.Immutable == nil {
		config.Immutable = &immutable
	}

	if config.OnlyInMemory == nil {
		config.OnlyInMemory = &onlyInMemory
	}

	if config.Path == "" && *config.OnlyInMemory {
		return nil, errors.New("it is impossible to disable onlyInMemory mode if the path is not specified")
	}

	documents := []map[string]interface{}{}

	database := Database{
		documents: documents,
		config:    *config,
	}

	if config.Path != "" && *config.Autoload {
		database.load()
	}

	return &database, nil
}

func (d *Database) InsertOne(document map[string]interface{}) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.documents = append(d.documents, document)

	if !*d.config.OnlyInMemory {
		err := d.save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) InsertMany(documents []map[string]interface{}) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.documents = append(d.documents, documents...)

	if !*d.config.OnlyInMemory {
		err := d.save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) FindOne(query map[string]interface{}) (map[string]interface{}, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	found := searchDocuments(query, d.documents)

	if len(found) == 0 {
		return nil, errors.New("could not find any matching documents")
	}

	return d.documents[found[0]], nil
}

func (d *Database) FindMany(query map[string]interface{}) ([]map[string]interface{}, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	found := searchDocuments(query, d.documents)

	if len(found) == 0 {
		return nil, errors.New("could not find any matching documents")
	}

	results := []map[string]interface{}{}

	for index := range found {
		results = append(results, d.documents[index])
	}

	return results, nil
}

func (d *Database) UpdateOne(query map[string]interface{}, document map[string]interface{}) (map[string]interface{}, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	found := searchDocuments(query, d.documents)

	if len(found) == 0 {
		return nil, errors.New("could not find document to update")
	}

	temp := d.documents[found[0]]

	for key, value := range document {
		temp[key] = value
	}

	d.documents[found[0]] = temp

	return temp, nil
}

func (d *Database) load() error {

	content := "[]"

	if d.config.Path != "" {
		f, err := os.Open(d.config.Path)
		if err != nil {
			return errors.New("failed to open file: " + d.config.Path)
		}

		bytes, err := ioutil.ReadFile(d.config.Path)
		if err != nil {
			return errors.New("failed to read file: " + d.config.Path)
		}

		content = string(bytes)
		f.Close()
	}

	documents := []map[string]interface{}{}
	err := json.Unmarshal([]byte(content), &documents)
	if err != nil {
		return err
	}

	d.documents = documents

	return nil
}

func (d *Database) save() error {
	bytes, err := json.MarshalIndent(d.documents, "", "\t")
	if err != nil {
		return errors.New("failed to marshal JSON")
	}

	temp := d.config.Path + ".temp"
	f, err := os.Create(temp)
	if err != nil {
		return errors.New("failed to create temporary file: " + temp)
	}
	defer f.Close()

	err = ioutil.WriteFile(temp, bytes, 0644)
	if err != nil {
		return errors.New("failed to write data to temporary file: " + temp)
	}

	err = os.Rename(temp, d.config.Path)
	if err != nil {
		return errors.New("failed to rename temporary file: " + temp + " to: " + d.config.Path)
	}

	return nil
}

func searchDocuments(query map[string]interface{}, documents []map[string]interface{}) []int {

	found := []int{}

	for index, document := range documents {

		include := true

		for key, queryValue := range query {
			if !include {
				break
			}

			documentValue := document[key]

			if !matchValues(queryValue, documentValue) {
				include = false
			}
		}

		if include {
			found = append(found, index)
		}
	}

	return found
}

func matchValues(queryValue interface{}, documentValue interface{}) bool {

	if IsString(queryValue) || IsNumber(queryValue) || IsBoolean(queryValue) || IsNil(queryValue) {
		return queryValue == documentValue
	}

	if IsFunction(queryValue) {
		queryFunc := reflect.ValueOf(queryValue).Interface().(QueryFunc)
		return queryFunc(documentValue)
	}

	return false
}
