package guard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Account struct {
	http.Client
	baseURL string
}

const (
	create = "account/create"
)

func NewAccount(baseURL string, client http.Client) *Account {
	return &Account{
		baseURL: baseURL,
		Client:  client,
	}
}

func (a Account) Create(email, password string) error {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := a.Client.Post(fmt.Sprintf("%s/%s", a.baseURL, create),
		"encoding/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
