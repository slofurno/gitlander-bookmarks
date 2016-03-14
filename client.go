package main

import (
	"fmt"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/slofurno/bookmarks/collection"
	"io/ioutil"
	"net/http"
	"os"
)

type ClusterClient struct {
	sock mangos.Socket
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

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
