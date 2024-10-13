package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
		"bootstrap.servers": "localhost:9092",
		"group.id":          "cop-workflow-intake",
		"auto.offset.reset": "latest"})

	if err != nil {
		log.Fatalln("Failed to create producer: ", err)
	}
	err = consumer.SubscribeTopics([]string{"OrderSubmitted"}, nil)

	msg_count := 0
	run := true

	for run == true {
		ev := consumer.Poll(5000)
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))

			var order *order.Order
			unmarshallErr := json.Unmarshal(e.Value, &order)

			if unmarshallErr != nil {
				//TODO: THis is not proper handling of the error at all, come back later with a "plan"
				log.Fatalln("Bad Payload Received: " + string(e.Value))
			}
			startWorkflow(c, order)

			msg_count += 1
			if msg_count%MESSAGE_COMMIT_COUNT == 0 {
				consumer.Commit()
				msg_count = 0
			}

		case kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		default:
			fmt.Printf("Nothing to do %v\n", e)
		}
	}
}

func startWorkflow(c client.Client, order *order.Order) error {
	options := client.StartWorkflowOptions{
		ID:                    "order-submitted-" + order.OrderNumber,
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
