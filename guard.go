package guard

import "net/http"

type Guard struct {
	Account *Account
	baseUrl string
}

func NewGuard(baseUrl string, client http.Client) *Guard {
	return &Guard{
		Account: NewAccount(baseUrl, &client),
		baseUrl: baseUrl,
	}
}
