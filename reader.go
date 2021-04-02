package AlgoeDB

import (
	"io/ioutil"
	"log"
	"os"
)

type Reader struct{}

func (r *Reader) read(path string) string {

	if path == "" {
		return "[]"
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	bytes, err := ioutil.ReadFile(path)
	return string(bytes)
}
