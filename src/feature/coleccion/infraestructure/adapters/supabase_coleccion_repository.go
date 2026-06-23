package adapters

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"

	"coleccionbackend/src/feature/coleccion/domain/entities"
)

type SupabaseColeccionRepository struct {
	db *sql.DB
}

func NewSupabaseColeccionRepository(db *sql.DB) *SupabaseColeccionRepository {
	return &SupabaseColeccionRepository{db: db}
}

func (r *SupabaseColeccionRepository) Create(ctx context.Context, v *entities.Vehiculo) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error iniciando transacción: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // no-op si ya se hizo commit

	query := `
		INSERT INTO vehiculos (id, nombre, marca, modelo, precio, imagenes, image_ids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = tx.ExecContext(ctx, query,
		v.ID, v.Nombre, v.Marca, v.Modelo, v.Precio,
		pq.Array(nonNil(v.Imagenes)), pq.Array(nonNil(v.ImageIDs)),
		v.CreatedAt, v.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creando vehículo: %w", err)
	}

	if err := setEtiquetas(ctx, tx, v.ID, v.EtiquetaIDs); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SupabaseColeccionRepository) Update(ctx context.Context, v *entities.Vehiculo) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error iniciando transacción: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	query := `
		UPDATE vehiculos
		SET nombre = $1, marca = $2, modelo = $3, precio = $4, imagenes = $5, image_ids = $6, updated_at = $7
		WHERE id = $8
	`
	res, err := tx.ExecContext(ctx, query,
		v.Nombre, v.Marca, v.Modelo, v.Precio,
		pq.Array(nonNil(v.Imagenes)), pq.Array(nonNil(v.ImageIDs)),
		v.UpdatedAt, v.ID,
	)
	if err != nil {
		return fmt.Errorf("error actualizando vehículo: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("vehículo no encontrado")
	}

	// Reemplazar las etiquetas: borrar las actuales y volver a insertar.
	if _, err := tx.ExecContext(ctx, `DELETE FROM vehiculo_etiquetas WHERE vehiculo_id = $1`, v.ID); err != nil {
		return fmt.Errorf("error limpiando etiquetas: %w", err)
	}
	if err := setEtiquetas(ctx, tx, v.ID, v.EtiquetaIDs); err != nil {
		return err
	}

	return tx.Commit()
}

// setEtiquetas inserta las relaciones vehículo-etiqueta dentro de una transacción.
func setEtiquetas(ctx context.Context, tx *sql.Tx, vehiculoID string, etiquetaIDs []string) error {
	for _, etID := range etiquetaIDs {
		if etID == "" {
			continue
		}
		_, err := tx.ExecContext(ctx,
			`INSERT INTO vehiculo_etiquetas (vehiculo_id, etiqueta_id) VALUES ($1, $2)
			 ON CONFLICT DO NOTHING`,
			vehiculoID, etID,
		)
		if err != nil {
			return fmt.Errorf("error asignando etiqueta: %w", err)
		}
	}
	return nil
}

// nonNil garantiza un arreglo vacío en vez de nil, para no insertar NULL
// en columnas TEXT[] con restricción NOT NULL (las imágenes son opcionales).
func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func (r *SupabaseColeccionRepository) GetAll(ctx context.Context, etiquetaFiltro string) ([]entities.Vehiculo, error) {
	query := `
		SELECT v.id, v.nombre, v.marca, v.modelo, v.precio, v.imagenes, v.image_ids, v.created_at, v.updated_at,
		       COALESCE(
		           json_agg(json_build_object('id', e.id, 'nombre', e.nombre))
		           FILTER (WHERE e.id IS NOT NULL), '[]'
		       ) AS etiquetas
		FROM vehiculos v
		LEFT JOIN vehiculo_etiquetas ve ON ve.vehiculo_id = v.id
		LEFT JOIN etiquetas e ON e.id = ve.etiqueta_id
	`
	var args []any
	if etiquetaFiltro != "" {
		// Filtra vehículos que tengan al menos una etiqueta que coincida,
		// pero sigue mostrando todas las etiquetas de esos vehículos.
		query += `
		WHERE v.id IN (
			SELECT ve2.vehiculo_id FROM vehiculo_etiquetas ve2
			JOIN etiquetas e2 ON e2.id = ve2.etiqueta_id
			WHERE e2.nombre ILIKE $1
		)`
		args = append(args, "%"+etiquetaFiltro+"%")
	}
	query += `
		GROUP BY v.id
		ORDER BY v.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error consultando vehículos: %w", err)
	}
	defer rows.Close()

	var vehiculos []entities.Vehiculo
	for rows.Next() {
		v, err := scanVehiculoConEtiquetas(rows)
		if err != nil {
			return nil, err
		}
		vehiculos = append(vehiculos, v)
	}
	return vehiculos, rows.Err()
}

func (r *SupabaseColeccionRepository) GetByID(ctx context.Context, id string) (*entities.Vehiculo, error) {
	query := `
		SELECT v.id, v.nombre, v.marca, v.modelo, v.precio, v.imagenes, v.image_ids, v.created_at, v.updated_at,
		       COALESCE(
		           json_agg(json_build_object('id', e.id, 'nombre', e.nombre))
		           FILTER (WHERE e.id IS NOT NULL), '[]'
		       ) AS etiquetas
		FROM vehiculos v
		LEFT JOIN vehiculo_etiquetas ve ON ve.vehiculo_id = v.id
		LEFT JOIN etiquetas e ON e.id = ve.etiqueta_id
		WHERE v.id = $1
		GROUP BY v.id
	`
	row := r.db.QueryRowContext(ctx, query, id)
	v, err := scanVehiculoConEtiquetas(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("vehículo no encontrado")
	}
	if err != nil {
		return nil, fmt.Errorf("error consultando vehículo: %w", err)
	}
	return &v, nil
}

func (r *SupabaseColeccionRepository) Delete(ctx context.Context, id string) error {
	// Las relaciones en vehiculo_etiquetas se borran en cascada.
	res, err := r.db.ExecContext(ctx, `DELETE FROM vehiculos WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error eliminando vehículo: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("vehículo no encontrado")
	}
	return nil
}

// Resumen retorna el conteo de vehículos, el valor total y el precio promedio.
func (r *SupabaseColeccionRepository) Resumen(ctx context.Context) (*entities.ResumenColeccion, error) {
	query := `
		SELECT COUNT(*)                         AS cantidad,
		       COALESCE(SUM(precio), 0)         AS total,
		       COALESCE(AVG(precio), 0)         AS promedio
		FROM vehiculos
	`
	var res entities.ResumenColeccion
	err := r.db.QueryRowContext(ctx, query).Scan(
		&res.CantidadVehiculos, &res.ValorTotal, &res.PrecioPromedio,
	)
	if err != nil {
		return nil, fmt.Errorf("error calculando resumen: %w", err)
	}
	return &res, nil
}

// scanner abstrae sql.Row y sql.Rows para reutilizar la lógica de scan.
type scanner interface {
	Scan(dest ...any) error
}

func scanVehiculoConEtiquetas(s scanner) (entities.Vehiculo, error) {
	var v entities.Vehiculo
	var imagenes, imageIDs pq.StringArray
	var etiquetasJSON []byte

	err := s.Scan(&v.ID, &v.Nombre, &v.Marca, &v.Modelo, &v.Precio, &imagenes, &imageIDs,
		&v.CreatedAt, &v.UpdatedAt, &etiquetasJSON)
	if err != nil {
		return v, err
	}

	v.Imagenes = []string(imagenes)
	v.ImageIDs = []string(imageIDs)

	if len(etiquetasJSON) > 0 {
		_ = json.Unmarshal(etiquetasJSON, &v.Etiquetas)
	}
	if v.Etiquetas == nil {
		v.Etiquetas = []entities.Etiqueta{}
	}
	return v, nil
}
