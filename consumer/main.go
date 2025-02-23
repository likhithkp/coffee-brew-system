package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "coffee_orders"

	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		log.Fatal("Failed to start consumer:", err)
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatal("Failed to get partitions:", err)
	}

	// Listen to all partitions
	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatal("Failed to start consuming:", err)
		}

		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("🔥 Received message: %s\n", string(msg.Value))
			}
		}(pc)
	}

	// Keep running until stopped
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan
}
