package application

import (
	"context"

	"coleccionbackend/src/feature/etiquetas/domain"
)

type DeleteEtiquetaUseCase struct {
	repo domain.EtiquetasRepository
}

func NewDeleteEtiquetaUseCase(r domain.EtiquetasRepository) *DeleteEtiquetaUseCase {
	return &DeleteEtiquetaUseCase{repo: r}
}

func (uc *DeleteEtiquetaUseCase) Execute(ctx context.Context, id string) error {
	// La relación en vehiculo_etiquetas se borra en cascada (ON DELETE CASCADE),
	// así que los vehículos simplemente pierden esta etiqueta.
	return uc.repo.Delete(ctx, id)
}
