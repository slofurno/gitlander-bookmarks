package main

import (
	"bufio"
	"fmt"
	"os"
)

//TODO:pls think of a better way to do this
type Filebase struct {
	Fd    *os.File
	Inbox chan []byte
	Pls   bool
}

type entry struct {
	User string
	Type string
	Data string
}

func NewFilebase(filename string) *Filebase {
	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	db := &Filebase{
		Fd:    fd,
		Inbox: make(chan []byte, 1024),
	}

	db.initWorker()
	return db
}

func (db *Filebase) initWorker() {
	go func() {
		defer db.Fd.Close()

		for {
			record := <-db.Inbox
			_, err := db.Fd.Write(record)
			if err != nil {
				fmt.Println(err.Error())
			}
			db.Fd.Write([]byte("\n"))
		}

	}()
}

func (db *Filebase) WriteRecord(b []byte) {
	db.Inbox <- b
}

func (db *Filebase) ReadRecords(f func([]byte)) {
	scanner := bufio.NewScanner(db.Fd)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		f(scanner.Bytes())
	}
}
