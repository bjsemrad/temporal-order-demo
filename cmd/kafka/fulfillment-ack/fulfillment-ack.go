package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	eventactivity "temporal-order-demo/pkg/order-activities/event"
	orderworkflowstep "temporal-order-demo/pkg/order-workflow/steps"

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
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false})

	if err != nil {
		log.Fatalln("Failed to create producer: ", err)
	}
	defer consumer.Close()

	err = consumer.Subscribe("OrderReadyForFulfillment", nil)
	if err != nil {
		log.Fatalln("Failed to subscribe: ", err)
	}
	for {
		message, kError := consumer.ReadMessage(3000)
		if kError == nil {
			fmt.Printf("%% Message on %s:\n%s\n", message.TopicPartition, string(message.Value))

			var orderFulfillEvent *eventactivity.OrderFulfillmentEvent
			unmarshallErr := json.Unmarshal(message.Value, &orderFulfillEvent)

			if unmarshallErr != nil {
				//TODO: THis is not proper handling of the error at all, come back later with a "plan"
				log.Fatalln("Bad Payload Received: " + string(message.Value))
			}
			err = c.SignalWorkflow(context.Background(), orderFulfillEvent.WorkflowID, orderFulfillEvent.RunID, orderworkflowstep.OrderFulfillmentConfirmedChannel, &orderworkflowstep.OrderConfirmedSignal{
				FulfillmentOrderNumber: strconv.FormatInt(rand.Int63n(9223372036854775807), 10),
			})
			//
			// if err != nil {
			// 	log.Fatalln("Unable to start the Workflow:", err)
			// }

			log.Printf("Signal Sent for Order order " + orderFulfillEvent.Order.OrderNumber + " run id " + orderFulfillEvent.RunID)

			consumer.Commit() //TODO: Should really consider moving this to message count based or duration based
		} else if kError.(kafka.Error).IsFatal() {
			log.Fatalln("Failed to read from kafka with a non-timeout error: ", err)
		} //TODO: Assume any other error is a retry for now
	}
}
