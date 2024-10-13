package pricingactivity

import (
	"context"
	"log"
	"temporal-order-demo/pkg/order"
	"temporal-order-demo/pkg/services/pricing"
)

func PriceOrder(ctx context.Context, order *order.Order) (pricing.OrderLinePricing, error) {
	log.Printf("Pricing order %s.\n\n", order.OrderNumber)
	pricingClient := pricing.InitializeClient()
	result, err := pricingClient.PriceOrder(order)
	log.Printf("Pricing for order %s complete.\n\n", order.OrderNumber)
	return result, err
}
