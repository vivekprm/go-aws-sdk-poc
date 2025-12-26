package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	event := "Hello, Kafka from AWS Lambda!"
	topic := "redundancy-events"
	bootstrapServers := "0.tcp.in.ngrok.io:13464"

	log.Printf("Producing to Kafka topic %s at %s\n", topic, bootstrapServers)

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"client.id":         "aws_lambda",
		"acks":              "all",
	})

	if err != nil {
		log.Fatalf("Failed to create producer: %s\n", err)
	}
	delivery_chan := make(chan kafka.Event, 10000)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(event)},
		delivery_chan,
	)

	if err != nil {
		log.Fatalf("Failed to produce message: %s\n", err)
	}
	e := <-delivery_chan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		log.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	close(delivery_chan)
}
