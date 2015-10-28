package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

	http.HandleFunc("/ws", authed(websocketHandler))
	http.HandleFunc("/api/follow", authed(subscriptionHandler))
	http.HandleFunc("/api/bookmarks", authed(bookmarkHandler))
	http.HandleFunc("/api/user", userHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":555", nil)
}

type RequestContext struct {
	isAuthed bool
	user     *User
}

func authed(h func(w http.ResponseWriter, r *http.Request, context *RequestContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		context := &RequestContext{}
		auth := r.Header.Get("Authorization")

		authParams := strings.Split(auth, ":")

		var userid string = ""
		var token string = ""

		if len(authParams) == 2 {
			userid = authParams[0]
			token = authParams[1]
		} else {
			qs := r.URL.Query()
			userid = qs.Get("user")
			token = qs.Get("token")
		}

		if userid != "" && token != "" {

			var user *User
			var ok bool

			if user, ok = CurrentUsers[userid]; ok {

				tb, _ := base32.StdEncoding.DecodeString(token)
				mac := hmac.New(sha256.New, secretKey)
				mac.Write([]byte(userid))
				expected := mac.Sum(nil)

				if hmac.Equal(expected, tb) {
					context.isAuthed = true
					context.user = user
					fmt.Println("user authed as: ", user.Id)
					//          w.WriteHeader(http.StatusForbidden)
				}
			}
		}

		h(w, r, context)
	}
}

func websocketHandler(w http.ResponseWriter, req *http.Request, context *RequestContext) {

	if !context.isAuthed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := context.user

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

func subscriptionHandler(w http.ResponseWriter, r *http.Request, context *RequestContext) {

	if !context.isAuthed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sub := r.URL.Query().Get("follow")

	var other *User
	var ok bool

	if other, ok = CurrentUsers[sub]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cant find sub"))
		return
	}

	other.AddSub(context.user)

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

func bookmarkHandler(w http.ResponseWriter, r *http.Request, context *RequestContext) {

	method := r.Method

	switch method {

	case "POST":
		b, _ := ioutil.ReadAll(r.Body)
		body := string(b)
		fmt.Println(body)
		if !context.isAuthed {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		bookmark := newBookmark(body)
		context.user.UpdateBookmark(bookmark)
		fmt.Fprintln(w, bookmark.Id)

	case "GET":
		w.Write(<-context.user.GetBookmarks())
	}

}
