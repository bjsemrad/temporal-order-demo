package main

import (
	"context"
	"log"
	orderworkflowstep "temporal-order-demo/pkg/order-workflow/steps"

	"go.temporal.io/sdk/client"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()
	orderNumber := "order-submitted-2145137998316979647" //999
	runID := "91d9a221-78da-431a-a046-9300d486cb3b"
	log.Printf("Sending Signal for Order order " + orderNumber + " run id " + runID)

	err = c.SignalWorkflow(context.Background(), orderNumber, runID, orderworkflowstep.OrderFulfillmentConfirmedChannel, &orderworkflowstep.OrderConfirmedSignal{
		FulfillmentOrderNumber: "order-submitted-2145137998316979647",
	})

	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("Signal Sent for Order order " + orderNumber + " run id " + runID)
}
