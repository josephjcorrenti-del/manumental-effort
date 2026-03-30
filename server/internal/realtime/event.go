package realtime

import "manumental-effort/server/internal/messages"

const (
	EventTypeSubscribe      = "subscribe"
	EventTypeUnsubscribe    = "unsubscribe"
	EventTypeMessageCreated = "message_created"
	EventTypeError          = "error"
)

type ClientEvent struct {
	Type      string `json:"type"`
	ChannelID string `json:"channel_id,omitempty"`
}

type ServerEvent struct {
	Type      string            `json:"type"`
	ChannelID string            `json:"channel_id,omitempty"`
	Message   *messages.Message `json:"message,omitempty"`
	Error     string            `json:"error,omitempty"`
}
