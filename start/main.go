package main

import (
	"context"
	"log"
	"temporal-order-demo/pkg/order"
	"temporal-order-demo/pkg/processing"

	"go.temporal.io/sdk/client"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()
	orderNumber := "85150787987"
	input := order.Order{
		OrderNumber: orderNumber,
		Lines: []*order.OrderLine{
			{LineNumber: 1, Product: "ABC123", Quantity: 10, Price: 10.32},
			{LineNumber: 2, Product: "XYZ", Quantity: 2, Price: 1.57},
		},
	}

	options := client.StartWorkflowOptions{
		ID:        "submit-order-" + orderNumber,
		TaskQueue: order.OrderIntakeTaskQueueName,
	}

	log.Printf("Starting order processing for order " + orderNumber)

	we, err := c.ExecuteWorkflow(context.Background(), options, processing.ProcessOrder, input)
	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	var result string

	err = we.Get(context.Background(), &result)

	if err != nil {
		log.Fatalln("Unable to get Workflow result:", err)
	}

	log.Println(result)
}
