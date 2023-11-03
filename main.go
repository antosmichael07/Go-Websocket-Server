package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func newServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) HandleWS(ws *websocket.Conn) {
	consoleLog(fmt.Sprintf("New connection from: %s", ws.RemoteAddr()), "SERVER")
	s.conns[ws] = true
	s.read(ws)
}

func (s *Server) read(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	charCount, _ := ws.Read(buf)
	data := fmt.Sprintf("%s", buf[:charCount])
	if data == "connectedToWS" {
		consoleLog(fmt.Sprintf("Connected to websocket from: %s", ws.RemoteAddr()), "SERVER")
	}
}

func sendHTMLServerHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := os.ReadFile("index.html")
	fmt.Fprintf(w, string(page))
}

func consoleLog(msg string, source string) {
	time := fmt.Sprintf("%s", time.Now())
	fmt.Printf("[%s] [%s] %s\n", time[0:19], source, msg)
}

func main() {
	server := newServer()
	http.Handle("/ws", websocket.Handler(server.HandleWS))
	http.HandleFunc("/", sendHTMLServerHandler)
	http.ListenAndServe(":8000", nil)
}
