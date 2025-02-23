package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/IBM/sarama"
)

type Order struct {
	CustomerName string `json:"customer_name"`
	CoffeeType   string `json:"coffee_type"`
}

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)
}

func PushOrderToQueue(topic string, message []byte) error {
	brokers := []string{"localhost:9092"}
	producer, err := ConnectProducer(brokers)
	if err != nil {
		return err
	}

	defer producer.Close()
	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(&msg)
	if err != nil {
		return err
	}

	log.Printf("Order is stored in topic (%s)/ partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

func produceCoffee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodOptions {
		http.Error(w, "Not a valid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	order := new(Order)

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	byteData, err := json.Marshal(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = PushOrderToQueue("Coffee Orders", byteData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"msg":     "Order for " + order.CustomerName + " placed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "Error placing orders", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/produce", produceCoffee)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
