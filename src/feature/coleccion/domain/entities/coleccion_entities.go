package entities

import "time"

type Vehiculo struct {
	ID        string     `json:"id"`
	Nombre    string     `json:"nombre"`
	Marca     string     `json:"marca"`
	Modelo    string     `json:"modelo"`
	Imagenes  []string   `json:"imagenes"`   // URLs públicas de Google Drive
	ImageIDs  []string   `json:"-"`          // IDs de Drive para poder eliminarlas
	Etiquetas []Etiqueta `json:"etiquetas"`  // etiquetas asignadas (para mostrar)
	EtiquetaIDs []string `json:"-"`          // IDs de etiquetas a asignar (entrada)
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Etiqueta es la referencia ligera a una etiqueta dentro de un vehículo.
type Etiqueta struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
}
