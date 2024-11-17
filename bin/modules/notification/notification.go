package notification

import (
	"context"

	"notification-service/bin/modules/notification/models"
	"notification-service/bin/pkg/utils"
)

type UsecaseQuery interface {
}

type UsecaseCommand interface {
	SendNotification(ctx context.Context, payload models.TripOrder) error
}

type MongodbRepositoryQuery interface {
	FindDriver(ctx context.Context, userId string) <-chan utils.Result
}

type MongodbRepositoryCommand interface {
	NewObjectID(ctx context.Context) string
	InsertNotification(ctx context.Context, data models.Notification) <-chan utils.Result
}
