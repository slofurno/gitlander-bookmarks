package main

import (
	"fmt"
	"github.com/go-mangos/mangos/protocol/bus"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/gorilla/mux"
	"github.com/slofurno/bookmarks/collection"
	"net/http"
	"os"
	"sync"
	"time"
)

type concurrentMap struct {
	m    map[string]collection.Collection
	lock *sync.Mutex
}

func (s *concurrentMap) Get(key string) collection.Collection {
	s.lock.Lock()
	defer s.lock.Unlock()

	ok, c := s.m[key]

	if !ok {
		c = collection.NewCollection()
		s.m[key] = c
	}

	return c
}

var collections *concurrentMap

func init() {
	collections = &concurrentMap{
		m:    make(map[string]collection.Collection),
		lock: &sync.Mutex{},
	}
}

func assert(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func insert(res http.ResponseWriter, req *http.Request) {

}

func get(res http.ResponseWriter, req *http.Request) {

}

func delete(res http.ResponseWriter, req *http.Request) {

}

func publish(outgoing chan []byte) {
	publish, err := pub.NewSocket()
	publish.AddTransport(tcp.NewTransport())

	assert(publish.Listen("tcp://*:11400"))

	for msg := range outgoing {
		assert(publish.Send(msg))
	}
}

func listen(args []string, outgoing chan []byte) {

	sock, err := bus.NewSocket()

	assert(err)

	sock.AddTransport(tcp.NewTransport())
	assert(sock.Listen(args[1]))

	time.Sleep(time.Second)

	for _, ad := range args[2:] {
		assert(sock.Dial(ad))
		fmt.Println("dialed", ad)
	}

	time.Sleep(time.Second)

	assert(sock.Send([]byte("hey from " + args[1])))

	for {
		msg, err := sock.Recv()
		assert(err)
		outgoing <- msg
	}
}

func main() {

	outgoing := make(chan []byte, 32)
	go listen(os.Args, outgoing)
	go publish(outgoing)

	r := mux.NewRouter()

	r.HandleFunc("/{key}", insert).Methods("POST")
	r.HandleFunc("/{key}", get).Methods("GET")
	r.HandleFunc("/{key}", delete).Methods("DELETE")

	http.ListenAndServe(":11411", nil)

}
