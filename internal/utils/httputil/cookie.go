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

// GetEmptyAuthSessionCookie returns a new http.Cookie configured for authentication sessions
// with an empty session ID to clear the cookie.
func GetEmptyAuthSessionCookie() *http.Cookie {
	sameSite := http.SameSiteLaxMode
	secure := false

	// If in production, enable cross-site cookie support
	if env.GetEnvOrDefault("APP_ENV", "dev") == "prod" {
		sameSite = http.SameSiteNoneMode
		secure = true
	}

	return &http.Cookie{
		Name:     constant.AuthSessionCookieName,
		Value:    "",
		Path:     constant.AuthSessionCookiePath,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   0,
	}
}

// GetCSRFTokenCookie returns a new http.Cookie configured for CSRF protection.
func GetCSRFTokenCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     constant.CSRFTokenCookieName,
		Value:    token,
		Path:     constant.CSRFTokenCookiePath,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   constant.CSRFTokenCookieMaxAge,
	}
}
