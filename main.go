package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var CurrentUsers = map[string]*User{}

func main() {

	http.HandleFunc("/ws", websocketHandler)
	http.HandleFunc("/api/follow", subscriptionHandler)
	http.HandleFunc("/api/bookmarks", bookmarkHandler)
	http.HandleFunc("/api/user", userHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":80", nil)
}

func websocketHandler(w http.ResponseWriter, req *http.Request) {
	userid := req.URL.Query().Get("user")
	fmt.Println("userid: ", userid)
	if userid == "" {
		return
	}

	if _, ok := CurrentUsers[userid]; !ok {
		return
	}

	ws := upgrade(w, req)

	user := &UserConnection{
		Socket: ws,
	}

	ws.WriteS("welcome")

	//lock this
	CurrentUsers[userid].listeners <- user

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

func subscriptionHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	userid := qs.Get("user")
	sub := qs.Get("sub")

	if userid == "" {
		fmt.Println("no userid")
		return
	}

	var user *User
	var other *User
	var ok bool

	if user, ok = CurrentUsers[userid]; !ok {
		fmt.Println("cant find user")
		return
	}

	if other, ok = CurrentUsers[sub]; !ok {
		fmt.Println("cant find sub")
		return
	}

	other.AddSub(user)

}

func userHandler(w http.ResponseWriter, r *http.Request) {

	user := newUser()
	CurrentUsers[user.Id] = user

	w.Write([]byte(user.Id))

}

func bookmarkHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	userid := qs.Get("user")

	if userid == "" {
		fmt.Println("no userid")
		return
	}

	var user *User
	var ok bool

	if user, ok = CurrentUsers[userid]; !ok {
		return
	}

	b, _ := ioutil.ReadAll(r.Body)
	body := string(b)

	bookmark := newBookmark(body)

	user.UpdateBookmark(bookmark)

	fmt.Fprintln(w, bookmark.Id)

}
