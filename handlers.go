package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/slofurno/ws"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

//TODO: this entire handler is a duplicate, as a workaround for restrictions on mixed content
func imgHandler(w http.ResponseWriter, r *http.Request, context *RequestContext) {

	var err error
	body := r.URL.Query().Get("body")

	if !context.isAuthed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	br := &BookmarkRequest{}
	err = json.Unmarshal([]byte(body), br)

	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	buf := []byte(br.Url)
	fmt.Println("req to node")
	resp, err := http.Post("http://localhost:8765", "text/plain", bytes.NewBuffer(buf))

	var content []byte

	if err == nil {
		defer resp.Body.Close()
		content, err = ioutil.ReadAll(resp.Body)
	} else {
		fmt.Println("scraper err")
	}

	bookmark := &Bookmark{
		Id:          makeUuid(),
		Url:         br.Url,
		Description: br.Description,
		Tags:        br.Tags,
		Owner:       context.userinfo.userid,
		Time:        getCurrentTime(),
		Summary:     string(content),
	}

	dataStore.UpsertBookmark(context.userinfo, bookmark)
	fmt.Fprintln(w, bookmark.Id)
}

func websocketHandler(w http.ResponseWriter, req *http.Request, context *RequestContext) {
	if !context.isAuthed {
		fmt.Println("not authed")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sock := ws.Upgrade(w, req)
	fmt.Println("user connected:", context.user)

	for x := range mergeFetch(context.user) {
		b := []byte(x.Value)
		next := &Bookmark{}
		err := json.Unmarshal(b, next)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		sock.Write(b)
	}

	func() {
		for {
			read, code, err := sock.Read()
			if err != nil || code == ws.Close {
				return
			}
			fmt.Println(read)
		}
	}()

	sock.Close()
	fmt.Println("disconnected")
}

func mergeFetch(user string) chan *Tuple {
	user = "ls" + user
	client := &ClusterClient{}
	subs := client.Fetch(user)
	res := make(chan *Tuple, 64)

	push := func(xs <-chan *Tuple) {
		for x := range xs {
			res <- x
		}
	}

	go func() {
		for s := range subs {
			c := client.Fetch(string(s.Key))
			go push(c)
		}
	}()

	return res
}

func subscriptionHandler(w http.ResponseWriter, r *http.Request, context *RequestContext) {

	if !context.isAuthed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sub := r.URL.Query().Get("follow")

	switch r.Method {
	case "DELETE":
		dataStore.DeleteSubscription(context.userinfo, sub)
	case "POST":
		client := &ClusterClient{}

		x := &Tuple{
			Time:  0,
			Key:   sub,
			Value: sub,
		}
		client.Post("ls"+context.user, x)
	}

	w.WriteHeader(http.StatusOK)
}

func userHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":

		code := r.URL.Query().Get("code")

		if code == "" {
			return
		}

		response, err := http.Post("https://github.com/login/oauth/access_token?client_id="+client_id+"&client_secret="+client_secret+"&code="+code, "application/json", nil)
		defer response.Body.Close()

		fmt.Println("github status code:", response.StatusCode)
		bb, _ := ioutil.ReadAll(response.Body)
		githubqs, err := url.ParseQuery(string(bb))

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		githubtoken := githubqs.Get("access_token")
		githubresponse, err := http.Get("https://api.github.com/user?access_token=" + githubtoken)

		ghb, _ := ioutil.ReadAll(githubresponse.Body)

		github_user := &GithubUserResponse{}
		err = json.Unmarshal(ghb, github_user)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		//we can get here with all failures from github...
		if github_user.Id == 0 {
			fmt.Println("github userid of 0...")
			return
		}

		fmt.Println("github user:", github_user)
		userid := strconv.FormatUint(github_user.Id, 10)

		userResponse := struct {
			User   string `json:"user"`
			Token  string `json:"token"`
			UserId string `json:"userid"`
		}{"", "", userid}

		jb, err := json.Marshal(userResponse)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(jb)

	case "GET":

		usersummaries := map[string]*userinfoDto{}
		for userid, userinfo := range userInfos {
			usersummaries[userid] = &userinfoDto{Summary: userinfo.summary, Name: userinfo.name}
		}

		j, _ := json.Marshal(usersummaries)
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		return
	}

}

func summaryHandler(w http.ResponseWriter, rep *http.Request) {
	j, _ := json.Marshal(globalSummary)
	fmt.Println(string(j))

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func bookmarkHandler(w http.ResponseWriter, r *http.Request, context *RequestContext) {

	method := r.Method
	switch method {

	case "GET":

		qs := r.URL.Query()
		terms := qs["t"]
		if len(terms) == 0 {
			w.WriteHeader(200)
			return
		}

		seen := map[string]int{}
		client := &ClusterClient{}

		for _, x := range terms {
			for _, user := range client.Get("t:" + x) {
				seen[user.Key] += 1
			}
		}

		matches := []string{}
		marks := map[string][]*Bookmark{}

		for k, v := range seen {
			if v == len(terms) {
				matches = append(matches, k)
				for _, tuple := range client.Get(k) {
					x := &Bookmark{}
					err := json.Unmarshal([]byte(tuple.Value), x)
					if err == nil {
						marks[k] = append(marks[k], x)
					}
				}
			}
		}

		b, _ := json.Marshal(marks)
		w.Write(b)

		//fmt.Fprint(w, marks)

	case "POST":

		if !context.isAuthed {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		b, _ := ioutil.ReadAll(r.Body)

		br := &BookmarkRequest{}
		json.Unmarshal(b, br)

		bookmark := &Bookmark{
			Id:          makeUuid(),
			Url:         br.Url,
			Description: br.Description,
			Tags:        br.Tags,
			Owner:       context.user,
			Time:        getCurrentTime(),
			//	Summary:     string(content),
		}

		buf, _ := json.Marshal(bookmark)

		client := &ClusterClient{}
		x := &Tuple{
			Key:   bookmark.Id,
			Time:  bookmark.Time,
			Value: string(buf),
		}

		client.Post(context.user, x)

		for _, tag := range bookmark.Tags {
			tag = "t:" + tag
			x := &Tuple{
				Key: context.user,
			}

			client.Post(tag, x)
		}

	}

}
