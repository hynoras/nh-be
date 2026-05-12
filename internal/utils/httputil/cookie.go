package httputil

import (
	"net/http"
	"nh-be/internal/constant"
)

// GetAuthSessionCookie returns a new http.Cookie configured for authentication sessions
// with the provided session ID.
func GetAuthSessionCookie(sessionId string) *http.Cookie {
	return &http.Cookie{
		Name:     constant.AuthSessionCookieName,
		Value:    sessionId,
		Path:     constant.AuthSessionCookiePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   constant.AuthSessionCookieMaxAge,
	}
}
