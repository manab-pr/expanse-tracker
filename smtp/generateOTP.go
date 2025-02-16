package smtp

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateOTP(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate OTP : %v", err)
	}
	return hex.EncodeToString(bytes)[:length], nil
}
