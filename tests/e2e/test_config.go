package e2e

import (
	"os"
	"strconv"
)

type TestConfig struct {
	MailpitImage     string
	MailcatcherImage string
	TestDomain       string
	TestTimeout      int
}

func LoadTestConfig() *TestConfig {
	config := &TestConfig{
		MailpitImage:     getEnvOrDefault("MAILPIT_IMAGE", "axllent/mailpit:latest"),
		MailcatcherImage: getEnvOrDefault("MAILCATCHER_IMAGE", ""), // Use local build if empty
		TestDomain:       getEnvOrDefault("TEST_DOMAIN", "test.example.com"),
		TestTimeout:      getEnvIntOrDefault("TEST_TIMEOUT", 30),
	}
	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
