package utils

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var IsJWT = validation.NewStringRule(isJWTToken, "must be a valid JWT token")

func isJWTToken(value string) bool {
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return false
	}
	for _, part := range parts {
		if len(part) == 0 {
			return false
		}
	}
	return true
}
