package main

import (
	"log"
	creditreviewactivity "temporal-order-demo/pkg/order-activities/creditreview"
	orderworkflowqueues "temporal-order-demo/pkg/order-workflow/queues"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, orderworkflowqueues.CreditReviewTaskQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterActivity(creditreviewactivity.ValidateAndReserveCredit)
	w.RegisterActivity(creditreviewactivity.SubmitCreditReview)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
