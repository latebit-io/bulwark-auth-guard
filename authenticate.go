package bulwark

import (
	"context"
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
	passwordUrl    = "authentication/authenticate"
	acknowledgeUrl = "authentication/acknowledge"
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
