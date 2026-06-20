package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"coleccionbackend/src/feature/etiquetas/domain"
	"coleccionbackend/src/feature/etiquetas/domain/entities"
)

type CreateEtiquetaUseCase struct {
	repo domain.EtiquetasRepository
}

func NewCreateEtiquetaUseCase(r domain.EtiquetasRepository) *CreateEtiquetaUseCase {
	return &CreateEtiquetaUseCase{repo: r}
}

func (uc *CreateEtiquetaUseCase) Execute(ctx context.Context, nombre string) (*entities.Etiqueta, error) {
	nombre = strings.TrimSpace(nombre)
	if nombre == "" {
		return nil, errors.New("el nombre de la etiqueta es requerido")
	}

	e := &entities.Etiqueta{
		ID:        uuid.New().String(),
		Nombre:    nombre,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.Create(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}
