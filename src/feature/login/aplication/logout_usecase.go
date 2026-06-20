package aplication

import (
	"context"
	"coleccionbackend/src/feature/login/domain"
)

type LogoutUseCase struct {
	repo domain.LoginRepository
}

func NewLogoutUseCase(r domain.LoginRepository) *LogoutUseCase {
	return &LogoutUseCase{repo: r}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, accessToken string) error {
	return uc.repo.Logout(ctx, accessToken)
}
