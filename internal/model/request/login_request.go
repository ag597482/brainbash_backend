package request

// GoogleLoginRequest is the request body for POST /auth/google.
type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}