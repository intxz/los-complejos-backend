// complejo_handler.go
package handlers

import (
	"los-complejos-backend/models"
	"los-complejos-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateComplejo creates a new Complejo and inserts it into the MongoDB collection.
//
// This function accepts a JSON payload to create a new Complejo document. It generates a unique ID for the Complejo,
// calculates its IMC (Body Mass Index) based on the weight and height provided, and generates a JWT token for authentication.
//
// HTTP Status Codes:
// - 201 Created: The Complejo was successfully created.
// - 400 Bad Request: Invalid JSON data was provided.
// - 500 Internal Server Error: There was an issue inserting the Complejo into the database or generating the token.
//
// Parameters:
// - collection (*mongo.Collection): The MongoDB collection where the Complejo documents are stored.
//
// Example JSON payload for creating a Complejo:
//
//	{
//	    "username": "test_user",
//	    "password": "securepassword",
//	    "role": "user",
//	    "weight": "75.5",
//	    "height": "1.78",
//	    "gender": "male",
//	    "bench": "100",
//	    "squad": "140",
//	    "dl": "180",
//	    "photo": "base64_encoded_photo"
//	}
//
// Example usage:
// r.POST("/complejo", CreateComplejo(collection))
func CreateComplejo(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var complejo models.Complejo

		// Parse the incoming JSON request into the Complejo model
		if err := c.ShouldBindJSON(&complejo); err != nil {
			// 400 Bad Request: The JSON is invalid
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"code":    http.StatusBadRequest,
				"message": "Invalid JSON format: " + err.Error(),
			})
			return
		}

		// Generate a unique ID and calculate the IMC
		complejo.ID = uuid.NewString()
		complejo.IMC = utils.CalcIMC(complejo.Weight, complejo.Height)

		// Prepare the document for MongoDB insertion
		document := bson.M{
			"_id":      complejo.ID,
			"username": complejo.Username,
			"password": complejo.Password,
			"role":     complejo.Role,
			"weight":   complejo.Weight,
			"height":   complejo.Height,
			"imc":      complejo.IMC,
			"gender":   complejo.Gender,
			"bench":    complejo.Bench,
			"squad":    complejo.Squad,
			"dl":       complejo.DL,
			"photo":    complejo.Photo,
		}

		// Insert the document into the MongoDB collection
		_, err := collection.InsertOne(c, document)
		if err != nil {
			// 500 Internal Server Error: Failed to insert the document
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Failed to create Complejo: " + err.Error(),
			})
			return
		}

		// Generate a token for the user (infinite or long-lived)
		token, err := utils.GenerateToken(complejo.ID, complejo.Role, complejo.Username)
		if err != nil {
			// 500 Internal Server Error: Failed to generate the token
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Failed to generate token: " + err.Error(),
			})
			return
		}

		// 201 Created: The Complejo was successfully created
		c.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"code":    http.StatusCreated,
			"message": "Complejo created successfully",
			"data":    complejo,
			"token":   token,
		})
	}
}

// GetComplejos retrieves all Complejos from the MongoDB collection.
//
// This function fetches all Complejo documents from the MongoDB collection. If no Complejos are found, it responds with a 404 status.
//
// HTTP Status Codes:
// - 200 OK: Successfully retrieved all Complejos.
// - 404 Not Found: No Complejos were found in the database.
// - 500 Internal Server Error: An issue occurred while fetching or processing the data.
//
// Parameters:
// - collection (*mongo.Collection): The MongoDB collection where the Complejo documents are stored.
//
// Example usage:
// r.GET("/complejo", GetComplejos(collection))
func GetComplejos(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Find all documents in the collection
		cursor, err := collection.Find(c, bson.M{})
		if err != nil {
			// 500 Internal Server Error: Database query failed
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Failed to fetch Complejos from the database: " + err.Error(),
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

		// Parse the cursor results into a slice of Complejos
		var complejos []models.Complejo
		if err := cursor.All(c, &complejos); err != nil {
			// 500 Internal Server Error: Failed to parse data
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Failed to parse Complejos data: " + err.Error(),
			})
			return
		}

		// Handle the case where no Complejos are found
		if len(complejos) == 0 {
			// 404 Not Found: No Complejos exist
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"code":    http.StatusNotFound,
				"message": "No Complejos found in the database",
			})
			return
		}

		// 200 OK: Successfully retrieved all Complejos
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"code":    http.StatusOK,
			"message": "Complejos retrieved successfully",
			"data":    complejos,
		})
	}
}

// GetComplejo retrieves a single Complejo by ID from the MongoDB collection.
//
// This function fetches a single Complejo document using its unique `_id`.
// If the document is not found, it responds with a 404 status.
//
// HTTP Status Codes:
// - 200 OK: Successfully retrieved the Complejo.
// - 404 Not Found: The Complejo with the specified ID was not found.
// - 500 Internal Server Error: Failed to fetch or process the Complejo.
//
// Parameters:
// - collection (*mongo.Collection): The MongoDB collection where the Complejo documents are stored.
//
// Example usage:
// r.GET("/complejo/:id", GetComplejo(collection))
func GetComplejo(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Find the document in the collection by "_id"
		var complejo models.Complejo
		err := collection.FindOne(c, bson.M{"_id": id}).Decode(&complejo)
		if err != nil {
			// 404 Not Found: Document not found
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  "error",
					"code":    http.StatusNotFound,
					"message": "Complejo not found",
				})
				return
			}
			// 500 Internal Server Error: Query error
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Failed to retrieve Complejo: " + err.Error(),
			})
			return
		}

		// 200 OK: Successfully retrieved the Complejo
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"code":    http.StatusOK,
			"message": "Complejo retrieved successfully",
			"data":    complejo,
		})
	}
}

// UpdateComplejoForUser updates specific fields of a Complejo, restricted to user role.
//
// This function allows users with the "user" role to update specific personal fields in their Complejo document.
// Only the fields listed as "allowed" are updated, and any invalid or unauthorized fields are ignored.
//
// HTTP Status Codes:
// - 200 OK: Successfully updated the Complejo.
// - 400 Bad Request: Invalid JSON data was provided or no valid fields were included in the payload.
// - 404 Not Found: The Complejo with the specified ID was not found or the role is not "user".
// - 500 Internal Server Error: An issue occurred while updating the Complejo in the database.
//
// Parameters:
// - collection (*mongo.Collection): The MongoDB collection where the Complejo documents are stored.
//
// Example JSON payload for updating a Complejo:
//
//	{
//	    "weight": 80.5,
//	    "height": 1.75,
//	    "bench": 120.0
//	}
//
// Example usage:
// r.PUT("/complejo/user", UpdateComplejoForUser(collection))
func UpdateComplejoForUser(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, idExist := c.Get("_id")
		role, roleExist := c.Get("role")

		if !idExist || !roleExist || role == "user" {
			// 403 Forbidden: Insufficient permissions
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"code":    http.StatusForbidden,
				"message": "You do not have permission to update this Complejo.",
			})
			return
		}

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

		allowedFields := []string{"username", "weight", "height", "bench", "squad", "deadlift", "photo"}
		filteredUpdate := bson.M{}
		for _, field := range allowedFields {
			if value, exists := updateData[field]; exists {
				filteredUpdate[field] = value
			}
		}

		// Ensure no invalid fields were sent
		if len(filteredUpdate) == 0 {
			// 400 Bad Request: No valid fields provided
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"code":    http.StatusBadRequest,
				"message": "No valid fields to update",
			})
			return
		}

		// Prepare the update payload
		update := bson.M{"$set": filteredUpdate}

		// Perform the update operation
		result, err := collection.UpdateOne(c, bson.M{"_id": id, "role": "user"}, update)
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
			// 404 Not Found: Document with the given ID does not exist or is not a user
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"code":    http.StatusNotFound,
				"message": "Complejo not found or insufficient permissions",
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

// UpdateComplejoForAdmin updates specific fields of a Complejo by ID, restricted to admin role.
//
// This function allows administrators with the "admin" role to update any field of a Complejo document.
// Unlike user updates, admin updates have no restrictions on the fields that can be modified.
//
// HTTP Status Codes:
// - 200 OK: Successfully updated the Complejo.
// - 400 Bad Request: Invalid JSON data was provided.
// - 403 Forbidden: The user does not have sufficient permissions to perform this action.
// - 404 Not Found: The Complejo with the specified ID was not found.
// - 500 Internal Server Error: An issue occurred while updating the Complejo in the database.
//
// Parameters:
// - collection (*mongo.Collection): The MongoDB collection where the Complejo documents are stored.
//
// Example JSON payload for updating a Complejo:
//
//	{
//	    "username": "admin_updated_user",
//	    "role": "superadmin",
//	    "weight": 85.0
//	}
//
// Example usage:
// r.PUT("/complejos/admin", UpdateComplejoForAdmin(collection))
func UpdateComplejoForAdmin(collection *mongo.Collection) gin.HandlerFunc {
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
