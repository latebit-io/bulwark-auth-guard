package bulwark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Authenticated returned on successful authentication and should be acknowledged
type Authenticated struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Authenticate is used for authentication bulwark-auth tasks, but it's preferable to use it via the Guard struct.
type Authenticate struct {
	client  *http.Client
	baseUrl string
}

const (
	passwordUrl         = "authentication/authenticate"
	acknowledgeUrl      = "authentication/acknowledge"
	requestMagicCodeUrl = "passwordless/magic/request"
	magicCodeUrl        = "passwordless/magic/authenticate"
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
func (a *Authenticate) Acknowledge(ctx context.Context, authenticated Authenticated, email, deviceId string) error {
	payload := struct {
		Email        string `json:"email"`
		DeviceId     string `json:"deviceId"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{
		Email:        email,
		DeviceId:     deviceId,
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
	resp, err := a.client.Get(fmt.Sprintf("%s/%s/%s", a.baseUrl, requestMagicCodeUrl, email))
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		jsonError := &JsonError{}
		if err := json.NewDecoder(resp.Body).Decode(jsonError); err != nil {
			return err
		}
		if jsonError != nil {
			return fmt.Errorf("%s - %s", jsonError.Title, jsonError.Detail)
		}
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
