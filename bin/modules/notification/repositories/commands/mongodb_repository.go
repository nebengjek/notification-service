package commands

import (
	"context"

	user "notification-service/bin/modules/notification"
	"notification-service/bin/modules/notification/models"
	"notification-service/bin/pkg/databases/mongodb"
	"notification-service/bin/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type commandMongodbRepository struct {
	mongoDb mongodb.MongoDBLogger
}

func NewCommandMongodbRepository(mongodb mongodb.MongoDBLogger) user.MongodbRepositoryCommand {
	return &commandMongodbRepository{
		mongoDb: mongodb,
	}
}

func (c commandMongodbRepository) NewObjectID(ctx context.Context) string {
	return primitive.NewObjectID().Hex()
}

func (c commandMongodbRepository) InsertNotification(ctx context.Context, data models.Notification) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		err := c.mongoDb.InsertOne(mongodb.InsertOne{
			CollectionName: "notification",
			Document:       data,
		}, ctx)

		if err != nil {
			output <- utils.Result{
				Error: err,
			}
			return
		}

		output <- utils.Result{
			Data: data.ID,
		}
	}()

	return output
}
