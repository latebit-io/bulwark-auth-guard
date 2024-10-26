package bulwark

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	gohog "github.com/latebitflip-io/go-hog"
	"net/http"
	"testing"
)

func TestAuthenticatePassword(t *testing.T) {
	client := &http.Client{}
	id := uuid.New()
	email := fmt.Sprintf("%s@bulwark.io", id.String())
	password := "password12!P"
	guard := NewGuard(baseUri, client)
	ctx := context.Background()
	err := createAndVerifyAccount(ctx, email, password, guard, client)
	if err != nil {
		t.Error(err)
	}

	authenticated, err := guard.Authenticate.Password(ctx, email, password)
	if err != nil {
		t.Error(err)
	}

	if authenticated.AccessToken == "" {
		t.Error("Token not returned")
	}

	err = guard.Authenticate.Acknowledge(ctx, authenticated, email, "testdevice")
	if err != nil {
		t.Error(err)
	}
}

func createAndVerifyAccount(ctx context.Context, email, password string, guard *Guard, client *http.Client) error {
	err := guard.Account.Create(ctx, email, password)
	if err != nil {
		return err
	}
	gohog := gohog.NewGoHogClient(mailHogUri, client)
	messages, err := gohog.Messages(0, 100)
	if err != nil {
		return err
	}
	message, err := findToMessage(messages, email)
	if err != nil {
		return err
	}
	err = guard.Account.Verify(ctx, email, message.Subject())
	if err != nil {
		return err
	}
	err = gohog.DeleteAll()
	if err != nil {
		return err
	}
	return nil
}
