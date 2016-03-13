package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-mangos/mangos/protocol/bus"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/gorilla/mux"
	"github.com/slofurno/bookmarks/collection"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

func getCurrentTime() int64 {
	nanos := time.Now().UnixNano()
	return nanos / 1000000
}

type concurrentMap struct {
	m    map[string]collection.Collection
	lock *sync.Mutex
}

func (s *concurrentMap) Get(key string) collection.Collection {
	s.lock.Lock()
	defer s.lock.Unlock()

	c, ok := s.m[key]

	if !ok {
		c = collection.NewCollection()
		s.m[key] = c
	}

	return c
}

var store *concurrentMap
var outbox chan []byte
var inbox chan []byte

func init() {
	store = &concurrentMap{
		m:    make(map[string]collection.Collection),
		lock: &sync.Mutex{},
	}

	outbox = make(chan []byte, 128)
	inbox = make(chan []byte, 128)
}

func assert(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func insert(res http.ResponseWriter, req *http.Request) {
	fmt.Println("insert called")
	vars := mux.Vars(req)
	col := vars["collection"]

	item := &collection.Tuple{}

	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, item)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	collection := store.Get(col)
	//TODO: use update for everything?
	collection.Update(item)

	ok := []byte("ADD")
	ok = append(ok, body...)
	outbox <- ok

	res.WriteHeader(200)
}

func getAll(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	key := vars["collection"]

	encoder := json.NewEncoder(res)
	col := store.Get(key)
	items := col.Get()

	for _, item := range items {
		fmt.Println(string(item.Key), string(item.Value))
	}

	encoder.Encode(items)
}

func delete(res http.ResponseWriter, req *http.Request) {
	//vars := mux.Vars(req)
	//key := vars["key"]
}

func publish(args []string) {

}

func listen(args []string) {

	publish, _ := pub.NewSocket()
	publish.AddTransport(tcp.NewTransport())

	assert(publish.Listen(args[2]))
	//"tcp://:11400"
	sock, err := bus.NewSocket()

	assert(err)

	sock.AddTransport(tcp.NewTransport())
	assert(sock.Listen(args[3]))

	time.Sleep(time.Second)

	for _, ad := range args[4:] {
		assert(sock.Dial(ad))
		fmt.Println("dialed", ad)
	}

	time.Sleep(time.Second)

	go func() {
		for {
			select {
			case next := <-outbox:
				fmt.Println("publishing ours")
				assert(sock.Send(next))
				assert(publish.Send(next))

			case forward := <-inbox:
				fmt.Println("forwarding from bus")
				assert(publish.Send(forward))
			}

		}
	}()

	for {
		msg, err := sock.Recv()
		assert(err)
		inbox <- msg
	}
}

//args: name, http, pub, [bus ips]
func main() {

	args := os.Args

	fmt.Println(args)
	go listen(args)
	go publish(args)

	r := mux.NewRouter()

	r.HandleFunc("/{collection}", insert).Methods("POST")
	r.HandleFunc("/{collection}", getAll).Methods("GET")
	r.HandleFunc("/{collection}", delete).Methods("DELETE")

	//11411
	fmt.Println("api running on", args[1])
	assert(http.ListenAndServe(args[1], r))

}
