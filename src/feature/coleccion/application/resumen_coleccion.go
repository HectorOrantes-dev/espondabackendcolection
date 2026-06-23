package application

import (
	"context"

	"coleccionbackend/src/feature/coleccion/domain"
	"coleccionbackend/src/feature/coleccion/domain/entities"
)

type ResumenColeccionUseCase struct {
	repo domain.ColeccionRepository
}

func NewResumenColeccionUseCase(r domain.ColeccionRepository) *ResumenColeccionUseCase {
	return &ResumenColeccionUseCase{repo: r}
}

func (uc *ResumenColeccionUseCase) Execute(ctx context.Context) (*entities.ResumenColeccion, error) {
	return uc.repo.Resumen(ctx)
}
