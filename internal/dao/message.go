package dao

import (
	"context"
	"time"
)

type Message struct {
	MessageID string    `bson:"message_id"`
	SenderID  string    `bson:"sender_id"`
	Text      string    `bson:"text"`
	CreatedAt time.Time `bson:"created_at"`
}

type MessagePagination struct {
	Count int
	After time.Time
}

type MessageDAO interface {
	GetByUserID(ctx context.Context, userID string, pagination MessagePagination) ([]*Message, error)
	Store(ctx context.Context, message *Message) error
}
