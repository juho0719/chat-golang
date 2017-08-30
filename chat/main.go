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
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func newRoom(UseAuthAvatar) *room {
	return &room{
		forward: make(chan *message),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
		tracer: trace.Off(),
		avatar: avatar,
	}
}
func main() {
	var addr = flag.String("addr", ":8081", "The addr of the application.")
	flag.Parse()	// 플래그 파싱
	//gomniauth 설정
	gomniauth.SetSecurityKey("PUT YOUR AUTH KEY HERE")
	gomniauth.WithProviders(
		facebook.New("key", "secret", "http://localhost:8081/auth/callback/facebook"),
		github.New("key", "secret", "http://localhost:8081/auth/callback/github"),
		google.New("key", "secret", "http://localhost:8081/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("/chat-golang/chat/assets/"))))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: "",
			Path: "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	// 방을 가져옴
	go r.run()
	// 웹 서버 시작
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
