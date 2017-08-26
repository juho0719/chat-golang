package main

import (
	"log"
	"net/http"
	"sync"
	"html/template"
	"path/filepath"
	"flag"
	"chat-golang/trace"
	"os"
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
	t.templ.Execute(w, r)
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
		tracer: trace.Off(),
	}
}
func main() {
	var addr = flag.String("addr", ":8081", "The addr of the application.")
	flag.Parse()	// 플래그 파싱
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// 방을 가져옴
	go r.run()
	// 웹 서버 시작
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
