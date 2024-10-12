package kafkaDB

import "github.com/segmentio/kafka-go"

func defaultKafkaParams(w *kafka.Writer) {
	if w.Balancer == nil {
		w.Balancer = &kafka.RoundRobin{}
	}
}

func NewKafkaConn(w *kafka.Writer) *kafka.Writer {
	defaultKafkaParams(w)
	conn := w
	return conn
}
