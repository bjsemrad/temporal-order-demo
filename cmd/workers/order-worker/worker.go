package main

import (
	"log"
	creditreviewactivity "temporal-order-demo/pkg/order-activities/creditreview"
	eventactivity "temporal-order-demo/pkg/order-activities/event"
	fraudactivity "temporal-order-demo/pkg/order-activities/fraud"
	orderworkflow "temporal-order-demo/pkg/order-workflow"
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

	w := worker.New(c, orderworkflowqueues.OrderIntakeTaskQueueName, worker.Options{})

	w.RegisterWorkflow(orderworkflow.ProcessOrder)
	var eventActivity *eventactivity.EventBroker
	var creditActivity *creditreviewactivity.CreditActivity
	w.RegisterActivity(eventActivity)
	w.RegisterActivity(fraudactivity.CheckOrderFraudulent)
	w.RegisterActivity(creditActivity)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
