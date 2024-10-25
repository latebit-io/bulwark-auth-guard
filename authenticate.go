package bulwark

import (
	"fmt"
	"net/http"
)

type Authenticated struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Authenticate struct {
	client  *http.Client
	baseUrl string
}

const (
	passwordUrl    = "authentication/authenticate"
	acknowledgeUrl = "authentication/acknowledge"
)

func NewAuthenticateClient(baseUrl string, client *http.Client) *Authenticate {
	return &Authenticate{
		client:  client,
		baseUrl: baseUrl,
	}
}

func (a *Authenticate) Password(email, password string) (Authenticated, error) {
	authenticated := Authenticated{}
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}
	err := doPost(fmt.Sprintf("%s/%s", a.baseUrl, passwordUrl), payload, &authenticated, a.client)

	if err != nil {
		return Authenticated{}, err
	}

	return authenticated, nil
}

func (a *Authenticate) Acknowledge(authenticated Authenticated, email, deviceId string) error {
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

	err := doPost(fmt.Sprintf("%s/%s", a.baseUrl, acknowledgeUrl), payload, nil, a.client)
	if err != nil {
		return err
	}

	return nil
}
