package guard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Account struct is the base for account tasks
type Account struct {
	client  *http.Client
	baseURL string
}

const (
	createUrl = "accounts/create"
	verifyUrl = "accounts/verify"
)

// NewAccountClient creates a client for account tasks
func NewAccountClient(baseURL string, client *http.Client) *Account {
	return &Account{
		baseURL: baseURL,
		client:  client,
	}
}

// Create will create a user account and send a verification email
func (a Account) Create(email, password string) error {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	err := doPost(fmt.Sprintf("%s/%s", a.baseURL, createUrl), payload,
		nil, a.client)

	if err != nil {
		return err
	}

	return nil
}

// Verify will verify a account
func (a Account) Verify(email, verificationToken string) error {
	payload := struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		Email: email,
		Token: verificationToken,
	}

	err := doPost(fmt.Sprintf("%s/%s", a.baseURL, verifyUrl), payload,
		nil, a.client)

	if err != nil {
		return err
	}

	return nil
}

func doPost(url string, payload interface{}, model interface{}, client *http.Client) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		jsonError := &JsonError{}
		if err := json.NewDecoder(resp.Body).Decode(jsonError); err != nil {
			return err
		}
		if jsonError != nil {
			return fmt.Errorf("%s - %s", jsonError.Title, jsonError.Detail)
		}
	}

	if resp.Body != http.NoBody {
		if err := json.NewDecoder(resp.Body).Decode(model); err != nil {
			return err
		}
	}

	defer resp.Body.Close()

	return nil
}
