package aplication

import (
	"context"
	"coleccionbackend/src/feature/login/domain"
	"coleccionbackend/src/feature/login/domain/entities"
)

type RefreshUseCase struct {
	repo domain.LoginRepository
}

func NewRefreshUseCase(r domain.LoginRepository) *RefreshUseCase {
	return &RefreshUseCase{repo: r}
}

func (uc *RefreshUseCase) Execute(ctx context.Context, refreshToken string) (*entities.TokenResponse, error) {
	return uc.repo.Refresh(ctx, refreshToken)
}
