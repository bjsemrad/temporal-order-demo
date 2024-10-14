package creditreviewactivity

import (
	"context"
	"log"
	"temporal-order-demo/pkg/order"
	"temporal-order-demo/pkg/services/creditreview"
)

type CreditActivity struct {
	CreditClient *creditreview.CreditReviewClient
}

func (c *CreditActivity) ValidateAndReserveCredit(ctx context.Context, order *order.Order) (creditreview.CreditReservationResult, error) {
	log.Printf("Checking available credit on order %s.\n\n", order.OrderNumber)
	creditResult, err := c.CreditClient.ReserveCredit(*order)
	if err != nil {
		return creditResult, err
	}
	return creditResult, nil
}

func (c *CreditActivity) SubmitCreditReview(ctx context.Context, order *order.Order, wfID string, runID string) error {
	log.Printf("Checking available credit on order %s.\n\n", order.OrderNumber)
	err := c.CreditClient.InitiateCreditReview(wfID, runID, order)
	return err
}
