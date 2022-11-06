package linebot

import (
	"github.com/amikai/line-go-prac/internal/dao"
)

type GetUserByIDRequest struct {
	Count int   `form:"count"`
	After int64 `form:"after"`
}

type GetUserByIDResponse struct {
	Messages []*Message
}

type Message struct {
	MessageID string `json:"message_id"`
	Text      string `json:"text"`
	CreatedAt int64  `json:"created_at"`
}

func (m *Message) FromMessageDAO(daoMessage *dao.Message) {
	m.MessageID = daoMessage.MessageID
	m.Text = daoMessage.Text
	m.CreatedAt = daoMessage.CreatedAt.Unix()
}
