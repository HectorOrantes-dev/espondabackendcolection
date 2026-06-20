package domain

import (
	"context"
	"coleccionbackend/src/feature/coleccion/domain/entities"
)

type ColeccionRepository interface {
	Create(ctx context.Context, v *entities.Vehiculo) error
	// GetAll lista los vehículos. Si etiquetaFiltro no está vacío, filtra por
	// el nombre de la etiqueta (búsqueda parcial, insensible a mayúsculas).
	GetAll(ctx context.Context, etiquetaFiltro string) ([]entities.Vehiculo, error)
	GetByID(ctx context.Context, id string) (*entities.Vehiculo, error)
	Update(ctx context.Context, v *entities.Vehiculo) error
	Delete(ctx context.Context, id string) error
}

// ImageService es el puerto para el servicio de almacenamiento de imágenes.
type ImageService interface {
	Upload(filename string, content []byte, mimeType string) (url string, fileID string, err error)
	Download(fileID string) ([]byte, error)
	Delete(fileID string) error
}
