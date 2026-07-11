package auth

import (
	"html/template"
	"net/http"
	"nh-be/internal/utils/httputil"

	"github.com/gin-gonic/gin"
)

// VerifyTokenHandler godoc
// @Summary Verify email token
// @Description Verify email token and activate user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token path string true "Token to verify"
// @Success 200 {object} httputil.SuccessResponse "User verified successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid token"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 500 {object} httputil.ErrorResponse "Failed to verify token"
// @Router /auth/verify/{token} [get]
func VerifyTokenHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		sessionId, err := s.VerifyEmail(c.Request.Context(), token)
		if err != nil {
			httputil.MakeErrorResponse(
				c,
				http.StatusUnauthorized,
				ErrVerifyTokenFailed,
				err.Error(),
			)
			return
		}

		http.SetCookie(c.Writer, httputil.GetAuthSessionCookie(sessionId))

		httputil.MakeSuccessResponse(c, http.StatusOK, "User verified successfully", nil)
	}
}

// LoginHandler godoc
// @Summary User login
// @Description Authenticate user with email and password, returns user information with permissions and creates a session
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginDto true "Login credentials (email and password)"
// @Success 200 {object} httputil.SuccessResponse{data=LoginResponseDto} "Successfully logged in"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request format"
// @Failure 401 {object} httputil.ErrorResponse "Invalid email or password"
// @Failure 500 {object} httputil.ErrorResponse "Session save failed"
// @Router /auth/login [post]
func LoginHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginDto
		if err := httputil.ValidateRequestFormat(c, &req); err != nil {
			return
		}

		userRes, sessionId, csrfToken, err := s.Login(c.Request.Context(), req.Email, req.Password)
		if httputil.MakeServiceErrorResponse(c, err, ErrLoginFailed) {
			return
		}

		http.SetCookie(c.Writer, httputil.GetAuthSessionCookie(sessionId))
		http.SetCookie(c.Writer, httputil.GetCSRFTokenCookie(csrfToken))

		httputil.MakeSuccessResponse(c, http.StatusOK, "User logged in successfully", userRes)
	}
}

// @Router /auth/:provider/login [get]
func ProviderLoginHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := c.Param("provider")
		if provider != "google" {
			httputil.MakeErrorResponse(c, http.StatusBadRequest, "Invalid provider", "Provider not supported")
			return
		}

		url, state, verifier, err := s.GenerateProviderLoginURL(c.Request.Context(), provider)
		if err != nil {
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to generate login URL", err.Error())
			return
		}

		// Store state and verifier temporarily in cookies to verify in the callback.
		// Attributes (Secure, SameSite) are set based on APP_ENV.
		http.SetCookie(c.Writer, httputil.GetOAuthStateCookie(state))
		http.SetCookie(c.Writer, httputil.GetOAuthVerifierCookie(verifier))

		// Handler executes the redirect
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// @Router /auth/:provider/callback [get]
func ProviderCallbackHandler(s Service, frontendURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := c.Param("provider")
		code := c.Query("code")
		state := c.Query("state")

		// Validate state
		cookieState, err := c.Cookie("oauth_state")
		fmt.Println("cookieState 2", cookieState)
		fmt.Println("state 2", state)
		if err != nil || state != cookieState {
			httputil.MakeErrorResponse(c, http.StatusBadRequest, "Invalid state parameter", "State validation failed")
			return
		}

		// Get verifier
		cookieVerifier, err := c.Cookie("oauth_verifier")
		if err != nil {
			httputil.MakeErrorResponse(c, http.StatusBadRequest, "Missing verifier", "PKCE verifier cookie not found")
			return
		}

		_, sessionId, csrfToken, err := s.ProviderCallback(c.Request.Context(), provider, code, cookieVerifier)
		if err != nil {
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to complete login", err.Error())
			return
		}

		// Clear OAuth cookies
		for _, cookie := range httputil.GetClearOAuthCookies() {
			http.SetCookie(c.Writer, cookie)
		}

		// Set login cookie
		http.SetCookie(c.Writer, httputil.GetAuthSessionCookie(sessionId))
		http.SetCookie(c.Writer, httputil.GetCSRFTokenCookie(csrfToken))

		tmpl, parseErr := template.ParseFiles("templates/provider_login.html")
		if parseErr != nil {
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to load template", parseErr.Error())
			return
		}

		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/html; charset=utf-8")
		if execErr := tmpl.Execute(c.Writer, map[string]string{"FrontendURL": frontendURL}); execErr != nil {
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to render template", execErr.Error())
			return
		}
	}
}

// LogoutHandler godoc
// @Summary User logout
// @Description Clear user session and log out from the system
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} httputil.SuccessResponse "Successfully logged out"
// @Failure 500 {object} httputil.ErrorResponse "Failed to logout"
// @Security SessionAuth
// @Router /auth/logout [post]
func LogoutHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("auth_session")
		if err != nil {
			httputil.MakeErrorResponse(
				c,
				http.StatusUnauthorized,
				"Unauthorized",
				err,
			)
			return
		}

		serviceErr := s.Logout(c.Request.Context(), cookie.Value)
		if httputil.MakeServiceErrorResponse(c, serviceErr, ErrLogoutFailed) {
			return
		}
		//set empty cookie upon logging out for frontend
		http.SetCookie(c.Writer, httputil.GetEmptyAuthSessionCookie())
		httputil.MakeSuccessResponse(c, http.StatusOK, "User logged out successfully", nil)
	}
}

// ChangePasswordHandler godoc
// @Summary Change user password
// @Description Update user password with new password and confirmation
// @Tags Authentication
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID format)"
// @Param request body ChangePasswordDto true "New password and confirmation"
// @Success 200 {object} httputil.SuccessResponse "Password changed successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid user ID or request body"
// @Security SessionAuth
// @Router /auth/users/{id}/change-password [put]
func ChangePasswordHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := httputil.ParseStringToUUID(c.Param("id"))
		if err != nil {
			httputil.MakeErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
			return
		}

		var req ChangePasswordDto
		valReqErr := httputil.ValidateRequestFormat(c, &req)
		if valReqErr != nil {
			return
		}

		serviceErr := s.ChangePassword(c.Request.Context(), userID, req)
		if httputil.MakeServiceErrorResponse(c, serviceErr, ErrChangePasswordFailed) {
			return
		}
		httputil.MakeSuccessResponse(c, http.StatusOK, "User password changed successfully", nil)
	}
}

// SignUpHandler godoc
// @Summary User sign up
// @Description Register a new user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body SignUpDto true "User registration details"
// @Success 201 {object} httputil.SuccessResponse "User signed up successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request format"
// @Failure 409 {object} httputil.ErrorResponse "User already exists"
// @Failure 500 {object} httputil.ErrorResponse "Failed to sign up"
// @Router /auth/signup [post]
func SignUpHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SignUpDto
		valReqErr := httputil.ValidateRequestFormat(c, &req)
		if valReqErr != nil {
			return
		}

		serviceErr := s.SignUp(c.Request.Context(), req)
		if httputil.MakeServiceErrorResponse(c, serviceErr, ErrSignUpFailed) {
			return
		}
		httputil.MakeSuccessResponse(c, http.StatusOK, "User signed up successfully", nil)
	}
}
