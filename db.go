package main

import (
	"fmt"
	"os"
)

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

func newFilebase(filename string) *Filebase {
	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err.Error())
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
	if !db.Pls {
		return
	}
	db.Inbox <- b
}
