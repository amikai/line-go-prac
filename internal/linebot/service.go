package linebot

import (
	"context"
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/amikai/line-go-prac/internal/dao"
	"github.com/amikai/line-go-prac/pkg/linebotkit"
)

type Service interface {
	GetUserByID(ctx context.Context, userID string, pagination dao.MessagePagination) ([]*dao.Message, error)
	HandleEvent(ctx context.Context, event *linebot.Event) error
}

var _ Service = &echoService{}

type echoService struct {
	messageDAO    dao.MessageDAO
	linebotClient *linebotkit.LinebotClient
}

func NewEchoService(messageDAO dao.MessageDAO, linebotClient *linebotkit.LinebotClient) Service {
	return &echoService{
		messageDAO:    messageDAO,
		linebotClient: linebotClient,
	}
}

func (s *echoService) HandleEvent(ctx context.Context, event *linebot.Event) error {
	switch event.Type {
	case linebot.EventTypeMessage:
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			// store message
			s.messageDAO.Store(ctx, &dao.Message{
				MessageID: message.ID,
				SenderID:  event.Source.UserID,
				Text:      message.Text,
				CreatedAt: event.Timestamp,
			})

			responseMessage := linebot.NewTextMessage(message.Text)
			_, err := s.linebotClient.ReplyMessage(event.ReplyToken, responseMessage).Do()
			if err != nil {
				return fmt.Errorf("failed to reply message: %w", err)
			}
		}
	}
	return nil
}

func (s *echoService) GetUserByID(ctx context.Context, userID string, pagination dao.MessagePagination) ([]*dao.Message, error) {
	return s.messageDAO.GetByUserID(ctx, userID, pagination)
}
