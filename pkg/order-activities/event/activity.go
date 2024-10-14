package eventactivity

import (
	"context"
	"encoding/json"
	"log"
	"temporal-order-demo/pkg/order"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type EventEmitOutput struct {
	Success bool
	Order   order.Order
}

type EventBroker struct {
	KafkaProducer *kafka.Producer
}

type OrderFulfillmentEvent struct {
	Order      *order.Order
	WorkflowID string
	RunID      string
}

func (eb *EventBroker) EmitStatusUpdateEvent(ctx context.Context, data order.Order) (EventEmitOutput, error) {
	log.Printf("Emitting Order %s Update Status: %s Event. \n\n", data.OrderNumber, data.Status)
	message, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to serialized order")
		return EventEmitOutput{Success: false}, err
	}

	topic := "OrderUpdated"
	eb.KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
	eb.KafkaProducer.Flush(3000)

	result := EventEmitOutput{Success: true, Order: data} //TODO: Make this talk to kafka
	return result, nil
}

func (eb *EventBroker) EmitFullfilmentEvent(ctx context.Context, wfID string, runID string, data *order.Order) (EventEmitOutput, error) {
	log.Printf("Emitting Order Fulfillment Event %s. \n\n", data.OrderNumber)
	event := OrderFulfillmentEvent{
		Order:      data,
		WorkflowID: wfID,
		RunID:      runID,
	}
	message, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to serialized order")
		return EventEmitOutput{Success: false}, err
	}

	topic := "OrderReadyForFulfillment"
	eb.KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
	eb.KafkaProducer.Flush(3000)
	result := EventEmitOutput{Success: true, Order: *data}
	return result, nil
}
