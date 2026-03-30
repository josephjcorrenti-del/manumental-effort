package realtime

import (
	"log"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	hub        *Hub
	conn       *websocket.Conn
	userID     primitive.ObjectID
	send       chan ServerEvent
	subscribed map[primitive.ObjectID]struct{}
}

func NewClient(hub *Hub, conn *websocket.Conn, userID primitive.ObjectID) *Client {
	return &Client{
		hub:        hub,
		conn:       conn,
		userID:     userID,
		send:       make(chan ServerEvent, 16),
		subscribed: make(map[primitive.ObjectID]struct{}),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		_ = c.conn.Close()
	}()

	for {
		var event ClientEvent
		if err := c.conn.ReadJSON(&event); err != nil {
			break
		}

		switch event.Type {
		case EventTypeSubscribe:
			c.hub.Subscribe(c, event.ChannelID)
		case EventTypeUnsubscribe:
			c.hub.Unsubscribe(c, event.ChannelID)
		default:
			c.sendError("unsupported event type")
		}
	}
}

func (c *Client) WritePump() {
	defer func() {
		_ = c.conn.Close()
	}()

	for event := range c.send {
		if err := c.conn.WriteJSON(event); err != nil {
			log.Printf("write websocket event: %v", err)
			return
		}
	}
}

func (c *Client) sendError(message string) {
	select {
	case c.send <- ServerEvent{
		Type:  EventTypeError,
		Error: message,
	}:
	default:
	}
}
