package util

import (
	"context"
	"fmt"

	"github.com/royanqodri/fms-ai-hackathon/config"
	"google.golang.org/api/idtoken"
)

// GoogleUserInfo Get data user
type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
}

// VerifyGoogleIDToken verifications ID Token from Google Sign-In
func VerifyGoogleIDToken(ctx context.Context, idToken string) (*GoogleUserInfo, error) {
	if idToken == "" {
		return nil, fmt.Errorf("id_token is required")
	}

	clientID := config.Get().Google.ClientId
	if clientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is not configured")
	}

	payload, err := idtoken.Validate(ctx, idToken, clientID)
	if err != nil {
		return nil, fmt.Errorf("invalid google id token: %w", err)
	}

	// GET claims
	email, _ := payload.Claims["email"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	givenName, _ := payload.Claims["given_name"].(string)
	familyName, _ := payload.Claims["family_name"].(string)

	if email == "" {
		return nil, fmt.Errorf("email not found in google token")
	}

	if !emailVerified {
		return nil, fmt.Errorf("email is not verified by google")
	}

	return &GoogleUserInfo{
		Sub:           payload.Subject,
		Email:         email,
		EmailVerified: emailVerified,
		Name:          name,
		Picture:       picture,
		GivenName:     givenName,
		FamilyName:    familyName,
	}, nil
}
