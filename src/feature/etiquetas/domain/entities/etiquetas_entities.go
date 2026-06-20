package entities

import "time"

type Etiqueta struct {
	ID        string    `json:"id"`
	Nombre    string    `json:"nombre"`
	Cantidad  int       `json:"cantidad"` // cuántos vehículos tienen esta etiqueta
	CreatedAt time.Time `json:"created_at"`
}
