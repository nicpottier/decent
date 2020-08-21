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
	"strings"
	"time"

	"go.bug.st/serial"

	"github.com/nicpottier/decent/hub"
	"github.com/nicpottier/decent/parser"
	_ "github.com/nicpottier/decent/types"
)

var server = flag.String("server", "0.0.0.0:8080", "what address and port to start the HTTP server on")
var de1 = flag.String("de1", "", "filesystem path (e.g. for a serial device)  or server ip/hostname:port for connecting to the de1")

func serveRoots(w http.ResponseWriter, r *http.Request) {
	// we only deal with roots
	if strings.Count(r.URL.Path, "/") > 1 {
		http.Error(w, "not found", http.StatusNotFound)
	}

	path := r.URL.Path
	if path == "/" {
		path = "index"
	}
	path = strings.TrimPrefix(path, "/")
	http.ServeFile(w, r, fmt.Sprintf("./static/%s.html", path))
}

var tcpRegexp = regexp.MustCompile(`[a-zA-Z0-9.]+:[0-9]+`)

// connects to the DE1 at the passed in location. loc can be either a domain/ip:port or
// a path to a file to serial port
func connectToDE1(loc string) (io.ReadCloser, error) {
	if tcpRegexp.MatchString(loc) {
		return net.Dial("tcp", loc)
	}

	m := &serial.Mode{
		BaudRate: 115200,
	}
	return serial.Open(loc, m)
}

func main() {
	flag.Parse()

	// start our hub
	h := hub.New()
	go h.Run()

	// then our webserver

	// convenience rewritter for "apps" with root html files.. /index will be rewritten to ./static/index.html
	http.HandleFunc("/", serveRoots)

	// server ./static as /static
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// expose web socket at /ws
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWS(h, w, r)
	})

	// we have a de1 to connect to, do so
	if *de1 != "" {
		go func() {
			for {
				fmt.Printf("connecting to: %s\n", *de1)
				conn, err := connectToDE1(*de1)
				defer conn.Close()

				if err != nil {
					fmt.Printf("error connecting: %s\n", err.Error())
				} else {
					fmt.Println("connected!")
					reader := bufio.NewReader(conn)
					for {
						ms, err := parser.ReadNextToken(reader)
						fmt.Printf("%s\n", ms)
						if err == io.EOF {
							fmt.Println("de1 connectioni closed")
							break
						}
						if err != nil {
							fmt.Printf("error: %s", err.Error())
							time.Sleep(100 * time.Millisecond)
							continue
						}

						h.Broadcast <- ms
					}
				}

				fmt.Println("sleeping before reconnecting..")
				time.Sleep(5 * time.Second)
			}
		}()
	}

	fmt.Printf("listening on %s\n", *server)
	err := http.ListenAndServe(*server, nil)
	if err != nil {
		log.Fatal("Error while serving: ", err)
	}
}
