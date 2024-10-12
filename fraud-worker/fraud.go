package main

import (
	"log"
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

	w := worker.New(c, processorqueue.FraudTaskQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterActivity(fraud.CheckOrderFraudulent)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
