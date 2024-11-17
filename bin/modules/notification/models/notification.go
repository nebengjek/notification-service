package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           string `json:"_id" bson:"_id"`
	FullName     string `json:"fullName" bson:"fullName" validate:"required,min=3,max=100"`
	MobileNumber string `json:"mobileNumber" bson:"mobileNumber" validate:"required"`
	Completed    bool   `'json:"completed" bson:"completed"`
}

type TripOrder struct {
	OrderID       string    `json:"orderId" bson:"orderId"`
	PassengerID   string    `json:"passengerId" bson:"passengerId"`
	DriverID      string    `json:"driverId,omitempty" bson:"driverId,omitempty"`
	Origin        Location  `json:"origin" bson:"origin"`
	Destination   Location  `json:"destination" bson:"destination"`
	Status        string    `json:"status" bson:"status"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
	EstimatedFare float64   `json:"estimatedFare" bson:"estimatedFare"`
	DistanceKm    float64   `json:"distanceKm" bson:"distanceKm"`
}

type Location struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
	Address   string  `json:"address" bson:"address"`
}

type Notification struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	RecipientID string             `json:"recipientId" bson:"recipientId"`
	Title       string             `json:"title" bson:"title"`
	Message     string             `json:"message" bson:"message"`
	Type        string             `json:"type" bson:"type"`
	Timestamp   time.Time          `json:"timestamp" bson:"timestamp"`
	Data        map[string]string  `json:"data" bson:"data,omitempty"`
	Status      string             `json:"status" bson:"status"`
}
