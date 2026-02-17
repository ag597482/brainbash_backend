package request

// GoogleLoginRequest is the request body for POST /auth/google.
// Accepts either an id_token (from mobile) or access_token (from web).
type GoogleLoginRequest struct {
	IDToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
}
