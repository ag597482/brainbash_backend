package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	googleTokenInfoURL = "https://oauth2.googleapis.com/tokeninfo?id_token=%s"
	googleUserInfoURL  = "https://www.googleapis.com/oauth2/v3/userinfo"
)

// GoogleUserInfo holds the user details extracted from a verified Google token.
type GoogleUserInfo struct {
	Sub        string `json:"sub"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

// GoogleAuthService handles Google token verification.
type GoogleAuthService struct {
	allowedClientIDs map[string]struct{} // set of accepted client IDs (web, Android, etc.)
	httpClient       *http.Client
}

// NewGoogleAuthService creates a new GoogleAuthService.
// clientIDs is one or more Google OAuth client IDs (comma-separated string split into list).
// Any of these are accepted as the id_token "aud" claim (e.g. web + Android).
func NewGoogleAuthService(clientIDs []string) *GoogleAuthService {
	allowed := make(map[string]struct{})
	for _, id := range clientIDs {
		id = strings.TrimSpace(id)
		if id != "" {
			allowed[id] = struct{}{}
		}
	}
	return &GoogleAuthService{
		allowedClientIDs: allowed,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// VerifyIDToken verifies a Google ID token and returns the user info.
// It calls Google's tokeninfo endpoint and checks that the audience (aud)
// matches the configured client ID.
func (s *GoogleAuthService) VerifyIDToken(idToken string) (*GoogleUserInfo, error) {
	resp, err := s.httpClient.Get(fmt.Sprintf(googleTokenInfoURL, idToken))
	if err != nil {
		return nil, fmt.Errorf("failed to verify token with Google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google token verification failed with status: %d", resp.StatusCode)
	}

	var payload struct {
		GoogleUserInfo
		Aud string `json:"aud"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode Google token response: %w", err)
	}

	if _, allowed := s.allowedClientIDs[payload.Aud]; !allowed {
		return nil, fmt.Errorf("token audience mismatch: token aud %q is not in allowed client IDs", payload.Aud)
	}

	return &payload.GoogleUserInfo, nil
}

// VerifyAccessToken verifies a Google access token by calling the userinfo endpoint.
// This is used for web clients where the Google Sign-In SDK provides an access_token
// instead of an id_token.
func (s *GoogleAuthService) VerifyAccessToken(accessToken string) (*GoogleUserInfo, error) {
	req, err := http.NewRequest("GET", googleUserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info from Google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google userinfo request failed with status: %d", resp.StatusCode)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode Google userinfo response: %w", err)
	}

	if userInfo.Sub == "" {
		return nil, fmt.Errorf("invalid Google access token: no user ID returned")
	}

	return &userInfo, nil
}
