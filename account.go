package bulwark

import (
	"context"
	"fmt"
	"net/http"
)

// Account is used for account bulwark-auth tasks, but it's preferable to use it via the Guard struct.
type Account struct {
	client  *http.Client
	baseURL string
}

const (
	createUrl         = "api/accounts"
	verifyUrl         = "api/accounts/verify"
	changePasswordUrl = "api/accounts/password"
)

// NewAccountClient creates a client for account tasks
func NewAccountClient(baseURL string, client *http.Client) *Account {
	return &Account{
		baseURL: baseURL,
		client:  client,
	}
}

// Create will create a user account and send a verification email
func (a Account) Create(ctx context.Context, tenantID, email, password string) error {
	payload := struct {
		TenantID string `json:"tenantId"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		TenantID: tenantID,
		Email:    email,
		Password: password,
	}

	err := doPost(ctx, fmt.Sprintf("%s/%s", a.baseURL, createUrl), payload, nil, a.client)

	if err != nil {
		return err
	}

	return nil
}

// Verify will verify a account with a verification token supplied via email
func (a Account) Verify(ctx context.Context, tenantID, email, verificationToken string) error {
	payload := struct {
		TenantID string `json:"tenantId"`
		Email    string `json:"email"`
		Token    string `json:"token"`
	}{
		TenantID: tenantID,
		Email:    email,
		Token:    verificationToken,
	}

	err := doPost(ctx, fmt.Sprintf("%s/%s", a.baseURL, verifyUrl), payload,
		nil, a.client)

	if err != nil {
		return err
	}

	return nil
}

// ChangePassword changes a password for an account, valid access token is required
func (a Account) ChangePassword(ctx context.Context, tenantID, email, newPassword, accessToken string) error {
	payload := struct {
		TenantID    string `json:"tenantId"`
		Email       string `json:"email"`
		NewPassword string `json:"newPassword"`
		AccessToken string `json:"accessToken"`
	}{
		TenantID:    tenantID,
		Email:       email,
		NewPassword: newPassword,
		AccessToken: accessToken,
	}
	err := doPut(ctx, fmt.Sprintf("%s/%s", a.baseURL, changePasswordUrl), payload, a.client)

	if err != nil {
		return err
	}

	return nil
}
