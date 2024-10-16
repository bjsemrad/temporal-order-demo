package main

import (
	"context"
	"log"
	orderworkflowstep "temporal-order-demo/pkg/order-workflow/steps"
	"time"

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
	runID := "49a400f4-585a-4515-ba3f-6abb72c08510"
	log.Printf("Sending Signal for Order order " + orderNumber + " run id " + runID)

	err = c.SignalWorkflow(context.Background(), "order-submitted-"+orderNumber, runID, orderworkflowstep.CreditReviewDecisionChannel, &orderworkflowstep.CreditReviewDecisionSignal{
		CreditDecision: orderworkflowstep.CreditExtended,
		Reviewier:      "xsed241",
		NewLimit:       100000.00,
		DecisionDate:   time.Now(),
	})

	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("Signal Sent for Order order " + orderNumber + " run id " + runID)
}
