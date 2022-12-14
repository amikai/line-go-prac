package bot

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/amikai/line-go-prac/config"
	"github.com/amikai/line-go-prac/internal/dao"
	"github.com/amikai/line-go-prac/internal/linebot"
	"github.com/amikai/line-go-prac/pkg/ginkit"
	"github.com/amikai/line-go-prac/pkg/linebotkit"
	"github.com/amikai/line-go-prac/pkg/logkit"
	"github.com/amikai/line-go-prac/pkg/mongokit"
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

	var mongoClient *mongokit.MongoClient
	{
		mongoConfig := &mongokit.MongoConfig{
			URL:      conf.Mongo.URL,
			Database: conf.Mongo.Database,
		}
		mongoClient = mongokit.NewMongoClient(ctx, mongoConfig)
	}

	messageDAO := dao.NewMongoMessageDAO(mongoClient)

	var linebotClient *linebotkit.LinebotClient
	{
		linebotConfig := &linebotkit.LinebotConfig{
			ChannelSecret: conf.Linebot.Channel.Secret,
			ChannelToken:  conf.Linebot.Channel.Token,
		}
		linebotClient = linebotkit.NewClient(ctx, linebotConfig)
	}

	linebotService := linebot.NewEchoService(messageDAO, linebotClient)
	linebotGinHandler := linebot.NewGinHandler(linebotClient, linebotService, logger)

	router := ginkit.Default(logger)
	router.GET("/healthz", ginkit.HeartbeatHandler)
	linebotGroup := router.Group("/linebot")
	{
		linebotGroup.POST("/webhook", linebotGinHandler.Webhook)
		linebotGroup.GET("/user/:userID", linebotGinHandler.GetUserByID)
	}

	// gracefully shutdown. See https://pkg.go.dev/net/http#Server.Shutdown
	srv := http.Server{
		Addr:    ":9999",
		Handler: router,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Error("Http server shutdown", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("HTTP server ListenAndServe", zap.Error(err))
	}

	<-idleConnsClosed

	return nil
}
