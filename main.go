package main

import (
	"websocket-chat/hub"
	"websocket-chat/server"
)

func main() {
	h := hub.NewHub()
	go h.Run()
	wsServer := server.NewDefaultServer(h)
	wsServer.Run()
	select {}
}
