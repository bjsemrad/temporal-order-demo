package event

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
	log.Printf("Emitting Order Update Event. %s \n", data.OrderNumber)
	result := EventEmitOutput{Success: true, Order: data}
	return result, nil
}
