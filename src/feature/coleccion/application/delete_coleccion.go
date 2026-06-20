package application

import (
	"context"

	"coleccionbackend/src/feature/coleccion/domain"
)

type DeleteColeccionUseCase struct {
	repo         domain.ColeccionRepository
	imageService domain.ImageService
}

func NewDeleteColeccionUseCase(r domain.ColeccionRepository, i domain.ImageService) *DeleteColeccionUseCase {
	return &DeleteColeccionUseCase{repo: r, imageService: i}
}

func (uc *DeleteColeccionUseCase) Execute(ctx context.Context, id string) error {
	v, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Eliminar imágenes de Google Drive antes de borrar el registro
	for _, imageID := range v.ImageIDs {
		_ = uc.imageService.Delete(imageID)
	}

	return uc.repo.Delete(ctx, id)
}
