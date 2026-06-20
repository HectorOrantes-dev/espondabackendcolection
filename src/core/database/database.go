package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func NewConnection() (*sql.DB, error) {
	dsn := os.Getenv("SUPABASE_DB_URL")
	if dsn == "" {
		return nil, fmt.Errorf("SUPABASE_DB_URL no está configurado")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error abriendo conexión: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error conectando a la base de datos: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	// El pooler de Supabase cierra conexiones inactivas, lo que provoca
	// "connection reset" al reutilizar una conexión muerta. Reciclamos las
	// conexiones antes de que el pooler las cierre del lado del servidor.
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	return db, nil
}
