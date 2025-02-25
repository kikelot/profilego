package mq

import (
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQConnection mantiene la conexión con RabbitMQ
type RabbitMQConnection struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

// NewRabbitMQ crea una nueva conexión
func NewRabbitMQ(url string) (*RabbitMQConnection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("❌ Error conectando a RabbitMQ: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Error abriendo canal RabbitMQ: %v", err)
		return nil, err
	}

	return &RabbitMQConnection{Conn: conn, Ch: ch}, nil
}

// Close cierra la conexión con RabbitMQ
func (r *RabbitMQConnection) Close() {
	r.Ch.Close()
	r.Conn.Close()
}
