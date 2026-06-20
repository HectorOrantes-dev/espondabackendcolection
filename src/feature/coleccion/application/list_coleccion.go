package application

import (
	"context"

	"coleccionbackend/src/feature/coleccion/domain"
	"coleccionbackend/src/feature/coleccion/domain/entities"
)

type ListColeccionUseCase struct {
	repo domain.ColeccionRepository
}

func NewListColeccionUseCase(r domain.ColeccionRepository) *ListColeccionUseCase {
	return &ListColeccionUseCase{repo: r}
}

// Execute lista los vehículos. etiquetaFiltro opcional filtra por nombre de etiqueta.
func (uc *ListColeccionUseCase) Execute(ctx context.Context, etiquetaFiltro string) ([]entities.Vehiculo, error) {
	return uc.repo.GetAll(ctx, etiquetaFiltro)
}

type GetByIDUseCase struct {
	repo domain.ColeccionRepository
}

func NewGetByIDUseCase(r domain.ColeccionRepository) *GetByIDUseCase {
	return &GetByIDUseCase{repo: r}
}

func (uc *GetByIDUseCase) Execute(ctx context.Context, id string) (*entities.Vehiculo, error) {
	return uc.repo.GetByID(ctx, id)
}
