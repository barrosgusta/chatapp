package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/barrosgusta/chatapp/services/gateway-websocket/sqs"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	name string
	send chan OutgoingMessage
}

type Hub struct {
	clients   map[*Client]struct{}
	names     map[string]*Client
	typing    map[string]bool
	mu        sync.Mutex
	register   chan *Client
	unregister chan *Client
	sqsProducer sqs.Producer
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		names:      make(map[string]*Client),
		typing:     make(map[string]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				h.mu.Lock()
				if client.name != "" {
					delete(h.names, client.name)
					delete(h.typing, client.name)
				}
				h.mu.Unlock()
				h.broadcastUserList()
				h.broadcastTyping()
				close(client.send)
			}
		}
	}
}

func (h *Hub) broadcastUserList() {
	users := make([]string, 0)
	h.mu.Lock()
	for name := range h.names {
		users = append(users, name)
	}
	h.mu.Unlock()
	msg := OutgoingMessage{Type: "USER_LIST", Users: users}
	for c := range h.clients {
		c.send <- msg
	}
}

func (h *Hub) broadcastTyping() {
	h.mu.Lock()
	typingUsers := make([]string, 0)
	for name, isTyping := range h.typing {
		if isTyping {
			typingUsers = append(typingUsers, name)
		}
	}
	h.mu.Unlock()
	msg := OutgoingMessage{Type: "TYPING", Typing: typingUsers}
	for c := range h.clients {
		c.send <- msg
	}
}

func (h *Hub) broadcastMessage(msg *ChatMessage) {
	out := OutgoingMessage{Type: "MESSAGE", Message: msg}
	for c := range h.clients {
		c.send <- out
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Dev only!
}

func wsHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
				log.Printf("WebSocket upgrade failed: %v", err)
				return
		}

		client := &Client{conn: conn, send: make(chan OutgoingMessage, 8)}
		hub.register <- client

		go writePump(client)
		readPump(hub, client)
	}
}

func writePump(c *Client) {
	for msg := range c.send {
		c.conn.WriteJSON(msg)
	}
	c.conn.Close()
}

func readPump(hub *Hub, c *Client) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	   for {
			   var rawMsg json.RawMessage
			   err := c.conn.ReadJSON(&rawMsg)
			   if err != nil {
					   break
			   }
			   log.Printf("Raw incoming: %s", string(rawMsg))

			   var in IncomingMessage
			   if err := json.Unmarshal(rawMsg, &in); err != nil {
					   log.Printf("Failed to unmarshal incoming message: %v", err)
					   continue
			   }
			   switch in.Type {
		case "SET_NAME":
			name := strings.TrimSpace(in.Name)
			hub.mu.Lock()
			if _, exists := hub.names[name]; name == "" || exists {
				c.send <- OutgoingMessage{Type: "NAME_REJECTED", Reason: "Name taken"}
				hub.mu.Unlock()
			} else {
				c.name = name
				hub.names[name] = c
				c.send <- OutgoingMessage{Type: "NAME_ACCEPTED"}
				hub.mu.Unlock()
				hub.broadcastUserList()
			}
		case "MESSAGE":
			if c.name == "" {
				continue
			}
			chatMsg := &ChatMessage{
				// Unique ID for the message
				ID:        time.Now().Format("20060102150405.000") + "-" + c.name,
				User:      c.name,
				Text:      in.Text,
				Timestamp: time.Now().Format(time.RFC3339),
			}
			hub.broadcastMessage(chatMsg)
			
			// Send to SQS
			log.Printf("Sending chat message to SQS: %s", chatMsg.Text)
			msgBytes, err := json.Marshal(chatMsg)
			if err != nil {
				log.Printf("Failed to marshal chat message: %v", err)
				break
			}
			go func(data string) {
				err := hub.sqsProducer.SendMessage(context.Background(), data)
				if err != nil {
					log.Printf("Failed to send message to SQS: %v", err)
				} else {
					log.Printf("Chat message sent to SQS")
				}
			}(string(msgBytes))
		case "TYPING_START":
			if c.name == "" {
				continue
			}
			hub.mu.Lock()
			hub.typing[c.name] = true
			hub.mu.Unlock()
			hub.broadcastTyping()
		case "TYPING_STOP":
			if c.name == "" {
				continue
			}
			hub.mu.Lock()
			hub.typing[c.name] = false
			hub.mu.Unlock()
			hub.broadcastTyping()
		}
	}
}

func main() {
	ctx := context.Background()
	sqsProducer := sqs.NewSQSProducer(ctx)
	hub := NewHub()
	hub.sqsProducer = sqsProducer
	go hub.Run()

	r := chi.NewRouter()
	r.Get("/ws", wsHandler(hub))
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
