package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-mangos/mangos/protocol/bus"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/gorilla/mux"
	"github.com/slofurno/bookmarks/collection"
	"github.com/slofurno/bookmarks/filebase"

	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

type concurrentMap struct {
	m    map[string]collection.Collection
	lock *sync.Mutex
}

type record struct {
	Topic string
	Tuple []byte
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
var db *filebase.Filebase

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

func insert(key string, payload []byte) error {

	item := &collection.Tuple{}
	err := json.Unmarshal(payload, item)

	if err != nil {
		return err
	}

	collection := store.Get(key)
	collection.Update(item)

	return nil
}

func postTuple(res http.ResponseWriter, req *http.Request) {
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

	r := &record{
		Topic: col,
		Tuple: body,
	}

	fb, _ := json.Marshal(r)
	db.WriteRecord(fb)

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

	assert(publish.Listen(args[3]))
	//"tcp://:11400"
	sock, err := bus.NewSocket()

	assert(err)

	sock.AddTransport(tcp.NewTransport())
	assert(sock.Listen(args[4]))

	time.Sleep(time.Second)

	for _, ad := range args[5:] {
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
	db = filebase.New(args[1] + ".log")

	reinsert := func(line []byte) {
		rec := &record{}
		err := json.Unmarshal(line, rec)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		insert(rec.Topic, rec.Tuple)
	}

	fmt.Println("reading log...")
	db.ReadRecords(reinsert)
	fmt.Println("done")

	fmt.Println(args)
	go listen(args)
	go publish(args)

	r := mux.NewRouter()

	r.HandleFunc("/{collection}", postTuple).Methods("POST")
	r.HandleFunc("/{collection}", getAll).Methods("GET")
	r.HandleFunc("/{collection}", delete).Methods("DELETE")

	//11411
	fmt.Println("api running on", args[2])
	assert(http.ListenAndServe(args[2], r))

}
