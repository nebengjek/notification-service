package usecases

import (
	"context"
	"fmt"
	"time"

	driver "notification-service/bin/modules/notification"
	"notification-service/bin/modules/notification/models"
	httpError "notification-service/bin/pkg/http-error"
	kafkaPkgConfluent "notification-service/bin/pkg/kafka/confluent"
	"notification-service/bin/pkg/log"
	"notification-service/bin/pkg/utils"

	"github.com/redis/go-redis/v9"
)

type commandUsecase struct {
	driverRepositoryQuery   driver.MongodbRepositoryQuery
	driverRepositoryCommand driver.MongodbRepositoryCommand
	redisClient             redis.UniversalClient
	kafkaProducer           kafkaPkgConfluent.Producer
}

func NewCommandUsecase(mq driver.MongodbRepositoryQuery, mc driver.MongodbRepositoryCommand, rc redis.UniversalClient, kp kafkaPkgConfluent.Producer) driver.UsecaseCommand {
	return &commandUsecase{
		driverRepositoryQuery:   mq,
		driverRepositoryCommand: mc,
		redisClient:             rc,
		kafkaProducer:           kp,
	}
}

func (c *commandUsecase) SendNotification(ctx context.Context, payload models.TripOrder) error {
	message := ""
	switch payload.Status {
	case "request-pickup":
		message = "Ride has been requested. A driver will be assigned shortly."
	case "ontheway":
		message = "Your driver is on the way to pick you up."
	case "completed":
		message = "Your trip has been completed. Thank you for using our service!"
	default:
		message = "Unknown status received. Please check the request."
	}

	driverInfo := <-c.driverRepositoryQuery.FindDriver(ctx, payload.DriverID)
	if driverInfo.Error != nil {
		errObj := httpError.BadRequest("Profile Driver not completed")
		return errObj
	}

	key := fmt.Sprintf("USER:ROUTE:%s", payload.PassengerID)
	newTTL := time.Duration(60) * time.Minute
	redisErr := c.redisClient.Expire(ctx, key, newTTL).Err()
	if redisErr != nil {
		log.GetLogger().Error("command_usecase", fmt.Sprintf("Error updating TTL in Redis: %v", redisErr), "UpdateTTL", utils.ConvertString(redisErr))
		errObj := httpError.InternalServerError(fmt.Sprintf("Error updating TTL: %v", redisErr))
		return errObj
	}

	driver, _ := driverInfo.Data.(models.User)
	notification := models.Notification{
		RecipientID: payload.PassengerID,
		Title:       "Trip Update",
		Message:     message,
		Type:        "info",
		Timestamp:   time.Now(),
		Data: map[string]string{
			"orderId":      payload.OrderID,
			"driverName":   driver.FullName,
			"mobileNumber": driver.MobileNumber,
		},
		Status: "sent",
	}
	sendNotif := <-c.driverRepositoryCommand.InsertNotification(ctx, notification)
	if sendNotif.Error != nil {
		errObj := httpError.BadRequest("send notif failed")
		return errObj
	}
	return nil
}
