package aplication

import (
	"context"
	"coleccionbackend/src/feature/login/domain"
	"coleccionbackend/src/feature/login/domain/entities"
)

type LoginUseCase struct {
	repo domain.LoginRepository
}

func NewLoginUseCase(r domain.LoginRepository) *LoginUseCase {
	return &LoginUseCase{repo: r}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req entities.LoginRequest) (*entities.TokenResponse, error) {
	return uc.repo.Login(ctx, req)
}
