package models

import "time"

// Event represents the structure of an event in the system
type Event struct {
	ID           string    `json:"_id" bson:"_id"`                                     // Unique identifier for the event
	Title        string    `json:"title" bson:"title" validate:"required"`             // Title of the event (required)
	Description  string    `json:"description" bson:"description" validate:"required"` // Description of the event (required)
	Participants []string  `json:"participants" bson:"participants" default:"[]"`      // List of participants (default: empty)
	Date         time.Time `json:"date" bson:"date" validate:"required"`               // Date of the event (required)
	Image        *string   `json:"image,omitempty" bson:"image,omitempty"`             // Optional image URL for the event
	Location     string    `json:"location" bson:"location" validate:"required"`       // Location of the event (required)
}
