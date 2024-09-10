package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	Status    int         `json:"status"`
	Error     string      `json:"error"`
	Code      string      `json:"code,omitempty"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// NewErrorResponse creates a new error response with the given parameters.
func NewErrorResponse(status int, message string, code string, details interface{}) []byte {
	// Restrict details to non-sensitive information
	if status == fasthttp.StatusUnauthorized || status == fasthttp.StatusForbidden {
		details = nil // Don't expose sensitive details
	}

	errorResponse := ErrorResponse{
		Status:    status,
		Error:     fasthttp.StatusMessage(status),
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	responseJSON, _ := json.Marshal(errorResponse)
	return responseJSON
}

// GetTokenFromHeader extracts the Bearer token from the Authorization header
func GetTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	return parts[1], nil
}

func GenerateLicenseKey(prefix string, mask string) string {
	key := []int32{}
	for _, char := range mask {
		if char == 'X' {
			key = append(key, randomChar())
		} else if char == '-' {
			key = append(key, '-')
		} else {
			key = append(key, char)
		}
	}
	return fmt.Sprintf("%s-%s", prefix, string(key))
}

func randomChar() int32 {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return int32(chars[rand.Intn(len(chars))])
}

func GetCurrentDatetime() string {
	return time.Now().Format("2006-01-02 @ 03:04 PM")
}

// CalculateExpiryDateFromText calculates the expiry date from a text like "1 days"
func CalculateExpiryDateFromText(durationText string) string {
	layout := "2006-01-02 @ 03:04 PM"
	currentDateTime := GetCurrentDatetime()
	currentTime, err := time.Parse(layout, currentDateTime)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	parts := strings.Split(durationText, " ")
	if len(parts) != 2 {
		fmt.Println("invalid duration format")
		return ""
	}

	duration, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Println(err)
		return ""
	}

	unit := parts[1]
	switch strings.ToLower(unit) {
	case "day", "days":
		currentTime = currentTime.AddDate(0, 0, duration)
	case "week", "weeks":
		currentTime = currentTime.AddDate(0, 0, duration*7)
	case "month", "months":
		currentTime = currentTime.AddDate(0, duration, 0)
	case "year", "years":
		currentTime = currentTime.AddDate(duration, 0, 0)
	case "hour", "hours":
		currentTime = currentTime.Add(time.Duration(duration) * time.Hour)
	case "minute", "minutes":
		currentTime = currentTime.Add(time.Duration(duration) * time.Minute)
	default:
		fmt.Println("unsupported unit")
		return ""
	}

	return currentTime.Format(layout)
}

func GetClientIP(ctx *gin.Context) string {
	ip := ctx.Request.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = ctx.Request.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = ctx.ClientIP()
	}
	return ip
}

func FormatDuration(duration int, unit string) string {
	return fmt.Sprintf("%d %s(s)", duration, unit)
}
