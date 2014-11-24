package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var connections map[*websocket.Conn]bool

func sendAll(msg []byte) {
	for conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			delete(connections, conn)
			return
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	connections[conn] = true
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		log.Println(string(msg))
		sendAll(msg)
	}
}

func main() {
	port := flag.Int("port", 80, "default port")
	dir := flag.String("directory", "web/", "files directory")
	flag.Parse()

	connections = make(map[*websocket.Conn]bool)

	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	http.HandleFunc("/ws", wsHandler)

	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)

	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
