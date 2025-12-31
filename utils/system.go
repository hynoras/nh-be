package utils

import (
	"log"
	"os"
	"strconv"
)

func MustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return v
}

func MustEnvInt(key string) int {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid integer value for env var %s: %v", key, err)
	}
	return i
}
