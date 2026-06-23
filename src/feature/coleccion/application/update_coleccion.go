package application

import (
	"context"
	"errors"
	"time"

	"coleccionbackend/src/feature/coleccion/domain"
)

type UpdateColeccionInput struct {
	ID     string
	Nombre string
	Marca  string
	Modelo string
	// Precio: nil = no se toca; puntero = nuevo valor (permite poner 0).
	Precio *float64
	Images []ImageInput // Si se envían imágenes nuevas, reemplazan las anteriores
	// EtiquetaIDs: nil = no se tocan las etiquetas; slice (aunque vacío) = reemplazo.
	EtiquetaIDs []string
}

type UpdateColeccionUseCase struct {
	repo         domain.ColeccionRepository
	imageService domain.ImageService
}

func NewUpdateColeccionUseCase(r domain.ColeccionRepository, i domain.ImageService) *UpdateColeccionUseCase {
	return &UpdateColeccionUseCase{repo: r, imageService: i}
}

func (uc *UpdateColeccionUseCase) Execute(ctx context.Context, input UpdateColeccionInput) error {
	if len(input.Images) > 3 {
		return errors.New("máximo 3 imágenes permitidas por vehículo")
	}

	existing, err := uc.repo.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}

	if input.Nombre != "" {
		existing.Nombre = input.Nombre
	}
	if input.Marca != "" {
		existing.Marca = input.Marca
	}
	if input.Modelo != "" {
		existing.Modelo = input.Modelo
	}
	if input.Precio != nil {
		existing.Precio = *input.Precio
	}

	if len(input.Images) > 0 {
		newURLs, newIDs, uploadErr := uploadImagesParallel(uc.imageService, input.Images)
		if uploadErr != nil {
			return uploadErr
		}

		for _, oldID := range existing.ImageIDs {
			_ = uc.imageService.Delete(oldID)
		}

		existing.Imagenes = newURLs
		existing.ImageIDs = newIDs
	}

	// Etiquetas: si vienen en el input (aunque sea lista vacía) se reemplazan;
	// si es nil, se conservan las que ya tenía el vehículo.
	if input.EtiquetaIDs != nil {
		existing.EtiquetaIDs = input.EtiquetaIDs
	} else {
		ids := make([]string, len(existing.Etiquetas))
		for i, e := range existing.Etiquetas {
			ids[i] = e.ID
		}
		existing.EtiquetaIDs = ids
	}

	existing.UpdatedAt = time.Now()
	return uc.repo.Update(ctx, existing)
}
