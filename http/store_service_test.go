package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/rafaeljesus/resilient-go-example"
	"github.com/rafaeljesus/resilient-go-example/mock"
)

func TestUserStore(t *testing.T) {
	client := new(mock.HTTPClientMock)
	client.GetFunc = func(url string) (*http.Response, error) {
		if url != usersAPIV2+"/foo@mail.com/de" {
			t.Fatal("unexpected url")
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: nopCloser{
				bytes.NewBufferString(`{"email": "foo@mail.com", "country": "de"}`),
			},
		}, nil
	}
	client.PostFunc = func(url, contentType string, body io.Reader) (*http.Response, error) {
		if url != usersAPIV2 {
			t.Fatal("unexpected url")
		}
		if contentType != "application/json" {
			t.Fatal("unexpected contentType")
		}
		u := new(srv.User)
		if err := json.NewDecoder(body).Decode(u); err != nil {
			t.Fatalf("unexpected body decode error: %v", err)
		}
		if u.Country != "de" {
			t.Fatal("unexpected country")
		}
		if u.Email != "foo@mail.com" {
			t.Fatal("unexpected email")
		}
		return &http.Response{StatusCode: http.StatusOK}, nil
	}
	storer := NewStoreService(client)
	u := srv.NewUser("foo@mail.com", "de")
	if err := storer.Store(u); err != nil {
		t.Fatalf("failed to store user: %v", err)
	}
	if !client.PostInvoked {
		t.Fatal("expected client.Post() to be invoked")
	}
	if !client.GetInvoked {
		t.Fatal("expected client.Get() to be invoked")
	}
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }
