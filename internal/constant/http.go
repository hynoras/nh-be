package constant

import "time"

const (
	// AuthSessionCookieName is the name of the cookie used for authentication sessions.
	AuthSessionCookieName = "auth_session"
	// AuthSessionCookiePath is the URL path that the authentication session cookie is valid for.
	AuthSessionCookiePath = "/"
	// AuthSessionCookieMaxAge is the maximum age of the authentication session cookie in seconds (8 hours).
	AuthSessionCookieMaxAge = 8 * 60 * 60
	// CSRFTokenCookieName is the name of the cookie used for CSRF protection.
	CSRFTokenCookieName = "csrf_token"
	// CSRFTokenCookiePath is the URL path that the CSRF token cookie is valid for.
	CSRFTokenCookiePath = "/"
	// CSRFTokenCookieMaxAge is the maximum age of the CSRF token cookie in seconds (8 hours).
	CSRFTokenCookieMaxAge = 8 * 60 * 60
)

var (
	CorsAllowMethods     = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	CorsAllowHeaders     = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	CorsExposeHeaders    = []string{"Content-Length"}
	CorsAllowCredentials = true
	CorsMaxAge           = 12 * time.Hour
)
