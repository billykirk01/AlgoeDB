package AlgoeDB

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sync"
)

type Database struct {
	config    DatabaseConfig
	mu        sync.Mutex
	documents []map[string]interface{}
	writer    Writer
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
		return nil, errors.New("It is impossible to disable onlyInMemory mode if the path is not specified")
	}

	documents := []map[string]interface{}{}

	writer := Writer{Path: config.Path}

	database := Database{
		documents: documents,
		writer:    writer,
		config:    *config,
	}

	if config.Path != "" && *config.Autoload {
		database.load()
	}

	return &database, nil
}

func (d *Database) InsertOne(document map[string]interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.documents = append(d.documents, document)

	if !*d.config.OnlyInMemory {
		d.save()
	}
}

func (d *Database) InsertMany(documents []map[string]interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, document := range documents {
		d.documents = append(d.documents, document)
	}

	if !*d.config.OnlyInMemory {
		d.save()
	}
}

func (d *Database) FindOne(query map[string]interface{}) map[string]interface{} {
	d.mu.Lock()
	defer d.mu.Unlock()

	found := searchDocuments(query, d.documents)

	if len(found) > 0 {
		return d.documents[found[0]]
	}

	return nil
}

func (d *Database) load() {

	content := "[]"

	if d.config.Path != "" {
		f, err := os.Open(d.config.Path)
		if err != nil {
			log.Fatal(err)
		}

		bytes, err := ioutil.ReadFile(d.config.Path)
		content = string(bytes)
		f.Close()
	}

	documents, err := parseDatabaseStorage(content)
	if err != nil {
		log.Fatal(err)
	}

	d.documents = documents
}

func (d *Database) save() {
	bytes, err := json.MarshalIndent(d.documents, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	d.writer.Write(string(bytes))
}

func (d *Database) drop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.documents = []map[string]interface{}{}

	if *d.config.OnlyInMemory {
		d.save()
	}
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

func parseDatabaseStorage(content string) ([]map[string]interface{}, error) {
	documents := []map[string]interface{}{}
	err := json.Unmarshal([]byte(content), &documents)
	return documents, err
}
