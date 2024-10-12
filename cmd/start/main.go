package main

import (
	"context"
	"log"
	"temporal-order-demo/pkg/order"
	orderworkflow "temporal-order-demo/pkg/order-workflow"
	orderworkflowqueues "temporal-order-demo/pkg/order-workflow/queues"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()
	orderNumber := "99985150787987"
	input := &order.Order{
		OrderNumber: orderNumber,
		Status:      order.Submitted,
		LastUpdated: time.Now(),
		Lines: []*order.OrderLine{
			{LineNumber: 1, Product: "ABC123", Quantity: 10, Price: 10.32},
			{LineNumber: 2, Product: "XYZ", Quantity: 2, Price: 1.57},
			{LineNumber: 3, Product: "3B", Quantity: 14, Price: 10.32},
			{LineNumber: 4, Product: "Z4", Quantity: 23, Price: 100.57},
			{LineNumber: 5, Product: "A5", Quantity: 13, Price: 200.32},
			// {LineNumber: 5, Product: "Y6", Quantity: 3, Price: 100.57},
		},
		Payment: &order.Payment{
			AccountNumber: "13676876876",
		},
	}

	options := client.StartWorkflowOptions{
		ID:        "submit-order-" + orderNumber + time.Now().String(),
		TaskQueue: orderworkflowqueues.OrderIntakeTaskQueueName,
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{"FraudDetectedError"},
		},
	}

	log.Printf("Starting order processing for order " + orderNumber)

	we, err := c.ExecuteWorkflow(context.Background(), options, orderworkflow.ProcessOrder, input)
	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	var result order.Order

	err = we.Get(context.Background(), &result)

	if err != nil {
		log.Fatalln("Unable to get Workflow result:", err)
	}

	log.Println(result)
}
