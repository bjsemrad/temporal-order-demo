package settingsactivity

import (
	"context"
	"log"
	"temporal-order-demo/pkg/order"
)

func ApplySettings(ctx context.Context, custOrder *order.Order) error {
	log.Printf("Applying Settings to order %s.\n\n", custOrder.OrderNumber)
	//Picking arbitrary means for controlling this obviously this doesn't make sense
	if custOrder.Total() > 1000 {
		custOrder.Settings = &order.OrderSettings{
			LiftGateRequired:     true,
			PackingListInEachBox: false,
		}
	}
	return nil
}
