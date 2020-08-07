package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/nicpottier/decent/types"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}

func serveDebug(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "debug.html")
}

func main() {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/debug", serveDebug)
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
