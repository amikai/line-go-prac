package mongokit

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"

	"github.com/amikai/line-go-prac/pkg/logkit"
)

type MongoConfig struct {
	URL      string
	Database string
}

type MongoClient struct {
	*mongo.Client
	database  *mongo.Database
	closeFunc func()
}

func (c *MongoClient) Database() *mongo.Database {
	return c.database
}

func (c *MongoClient) Close() error {
	if c.closeFunc != nil {
		c.closeFunc()
	}

	return c.Client.Disconnect(context.Background())
}

func NewMongoClient(ctx context.Context, conf *MongoConfig) *MongoClient {
	logger := logkit.FromContext(ctx).With(
		zap.String("url", conf.URL),
		zap.String("database", conf.Database),
	)

	o := options.Client()
	o.ApplyURI(conf.URL)

	client, err := mongo.NewClient(o)
	if err != nil {
		logger.Fatal("failed to create MongoDB client", zap.Error(err))
	}

	if err := client.Connect(ctx); err != nil {
		logger.Fatal("failed to connect to MongoDB", zap.Error(err))
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Fatal("failed to ping to MongoDB", zap.Error(err))
	}

	logger.Info("create MongoDB client successfully")

	return &MongoClient{
		Client:   client,
		database: client.Database(conf.Database),
	}
}
