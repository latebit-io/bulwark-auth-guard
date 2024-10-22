package guard_test

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/latebitflip-io/bulwark-auth-guard"
	"github.com/latebitflip-io/go-hog"
	"net/http"
	"testing"
)

const baseUri = "http://localhost:8080"
const mailHogUri = "http://localhost:8025"

func TestAccountCreate(t *testing.T) {
	client := &http.Client{}
	id := uuid.New()
	account := guard.NewAccount(baseUri, client)
	err := account.Create(fmt.Sprintf("%s@bulwark.io", id.String()), "password12!P")
	if err != nil {
		t.Error(err)
	}
}

func TestAccountCreateDuplicate(t *testing.T) {
	client := &http.Client{}
	id := uuid.New()
	account := guard.NewAccount(baseUri, client)
	err := account.Create(fmt.Sprintf("%s@bulwark.io", id.String()), "password12!P")
	if err != nil {
		t.Error(err)
	}

	err = account.Create(fmt.Sprintf("%s@bulwark.io", id.String()), "password12!P")
	if err == nil {
		t.Error(err)
	}
}

func TestAccountCreateAndVerify(t *testing.T) {
	client := &http.Client{}
	id := uuid.New()
	email := fmt.Sprintf("%s@bulwark.io", id.String())
	account := guard.NewAccount(baseUri, client)
	err := account.Create(email, "password12!P")
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
	err = account.Verify(email, message.Content.Headers.Subject[0])
	if err != nil {
		t.Error(err)
	}
	err = gohog.DeleteAll()
}

func findToMessage(messages gohog.Messages, to string) (gohog.Message, error) {
	for _, m := range messages.Items {
		if m.Content.Headers.To[0] == to {
			return m, nil
		}
	}
	return gohog.Message{}, errors.New("message not found")
}
