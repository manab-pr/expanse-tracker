package smtp

import (
	"bytes"
	"context"
	"expanse-tracker/constants"
	"expanse-tracker/db"
	"expanse-tracker/models"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"
)

func SendOTP(toEmail string) (string, error) {
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
		log.Fatal("BREVO_API_KEY is not set")
	}

	otp, err := GenerateOTP(6)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %v", err)
	}

	subject := "Your OTP for verification"
	tmpl, err := template.New("email").Parse(constants.EmailTemplateOTP)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %v", err)
	}
	var htmlContent bytes.Buffer
	data := struct {
		OTP string
	}{
		OTP: otp,
	}

	if err := tmpl.Execute(&htmlContent, data); err != nil {
		return "", fmt.Errorf("failed to execute email template: %v", err)
	}

	log.Printf("Sending OTP to %s: %s", toEmail, otp)

	err = SendEmail(apiKey, toEmail, subject, htmlContent.String())
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return "", fmt.Errorf("failed to send email: %v", err)
	}

	// Store OTP in MongoDB
	otpDoc := models.OTP{
		Email:     toEmail,
		OTP:       otp,
		ExpiresAt: time.Now().Add(5 * time.Minute), // OTP expires in 5 minutes
	}

	collection := db.OpenCollection(db.Client, "otps")
	_, err = collection.InsertOne(context.Background(), otpDoc)
	if err != nil {
		log.Printf("Failed to store OTP in MongoDB: %v", err)
		return "", fmt.Errorf("failed to store OTP in MongoDB: %v", err)
	}

	log.Println("OTP sent and stored successfully")
	return otp, nil
}
