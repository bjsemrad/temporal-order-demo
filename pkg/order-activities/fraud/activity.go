package fraudactivity

import (
	"log"
	"temporal-order-demo/pkg/order"
	"temporal-order-demo/pkg/services/fraud"

	"go.temporal.io/sdk/workflow"
)

func CheckOrderFraudulent(ctx workflow.Context, data order.Order) (fraud.FraudDecision, error) {
	log.Printf("Checking order %s for fraud.\n\n", data.OrderNumber)
	fraudClient := fraud.InitializeClient()
	result, err := fraudClient.ValidateOrder(data)
	log.Printf("Fraud check for order %s complete.\n\n", data.OrderNumber)
	return result, err
}
