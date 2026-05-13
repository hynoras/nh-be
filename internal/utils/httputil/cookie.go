package httputil

import (
	"net/http"
	"nh-be/internal/constant"
	"nh-be/pkg/env"
)

// GetAuthSessionCookie returns a new http.Cookie configured for authentication sessions
// with the provided session ID.
func GetAuthSessionCookie(sessionId string) *http.Cookie {
	sameSite := http.SameSiteLaxMode
	secure := false

	// If in production, enable cross-site cookie support
	if env.GetEnvOrDefault("APP_ENV", "dev") == "prod" {
		sameSite = http.SameSiteNoneMode
		secure = true
	}

	return &http.Cookie{
		Name:     constant.AuthSessionCookieName,
		Value:    sessionId,
		Path:     constant.AuthSessionCookiePath,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   constant.AuthSessionCookieMaxAge,
	}
}

