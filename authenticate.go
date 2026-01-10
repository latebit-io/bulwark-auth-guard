package bulwark

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Authenticated returned on successful authentication and should be acknowledged
type Authenticated struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AccessTokenClaims struct {
	TenantID  string    `json:"tenantId"`
	Roles     []string  `json:"roles"`
	Issuer    string    `json:"issuer"`
	Subject   string    `json:"subject"`
	Audience  string    `json:"audience"`
	ExpiresAt time.Time `json:"expiresAt"`
	NotBefore time.Time `json:"notBefore"`
	IssuedAt  time.Time `json:"issuedAt"`
	ID        string    `json:"id,omitempty"`
	ClientID  string    `json:"clientId,omitempty"`
}

// Authenticate is used for authentication bulwark-auth tasks, but it's preferable to use it via the Guard struct.
type Authenticate struct {
	client  *http.Client
	baseUrl string
}

const (
	passwordUrl            = "api/authenticate"
	acknowledgeUrl         = "api/authenticate/ack"
	requestMagicCodeUrl    = "api/authenticate/logon/request"
	magicCodeUrl           = "api/authenticate/code"
	validateAccessTokenUrl = "api/authenticate/token/validate"
	renewUrl               = "api/authenticate/renew"
	revokeUrl              = "api/authenticate/revoke"
)

// NewAuthenticateClient creates a client for account tasks
func NewAuthenticateClient(baseUrl string, client *http.Client) *Authenticate {
	return &Authenticate{
		client:  client,
		baseUrl: baseUrl,
	}
}

// Password traditional authentication by email and password
func (a *Authenticate) Password(ctx context.Context, tenantID, email, password, clientID string) (Authenticated, error) {
	authenticated := Authenticated{}
	payload := struct {
		TenantID string `json:"tenantId"`
		Email    string `json:"email"`
		Password string `json:"password"`
		ClientID string `json:"clientId"`
	}{
		TenantID: tenantID,
		Email:    email,
		Password: password,
		ClientID: clientID,
	}
	err := doPost(ctx, fmt.Sprintf("%s/%s", a.baseUrl, passwordUrl), payload, &authenticated, a.client)

	if err != nil {
		return Authenticated{}, err
	}

	return authenticated, nil
}

// Acknowledge notifies the server a token is in use, this should be done after each authentication
func (a *Authenticate) Acknowledge(ctx context.Context, tenantID string, authenticated Authenticated) error {
	payload := struct {
		TenantID     string `json:"tenantId"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{
		TenantID:     tenantID,
		AccessToken:  authenticated.AccessToken,
		RefreshToken: authenticated.RefreshToken,
	}

	err := doPost(ctx, fmt.Sprintf("%s/%s", a.baseUrl, acknowledgeUrl), payload, nil, a.client)
	if err != nil {
		return err
	}

	return nil
}

// RequestMagicCode will send an email with a magic code link
func (a *Authenticate) RequestMagicCode(ctx context.Context, tenantID, email string) error {
	payload := struct {
		TenantID string `json:"tenantId"`
		Email    string `json:"email"`
	}{
		Email: email,
	}

	err := doPost(ctx, fmt.Sprintf("%s/%s", a.baseUrl, requestMagicCodeUrl), payload, nil, a.client)
	if err != nil {
		return err
	}

	return nil
}

// MagicCode authenticates a user with email and a magic code
func (a *Authenticate) MagicCode(ctx context.Context, tenantID, email, magicCode, clientID string) (Authenticated, error) {
	authenticated := Authenticated{}
	payload := struct {
		TenantID string `json:"tenantId"`
		Email    string `json:"email"`
		Code     string `json:"code"`
		ClientID string `json:"clientId"`
	}{
		TenantID: tenantID,
		Email:    email,
		Code:     magicCode,
		ClientID: clientID,
	}

	err := doPost(ctx, fmt.Sprintf("%s/%s", a.baseUrl, magicCodeUrl), payload, &authenticated, a.client)
	if err != nil {
		return Authenticated{}, err
	}

	return authenticated, nil
}

func (a *Authenticate) ValidateAccessToken(ctx context.Context, tenantID string, accessToken string) (AccessTokenClaims, error) {
	claims := AccessTokenClaims{}
	payload := struct {
		TenantID string `json:"tenantId"`
		Token    string `json:"token"`
	}{
		TenantID: tenantID,
		Token:    accessToken,
	}

	err := doPost(ctx, fmt.Sprintf("%s/%s", a.baseUrl, validateAccessTokenUrl), payload, &claims, a.client)
	if err != nil {
		return claims, err
	}
	return claims, nil
}

func (a *Authenticate) Renew(ctx context.Context, tenantID, email, refreshToken string) (Authenticated, error) {
	authenticated := Authenticated{}
	payload := struct {
		TenantID     string `json:"tenantId"`
		Email        string `json:"email"`
		RefreshToken string `json:"refreshToken"`
	}{
		TenantID:     tenantID,
		Email:        email,
		RefreshToken: refreshToken,
	}
	err := doPost(ctx, fmt.Sprintf("%s/%s", a.baseUrl, renewUrl), payload, &authenticated, a.client)
	if err != nil {
		return Authenticated{}, err
	}
	return authenticated, nil
}

func (a *Authenticate) Revoke(ctx context.Context, tenantID, email, accessToken, clientID string) error {
	payload := struct {
		TenantID    string `json:"tenantId"`
		Email       string `json:"email"`
		ClientID    string `json:"clientId"`
		AccessToken string `json:"accessToken"`
	}{
		TenantID:    tenantID,
		Email:       email,
		ClientID:    clientID,
		AccessToken: accessToken,
	}

	return doDelete(ctx, fmt.Sprintf("%s/%s", a.baseUrl, revokeUrl), payload, a.client)
}
