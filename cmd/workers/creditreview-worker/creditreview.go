package main

import (
	"log"
	creditreviewactivity "temporal-order-demo/pkg/order-activities/creditreview"
	orderworkflowqueues "temporal-order-demo/pkg/order-workflow/queues"
	"temporal-order-demo/pkg/services/creditreview"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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

	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}
	defer kafkaProducer.Close()

	// This worker hosts both Workflow and Activity functions.
	w.RegisterActivity(&creditreviewactivity.CreditActivity{
		CreditClient: &creditreview.CreditReviewClient{
			KafkaProducer: kafkaProducer,
		},
	})

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
