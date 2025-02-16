package handlers

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"expanse-tracker/constants"
	"expanse-tracker/db"
	"expanse-tracker/models"
	"expanse-tracker/smtp"
)

func ForgotPasswordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Email string `json:"email"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Check if the user exists in the database
		userCollection := db.OpenCollection(db.Client, "user")
		filter := bson.M{"email": request.Email}
		var user models.User
		err := userCollection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
			}
			return
		}

		// Generate a password reset token
		resetToken, err := smtp.GenerateOTP(32) // Use a longer token for password reset
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
			return
		}

		// Store the reset token in MongoDB
		resetDoc := models.PasswordReset{
			Email:     request.Email,
			Token:     resetToken,
			ExpiresAt: time.Now().Add(15 * time.Minute), // Token expires in 15 minutes
		}

		resetCollection := db.OpenCollection(db.Client, "password_resets")
		_, err = resetCollection.InsertOne(context.Background(), resetDoc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store reset token"})
			return
		}

		// Send the password reset email
		resetLink := fmt.Sprintf("https://expanse-tracker/reset-password?token=%s", resetToken)
		subject := "Password Reset Request"

		// Render the email template
		tmpl, err := template.New("email").Parse(constants.EmailTemplateFP)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse email template"})
			return
		}

		var htmlContent bytes.Buffer
		data := struct {
			OTP string
		}{
			OTP: resetLink, // Use the reset link as the "OTP" in the template
		}

		if err := tmpl.Execute(&htmlContent, data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render email template"})
			return
		}

		err = smtp.SendEmail(os.Getenv("BREVO_API_KEY"), request.Email, subject, htmlContent.String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send password reset email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent successfully"})
	}
}
