package bulwark

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	gohog "github.com/latebit-io/go-hog"
	"net/http"
	"testing"
)

const baseUri = "http://localhost:8080"
const mailHogUri = "http://localhost:8025"

func TestAccountCreate(t *testing.T) {
	client := &http.Client{}
	id := uuid.New()
	guard := NewGuard(baseUri, client)
	ctx := context.Background()
	err := guard.Account.Create(ctx, fmt.Sprintf("%s@bulwark.io", id.String()), "password12!P")
	if err != nil {
		t.Error(err)
	}
}

func TestAccountCreateDuplicate(t *testing.T) {
	client := &http.Client{}
	id := uuid.New()
	guard := NewGuard(baseUri, client)
	ctx := context.Background()
	err := guard.Account.Create(ctx, fmt.Sprintf("%s@bulwark.io", id.String()), "password12!P")
	if err != nil {
		t.Error(err)
	}

	err = guard.Account.Create(ctx, fmt.Sprintf("%s@bulwark.io", id.String()), "password12!P")
	if err == nil {
		t.Error(err)
	}
}

func TestAccountCreateAndVerify(t *testing.T) {
	client := &http.Client{}
	id := uuid.New()
	email := fmt.Sprintf("%s@bulwark.io", id.String())
	guard := NewGuard(baseUri, client)
	ctx := context.Background()
	err := guard.Account.Create(ctx, email, "password12!P")
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
	err = guard.Account.Verify(ctx, email, message.Subject())
	if err != nil {
		t.Error(err)
	}

	authenticated, err := guard.Authenticate.Password(ctx, email, "password12!P")
	if err != nil {
		t.Error(err)
	}

	err = guard.Authenticate.Acknowledge(ctx, authenticated, email, "deviceId")
	if err != nil {
		t.Error(err)
	}

	err = guard.Account.ChangePassword(ctx, email, "newPassword12!P", authenticated.AccessToken)
	if err != nil {
		t.Error(err)
	}

	authenticated, err = guard.Authenticate.Password(ctx, email, "newPassword12!P")
	if err != nil {
		t.Error(err)
	}

	err = gohog.DeleteAll()
}

func findToMessage(messages gohog.Messages, to string) (gohog.Message, error) {
	for _, m := range messages.Items {
		if m.ToAddresses()[0] == to {
			return m, nil
		}
	}
	return gohog.Message{}, errors.New("message not found")
}
