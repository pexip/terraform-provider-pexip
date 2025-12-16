package helpers

import (
	"os"
)

func GetEnvStringOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
