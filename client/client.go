package client

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	connection     *websocket.Conn
	sendChan       chan []byte
	unregisterChan chan *Client
	broadcastChan  chan []byte
}

func NewClient(connection *websocket.Conn, sendChan chan []byte, unregister chan *Client, broadcastChan chan []byte) *Client {
	return &Client{connection: connection, sendChan: sendChan, unregisterChan: unregister, broadcastChan: broadcastChan}
}

func (c *Client) ReadPump() {
	defer func() {
		c.unregisterChan <- c
		c.connection.Close()
	}()

	for {
		_, message, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.broadcastChan <- message
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.connection.Close()
	}()
	for {
		select {
		case message, ok := <-c.sendChan:
			if !ok {
				c.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (c *Client) SendMessage(m []byte) {
	c.sendChan <- m
}

func (c *Client) SendChannel() chan []byte {
	return c.sendChan
}
