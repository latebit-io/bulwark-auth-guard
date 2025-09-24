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
	Roles     []string  `json:"roles"`
	Issuer    string    `json:"issuer"`
	Subject   string    `json:"subject"`
	Audience  string    `json:"audience"`
	ExpiresAt time.Time `json:"expiresAt"`
	NotBefore time.Time `json:"notBefore"`
	IssuedAt  time.Time `json:"issuedAt"`
	ID        string    `json:"Id,omitempty"`
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
)

// NewAuthenticateClient creates a client for account tasks
func NewAuthenticateClient(baseUrl string, client *http.Client) *Authenticate {
	return &Authenticate{
		client:  client,
		baseUrl: baseUrl,
	}
}

// Password traditional authentication by email and password
func (a *Authenticate) Password(ctx context.Context, email, password string) (Authenticated, error) {
	authenticated := Authenticated{}
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}
	err := doPost(ctx, fmt.Sprintf("%s/%s", a.baseUrl, passwordUrl), payload, &authenticated, a.client)

	if err != nil {
		return Authenticated{}, err
	}

	return authenticated, nil
}

// Acknowledge notifies the server a token is in use, this should be done after each authentication
func (a *Authenticate) Acknowledge(ctx context.Context, authenticated Authenticated, email, clientID string) error {
	payload := struct {
		Email        string `json:"email"`
		ClientId     string `json:"clientId"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{
		Email:        email,
		ClientId:     clientID,
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
func (a *Authenticate) RequestMagicCode(ctx context.Context, email string) error {
	payload := struct {
		Email string `json:"email"`
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
func (a *Authenticate) MagicCode(ctx context.Context, email, magicCode string) (Authenticated, error) {
	authenticated := Authenticated{}
	payload := struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}{
		Email: email,
		Code:  magicCode,
	}

	err := doPost(ctx, fmt.Sprintf("%s/%s", a.baseUrl, magicCodeUrl), payload, &authenticated, a.client)
	if err != nil {
		return Authenticated{}, err
	}

	return authenticated, nil
}

func (a *Authenticate) ValidateAccessToken(ctx context.Context, email, accessToken, deviceId string) (AccessTokenClaims, error) {
	claims := AccessTokenClaims{}
	payload := struct {
		Email    string `json:"email"`
		ClientId string `json:"clientId"`
		Token    string `json:"token"`
	}{
		Email:    email,
		ClientId: deviceId,
		Token:    accessToken,
	}

	err := doPost(ctx, fmt.Sprintf("%s/%s", a.baseUrl, validateAccessTokenUrl), payload, &claims, a.client)
	if err != nil {
		return claims, err
	}
	return claims, nil
}

func (a *Authenticate) Renew(ctx context.Context, email, refreshToken string) (Authenticated, error) {
	authenticated := Authenticated{}
	payload := struct {
		Email        string `json:"email"`
		RefreshToken string `json:"refreshToken"`
	}{
		Email:        email,
		RefreshToken: refreshToken,
	}
	err := doPost(ctx, fmt.Sprintf("%s/%s", a.baseUrl, renewUrl), payload, &authenticated, a.client)
	if err != nil {
		return Authenticated{}, err
	}
	return authenticated, nil
}
