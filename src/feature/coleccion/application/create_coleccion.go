package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"coleccionbackend/src/feature/coleccion/domain"
	"coleccionbackend/src/feature/coleccion/domain/entities"
)

type ImageInput struct {
	Filename string
	Content  []byte
	MimeType string
}

type CreateColeccionInput struct {
	Nombre      string
	Marca       string
	Modelo      string
	Precio      float64
	Images      []ImageInput
	EtiquetaIDs []string
}

type CreateColeccionUseCase struct {
	repo         domain.ColeccionRepository
	imageService domain.ImageService
}

func NewCreateColeccionUseCase(r domain.ColeccionRepository, i domain.ImageService) *CreateColeccionUseCase {
	return &CreateColeccionUseCase{repo: r, imageService: i}
}

func (uc *CreateColeccionUseCase) Execute(ctx context.Context, input CreateColeccionInput) (*entities.Vehiculo, error) {
	if len(input.Images) > 3 {
		return nil, errors.New("máximo 3 imágenes permitidas por vehículo")
	}

	urls, ids, err := uploadImagesParallel(uc.imageService, input.Images)
	if err != nil {
		return nil, err
	}

	v := &entities.Vehiculo{
		ID:        uuid.New().String(),
		Nombre:    input.Nombre,
		Marca:     input.Marca,
		Modelo:    input.Modelo,
		Precio:      input.Precio,
		Imagenes:    urls,
		ImageIDs:    ids,
		EtiquetaIDs: input.EtiquetaIDs,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.repo.Create(ctx, v); err != nil {
		for _, id := range ids {
			_ = uc.imageService.Delete(id)
		}
		return nil, err
	}

	return v, nil
}
