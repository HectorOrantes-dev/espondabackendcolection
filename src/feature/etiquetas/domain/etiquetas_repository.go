package domain

import (
	"context"

	"coleccionbackend/src/feature/etiquetas/domain/entities"
)

type EtiquetasRepository interface {
	Create(ctx context.Context, e *entities.Etiqueta) error
	GetAll(ctx context.Context) ([]entities.Etiqueta, error)
	GetByID(ctx context.Context, id string) (*entities.Etiqueta, error)
	Update(ctx context.Context, e *entities.Etiqueta) error
	Delete(ctx context.Context, id string) error
}
