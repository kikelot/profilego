package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"profilego/internal/client"

	"profilego/internal/repository"
	"profilego/internal/service"
	"profilego/internal/transport/http"
	"profilego/internal/transport/mq"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"

	"profilego/internal/middleware" //
)

func main() {
	// Configuración de PostgreSQL
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

	fmt.Println("✅ Conexión exitosa a PostgreSQL")

	//===========RABBITMQ=======================================================
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("No se pudo conectar a RabbitMQ: %s", err)
	}
	defer conn.Close()

	fmt.Println("✅ Conexión exitosa a RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Error abriendo canal RabbitMQ: %v", err)
	}
	defer ch.Close()

	//===========================================================================

	// Crear repositorios
	profileRepo := repository.NewProfileRepository(db)
	addressRepo := repository.NewAddressRepository(db)

	// Crear publisher para RabbitMQ
	rabbitPublisher := mq.NewPublisher(conn)

	// Crear servicios
	profileService := service.NewProfileRabbitService(*profileRepo, rabbitPublisher) // ✅ Ahora con RabbitMQ
	addressService := service.NewAddressService(*addressRepo, *profileRepo)

	// Crear conexión a RabbitMQ
	rabbitConn, err := mq.NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("❌ Error conectando a RabbitMQ: %v", err)
	}
	// Crear consumidor
	consumer := mq.NewConsumer(rabbitConn, profileService)

	// Iniciar consumidor en un Goroutine para que no bloquee el servidor HTTP
	go consumer.StartListening("direct_profile")

	// Configurar router con Gin
	router := gin.Default()

	// Inicializar handlers con los servicios correctos
	profileHandler := http.NewProfileHandler(*profileService) // ✅ Se usa el `profileService` con RabbitMQ
	addressHandler := http.NewAdressHandler(*addressService)

	authClient := client.NewAuthClient("http://localhost:3000")

	// Definir rutas
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(authClient)) // ⬅️ Aplica autenticación a todas las rutas dentro de /v1
	profileHandler.RegisterRoutes(api)
	addressHandler.RegisterRoutes(api)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("🚀 Servidor iniciado en el puerto %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Error al iniciar el servidor: %v", err)
	}
}
