package bot

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/amikai/line-go-prac/config"
	"github.com/amikai/line-go-prac/pkg/linebotkit"
	"github.com/amikai/line-go-prac/pkg/logkit"
)

func newBotCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "starts line bot server",
		RunE:  runBot,
	}
}

func runBot(_ *cobra.Command, _ []string) error {
	ctx := context.Background()

	conf, err := config.LoadConf()
	if err != nil {
		log.Fatal("failed to load config file", err.Error())
	}

	var logger *logkit.Logger
	{
		loggerConfig := &logkit.LoggerConfig{
			Level:       logkit.LoggerLevel(conf.Logger.Level),
			Development: conf.Logger.Developement,
		}
		logger = logkit.NewLogger(loggerConfig)
	}
	defer func() {
		_ = logger.Sync()
	}()
	ctx = logger.WithContext(ctx)

	var _ *linebotkit.LinebotClient
	{
		linebotConfig := &linebotkit.LinebotConfig{
			ChannelSecret: conf.Linebot.Channel.Secret,
			ChannelToken:  conf.Linebot.Channel.Token,
		}
		_ = linebotkit.NewClient(ctx, linebotConfig)
	}

	// TODO: create gin server to reply message from line
	return nil
}
