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
	"expanse-tracker/smtp"
)

func SendOTPHandler(c *gin.Context) {

	var request struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	_, err := smtp.SendOTP(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func VerifyOTPHandler(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Retrieve and verify OTP from MongoDB
	collection := db.OpenCollection(db.Client, "otps")
	filter := bson.M{
		"email": request.Email,
		"otp":   request.OTP,
		"expires_at": bson.M{
			"$gt": time.Now(), // Check if OTP is not expired
		},
	}

	var otpDoc models.OTP
	err := collection.FindOne(context.Background(), filter).Decode(&otpDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify OTP"})
		}
		return
	}

	// Delete the OTP after successful verification
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}
