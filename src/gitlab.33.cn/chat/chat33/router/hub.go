package router

import (
	"github.com/inconshreveable/log15"
)

var hubModelLog = log15.New("module", "chat33/router/hub")

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[Client]bool

	// Inbound messages from the clients.
	broadcast chan interface{}

	// Register requests from the clients.
	register chan Client

	// Unregister requests from clients.
	unregister chan Client
}

func NewHub() *Hub {
	hub := &Hub{
		broadcast:  make(chan interface{}),
		register:   make(chan Client),
		unregister: make(chan Client),
		clients:    make(map[Client]bool),
	}
	go hub.Run()
	return hub
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			delete(h.clients, client)
		case message, ok := <-h.broadcast:
			if !ok {
				hubModelLog.Error("channel closed")
				break
			}
			for client := range h.clients {
				err := client.Send(message)
				if err != nil {
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) IsExist(key Client) bool {
	if value, ok := h.clients[key]; ok && value {
		return true
	} else {
		return false
	}
}
