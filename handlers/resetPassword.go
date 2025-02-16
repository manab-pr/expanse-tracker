package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"expanse-tracker/db"
	"expanse-tracker/models"
	"expanse-tracker/utils"
)

func ResetPasswordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		var request struct {
			Token       string `json:"token"`
			NewPassword string `json:"new_password"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Retrieve the reset token from MongoDB
		collection := db.OpenCollection(db.Client, "password_resets")
		filter := bson.M{
			"token": request.Token,
			"expires_at": bson.M{
				"$gt": time.Now(), // Check if the token is not expired
			},
		}

		var resetDoc models.PasswordReset
		err := collection.FindOne(context.Background(), filter).Decode(&resetDoc)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token"})
			}
			return
		}

		// Update the user's password in the database (pseudo-code)
		userCollection := db.OpenCollection(db.Client, "user")
		updateFilter := bson.M{"email": resetDoc.Email}
		update := bson.M{"$set": bson.M{"password": utils.HashPassword(request.NewPassword)}}
		_, err = userCollection.UpdateOne(context.Background(), updateFilter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}

		// Delete the reset token after successful password reset
		_, err = collection.DeleteOne(context.Background(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reset token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
	}
}
