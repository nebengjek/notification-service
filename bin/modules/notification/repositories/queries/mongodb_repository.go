package queries

import (
	"context"
	driver "notification-service/bin/modules/notification"
	"notification-service/bin/modules/notification/models"
	"notification-service/bin/pkg/databases/mongodb"
	"notification-service/bin/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type queryMongodbRepository struct {
	mongoDb mongodb.MongoDBLogger
}

func NewQueryMongodbRepository(mongodb mongodb.MongoDBLogger) driver.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
	}
}

func (q queryMongodbRepository) FindDriver(ctx context.Context, userId string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		objId, _ := primitive.ObjectIDFromHex(userId)

		var driver models.User
		err := q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &driver,
			CollectionName: "user",
			Filter: bson.M{
				"_id": objId,
			},
		}, ctx)
		if err != nil {
			output <- utils.Result{
				Error: err,
			}
		}

		output <- utils.Result{
			Data: driver,
		}

	}()

	return output
}
