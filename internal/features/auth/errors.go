package auth

import (
	"errors"
	"net/http"
	"nh-be/internal/utils/httputil"
)

const (
	ErrVerifyTokenFailed    = "Failed to verify token"
	ErrSignUpFailed         = "Failed to sign up"
	ErrLoginFailed          = "Failed to login"
	ErrLogoutFailed         = "Failed to logout"
	ErrChangePasswordFailed = "Failed to change password"
)

var (
	ErrInvalidCredentials                      = errors.New("invalid credentials")
	ErrVerificationTokenNotFound               = errors.New("token not found")
	ErrUnauthenticated                         = errors.New("unauthenticated")
	ErrEmailAlreadyExists                      = errors.New("email already exists")
	ErrVerificationTokenExpired                = errors.New("token expired")
	ErrInvalidVerificationToken                = errors.New("invalid token, please check token type")
	ErrSessionNotFound                         = errors.New("session not found")
	ErrNewPasswordAndConfirmPasswordDoNotMatch = errors.New("new password and confirm password do not match")
	ErrNewPasswordIsTheSameAsOldPassword       = errors.New("new password is the same as the old password")
)

func init() {
	httputil.RegisterError(ErrInvalidCredentials, http.StatusUnauthorized, "Invalid credentials")
	httputil.RegisterError(ErrVerificationTokenNotFound, http.StatusNotFound, "Verification token not found")
	httputil.RegisterError(ErrUnauthenticated, http.StatusUnauthorized, "Unauthenticated")
	httputil.RegisterError(ErrEmailAlreadyExists, http.StatusConflict, "Email already exists")
	httputil.RegisterError(ErrVerificationTokenExpired, http.StatusUnauthorized, "Verification token expired")
	httputil.RegisterError(ErrInvalidVerificationToken, http.StatusUnauthorized, "Invalid verification token")
	httputil.RegisterError(ErrSessionNotFound, http.StatusUnauthorized, "Session not found")
	httputil.RegisterError(ErrNewPasswordAndConfirmPasswordDoNotMatch, http.StatusBadRequest, "Passwords do not match")
	httputil.RegisterError(ErrNewPasswordIsTheSameAsOldPassword, http.StatusBadRequest, "Password must be different")
}
