package utils

import (
	"github.com/google/uuid"
)

func GenerateRandomString(n int) string {
	return uuid.NewString()
}
