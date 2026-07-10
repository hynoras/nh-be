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

// getOAuthCookieAttrs returns the Secure flag and SameSite mode appropriate for
// the current environment. OAuth cookies must use SameSite=None + Secure=true in
// production because the login popup is served from the backend domain while the
// parent window lives on the frontend domain (cross-site).
func getOAuthCookieAttrs() (secure bool, sameSite http.SameSite) {
	if env.GetEnvOrDefault("APP_ENV", "dev") == "prod" {
		return true, http.SameSiteNoneMode
	}
	return false, http.SameSiteLaxMode
}

// GetOAuthStateCookie returns an http.Cookie that stores the CSRF state value
// for an OAuth flow. It expires after 10 minutes.
func GetOAuthStateCookie(state string) *http.Cookie {
	secure, sameSite := getOAuthCookieAttrs()
	return &http.Cookie{
		Name:     constant.OAuthStateCookieName,
		Value:    state,
		Path:     constant.OAuthStateCookiePath,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   constant.OAuthStateCookieMaxAge,
	}
}

// GetOAuthVerifierCookie returns an http.Cookie that stores the PKCE verifier
// for an OAuth flow. It expires after 10 minutes.
func GetOAuthVerifierCookie(verifier string) *http.Cookie {
	secure, sameSite := getOAuthCookieAttrs()
	return &http.Cookie{
		Name:     constant.OAuthVerifierCookieName,
		Value:    verifier,
		Path:     constant.OAuthVerifierCookiePath,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   constant.OAuthVerifierCookieMaxAge,
	}
}

// GetClearOAuthCookies returns two expired cookies that instruct the browser to
// delete the oauth_state and oauth_verifier cookies.
func GetClearOAuthCookies() []*http.Cookie {
	secure, sameSite := getOAuthCookieAttrs()
	clear := func(name string) *http.Cookie {
		return &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   secure,
			SameSite: sameSite,
			MaxAge:   -1,
		}
	}
	return []*http.Cookie{clear("oauth_state"), clear("oauth_verifier")}
}
