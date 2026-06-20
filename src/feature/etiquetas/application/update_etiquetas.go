package application

import (
	"context"
	"errors"
	"strings"

	"coleccionbackend/src/feature/etiquetas/domain"
)

type UpdateEtiquetaUseCase struct {
	repo domain.EtiquetasRepository
}

func NewUpdateEtiquetaUseCase(r domain.EtiquetasRepository) *UpdateEtiquetaUseCase {
	return &UpdateEtiquetaUseCase{repo: r}
}

func (uc *UpdateEtiquetaUseCase) Execute(ctx context.Context, id, nombre string) error {
	nombre = strings.TrimSpace(nombre)
	if nombre == "" {
		return errors.New("el nombre de la etiqueta es requerido")
	}

	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	existing.Nombre = nombre
	return uc.repo.Update(ctx, existing)
}
