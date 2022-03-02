package client

import (
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Clients struct {
	mut     sync.Mutex
	clients map[string]*Client
}

type Client struct {
	ID   uuid.UUID
	Addr string
}

func NewClients() *Clients {
	return &Clients{
		clients: make(map[string]*Client),
	}
}

func (c *Clients) Get(ip string) *Client {
	c.mut.Lock()
	defer c.mut.Unlock()
	client, ok := c.clients[ip]
	if ok {
		return client
	} else {
		client := &Client{
			ID:   uuid.New(),
			Addr: ip,
		}
		c.clients[ip] = client
		return client
	}
}

var clients *Clients

func init() {
	clients = NewClients()
}

func ClientFromRequest(r *http.Request) *Client {
	var ip string
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip = strings.Split(xff, ",")[0] // might need to handle multiple IP addresses in the header. Likey from load balancer
	} else {
		ip = strings.Split(r.RemoteAddr, ":")[0] // Fallback to RemoteAddr if X-Forwarded-For is not set
	}
	return clients.Get(ip)
}
