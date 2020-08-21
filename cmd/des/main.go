package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"time"

	"go.bug.st/serial"

	"github.com/nicpottier/decent/hub"
	"github.com/nicpottier/decent/parser"
	_ "github.com/nicpottier/decent/types"
)

var server = flag.String("server", "0.0.0.0:8080", "what address and port to start the HTTP server on")
var de1 = flag.String("de1", "", "filesystem path (e.g. for a serial device)  or server ip/hostname:port for connecting to the de1")

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/html/home.html")
}

func serveDebug(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/html/debug.html")
}

func main() {
	flag.Parse()

	h := hub.New()
	go h.Run()
	http.HandleFunc("/debug", serveDebug)
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWS(h, w, r)
	})

	if *de1 != "" {
		fmt.Printf("connecting to %s\n", *de1)
		tcpRegexp := regexp.MustCompile(`[a-zA-Z0-9.]+:[0-9]+`)
		var conn io.Reader
		var err error
		if tcpRegexp.MatchString(*de1) {
			conn, err = net.Dial("tcp", *de1)
		} else {
			m := &serial.Mode{
				BaudRate: 115200,
			}
			conn, err = serial.Open(*de1, m)
		}
		if err != nil {
			panic(fmt.Sprintf("unable to connect to de1: %s", err.Error()))
		}
		fmt.Printf("connected!\n")

		reader := bufio.NewReader(conn)
		go func() {
			for {
				ms, err := parser.ReadNextToken(reader)
				fmt.Printf("%s\n", ms)
				if err == io.EOF {
					fmt.Println("de1 connectioni closed")
					return
				}
				if err != nil {
					fmt.Printf("error: %s", err.Error())
					time.Sleep(100 * time.Millisecond)
					continue
				}

				h.Broadcast <- ms
			}
		}()
	}

	fmt.Printf("listening on %s\n", *server)
	err := http.ListenAndServe(*server, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
