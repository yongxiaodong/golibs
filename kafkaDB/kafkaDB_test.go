package kafkaDB

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"testing"
)

func TestNewKafkaConn(t *testing.T) {
	c := NewKafkaConn(&kafka.Writer{
		Addr:  kafka.TCP("127.0.0.1:9092"),
		Topic: "test",
	})
	err := c.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte("sdjfksd"),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
}
