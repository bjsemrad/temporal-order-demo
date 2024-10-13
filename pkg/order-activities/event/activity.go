package eventactivity

import (
	"context"
	"log"
	"temporal-order-demo/pkg/order"
)

type EventEmitOutput struct {
	Success bool
	Order   order.Order
}

func EmitEvent(ctx context.Context, data order.Order) (EventEmitOutput, error) {
	log.Printf("Emitting Order %s Update Status: %s Event. \n\n", data.OrderNumber, data.Status)
	result := EventEmitOutput{Success: true, Order: data}
	return result, nil
}
