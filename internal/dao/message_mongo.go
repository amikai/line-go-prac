package dao

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/amikai/line-go-prac/pkg/mongokit"
)

const MessageCollection string = "message"

const MaximunMessageCount int = 100

type mongoMessageDAO struct {
	client *mongokit.MongoClient
}

var _ MessageDAO = (*mongoMessageDAO)(nil)

func NewMongoMessageDAO(mongoClient *mongokit.MongoClient) *mongoMessageDAO {
	return &mongoMessageDAO{
		client: mongoClient,
	}
}

func (dao *mongoMessageDAO) GetByUserID(ctx context.Context, userID string, pagination MessagePagination) ([]*Message, error) {
	findOpts := options.Find()
	findOpts.SetSort(bson.D{{"created_at", 1}})

	limit := pagination.Count
	if pagination.Count > MaximunMessageCount {
		limit = MaximunMessageCount
	}
	findOpts.SetLimit(int64(limit))

	mongoTime := primitive.NewDateTimeFromTime(pagination.After)
	filter := bson.D{
		{"created_at", bson.D{{"$gt", mongoTime}}},
		{"sender_id", userID},
	}

	cursor, err := dao.client.Database().Collection(MessageCollection).Find(
		ctx, filter, findOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}

	var messages []*Message
	for cursor.Next(ctx) {
		var message Message
		err := cursor.Decode(&message)
		if err != nil {
			return nil, fmt.Errorf("failed to decode message: %w", err)
		}
		messages = append(messages, &message)
	}
	return messages, nil
}

func (dao *mongoMessageDAO) Store(ctx context.Context, message *Message) error {
	if _, err := dao.client.Database().Collection(MessageCollection).InsertOne(ctx, message); err != nil {
		return fmt.Errorf("failed to store message: %w", err)
	}
	return nil
}
