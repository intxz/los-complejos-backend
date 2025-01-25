// complejo.go
package models

// Complejo represents a user in the system with optional fitness-related attributes.
type Complejo struct {
	ID       string `json:"_id" bson:"_id" validate:"required"`           // Unique identifier
	Username string `json:"username" bson:"username" validate:"required"` // User's username (required)
	Password string `json:"password" bson:"password" validate:"required"` // User's password (required)
	Role     string `json:"role" bson:"role" validate:"required"`         // Role of the user (e.g., "user" or "admin") (required)
	Weight   string `json:"weight" bson:"weight"`                         // Weight in kilograms (optional)
	Height   string `json:"height" bson:"height"`                         // Height in meters (optional)
	IMC      string `json:"imc" bson:"imc"`                               // Calculated IMC based on weight and height
	Gender   string `json:"gender" bson:"gender" validate:"required"`     // User's gender (required)
	Bench    string `json:"bench" bson:"bench"`                           // Bench press weight in kilograms (optional)
	Squad    string `json:"squad" bson:"squad"`                           // Squat weight in kilograms (optional)
	DL       string `json:"dl" bson:"dl"`                                 // Deadlift weight in kilograms (optional)
	Photo    string `json:"photo" bson:"photo"`                           // Base64-encoded profile photo (optional)
}
