package main

import (
	"encoding/json"
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

	if userid == "" {
		return
	}

	var user *User
	var ok bool

	if user, ok = CurrentUsers[userid]; !ok {
		return
	}

	ws := upgrade(w, req)

	connection := &UserConnection{
		Socket: ws,
	}

	user.listeners <- connection

	func() {
		for {
			read, code, err := ws.Read()
			if err != nil || code == Close {
				return
			}
			fmt.Println(read)
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

	switch r.Method {
	case "POST":
		user := newUser()
		CurrentUsers[user.Id] = user

		w.Write([]byte(user.Id))
	case "GET":

		userids := []string{}
		for id, _ := range CurrentUsers {
			userids = append(userids, id)
		}

		j, _ := json.Marshal(userids)

		w.Write(j)

	}

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

	method := r.Method

	switch method {

	case "POST":
		b, _ := ioutil.ReadAll(r.Body)
		body := string(b)
		bookmark := newBookmark(body)
		user.UpdateBookmark(bookmark)
		fmt.Fprintln(w, bookmark.Id)

	case "GET":
		w.Write(<-user.GetBookmarks())
	}

}
