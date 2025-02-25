package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// InitDB inicializa la conexión con PostgreSQL
func InitDB() (*sql.DB, error) {
	// Obtener variables de entorno
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Crear la cadena de conexión
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Abrir conexión con PostgreSQL
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error al abrir la conexión con la base de datos: %v", err)
		return nil, err
	}

	// Verificar la conexión
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
		return nil, err
	}

	log.Println("✅ Conexión exitosa a PostgreSQL")
	return db, nil
}
