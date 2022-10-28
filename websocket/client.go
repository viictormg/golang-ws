package websocket

import "github.com/gorilla/websocket"

type Client struct {
	hub       *Hub
	id        string
	socket    *websocket.Conn
	outhbound chan []byte
}

func NewCliente(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:       hub,
		socket:    socket,
		outhbound: make(chan []byte),
	}
}

func (c *Client) Write() {
	for {
		select {
		case message, ok := <-c.outhbound:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
