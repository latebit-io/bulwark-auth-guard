package bulwark

import "net/http"

type Guard struct {
	Account      *Account
	Authenticate *Authenticate
	baseUrl      string
}

func NewGuard(baseUrl string, client *http.Client) *Guard {
	return &Guard{
		Account:      NewAccountClient(baseUrl, client),
		Authenticate: NewAuthenticateClient(baseUrl, client),
		baseUrl:      baseUrl,
	}
}
