package application

import (
	"context"

	"coleccionbackend/src/feature/etiquetas/domain"
	"coleccionbackend/src/feature/etiquetas/domain/entities"
)

type ListEtiquetasUseCase struct {
	repo domain.EtiquetasRepository
}

func NewListEtiquetasUseCase(r domain.EtiquetasRepository) *ListEtiquetasUseCase {
	return &ListEtiquetasUseCase{repo: r}
}

func (uc *ListEtiquetasUseCase) Execute(ctx context.Context) ([]entities.Etiqueta, error) {
	return uc.repo.GetAll(ctx)
}
