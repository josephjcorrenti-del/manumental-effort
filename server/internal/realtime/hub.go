package realtime

import (
	"context"
	"log"
	"sync"

	"manumental-effort/server/internal/channels"
	"manumental-effort/server/internal/memberships"
	"manumental-effort/server/internal/messages"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hub struct {
	channelRepository    *channels.Repository
	membershipRepository *memberships.Repository

	mu            sync.RWMutex
	clients       map[*Client]struct{}
	subscriptions map[primitive.ObjectID]map[*Client]struct{}
}

func NewHub(
	channelRepository *channels.Repository,
	membershipRepository *memberships.Repository,
) *Hub {
	return &Hub{
		channelRepository:    channelRepository,
		membershipRepository: membershipRepository,
		clients:              make(map[*Client]struct{}),
		subscriptions:        make(map[primitive.ObjectID]map[*Client]struct{}),
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = struct{}{}
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, client)

	for channelID := range client.subscribed {
		if subscribers, ok := h.subscriptions[channelID]; ok {
			delete(subscribers, client)
			if len(subscribers) == 0 {
				delete(h.subscriptions, channelID)
			}
		}
	}

	close(client.send)
}

func (h *Hub) Subscribe(client *Client, channelIDHex string) {
	channelID, err := primitive.ObjectIDFromHex(channelIDHex)
	if err != nil {
		client.sendError("invalid channel id")
		return
	}

	channel, err := h.channelRepository.GetByID(context.Background(), channelID)
	if err != nil {
		if err == channels.ErrChannelNotFound {
			client.sendError("channel not found")
			return
		}
		client.sendError("failed to load channel")
		return
	}

	isMember, err := h.membershipRepository.Exists(context.Background(), channel.SpaceID, client.userID)
	if err != nil {
		client.sendError("failed to verify membership")
		return
	}
	if !isMember {
		client.sendError("membership required")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.subscriptions[channelID]; !ok {
		h.subscriptions[channelID] = make(map[*Client]struct{})
	}

	h.subscriptions[channelID][client] = struct{}{}
	client.subscribed[channelID] = struct{}{}
}

func (h *Hub) Unsubscribe(client *Client, channelIDHex string) {
	channelID, err := primitive.ObjectIDFromHex(channelIDHex)
	if err != nil {
		client.sendError("invalid channel id")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if subscribers, ok := h.subscriptions[channelID]; ok {
		delete(subscribers, client)
		if len(subscribers) == 0 {
			delete(h.subscriptions, channelID)
		}
	}

	delete(client.subscribed, channelID)
}

func (h *Hub) BroadcastMessageCreated(message *messages.Message) {
	h.mu.RLock()
	subscribers := h.subscriptions[message.ChannelID]
	h.mu.RUnlock()

	if len(subscribers) == 0 {
		return
	}

	event := ServerEvent{
		Type:      EventTypeMessageCreated,
		ChannelID: message.ChannelID.Hex(),
		Message:   message,
	}

	for client := range subscribers {
		select {
		case client.send <- event:
		default:
			log.Printf("dropping websocket event for user=%s channel=%s", client.userID.Hex(), message.ChannelID.Hex())
		}
	}
}
