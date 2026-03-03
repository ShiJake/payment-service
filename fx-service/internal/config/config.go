package config

import (
	"os"
)

// GetPort returns the server port from environment or default
func GetPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "4000"
}

