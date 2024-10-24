package bulwark

import (
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
