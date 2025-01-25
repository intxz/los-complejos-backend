// event_handler.go
package handlers

import (
	"fmt"
	"los-complejos-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateEvent allows only admin users to create a new event and insert it into the MongoDB collection.
//
// This function:
// 1. Validates the user's role to ensure they are an admin.
// 2. Parses the incoming JSON payload to create a new Event document.
// 3. Inserts the Event into the MongoDB collection.
//
// HTTP Status Codes:
// - 201 Created: The Event was successfully created.
// - 400 Bad Request: Invalid JSON data was provided.
// - 403 Forbidden: The user does not have sufficient permissions to create an event.
// - 500 Internal Server Error: An issue occurred while inserting the Event into the database.
//
// Example JSON payload:
//
//	{
//	    "title": "Gym Meetup",
//	    "description": "A gathering of fitness enthusiasts.",
//	    "date": "2025-02-01T10:00:00Z",
//	    "location": "Local Gym, Main Street"
//	}
//
// Example usage:
// r.POST("/event", CreateEvent(collection))
func CreateEvent(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the role from the context (set by the JWT middleware)
		role, exists := c.Get("role")
		if !exists {
			// Log a message if the token is missing or invalid
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"code":    http.StatusUnauthorized,
				"message": "Authorization token is missing or invalid",
			})
			return
		}

		// Debug: Log the role extracted from the token
		fmt.Println("Token validated successfully. Role:", role)

		if role != "admin" {
			// 403 Forbidden: Insufficient permissions
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"code":    http.StatusForbidden,
				"message": "You do not have permission to create events.",
			})
			return
		}

		// Parse the incoming JSON request into the Event model
		var event models.Event
		if err := c.ShouldBindJSON(&event); err != nil {
			// 400 Bad Request: Invalid JSON format
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"code":    http.StatusBadRequest,
				"message": "Invalid JSON format: " + err.Error(),
			})
			return
		}

		// Generate a unique ID for the event
		event.ID = uuid.NewString()
		document := bson.M{
			"_id":          event.ID,
			"title":        event.Title,
			"description":  event.Description,
			"participants": event.Participants,
			"date":         event.Date,
			"image":        event.Image,
			"location":     event.Location,
		}

		// Insert the event into the MongoDB collection
		_, err := collection.InsertOne(c, document)
		if err != nil {
			// 500 Internal Server Error: Database insertion failed
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Failed to create event: " + err.Error(),
			})
			return
		}

		// 201 Created: The Event was successfully created
		c.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"code":    http.StatusCreated,
			"message": "Event created successfully",
			"data":    event,
		})
	}
}

// GetEvents retrieves all Event documents from the MongoDB collection.
//
// This function fetches all Event documents from the MongoDB collection.
// If no Events are found, it responds with a 404 status.
//
// HTTP Status Codes:
// - 200 OK: Successfully retrieved all Events.
// - 404 Not Found: No Events were found in the database.
// - 500 Internal Server Error: An issue occurred while fetching or processing the data.
//
// Parameters:
// - collection (*mongo.Collection): The MongoDB collection where the Event documents are stored.
//
// Example usage:
// r.GET("/events", GetEvents(collection))

// GetEvent retrieves a single Event by ID from the MongoDB collection.
//
// This function fetches a single Event document using its unique `_id`.
// If the document is not found, it responds with a 404 status.
//
// HTTP Status Codes:
// - 200 OK: Successfully retrieved the Event.
// - 404 Not Found: The Event with the specified ID was not found.
// - 500 Internal Server Error: Failed to fetch or process the Event.
//
// Parameters:
// - collection (*mongo.Collection): The MongoDB collection where the Event documents are stored.
//
// Example usage:
// r.GET("/event/:id", GetEvent(collection))
func GetEvents(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Find all documents in the collection
		cursor, err := collection.Find(c, bson.M{})
		if err != nil {
			// 500 Internal Server Error: Database query failed
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Failed to fetch Event from the database: " + err.Error(),
			})
			return
		}
		defer func() {
			if err := cursor.Close(c); err != nil {
				// 500 Internal Server Error: Failed to close the cursor
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"code":    http.StatusInternalServerError,
					"message": "Failed to close the database cursor: " + err.Error(),
				})
			}
		}()

		// Parse the cursor results into a slice of Event
		var events []models.Event
		if err := cursor.All(c, &events); err != nil {
			// 500 Internal Server Error: Failed to parse data
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Failed to parse Events data: " + err.Error(),
			})
			return
		}

		// Handle the case where no Event are found
		if len(events) == 0 {
			// 404 Not Found: No Event exist
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"code":    http.StatusNotFound,
				"message": "No Event found in the database",
			})
			return
		}

		// 200 OK: Successfully retrieved all Event
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"code":    http.StatusOK,
			"message": "Event retrieved successfully",
			"data":    events,
		})
	}
}

// GetEvent retrieves a single Event by ID from the MongoDB collection.
//
// This function fetches a single Event document using its unique `_id`.
// If the document is not found, it responds with a 404 status.
//
// HTTP Status Codes:
// - 200 OK: Successfully retrieved the Event.
// - 404 Not Found: The Event with the specified ID was not found.
// - 500 Internal Server Error: Failed to fetch or process the Event.
//
// Parameters:
// - collection (*mongo.Collection): The MongoDB collection where the Event documents are stored.
//
// Example usage:
// r.GET("/event/:id", GetEvent(collection))
func GetEvent(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Find the document in the collection by "_id"
		var event models.Event
		err := collection.FindOne(c, bson.M{"_id": id}).Decode(&event)
		if err != nil {
			// 404 Not Found: Document not found
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  "error",
					"code":    http.StatusNotFound,
					"message": "Event not found",
				})
				return
			}
			// 500 Internal Server Error: Query error
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Failed to retrieve Event: " + err.Error(),
			})
			return
		}

		// 200 OK: Successfully retrieved the Event
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"code":    http.StatusOK,
			"message": "Event retrieved successfully",
			"data":    event,
		})
	}
}

// UpdateEventForAdmin updates specific fields of an Event by ID, restricted to admin role.
//
// This function allows administrators with the "admin" role to update any field of an Event document.
// Unlike user updates, admin updates have no restrictions on the fields that can be modified.
//
// HTTP Status Codes:
// - 200 OK: Successfully updated the Event.
// - 400 Bad Request: Invalid JSON data was provided.
// - 403 Forbidden: The user does not have sufficient permissions to perform this action.
// - 404 Not Found: The Event with the specified ID was not found.
// - 500 Internal Server Error: An issue occurred while updating the Event in the database.
//
// Parameters:
// - collection (*mongo.Collection): The MongoDB collection where the Event documents are stored.
//
// Example JSON payload for updating an Event:
//
//	{
//	    "title": "Updated Gym Meetup",
//	    "location": "Updated Location"
//	}
//
// Example usage:
// r.PUT("/event/admin", UpdateEventForAdmin(collection))
func UpdateEventForAdmin(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Retrieve the id and role from the context (set by the JWT middleware)
		role, roleExists := c.Get("role")
		id, idExist := c.Get("_id")
		if !roleExists || role != "admin" && !idExist || id != "_id" {
			// 403 Forbidden: Insufficient permissions
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"code":    http.StatusForbidden,
				"message": "You do not have permission to update this Complejo.",
			})
			return
		}

		// Parse the incoming JSON to a map for flexible updates
		var updateData map[string]interface{}
		if err := c.ShouldBindJSON(&updateData); err != nil {
			// 400 Bad Request: Invalid JSON format
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"code":    http.StatusBadRequest,
				"message": "Invalid JSON format: " + err.Error(),
			})
			return
		}

		// Remove `_id` to avoid overwriting the document ID
		delete(updateData, "_id")

		// Prepare the update payload
		update := bson.M{"$set": updateData}

		// Perform the update operation
		result, err := collection.UpdateOne(c, bson.M{"_id": id}, update)
		if err != nil {
			// 500 Internal Server Error: Database update failed
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Failed to update Complejo: " + err.Error(),
			})
			return
		}

		// Handle the case where no document was updated
		if result.MatchedCount == 0 {
			// 404 Not Found: Document with the given ID does not exist
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"code":    http.StatusNotFound,
				"message": "Complejo not found",
			})
			return
		}

		// 200 OK: Successfully updated the Complejo
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"code":    http.StatusOK,
			"message": "Complejo updated successfully",
		})
	}
}

// SubscribeEvent allows a user to subscribe to an Event by adding their username to the Event's participants.
//
// This function:
// 1. Extracts the username from the JWT token.
// 2. Adds the username to the Event's participants list using MongoDB's `$addToSet` operator.
//
// HTTP Status Codes:
// - 200 OK: Successfully subscribed to the Event.
// - 403 Forbidden: The user does not have a valid username.
// - 404 Not Found: The Event with the specified ID was not found.
// - 409 Conflict: The user is already subscribed to the Event.
// - 500 Internal Server Error: An issue occurred while subscribing to the Event.
//
// Parameters:
// - collection (*mongo.Collection): The MongoDB collection where the Event documents are stored.
//
// Example usage:
// r.PUT("/event/:id/subscribe", SubscribeEvent(collection))
func SubscribeEvent(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID := c.Param("_id")
		username, exist := c.Get("username")
		if !exist || username == "username" {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"code":    http.StatusForbidden,
				"message": "You do not have a valid username.",
			})
			return
		}

		update := bson.M{
			"$addToSet": bson.M{"participants": username},
		}

		result, err := collection.UpdateOne(c, bson.M{"_id": eventID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to subscribe to the event: " + err.Error(),
			})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Event not found",
			})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Complejo is already subscribed to the event.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Successfully subscribed to the event",
		})
	}
}

// UnsuscribeEvent allows a user to unsubscribe from an Event by removing their username from the Event's participants.
//
// This function:
// 1. Extracts the username from the JWT token.
// 2. Removes the username from the Event's participants list using MongoDB's `$pull` operator.
//
// HTTP Status Codes:
// - 200 OK: Successfully unsubscribed from the Event.
// - 403 Forbidden: The user does not have a valid username.
// - 404 Not Found: The Event with the specified ID was not found or the user is not subscribed.
// - 409 Conflict: The user is not subscribed to the Event.
// - 500 Internal Server Error: An issue occurred while unsubscribing from the Event.
//
// Parameters:
// - collection (*mongo.Collection): The MongoDB collection where the Event documents are stored.
//
// Example usage:
// r.PUT("/event/:id/unsubscribe", UnsuscribeEvent(collection))
func UnsuscribeEvent(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID := c.Param("_id")
		username, exist := c.Get("username")
		if !exist || username == "username" {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"code":    http.StatusForbidden,
				"message": "You do not have a valid username.",
			})
			return
		}

		update := bson.M{
			"$pull": bson.M{
				"participants": username,
			},
		}

		result, err := collection.UpdateOne(c, bson.M{"_id": eventID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to unsubscribe from event: " + err.Error(),
			})
			return
		}

		// Comprobar si se encontró y actualizó el documento
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Event not found or user not subscribed",
			})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Complejo is not already subscribed to the event.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Successfully unsubscribed from event",
		})
	}
}
