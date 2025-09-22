package bulwark

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	gohog "github.com/latebit-io/go-hog"
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

	claims, err := guard.Authenticate.ValidateAccessToken(ctx, email, authenticated.AccessToken, "testdevice")
	if err != nil {
		t.Error(err)
	}

	if claims.Subject != email {
		t.Error("Subject does not match email")
	}

	authenticated, err = guard.Authenticate.Renew(ctx, claims.Subject, authenticated.RefreshToken)
	claims, err = guard.Authenticate.ValidateAccessToken(ctx, email, authenticated.AccessToken, "testdevice")

	if claims.Subject != email {
		t.Error("Refresh Token does not match email")
	}
}

func TestAuthenticateMagicCode(t *testing.T) {
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

	err = guard.Authenticate.RequestMagicCode(ctx, email)
	if err != nil {
		t.Error(err)
	}
	gohog := gohog.NewGoHogClient(mailHogUri, client)
	messages, err := gohog.Messages(0, 100)
	if err != nil {
		t.Error(err)
	}
	message, err := findToMessage(messages, email)
	if err != nil {
		t.Error(err)
	}
	code := message.Subject()
	fmt.Printf("Message: %s\n", message.Subject())

	authenticated, err := guard.Authenticate.MagicCode(ctx, email, code)
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

func TestAuthenticateMagicCodeFail(t *testing.T) {
	email := "test9090909@latebit.io"
	client := &http.Client{}
	guard := NewGuard(baseUri, client)
	ctx := context.Background()
	err := guard.Authenticate.RequestMagicCode(ctx, email)
	if err == nil {
		t.Error("should throw an magic link error")
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
