package bulwark

import "net/http"

// Guard is the main client struct to use with bulwark-auth
type Guard struct {
	Account      *Account
	Authenticate *Authenticate
	baseUrl      string
}

// NewGuard constructor
func NewGuard(baseUrl string, client *http.Client) *Guard {
	return &Guard{
		Account:      NewAccountClient(baseUrl, client),
		Authenticate: NewAuthenticateClient(baseUrl, client),
		baseUrl:      baseUrl,
	}
}
