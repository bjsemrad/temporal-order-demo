package orderworkflowstep

import (
	"log"
	"temporal-order-demo/pkg/order"
	pricingactivity "temporal-order-demo/pkg/order-activities/pricing"
	orderworkflowqueues "temporal-order-demo/pkg/order-workflow/queues"
	"temporal-order-demo/pkg/services/pricing"

	"go.temporal.io/sdk/workflow"
)

func EnsureOrderIsPriced(ctx workflow.Context, order *order.Order) error {
	var pricedLines pricing.OrderLinePricing
	// Execute Order Pricing

	if containsUnPriceLines(order) {
		log.Printf("Order requires pricing, pricing order.")
		priceErr := workflow.ExecuteActivity(
			workflow.WithTaskQueue(ctx, orderworkflowqueues.PricingTaskQueueName),
			pricingactivity.PriceOrder,
			order).Get(ctx, &pricedLines)

		if priceErr != nil {
			return priceErr
		}
		updateOrderPricing(order, pricedLines)
	}
	return nil
}

func containsUnPriceLines(order *order.Order) bool {
	for _, line := range order.Lines {
		if line.Price == 0 {
			return true
		}
	}
	return false
}

func updateOrderPricing(order *order.Order, pricedLines pricing.OrderLinePricing) {
	for _, line := range order.Lines {
		if value, ok := pricedLines.LinePricing[line.Product]; ok {
			line.Price = value
		}
	}

}
