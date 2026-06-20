package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"coleccionbackend/src/feature/etiquetas/domain/entities"
)

type SupabaseEtiquetasRepository struct {
	db *sql.DB
}

func NewSupabaseEtiquetasRepository(db *sql.DB) *SupabaseEtiquetasRepository {
	return &SupabaseEtiquetasRepository{db: db}
}

func (r *SupabaseEtiquetasRepository) Create(ctx context.Context, e *entities.Etiqueta) error {
	query := `INSERT INTO etiquetas (id, nombre, created_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, e.ID, e.Nombre, e.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return fmt.Errorf("ya existe una etiqueta con ese nombre")
		}
		return fmt.Errorf("error creando etiqueta: %w", err)
	}
	return nil
}

// GetAll retorna las etiquetas junto con la cantidad de vehículos que tiene cada una.
func (r *SupabaseEtiquetasRepository) GetAll(ctx context.Context) ([]entities.Etiqueta, error) {
	query := `
		SELECT e.id, e.nombre, e.created_at, COUNT(ve.vehiculo_id) AS cantidad
		FROM etiquetas e
		LEFT JOIN vehiculo_etiquetas ve ON ve.etiqueta_id = e.id
		GROUP BY e.id, e.nombre, e.created_at
		ORDER BY e.nombre ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error consultando etiquetas: %w", err)
	}
	defer rows.Close()

	var etiquetas []entities.Etiqueta
	for rows.Next() {
		var e entities.Etiqueta
		if err := rows.Scan(&e.ID, &e.Nombre, &e.CreatedAt, &e.Cantidad); err != nil {
			return nil, err
		}
		etiquetas = append(etiquetas, e)
	}
	return etiquetas, rows.Err()
}

func (r *SupabaseEtiquetasRepository) GetByID(ctx context.Context, id string) (*entities.Etiqueta, error) {
	query := `SELECT id, nombre, created_at FROM etiquetas WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var e entities.Etiqueta
	err := row.Scan(&e.ID, &e.Nombre, &e.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("etiqueta no encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("error consultando etiqueta: %w", err)
	}
	return &e, nil
}

func (r *SupabaseEtiquetasRepository) Update(ctx context.Context, e *entities.Etiqueta) error {
	res, err := r.db.ExecContext(ctx, `UPDATE etiquetas SET nombre = $1 WHERE id = $2`, e.Nombre, e.ID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return fmt.Errorf("ya existe una etiqueta con ese nombre")
		}
		return fmt.Errorf("error actualizando etiqueta: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("etiqueta no encontrada")
	}
	return nil
}

func (r *SupabaseEtiquetasRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM etiquetas WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error eliminando etiqueta: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("etiqueta no encontrada")
	}
	return nil
}
