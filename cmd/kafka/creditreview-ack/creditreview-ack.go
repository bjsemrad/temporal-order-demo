package main

import (
	"context"
	"encoding/json"
	"log"
	orderworkflowstep "temporal-order-demo/pkg/order-workflow/steps"
	"temporal-order-demo/pkg/services/creditreview"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.temporal.io/sdk/client"
)

const (
	MESSAGE_COMMIT_COUNT = 100
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092",
		"group.id":           "cop-workflow-fulfillmentack",
		"auto.offset.reset":  "latest",
		"enable.auto.commit": false})

	if err != nil {
		log.Fatalln("Failed to create producer: ", err)
	}
	defer consumer.Close()

	err = consumer.Subscribe("OrderCreditReview", nil)
	if err != nil {
		log.Fatalln("Failed to subscribe: ", err)
	}
	for {
		message, kError := consumer.ReadMessage(3000)
		if kError == nil {
			// fmt.Printf("%% Message on %s:\n%s\n", message.TopicPartition, string(message.Value))

			var creditReviewEvent *creditreview.CreditReviewEvent
			unmarshallErr := json.Unmarshal(message.Value, &creditReviewEvent)

			if unmarshallErr != nil {
				//TODO: THis is not proper handling of the error at all, come back later with a "plan"
				log.Fatalln("Bad Payload Received: " + string(message.Value))
			}

			var decision *orderworkflowstep.CreditReviewDecisionSignal
			if creditReviewEvent.Order.Total() > 2000 {
				decision = &orderworkflowstep.CreditReviewDecisionSignal{
					CreditDecision: orderworkflowstep.CreditExtensionDenied,
					Reviewier:      "xsed241",
					DecisionDate:   time.Now(),
				}
			} else {
				decision = &orderworkflowstep.CreditReviewDecisionSignal{
					CreditDecision: orderworkflowstep.CreditExtended,
					Reviewier:      "xsed241",
					NewLimit:       100000.00,
					DecisionDate:   time.Now(),
				}

			}
			err = c.SignalWorkflow(context.Background(), creditReviewEvent.WorkflowID, creditReviewEvent.RunID, orderworkflowstep.CreditReviewDecisionChannel, decision)

			//
			// if err != nil {
			// 	log.Fatalln("Unable to start the Workflow:", err)
			// }

			log.Printf("Signal Sent for CreditReviw " + string(decision.CreditDecision) + " order " + creditReviewEvent.Order.OrderNumber + " run id " + creditReviewEvent.RunID)

			consumer.Commit() //TODO: Should really consider moving this to message count based or duration based
		} else if kError.(kafka.Error).IsFatal() {
			log.Fatalln("Failed to read from kafka with a non-timeout error: ", err)
		} //TODO: Assume any other error is a retry for now
	}
}
