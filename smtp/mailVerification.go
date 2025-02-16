package smtp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Sender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ToEmail struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type EmailRequest struct {
	Sender      Sender    `json:"sender"`
	To          []ToEmail `json:"to"`
	Subject     string    `json:"subject"`
	HTMLContent string    `json:"htmlContent"`
}

var brevoUrl string = os.Getenv("BREVO_URL")

func SendEmail(apiKey string, toEmail string, subject string, htmlContent string) error {
	emailRequest := EmailRequest{
		Sender: Sender{
			Name:  "Expanse Tracker",
			Email: "manab2001maity@gmail.com",
		},
		To: []ToEmail{
			{
				Email: toEmail,
			},
		},
		Subject:     subject,
		HTMLContent: htmlContent,
	}

	requestBody, err := json.Marshal(emailRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %v", err)
	}

	log.Printf("Sending email to %s with subject: %s", toEmail, subject)

	req, err := http.NewRequest("POST", brevoUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Brevo API response status: %s", resp.Status)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send email, status code: %d", resp.StatusCode)
	}

	return nil
}
