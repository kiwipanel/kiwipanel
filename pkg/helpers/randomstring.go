package helpers

import (
	"math/rand"
	"time"
)

func GenerateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789" //Can add uppercase if needed
	rand.NewSource(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
