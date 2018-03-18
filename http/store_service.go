package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rafaeljesus/resilient-go-example"
)

var (
	usersAPIV1 = "http://localhost:8081/v1/users"
	usersAPIV2 = "http://localhost:8081/v2/users"
)

type (
	// StoreService holds fields and dependencies for storing users
	StoreService struct {
		client srv.HTTPClient
	}
)

// NewStoreService returns a configured StoreService.
func NewStoreService(hc srv.HTTPClient) *StoreService {
	return &StoreService{hc}
}

// Store responsible for storing the user.
func (s *StoreService) Store(user *srv.User) error {
	userResourceV2 := fmt.Sprintf("%s/%s/%s", usersAPIV2, user.Email, user.Country)
	res, err := s.client.Get(userResourceV2)
	if err != nil {
		if err == srv.ErrCircuitBreakerOpen {
			// fallback to apiv1
			userResourceV1 := fmt.Sprintf("%s?%s&%s", usersAPIV1, user.Email, user.Country)
			res, err = s.client.Get(userResourceV1)
			if err != nil {
				return fmt.Errorf("failed to fetch user: user service is unavailable: %v", err)
			}
		} else {
			return fmt.Errorf("failed to fetch user: %v", err)
		}
	}

	if res.StatusCode == http.StatusNotFound {
		return srv.ErrUserNotFound
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(user); err != nil {
		return fmt.Errorf("failed to marshal user payload: %v", err)
	}

	res, err = s.client.Post(usersAPIV2, "application/json", body)
	if err != nil {
		return fmt.Errorf("failed to store user: %v", err)
	}

	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return nil
	case http.StatusConflict:
		return srv.ErrUserAlreadyExists
	default:
		return errors.New("failed to store user")
	}
}
