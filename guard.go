package guard

import "net/http"

type Guard struct {
	Account *Account
	baseUrl string
}

func NewGuard(baseUrl string, client http.Client) *Guard {
	return &Guard{
		Account: NewAccountClient(baseUrl, &client),
		baseUrl: baseUrl,
	}
}
