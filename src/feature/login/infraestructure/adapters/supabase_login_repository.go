package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"coleccionbackend/src/feature/login/domain/entities"
)

type SupabaseLoginRepository struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewSupabaseLoginRepository() *SupabaseLoginRepository {
	return &SupabaseLoginRepository{
		baseURL: os.Getenv("SUPABASE_URL"),
		apiKey:  os.Getenv("SUPABASE_ANON_KEY"),
		client:  &http.Client{},
	}
}

func (r *SupabaseLoginRepository) Login(_ context.Context, req entities.LoginRequest) (*entities.TokenResponse, error) {
	body, _ := json.Marshal(map[string]string{
		"email":    req.Email,
		"password": req.Password,
	})

	url := fmt.Sprintf("%s/auth/v1/token?grant_type=password", r.baseURL)
	resp, err := r.doPost(url, body, "")
	if err != nil {
		return nil, err
	}

	var token entities.TokenResponse
	if err := json.Unmarshal(resp, &token); err != nil {
		return nil, fmt.Errorf("error parseando respuesta de login: %w", err)
	}
	return &token, nil
}

func (r *SupabaseLoginRepository) Logout(_ context.Context, accessToken string) error {
	url := fmt.Sprintf("%s/auth/v1/logout", r.baseURL)
	// El logout es idempotente: aunque Supabase rechace el token (expirado o
	// inválido), consideramos el logout exitoso para que el frontend pueda
	// limpiar la sesión sin recibir un 500.
	_, _ = r.doPost(url, nil, accessToken)
	return nil
}

func (r *SupabaseLoginRepository) Refresh(_ context.Context, refreshToken string) (*entities.TokenResponse, error) {
	body, _ := json.Marshal(map[string]string{
		"refresh_token": refreshToken,
	})

	url := fmt.Sprintf("%s/auth/v1/token?grant_type=refresh_token", r.baseURL)
	resp, err := r.doPost(url, body, "")
	if err != nil {
		return nil, err
	}

	var token entities.TokenResponse
	if err := json.Unmarshal(resp, &token); err != nil {
		return nil, fmt.Errorf("error parseando respuesta de refresh: %w", err)
	}
	return &token, nil
}

func (r *SupabaseLoginRepository) doPost(url string, body []byte, bearerToken string) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", r.apiKey)
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en solicitud a Supabase Auth: %w", err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error de Supabase Auth (%d): %s", resp.StatusCode, string(data))
	}
	return data, nil
}
