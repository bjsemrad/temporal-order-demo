package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"temporal-order-demo/pkg/order"
	orderworkflow "temporal-order-demo/pkg/order-workflow"
	orderworkflowqueues "temporal-order-demo/pkg/order-workflow/queues"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.temporal.io/api/enums/v1"
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
		"group.id":           "cop-workflow-intake",
		"auto.offset.reset":  "latest",
		"enable.auto.commit": false})

	if err != nil {
		log.Fatalln("Failed to create producer: ", err)
	}
	defer consumer.Close()

	err = consumer.Subscribe("OrderSubmitted", nil)
	if err != nil {
		log.Fatalln("Failed to subscribe: ", err)
	}
	for {
		message, kError := consumer.ReadMessage(3000)
		if kError == nil {
			fmt.Printf("%% Message on %s:\n%s\n", message.TopicPartition, string(message.Value))

			var order *order.Order
			unmarshallErr := json.Unmarshal(message.Value, &order)

			if unmarshallErr != nil {
				//TODO: THis is not proper handling of the error at all, come back later with a "plan"
				log.Fatalln("Bad Payload Received: " + string(message.Value))
			}
			startWorkflow(c, order)

			consumer.Commit() //TODO: Should really consider moving this to message count based or duration based
		} else if kError.(kafka.Error).IsFatal() {
			log.Fatalln("Failed to read from kafka with a non-timeout error: ", err)
		} //TODO: Assume any other error is a retry for now
	}
}

func startWorkflow(c client.Client, order *order.Order) error {
	options := client.StartWorkflowOptions{
		ID:                    "order-" + order.Channel + "-" + order.OrderNumber,
		TaskQueue:             orderworkflowqueues.OrderIntakeTaskQueueName,
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY,
	}

	log.Printf("Starting order processing for order " + order.OrderNumber)

	we, err := c.ExecuteWorkflow(context.Background(), options, orderworkflow.ProcessOrder, order)
	if err != nil {
		return err
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())
	return nil
}
