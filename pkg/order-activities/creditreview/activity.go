package creditreviewactivity

import (
	"context"
	"log"
	"temporal-order-demo/pkg/order"
	"temporal-order-demo/pkg/services/creditreview"
)

func ValidateAndReserveCredit(ctx context.Context, order *order.Order) (creditreview.CreditReservationResult, error) {
	log.Printf("Checking available credit on order %s.\n\n", order.OrderNumber)
	creditClient := creditreview.InitializeClient()
	creditResult, err := creditClient.ReserveCredit(*order)
	if err != nil {
		return creditResult, err
	}
	return creditResult, nil
}

func SubmitCreditReview(ctx context.Context, order *order.Order, wfID string, runID string) error {
	log.Printf("Checking available credit on order %s.\n\n", order.OrderNumber)
	creditClient := creditreview.InitializeClient()
	err := creditClient.InitiateCreditReview(wfID, runID, *order)
	return err
}
