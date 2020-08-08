package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/nicpottier/decent/parser"
	_ "github.com/nicpottier/decent/types"
)

var server = flag.String("server", ":8080", "what address and port to start the HTTP server on")
var de1 = flag.String("de1", "", "what ip and port the de1 is on")

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}

func serveDebug(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "debug.html")
}

func main() {
	flag.Parse()

	hub := newHub()
	go hub.run()
	http.HandleFunc("/debug", serveDebug)
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	if *de1 != "" {
		fmt.Printf("connecting to %s\n", *de1)
		conn, err := net.Dial("tcp", *de1)
		if err != nil {
			panic(fmt.Sprintf("unable to connect to de1: %s", err.Error()))
		}
		fmt.Printf("connected!\n")

		reader := bufio.NewReader(conn)
		go func() {
			for {
				mt, mb, err := parser.ReadNextToken(reader)
				fmt.Printf("[%s]%s\n", mt, strings.ToUpper(hex.EncodeToString(mb)))
				if err == io.EOF {
					fmt.Println("de1 connectioni closed")
					return
				}
				if err != nil {
					fmt.Printf("error: %s", err.Error())
					time.Sleep(100 * time.Millisecond)
					continue
				}

				hub.broadcast <- []byte(fmt.Sprintf("[%s]%s", mt, strings.ToUpper(hex.EncodeToString(mb))))
			}
		}()
	}

	err := http.ListenAndServe(*server, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
