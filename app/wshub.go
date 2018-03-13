package app

import (
	"bytes"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = pongWait / 2
	maxMessageSize = 512
)

type wsHub struct {
	mutex      sync.Mutex
	clients    map[*wsClient]struct{}
	broadcast  chan []byte
	register   chan *wsClient
	unregister chan *wsClient
	getPeers   chan struct{}
	peers      chan []*wsClient
}

type wsClient struct {
	hub  *wsHub
	conn *websocket.Conn
	send chan []byte
}

func newHub() *wsHub {
	h := &wsHub{
		clients:    make(map[*wsClient]struct{}),
		broadcast:  make(chan []byte),
		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
		getPeers:   make(chan struct{}, 1),
		peers:      make(chan []*wsClient, 1),
	}
	return h
}

func (h *wsHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		case <-h.getPeers:
			res := make([]*wsClient, len(h.clients))
			for client := range h.clients {
				res = append(res, client)
			}
			h.peers <- res
		}
	}
}

func (h *wsHub) Broadcast(msg string) {
	h.broadcast <- []byte(msg)
}

func (h *wsHub) Connect(client *wsClient) {
	h.register <- client
	h.Broadcast("+peer: " + client.conn.RemoteAddr().String())
}

func (h *wsHub) Disconnect(client *wsClient) {
	h.unregister <- client
	h.Broadcast("-peer: " + client.conn.RemoteAddr().String())
}

func (h *wsHub) GetPeers() []*wsClient {
	h.getPeers <- struct{}{}
	return <-h.peers
}

func (c *wsClient) readPump() {
	defer func() {
		c.hub.Disconnect(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		c.hub.broadcast <- message
	}
}

func (c *wsClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
