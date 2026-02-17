package response


// LoginResponse is the response body for POST /auth/google.
type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	User        UserInfo `json:"user"`
}

// UserInfo represents user details returned in auth responses.
type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}
