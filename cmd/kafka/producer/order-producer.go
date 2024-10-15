package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"temporal-order-demo/pkg/order"
	orderstatus "temporal-order-demo/pkg/order/status"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	topic := "OrderSubmitted"
	for i := 0; i < 20; i++ {
		orderNumber := strconv.FormatInt(rand.Int63n(9223372036854775807), 10)
		var channel = "gcom"
		if i%2 == 0 {
			channel = "custservice"
		}
		input := &order.Order{
			Channel:     channel,
			OrderNumber: orderNumber,
			Status: &orderstatus.OrderStatus{
				Code:   orderstatus.Submitted,
				Reason: "submitted",
			},
			LastUpdated: time.Now(),
			Lines: []*order.OrderLine{
				{LineNumber: 1, Product: "ABC123", Quantity: rand.Intn(10), Price: 10.32},
				{LineNumber: 2, Product: "XYZ", Quantity: rand.Intn(10), Price: 1.57},
				{LineNumber: 3, Product: "3B", Quantity: rand.Intn(10), Price: 10.32},
				{LineNumber: 4, Product: "Z4", Quantity: rand.Intn(10), Price: 100.57},
				{LineNumber: 5, Product: "A5", Quantity: rand.Intn(10), Price: 200.32},
				// {LineNumber: 5, Product: "Y6", Quantity: 3,  Price: 100.7},
			},
			Payment: &order.Payment{
				AccountNumber: "13676876876",
			},
		}

		message, err := json.Marshal(input)
		if err != nil {
			log.Printf("Failed to serialized order")
		}
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          message,
		}, nil)
	}

	// Wait for message deliveries before shutting down
	p.Flush(5000)
}
