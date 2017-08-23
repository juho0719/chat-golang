package main

import (
	"log"
	"net/http"
	"sync"
	"html/template"
	"path/filepath"
)

// templ은 하나의 템플릿
type templateHandler struct {
	once		sync.Once
	filename	string
	templ		*template.Template
}

// ServeHttp가 Http요청을 처리
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil);
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
	}
}
func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}

	if err := http.ListenAndServe(":8089", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	// webserver start
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
