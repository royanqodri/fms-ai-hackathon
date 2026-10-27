package util

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// SSEClient represents a single SSE client
type SSEClient struct {
	Channel chan string
	Level   string
}

// SSEServer holds the list of clients
type SSEServer struct {
	Clients map[*SSEClient]bool
	Add     chan *SSEClient
	Remove  chan *SSEClient
	Message chan Message
}

// Message struct to include message content and level
type Message struct {
	Content string
	Level   string
}

// NewSSEServer initializes a new SSE server
func NewSSEServer() *SSEServer {
	return &SSEServer{
		Clients: make(map[*SSEClient]bool),
		Add:     make(chan *SSEClient),
		Remove:  make(chan *SSEClient),
		Message: make(chan Message),
	}
}

// NewSSEClient initializes a new SSE client
func NewSSEClient(level string) *SSEClient {
	return &SSEClient{
		Channel: make(chan string),
		Level:   level,
	}
}

// Run starts the SSE server
func (server *SSEServer) Run() {
	semaphore := make(chan struct{}, 10) // Limit to 10 concurrent goroutines
	var wg sync.WaitGroup

	for {
		select {
		case client := <-server.Add:
			server.Clients[client] = true
		case client := <-server.Remove:
			delete(server.Clients, client)
			close(client.Channel)
		case message := <-server.Message:
			for client := range server.Clients {
				if client.Level == message.Level {
					client.Channel <- message.Content

					// Relay message to other clients with the same level
					wg.Add(1)
					semaphore <- struct{}{} // Acquire a semaphore slot
					go func(client *SSEClient, msg Message) {
						defer wg.Done()
						defer func() { <-semaphore }() // Release the semaphore slot
					}(client, message)
				}
			}
			wg.Wait() // Wait for all relaying goroutines to finish before processing the next message
		}
	}
}

// SSEHandler returns a gin.HandlerFunc that initializes a new SSE client and streams messages
func SSEHandler(server *SSEServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		level := c.Query("level")
		client := &SSEClient{Channel: make(chan string), Level: level}

		// Add client to SSE server
		server.Add <- client
		defer func() {
			server.Remove <- client
		}()

		// Set headers for SSE
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		// Stream messages to client
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			http.Error(c.Writer, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		for {
			select {
			case msg, ok := <-client.Channel:
				if !ok {
					return
				}
				fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
				flusher.Flush()
			case <-c.Writer.CloseNotify():
				server.Remove <- client
				return
			}
		}
	}
}

// SendMessageHandler returns a gin.HandlerFunc that sends messages to clients
func SendMessageHandler(server *SSEServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json struct {
			Message string `json:"message"`
			Level   string `json:"level"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		server.Message <- Message{Content: json.Message, Level: json.Level}
		c.JSON(http.StatusOK, gin.H{"status": "message sent"})
	}
}
