package creditcheck

import (
	"context"
	"log"
	"temporal-order-demo/pkg/order"
	"temporal-order-demo/pkg/services/creditcheck"
)

func SubmitCreditReview(ctx context.Context, data order.Order) (creditcheck.CreditApprovalDecision, error) {
	log.Printf("Checking available credit on order %s.\n\n", data.OrderNumber)
	creditClient := creditcheck.InitializeClient()
	result, err := creditClient.ProcessOrder(data)
	log.Printf("Credit Review for order %s complete.\n\n", data.OrderNumber)
	return result, err
}
