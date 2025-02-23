# Kafka Coffee Brew System ☕

This project is a simple Kafka-based coffee order system using Go. It consists of a producer that sends coffee orders to Kafka and a consumer that processes these orders.

## Prerequisites

- Kafka and Zookeeper installed and running locally.
- Go installed on your system.
- Sarama library installed (`go get github.com/IBM/sarama`).

## Installation

1. Clone this repository:
   ```sh
   git clone https://github.com/likhithkp/coffee-brew-system
   cd coffee-brew-system
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

3. Start Kafka and Zookeeper:
   ```sh
   # Start Zookeeper
   bin/zookeeper-server-start.sh config/zookeeper.properties

   # Start Kafka
   bin/kafka-server-start.sh config/server.properties
   ```

4. Create a Kafka topic:
   ```sh
   bin/kafka-topics.sh --create --topic coffee_orders --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
   ```

## Running the Coffee Order Producer

1. Start the producer service:
   ```sh
   go run producer.go
   ```

2. Send a POST request to place a coffee order:
   ```sh
   curl -X POST http://localhost:3000/produce -H "Content-Type: application/json" -d '{"customer_name":"John Doe", "coffee_type":"Espresso"}'
   ```

3. The producer will send the order to Kafka and log a success message.

## Running the Coffee Order Consumer

1. Start the consumer service:
   ```sh
   go run consumer.go
   ```

2. The consumer will receive coffee orders and process them by printing messages to the console.

## Troubleshooting

- Ensure Kafka is running on `localhost:9092`.
- Verify the topic exists using:
  ```sh
  bin/kafka-topics.sh --list --bootstrap-server localhost:9092
  ```
- Check logs for any errors in producer or consumer.

## License

MIT License
