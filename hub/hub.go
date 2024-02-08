package hub

import (
	"log"
	"websocket-chat/client"
)

type Hub struct {
	clients        map[*client.Client]bool
	register       chan *client.Client
	unregisterChan chan *client.Client
	broadcastChan  chan []byte
}

func (h *Hub) Register(c *client.Client) {
	h.register <- c
}

func NewHub() *Hub {
	return &Hub{
		clients:        make(map[*client.Client]bool),
		register:       make(chan *client.Client),
		unregisterChan: make(chan *client.Client),
		broadcastChan:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			log.Printf("Registering client in the hub")
			h.clients[c] = true
		case c := <-h.unregisterChan:
			log.Printf("Unregistering client")
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.SendChannel())
			}
		case m := <-h.broadcastChan:
			log.Printf("Broadcasting message %s", m)
			for c := range h.clients {
				c.SendMessage(m)
			}
		}
	}
}

func (h *Hub) UnregisterChan() chan *client.Client {
	return h.unregisterChan
}

func (h *Hub) BroadcastChan() chan []byte {
	return h.broadcastChan
}
