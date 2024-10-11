package event

import (
	"context"
	"log"
	"temporal-order-demo/pkg/order"
)

type EventEmitOutput struct {
}

func EmitEvent(ctx context.Context, data order.Order) (bool, error) {
	log.Printf("Emitting Order Update Event $s.\n\n", data.OrderNumber)

	return true, nil
}
