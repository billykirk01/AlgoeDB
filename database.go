package AlgoeDB

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Database struct {
	config    DatabaseConfig
	mutex     sync.Mutex
	documents []map[string]interface{}
}

type DatabaseConfig struct {
	Path            string
	OnlyInMemory    *bool
	SchemaValidator SchemaValidator
}

type SchemaValidator func(document interface{}) bool

type QueryFunc func(value interface{}) bool

func NewDatabase(config *DatabaseConfig) (*Database, error) {
	onlyInMemory := false

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

	if config.Path != "" {
		err := database.load()
		if err != nil {
			log.Fatal(err)
		}
	}

	return &database, nil
}

func (d *Database) InsertOne(document map[string]interface{}) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.config.SchemaValidator(document) {
		return errors.New("document failed schema validtion: " + fmt.Sprint(document))
	}

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
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for _, document := range documents {
		if d.config.SchemaValidator != nil && !d.config.SchemaValidator(document) {
			return errors.New("document failed schema validtion: " + fmt.Sprint(document))
		}
	}

	d.documents = append(d.documents, documents...)

	if !*d.config.OnlyInMemory {
		err := d.save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) FindOne(query map[string]interface{}) map[string]interface{} {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	found := d.searchDocuments(query)

	if len(found) == 0 {
		return nil
	}

	return d.documents[found[0]]
}

func (d *Database) FindMany(query map[string]interface{}) []map[string]interface{} {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	found := d.searchDocuments(query)

	if len(found) == 0 {
		return nil
	}

	results := []map[string]interface{}{}

	for _, index := range found {
		results = append(results, d.documents[index])
	}

	return results
}

func (d *Database) UpdateOne(query map[string]interface{}, document map[string]interface{}) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	found := d.searchDocuments(query)

	if len(found) == 0 {
		return errors.New("could not find document to update")
	}

	temp := d.documents[found[0]]

	for key, value := range document {
		temp[key] = value
	}

	d.documents[found[0]] = temp

	if !*d.config.OnlyInMemory {
		err := d.save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) UpdateMany(query map[string]interface{}, document map[string]interface{}) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	found := d.searchDocuments(query)

	if len(found) == 0 {
		return errors.New("could not find document(s) to update")
	}

	for _, index := range found {
		temp := d.documents[index]
		for key, value := range document {
			temp[key] = value
		}

		d.documents[index] = temp
	}

	if !*d.config.OnlyInMemory {
		err := d.save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) DeleteOne(query map[string]interface{}) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	found := d.searchDocuments(query)

	if len(found) == 0 {
		return errors.New("could not find document to update")
	}

	d.documents = append(d.documents[:found[0]], d.documents[found[0]+1:]...)

	if !*d.config.OnlyInMemory {
		err := d.save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) DeleteMany(query map[string]interface{}) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	found := d.searchDocuments(query)
	if len(found) == 0 {
		return errors.New("could not find document(s) to update")
	}

	foundMap := map[int]bool{}
	for _, value := range found {
		foundMap[value] = true
	}

	temp := []map[string]interface{}{}
	for index, value := range d.documents {
		if foundMap[index] {
			temp = append(temp, value)
		}
	}

	d.documents = temp

	if !*d.config.OnlyInMemory {
		err := d.save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) load() error {
	content := "[]"

	if d.config.Path != "" {
		if _, err := os.Stat(d.config.Path); os.IsNotExist(err) {
			_, err := os.Create(d.config.Path)
			if err != nil {
				return errors.New("failed to create file: " + d.config.Path)
			}
		}

		bytes, err := ioutil.ReadFile(d.config.Path)
		if err != nil {
			return errors.New("failed to read file: " + d.config.Path)
		}

		if len(bytes) != 0 {
			content = string(bytes)
		}
	}

	documents := []map[string]interface{}{}
	err := json.Unmarshal([]byte(content), &documents)
	if err != nil {
		return err
	}

	d.documents = documents

	err = d.save()
	if err != nil {
		return err
	}

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

func (d *Database) searchDocuments(query map[string]interface{}) []int {
	found := []int{}

	for index, document := range d.documents {

		include := true

		for key, queryValue := range query {
			if !include {
				break
			}

			documentValue := document[key]
			include = matchValues(queryValue, documentValue)
		}

		if include {
			found = append(found, index)
		}
	}

	return found
}

func matchValues(queryValue interface{}, documentValue interface{}) bool {

	if queryValue == documentValue {
		return true
	}

	switch x := queryValue.(type) {
	case QueryFunc:
		return x(documentValue)
	}

	return false
}
