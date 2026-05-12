package constant

const (
	// AuthSessionCookieName is the name of the cookie used for authentication sessions.
	AuthSessionCookieName = "auth_session"
	// AuthSessionCookiePath is the URL path that the authentication session cookie is valid for.
	AuthSessionCookiePath = "/"
	// AuthSessionCookieMaxAge is the maximum age of the authentication session cookie in seconds (8 hours).
	AuthSessionCookieMaxAge = 8 * 60 * 60
)
