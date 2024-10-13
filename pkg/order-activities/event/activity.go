package eventactivity

import (
	"log"
	"temporal-order-demo/pkg/order"

	"go.temporal.io/sdk/workflow"
)

type EventEmitOutput struct {
	Success bool
	Order   order.Order
}

func EmitStatusUpdateEvent(ctx workflow.Context, data order.Order) (EventEmitOutput, error) {
	log.Printf("Emitting Order %s Update Status: %s Event. \n\n", data.OrderNumber, data.Status)
	result := EventEmitOutput{Success: true, Order: data} //TODO: Make this talk to kafka
	return result, nil
}

func EmitFullfilmentEvent(ctx workflow.Context, wfID string, runID string, data order.Order) (EventEmitOutput, error) {
	log.Printf("Emitting Order Fulfillment Event %s. \n\n", data.OrderNumber)
	result := EventEmitOutput{Success: true, Order: data} //TODO: Make this talk to kafka
	return result, nil
}
