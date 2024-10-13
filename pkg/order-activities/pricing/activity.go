package pricingactivity

import (
	"log"
	"temporal-order-demo/pkg/order"
	"temporal-order-demo/pkg/services/pricing"

	"go.temporal.io/sdk/workflow"
)

func PriceOrder(ctx workflow.Context, order *order.Order) (pricing.OrderLinePricing, error) {
	log.Printf("Pricing order %s.\n\n", order.OrderNumber)
	pricingClient := pricing.InitializeClient()
	result, err := pricingClient.PriceOrder(order)
	log.Printf("Pricing for order %s complete.\n\n", order.OrderNumber)
	return result, err
}
