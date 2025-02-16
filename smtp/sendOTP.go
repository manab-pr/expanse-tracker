package smtp

import (
	"bytes"
	"context"
	"expanse-tracker/db"
	"expanse-tracker/models"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"
)

const emailTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>OTP Email</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f4f4f4;
      margin: 0;
      padding: 0;
    }
    .email-container {
      max-width: 600px;
      margin: 0 auto;
      background-color: #ffffff;
      border-radius: 8px;
      overflow: hidden;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
    .header img {
      width: 100%;
      height: auto;
    }
    .content {
      padding: 20px;
      text-align: center;
    }
    .content h1 {
      color: #333333;
      font-size: 24px;
      margin-bottom: 20px;
    }
    .otp {
      font-size: 32px;
      font-weight: bold;
      color: #007bff;
      margin: 20px 0;
    }
    .footer {
      background-color: #f8f9fa;
      padding: 10px;
      text-align: center;
      font-size: 14px;
      color: #666666;
    }
    .footer a {
      color: #007bff;
      text-decoration: none;
    }
  </style>
</head>
<body>
  <div class="email-container">
    <div class="header">
      <img src="https://images.pexels.com/photos/1761279/pexels-photo-1761279.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1" alt="Header Image">
    </div>
    <div class="content">
      <h1>Your One-Time Password (OTP)</h1>
      <p>Please use the following OTP to verify your account:</p>
      <div class="otp">{{.OTP}}</div>
      <p>This OTP is valid for <strong>5 minutes</strong>.</p>
    </div>
    <div class="footer">
      <p>If you did not request this OTP, please ignore this email.</p>
      <p>© 2025 Expanse Tracker. All rights reserved.</p>
      <p><a href="https://expanse-tracker.com">Visit our website</a></p>
    </div>
  </div>
</body>
</html>
`

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
	tmpl, err := template.New("email").Parse(emailTemplate)
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
