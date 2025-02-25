package mq

import (
	"context"
	"encoding/json"
	"log"
	"profilego/internal/domain"
	"profilego/internal/service"
)

// Consumer procesa los mensajes de RabbitMQ
type Consumer struct {
	Conn           *RabbitMQConnection
	ProfileService *service.ProfileService
}

// NewConsumer crea un nuevo consumidor
func NewConsumer(conn *RabbitMQConnection, profileService *service.ProfileService) *Consumer {
	return &Consumer{Conn: conn, ProfileService: profileService}
}

// StartListening inicia la escucha de eventos de RabbitMQ
func (c *Consumer) StartListening(queueName string) {
	msgs, err := c.Conn.Ch.Consume(
		queueName, // Cola
		"",        // Consumer
		true,      // Auto-ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Args
	)
	if err != nil {
		log.Fatalf("❌ Error consumiendo mensajes: %v", err)
	}

	log.Printf("📢 Escuchando mensajes en cola '%s'...", queueName)

	for msg := range msgs {
		var event struct {
			UserID        string `json:"userId"`
			ProfileID     string `json:"profileId"`
			ProfilePoints int    `json:"profilePoints"`
		}

		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("❌ Error procesando mensaje: %v", err)
			continue
		}

		/*log.Printf("📩 Mensaje recibido -> userId: %s, profileId: %s, profilePoints: %d",
		event.UserID, event.ProfileID, event.ProfilePoints)*/

		// 🔎 Obtener el perfil actual del usuario
		profile, err := c.ProfileService.Repo.GetByUserID(context.TODO(), event.UserID)
		if err != nil {
			log.Printf("❌ Error obteniendo perfil: %v", err)
			continue
		}

		// 📊 Calcular los nuevos puntos sumando los actuales con los recibidos
		//var totalPoints int
		//log.Println("PUNTOS", profile.ProfilePoints, event.ProfilePoints, totalPoints)
		//totalPoints = profile.ProfilePoints + event.ProfilePoints
		//log.Println("PUNTOS", profile.ProfilePoints, event.ProfilePoints, totalPoints)
		// 🔄 Mantener los valores actuales antes de actualizar
		profileUpdate := &domain.Profile{
			ProfilePoints: profile.ProfilePoints,
			ProfileLevel:  profile.ProfileLevel, // 🔹 Mantiene el nivel actual
		}

		// 🚀 Verificar si el total de puntos supera los 1000 y aún no ha subido de nivel
		if profile.ProfilePoints > 1000 {
			log.Printf("🎯 userId %s alcanzó más de 1000 puntos. Evaluando nivel...", event.UserID)

			// 📊 Subir el nivel a 1 solo si sigue en 0
			newLevel := profile.ProfileLevel + 1 // 🚀 Aquí forzamos el nivel a 1

			if newLevel != profile.ProfileLevel { // ✅ Evita actualizar si ya estaba en el nivel correcto
				profileUpdate.ProfileLevel = newLevel
				log.Printf("🎉 userId %s sube a nivel %d", event.UserID, newLevel)
			}
		}

		// 📝 Guardar la actualización solo si hubo cambios
		if profileUpdate.ProfilePoints != profile.ProfilePoints || profileUpdate.ProfileLevel != profile.ProfileLevel {
			err = c.ProfileService.UpdateProfileLevel(context.TODO(), event.UserID, profileUpdate)
			if err != nil {
				log.Printf("❌ Error actualizando nivel: %v", err)
			} else {
				log.Printf("✅ Nivel actualizado a %d para userId %s", profileUpdate.ProfileLevel, event.UserID)
			}
		} else {
			log.Printf("🔵 No hay cambios en puntos o nivel. No se actualiza el perfil.")
		}
	}
}
