package main

import (
	"log"
	"temporal-order-demo/pkg/processing"
	"temporal-order-demo/pkg/processing/event"
	"temporal-order-demo/pkg/processing/fraud"
	processorqueue "temporal-order-demo/pkg/queue"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, processorqueue.OrderIntakeTaskQueueName, worker.Options{})

	w.RegisterWorkflow(processing.ProcessOrder)
	w.RegisterActivity(event.EmitEvent)
	w.RegisterActivity(fraud.CheckOrderFraudulent)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
