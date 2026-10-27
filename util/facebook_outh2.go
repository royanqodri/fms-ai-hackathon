package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/royanqodri/fms-ai-hackathon/config"
)

// FacebookUserInfo berisi data yang didapat dari Facebook Graph API
type FacebookUserInfo struct {
	ID      string `json:"id"` // Facebook User ID (unik & stabil)
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

// facebookDebugTokenResponse merepresentasikan response dari endpoint debug_token
type facebookDebugTokenResponse struct {
	Data struct {
		AppID     string `json:"app_id"`
		IsValid   bool   `json:"is_valid"`
		UserID    string `json:"user_id"`
		ExpiresAt int64  `json:"expires_at"`
		Error     *struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		} `json:"error,omitempty"`
	} `json:"data"`
}

// VerifyFacebookAccessToken memverifikasi Access Token dari Facebook Login
// dan mengembalikan data user yang sudah tervalidasi.
func VerifyFacebookAccessToken(ctx context.Context, accessToken string) (*FacebookUserInfo, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access_token is required")
	}

	appID := config.Get().Facebook.AppId
	appSecret := config.Get().Facebook.AppSecret
	if appID == "" || appSecret == "" {
		return nil, fmt.Errorf("facebook app id/secret is not configured")
	}

	// 1. Validasi token via debug_token (memastikan token milik app kita & masih valid)
	if err := debugFacebookToken(ctx, accessToken, appID, appSecret); err != nil {
		return nil, err
	}

	// 2. Ambil data user dari endpoint /me
	userInfo, err := fetchFacebookUserInfo(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	if userInfo.Email == "" {
		return nil, fmt.Errorf("email not found in facebook profile (pastikan permission 'email' diminta saat login)")
	}

	return userInfo, nil
}

// debugFacebookToken memvalidasi access token menggunakan endpoint /debug_token
func debugFacebookToken(ctx context.Context, accessToken, appID, appSecret string) error {
	appAccessToken := fmt.Sprintf("%s|%s", appID, appSecret)

	endpoint := "https://graph.facebook.com/debug_token"
	params := url.Values{}
	params.Set("input_token", accessToken)
	params.Set("access_token", appAccessToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return fmt.Errorf("failed to build debug_token request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call facebook debug_token: %w", err)
	}
	defer resp.Body.Close()

	var result facebookDebugTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode debug_token response: %w", err)
	}

	if result.Data.Error != nil {
		return fmt.Errorf("invalid facebook access token: %s", result.Data.Error.Message)
	}
	if !result.Data.IsValid {
		return fmt.Errorf("facebook access token is not valid")
	}
	if result.Data.AppID != appID {
		return fmt.Errorf("facebook access token does not belong to this app")
	}

	return nil
}

// fetchFacebookUserInfo mengambil profil user dari endpoint /me
func fetchFacebookUserInfo(ctx context.Context, accessToken string) (*FacebookUserInfo, error) {
	endpoint := "https://graph.facebook.com/me"
	params := url.Values{}
	params.Set("fields", "id,name,email,picture")
	params.Set("access_token", accessToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build /me request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call facebook graph api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("facebook graph api returned status %d", resp.StatusCode)
	}

	var userInfo FacebookUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode facebook user info: %w", err)
	}

	if userInfo.ID == "" {
		return nil, fmt.Errorf("facebook user id not found")
	}

	return &userInfo, nil
}
