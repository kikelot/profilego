/*package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "root"
	password = "root"
	dbname   = "profileDB"
)

func main() {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Conexión exitosa a PostgreSQL!")
}*/

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	//"profilego/pkg/db"
	//"profilego/pkg/logger"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"profilego/internal/repository"
	"profilego/internal/service"
	"profilego/internal/transport/http"
	"profilego/internal/transport/mq"
)

func main() {
	// Configuración de PostgreSQL
	/*dbHost := os.Getenv("127.0.0.1")
	dbPort := os.Getenv("5432")
	dbUser := os.Getenv("root")
	dbPassword := os.Getenv("root")
	dbName := os.Getenv("profileDB")*/

	const (
		dbHost     = "127.0.0.1"
		dbPort     = "5432"
		dbUser     = "root"
		dbPassword = "root"
		dbName     = "profileDB"
	)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	fmt.Println("Conexión exitosa a PostgreSQL")

	//===========RABBIT=======================================================
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("No se pudo conectar a RabbitMQ: %s", err)
	}
	defer conn.Close()

	log.Println("✅ Conexión exitosa a RabbitMQ")

	// Inicializar Logger
	//logger.InitLogger()

	// Crear repositorios
	profileRepo := repository.NewProfileRepository(db)
	addressRepo := repository.NewAddressRepository(db)

	// Crear servicios
	profileService := service.NewProfileService(*profileRepo)
	addressService := service.NewAddressService(*addressRepo)

	// Configurar router con Gin
	router := gin.Default()

	// Inicializar handlers
	profileHandler := http.NewProfileHandler(*profileService)
	addressHandler := http.NewAdressHandler(*addressService)

	// Crear publisher
	rabbitPublisher := mq.NewPublisher(conn)

	// Crear servicio con RabbitMQ
	profileRabbitService := service.NewProfileRabbitService(*profileRepo, rabbitPublisher) // Pasamos correctamente los valores

	// Usamos profileRabbitService en algún lugar para evitar el error de variable no usada
	_ = profileRabbitService

	// Definir rutas
	api := router.Group("/api")
	profileHandler.RegisterRoutes(api)
	addressHandler.RegisterRoutes(api)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Servidor iniciado en el puerto %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
