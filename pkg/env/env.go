package env

import (
	"log/slog"
	"os"
	"strconv"
)

func MustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		slog.Error("missing required env var", "key", key)
		panic("missing required env var: " + key)
	}
	return v
}

func MustEnvInt(key string) int {
	v := os.Getenv(key)
	if v == "" {
		slog.Error("missing required env var", "key", key)
		panic("missing required env var: " + key)
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		slog.Error("invalid integer value for env var", "key", key, "error", err)
		panic("invalid integer value for env var: " + key)
	}
	return i
}
