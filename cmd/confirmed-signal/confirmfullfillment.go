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
	orderNumber := "85150787987" //999
	runID := "cc91e9ab-4076-4905-8bfb-3a0350dcfe09"
	log.Printf("Sending Signal for Order order " + orderNumber + " run id " + runID)

	err = c.SignalWorkflow(context.Background(), "order-submitted-"+orderNumber, runID, orderworkflowstep.OrderFulfillmentConfirmedChannel, &orderworkflowstep.OrderConfirmedSignal{
		FulfillmentOrderNumber: "879879879879",
	})

	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("Signal Sent for Order order " + orderNumber + " run id " + runID)
}
