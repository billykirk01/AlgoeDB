package AlgoeDB

import (
	"io/ioutil"
	"log"
	"os"
)

type Writer struct {
	Path   string
	next   string
	locked bool
}

func NewWriter(path string) *Writer {
	return &Writer{Path: path}
}

func (w *Writer) Write(data string) {
	// Add writing to the queue if writing is locked
	if w.locked {
		w.next = data
		return
	}

	// Lock writing
	w.locked = true

	// Write data
	temp := w.Path + ".temp"
	f, err := os.Create(temp)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = ioutil.WriteFile(temp, []byte(data), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename(temp, w.Path)
	if err != nil {
		log.Fatal(err)
	}

	// Unlock writing
	w.locked = false

	// Start next writing if there is data in the queue
	if w.next != "" {
		nextTmp := w.next
		w.next = ""
		w.Write(nextTmp)
	}
}
