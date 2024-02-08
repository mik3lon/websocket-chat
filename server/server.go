package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"websocket-chat/client"
	"websocket-chat/hub"
)

type WSServer struct {
	address string
	hub     *hub.Hub
}

func (s *WSServer) Run() {
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		s.handleWS(writer, request)
	})

	server := http.Server{
		Addr: s.address,
	}

	log.Printf("Server running in address %s", s.address)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error running server %s", err)
	}

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *WSServer) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("error upgrading to websocket %s", err)
	}

	c := client.NewClient(
		conn,
		make(chan []byte, 256),
		s.hub.UnregisterChan(),
		s.hub.BroadcastChan(),
	)

	s.hub.Register(c)

	go c.WritePump()
	go c.ReadPump()
}

type Config func(server *WSServer)

func NewServer(configs ...Config) *WSServer {
	s := &WSServer{}

	for _, o := range configs {
		o(s)
	}

	return s
}

func WithAddress(address string) Config {
	return func(server *WSServer) {
		server.address = address
	}
}

func NewDefaultServer(h *hub.Hub) *WSServer {
	return NewServer(
		WithAddress(":8000"),
		WithHub(h),
	)
}

func WithHub(h *hub.Hub) Config {
	return func(server *WSServer) {
		server.hub = h
	}
}
