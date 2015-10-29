package main

import (
	"bufio"
	"fmt"
	"os"
)

type Filebase struct {
	Fd      *os.File
	Inbox   chan []byte
	writer  *bufio.Writer
	scanner *bufio.Scanner
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

	scanner := bufio.NewScanner(db.Fd)
	writer := bufio.NewWriter(db.Fd)

	db.scanner = scanner
	db.writer = writer

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	fmt.Fprintln(writer, "lets add a line")
	fmt.Fprintln(writer, "lets add a line2")

	writer.Flush()

	return db
}

func (db *Filebase) AddUser(user *User) {

}
