package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/slofurno/bookmarks/collection"
	"io/ioutil"
	"net/http"
	"os"
)

var shards = []string{
	"http://127.0.0.1:11411/",
	"http://127.0.0.1:11412/",
	"http://127.0.0.1:11413/",
	"http://127.0.0.1:11414/",
}

var pubshards = []string{
	"tcp://127.0.0.1:11400",
	"tcp://127.0.0.1:11402",
	"tcp://127.0.0.1:11403",
	"tcp://127.0.0.1:11404",
}

type Tuple struct {
	Time  int64
	Key   string
	Value string
}

type ClusterIterator struct {
	items chan *Tuple
}

func (s *ClusterIterator) Next() *Tuple {
	return <-s.items
}

type ClusterClient struct {
}

func check(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (s *ClusterClient) Fetch(key string) <-chan *Tuple {
	sock, _ := sub.NewSocket()
	sock.AddTransport(tcp.NewTransport())

	bkey := []byte(key)
	si := collection.Crc16(bkey) % 4

	fmt.Println("index", si)
	uri := shards[si]

	check(sock.Dial(pubshards[si]))
	check(sock.SetOption(mangos.OptionSubscribe, bkey))

	items := make(chan *Tuple, 64)

	go func() {
		res, _ := http.Get(uri + key)
		xs := []*Tuple{}
		b, _ := ioutil.ReadAll(res.Body)
		if check(json.Unmarshal(b, &xs)) {
			for i := 0; i < len(xs); i++ {
				items <- xs[i]
			}
		}

		for {
			msg, err := sock.Recv()
			if !check(err) {
				return
			}
			x := &Tuple{}
			body := msg[len(bkey):]
			json.Unmarshal(body, x)
			items <- x
		}
	}()

	return items
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

/*
func NewClient() *ClusterClient {
	sock, _ := sub.NewSocket()
	sock.AddTransport(tcp.NewTransport())
	err := sock.Dial("tcp://127.0.0.1:11400")

	if err != nil {
		die("cannot subscribe: %s", err.Error())
	}

	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))

	if err != nil {
		die("cannot subscribe: %s", err.Error())
	}

	return &ClusterClient{
		sock: sock,
	}
}

func (s *ClusterClient) Get(key string) ([]byte, error) {

	fmt.Println("key:", key)
	server := collection.Crc16([]byte(key)) % 4
	fmt.Println("server:", server)
	res, err := http.Get("http://127.0.0.1:11411/" + key)

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(res.Body)
}
*/

func (s *ClusterClient) Post(key string, item *Tuple) (*http.Response, error) {
	si := collection.Crc16([]byte(key)) % 4
	b, _ := json.Marshal(item)
	buf := bytes.NewBuffer(b)

	return http.Post(shards[si]+key, "application/json", buf)
}
