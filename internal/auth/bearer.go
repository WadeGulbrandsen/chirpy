package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", errors.New("No Bearer token in headers")
	}
	return strings.TrimSpace(strings.TrimPrefix(auth, "Bearer ")), nil
}
