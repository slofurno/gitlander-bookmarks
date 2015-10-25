package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/ws", websocketHandler)
	http.HandleFunc("/api/bookmarks", bookmarkHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":80", nil)
}

func websocketHandler(w http.ResponseWriter, req *http.Request) {
	ws := upgrade(w, req)

	ws.Inbox <- []byte("WELCOME")
	ws.Inbox <- []byte("WELCOME SDFSDFSD")

	func() {
		for {

			read, code, err := ws.Read()

			if err != nil || code == Close {
				return
			}

			fmt.Println(read)
			ws.Inbox <- []byte("thanks for msg")
		}
	}()

	fmt.Println("disconnected")
}

func bookmarkHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	method := r.Method

	for k, v := range qs {
		fmt.Fprintln(w, method, k, v)
	}
}
