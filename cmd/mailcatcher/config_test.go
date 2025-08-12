package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfiguration_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      Configuration
		expectError bool
	}{
		{
			name: "valid configuration",
			config: Configuration{
				MC_PORT:          1025,
				MC_HOST:          "test.example.com",
				MC_REDIRECT_TO:   "test@example.com",
				MC_SENDER_MAIL:   "sender@example.com",
				MC_SMTP_HOST:     "smtp.example.com",
				MC_SMTP_PORT:     587,
				MC_SMTP_USER:     "user",
				MC_SMTP_PASSWORD: "password",
			},
			expectError: false,
		},
		{
			name: "invalid port - negative",
			config: Configuration{
				MC_PORT:          -1,
				MC_HOST:          "test.example.com",
				MC_REDIRECT_TO:   "test@example.com",
				MC_SENDER_MAIL:   "sender@example.com",
				MC_SMTP_HOST:     "smtp.example.com",
				MC_SMTP_PORT:     587,
				MC_SMTP_USER:     "user",
				MC_SMTP_PASSWORD: "password",
			},
			expectError: true,
		},
		{
			name: "invalid port - too high",
			config: Configuration{
				MC_PORT:          70000,
				MC_HOST:          "test.example.com",
				MC_REDIRECT_TO:   "test@example.com",
				MC_SENDER_MAIL:   "sender@example.com",
				MC_SMTP_HOST:     "smtp.example.com",
				MC_SMTP_PORT:     587,
				MC_SMTP_USER:     "user",
				MC_SMTP_PASSWORD: "password",
			},
			expectError: true,
		},
		{
			name: "invalid hostname",
			config: Configuration{
				MC_PORT:          1025,
				MC_HOST:          "invalid hostname with spaces",
				MC_REDIRECT_TO:   "test@example.com",
				MC_SENDER_MAIL:   "sender@example.com",
				MC_SMTP_HOST:     "smtp.example.com",
				MC_SMTP_PORT:     587,
				MC_SMTP_USER:     "user",
				MC_SMTP_PASSWORD: "password",
			},
			expectError: true,
		},
		{
			name: "invalid email - redirect_to",
			config: Configuration{
				MC_PORT:          1025,
				MC_HOST:          "test.example.com",
				MC_REDIRECT_TO:   "invalid-email",
				MC_SENDER_MAIL:   "sender@example.com",
				MC_SMTP_HOST:     "smtp.example.com",
				MC_SMTP_PORT:     587,
				MC_SMTP_USER:     "user",
				MC_SMTP_PASSWORD: "password",
			},
			expectError: true,
		},
		{
			name: "invalid email - sender_mail",
			config: Configuration{
				MC_PORT:          1025,
				MC_HOST:          "test.example.com",
				MC_REDIRECT_TO:   "test@example.com",
				MC_SENDER_MAIL:   "invalid-email",
				MC_SMTP_HOST:     "smtp.example.com",
				MC_SMTP_PORT:     587,
				MC_SMTP_USER:     "user",
				MC_SMTP_PASSWORD: "password",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfiguration_String(t *testing.T) {
	config := Configuration{
		MC_PORT:          1025,
		MC_HOST:          "test.example.com",
		MC_REDIRECT_TO:   "test@example.com",
		MC_SENDER_MAIL:   "sender@example.com",
		MC_SMTP_HOST:     "smtp.example.com",
		MC_SMTP_PORT:     587,
		MC_SMTP_USER:     "testuser",
		MC_SMTP_PASSWORD: "testpass",
	}

	result := config.String()

	// Verify all fields are present in the string representation
	assert.Contains(t, result, "1025")
	assert.Contains(t, result, "test.example.com")
	assert.Contains(t, result, "test@example.com")
	assert.Contains(t, result, "sender@example.com")
	assert.Contains(t, result, "smtp.example.com")
	assert.Contains(t, result, "587")
	assert.Contains(t, result, "testuser")

	// Verify password is not exposed
	assert.NotContains(t, result, "testpass")
}
