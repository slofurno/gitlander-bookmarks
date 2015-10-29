package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"net/http"
	"strings"
)

var userSubscriptions = map[string]*Collection{}
var userBookmarks = map[string]*Collection{}

var userInfos = map[string]*userInfo{}

var secretKey []byte
var database = newFilebase("test.db")

type userInfo struct {
	subscriptions *Collection
	bookmarks     *Collection
}

type RequestContext struct {
	isAuthed bool
	userinfo *userInfo
	userid   string
}

func newUserInfo() *userInfo {
	return &userInfo{subscriptions: newCollection(), bookmarks: newCollection()}
}

func makeUuid() string {
	u, _ := uuid.NewV4()
	return u.String()
}

func init() {
	var err error
	secretKey, err = ioutil.ReadFile("secret/secret.key")
	if err != nil {
		panic("i need the secret key")
	}
}

func main() {

	http.HandleFunc("/api/img/", authed(imgHandler))
	//http.HandleFunc("/api/img", imgHandler)
	http.HandleFunc("/ws", authed(websocketHandler))
	http.HandleFunc("/api/follow", authed(subscriptionHandler))
	http.HandleFunc("/api/bookmarks", authed(bookmarkHandler))
	http.HandleFunc("/api/user", userHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":555", nil)
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

			var user *userInfo
			var ok bool

			if user, ok = userInfos[userid]; ok {

				tb, _ := base32.StdEncoding.DecodeString(token)
				mac := hmac.New(sha256.New, secretKey)
				mac.Write([]byte(userid))
				expected := mac.Sum(nil)

				if hmac.Equal(expected, tb) {
					context.isAuthed = true
					context.userinfo = user
					context.userid = userid
					fmt.Println("user authed as: ", userid)
					//          w.WriteHeader(http.StatusForbidden)
				}
			}
		}

		h(w, r, context)
	}
}

func imgHandler(w http.ResponseWriter, r *http.Request, context *RequestContext) {

	body := r.URL.Query().Get("body")

	if !context.isAuthed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bookmark := &Bookmark{Id: makeUuid(), Url: body, Owner: context.userid}
	context.userinfo.bookmarks.Add(bookmark.Id, bookmark)
	fmt.Println("adding bookmark", body)
}

func websocketHandler(w http.ResponseWriter, req *http.Request, context *RequestContext) {

	if !context.isAuthed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ws := upgrade(w, req)

	connection := newUserConnection(context.userinfo.subscriptions, ws)

	func() {
		for {
			read, code, err := ws.Read()
			if err != nil || code == Close {
				return
			}
			fmt.Println(read)
		}
	}()

	connection.onstop()

	fmt.Println("disconnected")
}

func subscriptionHandler(w http.ResponseWriter, r *http.Request, context *RequestContext) {

	if !context.isAuthed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sub := r.URL.Query().Get("follow")

	//var other *userInfo
	var ok bool

	if _, ok = userInfos[sub]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cant find sub"))
		return
	}

	context.userinfo.subscriptions.Add(sub, true)

	//other.AddSub(context.user)

}

func userHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":

		userid := makeUuid()
		userinfo := newUserInfo()
		userinfo.subscriptions.Add(userid, userid)
		userInfos[userid] = userinfo

		mac := hmac.New(sha256.New, secretKey)
		mac.Write([]byte(userid))
		b := mac.Sum(nil)

		token := base32.StdEncoding.EncodeToString(b)

		userResponse := struct {
			User  string `json:"user"`
			Token string `json:"token"`
		}{userid, token}

		jb, err := json.Marshal(userResponse)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(jb)

	case "GET":

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

		bookmark := &Bookmark{Id: makeUuid(), Url: body, Owner: context.userid}
		context.userinfo.bookmarks.Add(bookmark.Id, bookmark)
		fmt.Fprintln(w, bookmark.Id)

	case "GET":
		//w.Write(<-context.user.GetBookmarks())
	}

}
