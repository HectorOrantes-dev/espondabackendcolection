package domain

import (
	"context"
	"coleccionbackend/src/feature/login/domain/entities"
)

type LoginRepository interface {
	Login(ctx context.Context, req entities.LoginRequest) (*entities.TokenResponse, error)
	Logout(ctx context.Context, accessToken string) error
	Refresh(ctx context.Context, refreshToken string) (*entities.TokenResponse, error)
}
