package main

import (
	"net/http"
	"os"
	"bufio"
	"container/list"
	"time"
	"io"
)

func handler (w http.ResponseWriter, r *http.Request, clients *list.List) {
	res := clients.PushBack(make(chan string))
	select {
	case <-time.After(60e9):
		clients.Remove(res)
		io.WriteString(w, "timeout")
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

func getString (c chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		thing, _ := reader.ReadString('\n')
		c <- thing
	}
}

func main() {
	responses := make(chan string)
	clients := list.New()
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/thing", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, clients)
	})
	go getString(responses)
	go responder(clients, responses)
	http.ListenAndServe("0.0.0.0:8080", nil)
}