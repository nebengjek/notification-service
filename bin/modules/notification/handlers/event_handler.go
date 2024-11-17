package handlers

import (
	driver "notification-service/bin/modules/notification"
	kafkaPkgConfluent "notification-service/bin/pkg/kafka/confluent"
)

func InitNotificationEventHandler(driver driver.UsecaseCommand, kc kafkaPkgConfluent.Consumer) {

	kc.SetHandler(NewNotificationConsumer(driver))
	kc.Subscribe("trip-created")

}
