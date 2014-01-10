package main

import (
	"net/http"
	"container/list"
	"time"
	"io"
	"io/ioutil"
)

func handler (w http.ResponseWriter, r *http.Request, clients *list.List) {
	res := clients.PushBack(make(chan string))
	select {
	case <-time.After(60e9):
		clients.Remove(res)
		http.Error(w, http.StatusText(408), 408)
	case r := <-res.Value.(chan string):
		clients.Remove(res)
		io.WriteString(w, r)
	}
}

func responder (clients *list.List, c chan string) {
	for {
		text := <-c
		for e := clients.Front(); e != nil; e = e.Next() {
			e.Value.(chan string) <- text
		}
	}
}

func enterChat (w http.ResponseWriter, r *http.Request, c chan string) {
	body, _ := ioutil.ReadAll(r.Body)
	c <- string(body)
}

func main () {
	responses := make(chan string)
	clients := list.New()
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, clients)
	})
	http.HandleFunc("/chatentry", func(w http.ResponseWriter, r *http.Request) {
		enterChat(w, r, responses)
	})
	go responder(clients, responses)
	http.ListenAndServe("0.0.0.0:8080", nil)
}