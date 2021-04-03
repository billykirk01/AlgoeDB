package AlgoeDB

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
)

type Database struct {
	documents []interface{}
	writer    Writer
	config    DatabaseConfig
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

	documents := []interface{}{}

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

func (d *Database) InsertOne(document interface{}) {
	d.documents = append(d.documents, document)
	if *d.config.OnlyInMemory == false {
		d.save()
	}
}

func (d *Database) InsertMany(documents []interface{}) {
	for _, document := range documents {
		d.documents = append(d.documents, document)
	}

	if *d.config.OnlyInMemory == false {
		d.save()
	}
}

func (d *Database) load() {
	if d.config.Path == "" {
		return
	}

	reader := Reader{}
	content := reader.read(d.config.Path)

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
	d.documents = []interface{}{}
	if *d.config.OnlyInMemory {
		d.save()
	}
}

func SearchDocuments(query map[string]interface{}, documents []interface{}) {

	// TODO
	for _, v := range query {

		if reflect.TypeOf(v).Kind() == reflect.Func {
			queryFunc := reflect.ValueOf(v).Interface().(QueryFunc)
			fmt.Println(queryFunc(1))
		}

	}
}

func matchValues(queryValue interface{}, documentValue interface{}) {
	// TODO
}

func parseDatabaseStorage(content string) ([]interface{}, error) {
	documents := []interface{}{}
	err := json.Unmarshal([]byte(content), &documents)
	return documents, err
}
