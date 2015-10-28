package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var CurrentUsers = map[string]*User{}
var secretKey []byte

func init() {
	var err error
	secretKey, err = ioutil.ReadFile("secret/secret.key")
	if err != nil {
		panic("i need the secret key")
	}
}

func main() {

	http.HandleFunc("/ws", websocketHandler)
	http.HandleFunc("/api/follow", subscriptionHandler)
	http.HandleFunc("/api/bookmarks", authed(bookmarkHandler))
	http.HandleFunc("/api/user", userHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":80", nil)
}

func authed(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		userid := r.URL.Query().Get("user")

		if token == "" || userid == "" {
			return
		}

		tb, _ := base32.StdEncoding.DecodeString(token)
		mac := hmac.New(sha256.New, secretKey)
		mac.Write([]byte(userid))
		expected := mac.Sum(nil)

		if !hmac.Equal(expected, tb) {
			fmt.Println("auth failed")
			return
		}

		fmt.Println("auth ok")

		h(w, r)
	}
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

		mac := hmac.New(sha256.New, secretKey)
		mac.Write([]byte(user.Id))
		b := mac.Sum(nil)

		token := base32.StdEncoding.EncodeToString(b)

		userResponse := struct {
			User  string `json:"user"`
			Token string `json:"token"`
		}{user.Id, token}

		jb, err := json.Marshal(userResponse)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(jb)

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
