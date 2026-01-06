package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func GenerateResetToken(userID string) (string, time.Time, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", time.Time{}, err
	}

	token := hex.EncodeToString(bytes)
	exp := time.Now().Add(1 * time.Hour)

	return token, exp, nil
}