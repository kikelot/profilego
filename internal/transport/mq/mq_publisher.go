package mq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// Publisher define el publicador de RabbitMQ
type Publisher struct {
	conn *amqp.Connection
}

// NewPublisher crea una nueva instancia de Publisher
func NewPublisher(conn *amqp.Connection) *Publisher {
	return &Publisher{conn: conn}
}

// PublishMessage publica un mensaje en una cola de RabbitMQ
func (p *Publisher) PublishMessage(queueName string, message []byte) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)

	if err != nil {
		log.Println("Error publicando mensaje en RabbitMQ:", err)
		return err
	}

	log.Printf("✅ Mensaje publicado en la cola: %s", queueName)
	return nil
}

// PublishProfileLevelUpdate - Método faltante para implementar la interfaz service.Publisher
// func (p *Publisher) PublishProfileLevelUpdate(userID string, level int) error { //kike modifico parcialmente esto (lo comentado en el cuerpo)
func (p *Publisher) PublishProfilePoints(profileid string, profilepoints int) error {
	message := map[string]interface{}{
		//"userID": userID,
		//"level":  level, // Asegurarse de que es un entero
		"profileid":     profileid,
		"profilepoints": profilepoints, // Asegurarse de que es un entero
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.PublishMessage("direct_profile", body)
}

func (p *Publisher) PublishProfileLevelUpdate(profileid string, profileLevel int) error {
	message := map[string]interface{}{
		//"userID": userID,
		//"level":  level, // Asegurarse de que es un entero
		"profileid":    profileid,
		"profilelevel": profileLevel, // Asegurarse de que es un entero
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.PublishMessage("direct_profile", body)
}
