package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	return getAuthToken(headers, "Bearer")
}

func GetAPIKey(headers http.Header) (string, error) {
	return getAuthToken(headers, "ApiKey")
}

func getAuthToken(headers http.Header, keyType string) (string, error) {
	prefix := strings.TrimSpace(keyType) + " "
	auth := headers.Get("Authorization")
	if !strings.HasPrefix(auth, prefix) {
		return "", fmt.Errorf("No %sin headers", prefix)
	}
	return strings.TrimSpace(strings.TrimPrefix(auth, prefix)), nil
}
