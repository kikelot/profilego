package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"profilego/internal/domain"
)

type AuthClient struct {
	BaseURL string
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{BaseURL: baseURL}
}

func (ac *AuthClient) GetCurrentUser(ctx context.Context, token string) (*domain.AuthUser, error) {
	url := fmt.Sprintf("%s/users/current", ac.BaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Agregar el Bearer Token en el header de autorización
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: código de estado %d", resp.StatusCode)
	}

	var user domain.AuthUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
