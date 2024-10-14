package creditreview

import (
	"encoding/json"
	"log"
	"temporal-order-demo/pkg/order"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type CreditReviewClient struct {
	KafkaProducer *kafka.Producer
}

type CreditApprovalDecision struct {
	Approved bool
	Limit    float64
}

func InitializeClient() *CreditReviewClient {
	return &CreditReviewClient{}
}

type CreditReservationResult struct {
	CreditReserved  bool
	AvailableCredit float64
}

type CreditReviewEvent struct {
	Order      *order.Order
	WorkflowID string
	RunID      string
}

func (c *CreditReviewClient) ReserveCredit(order order.Order) (CreditReservationResult, error) {
	if order.Total() > 1000 {
		return CreditReservationResult{
			CreditReserved:  false,
			AvailableCredit: 100.00,
		}, nil
	}
	return CreditReservationResult{
		CreditReserved:  true,
		AvailableCredit: 5000.00,
	}, nil
}

func (c *CreditReviewClient) InitiateCreditReview(workflowId string, runId string, order *order.Order) error {
	event := CreditReviewEvent{
		Order:      order,
		WorkflowID: workflowId,
		RunID:      runId,
	}
	message, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to serialized order")
		return err
	}

	topic := "OrderCreditReview"
	c.KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
	c.KafkaProducer.Flush(3000)

	return nil
}
