package utils

import "os"

// EnvOrDefault returns value of environment with `name`
// if the value is empty, returns `fallback`
func EnvOrDefault(name, fallback string) string {
	if envVal := os.Getenv(name); envVal != "" {
		return envVal
	}
	return fallback
}
