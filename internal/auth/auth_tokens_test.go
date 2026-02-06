package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		headerName  string
		headerValue string
		expected    string
		wantErr     bool
	}{
		{
			name:        "Valid bearer token",
			headerName:  "Authorization",
			headerValue: "Bearer FooBarBazz",
			expected:    "FooBarBazz",
			wantErr:     false,
		},
		{
			name:        "Valid bearer token with extra spaces",
			headerName:  "Authorization",
			headerValue: "Bearer     FooBarBazz    ",
			expected:    "FooBarBazz",
			wantErr:     false,
		},
		{
			name:        "Bearer in wrong header",
			headerName:  "Bearer",
			headerValue: "Bearer FooBarBazz",
			expected:    "",
			wantErr:     true,
		},
		{
			name:        "Bare token",
			headerName:  "Authorization",
			headerValue: "FooBarBazz",
			expected:    "",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := http.Header{}
			header.Add(tt.headerName, tt.headerValue)
			actual, err := GetBearerToken(header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if actual != tt.expected {
				t.Errorf("GetBearerToken() expected: %v, got: %v", tt.expected, actual)
			}
		})
	}
}

func TestGetApiKey(t *testing.T) {
	tests := []struct {
		name        string
		headerName  string
		headerValue string
		expected    string
		wantErr     bool
	}{
		{
			name:        "Valid ApiKey token",
			headerName:  "Authorization",
			headerValue: "ApiKey FooBarBazz",
			expected:    "FooBarBazz",
			wantErr:     false,
		},
		{
			name:        "Valid ApiKey with extra spaces",
			headerName:  "Authorization",
			headerValue: "ApiKey     FooBarBazz    ",
			expected:    "FooBarBazz",
			wantErr:     false,
		},
		{
			name:        "Bearer in wrong header",
			headerName:  "ApiKey",
			headerValue: "ApiKey FooBarBazz",
			expected:    "",
			wantErr:     true,
		},
		{
			name:        "Bare token",
			headerName:  "Authorization",
			headerValue: "FooBarBazz",
			expected:    "",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := http.Header{}
			header.Add(tt.headerName, tt.headerValue)
			actual, err := GetAPIKey(header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAPIKey() error = %v, wantErr %v", err, tt.wantErr)
			}
			if actual != tt.expected {
				t.Errorf("GetAPIKey() expected: %v, got: %v", tt.expected, actual)
			}
		})
	}
}
