package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	staticDir = flag.String("dir", "/static", "static dir")
)

func main() {
	flag.Parse()
	fs := http.FileServer(http.Dir(*staticDir))
	http.Handle("/", fs)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
