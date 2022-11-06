package linebotkit

import (
	"context"

	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"

	"github.com/amikai/line-go-prac/pkg/logkit"
)

type LinebotConfig struct {
	ChannelSecret string
	ChannelToken  string
}

type LinebotClient struct {
	*linebot.Client
}

func NewClient(ctx context.Context, conf *LinebotConfig) *LinebotClient {
	logger := logkit.FromContext(ctx)

	bot, err := linebot.New(conf.ChannelSecret, conf.ChannelToken)
	if err != nil {
		logger.Fatal("failed to create Linebot client", zap.Error(err))
	}

	logger.Info("create Linebot client successfully")
	return &LinebotClient{
		Client: bot,
	}
}
