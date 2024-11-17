package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	notification "notification-service/bin/modules/notification"
	"notification-service/bin/modules/notification/models"
	kafkaPkgConfluent "notification-service/bin/pkg/kafka/confluent"
	"notification-service/bin/pkg/log"

	k "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type notificationHandler struct {
	driverUsecaseCommand notification.UsecaseCommand
}

func NewNotificationConsumer(dc notification.UsecaseCommand) kafkaPkgConfluent.ConsumerHandler {
	return &notificationHandler{
		driverUsecaseCommand: dc,
	}
}

func (i notificationHandler) HandleMessage(message *k.Message) {
	log.GetLogger().Info("consumer", fmt.Sprintf("Partition: %v - Offset: %v", message.TopicPartition.Partition, message.TopicPartition.Offset.String()), *message.TopicPartition.Topic, string(message.Value))

	var msg models.TripOrder
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.GetLogger().Error("consumer", "unmarshal-data", err.Error(), string(message.Value))
		return
	}

	if err := i.driverUsecaseCommand.SendNotification(context.Background(), msg); err != nil {
		log.GetLogger().Error("consumer", "BroadcastPickupnotification", err.Error(), string(message.Value))
		return
	}

	return
}
