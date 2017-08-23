package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
		<html>
		<head>
			<title>Chat</title>
		</head>
		<body>
			gogogo
		</body>
		</html>
		`))
	})

	// webserver start
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndSerfve:", err)
	}
}
