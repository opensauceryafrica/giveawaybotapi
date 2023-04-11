package config

import (
	"os"
)

func MustGet(key string, def string) string {
	value := os.Getenv(key)
	if def != "" && value == "" {
		return def
	} else if def == "" && value == "" {
		panic(key + ": not found in env")
	}
	return value
}
