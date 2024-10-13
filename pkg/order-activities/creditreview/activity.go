package creditreviewactivity

import (
	"log"
	"temporal-order-demo/pkg/order"
	"temporal-order-demo/pkg/services/creditreview"

	"go.temporal.io/sdk/workflow"
)

func ValidateAndReserveCredit(ctx workflow.Context, order *order.Order) (creditreview.CreditReservationResult, error) {
	log.Printf("Checking available credit on order %s.\n\n", order.OrderNumber)
	creditClient := creditreview.InitializeClient()
	creditResult, err := creditClient.ReserveCredit(*order)
	if err != nil {
		return creditResult, err
	}
	return creditResult, nil
}

func SubmitCreditReview(ctx workflow.Context, order *order.Order) error {
	log.Printf("Checking available credit on order %s.\n\n", order.OrderNumber)
	creditClient := creditreview.InitializeClient()
	wfID := workflow.GetInfo(ctx).WorkflowExecution.ID
	runID := workflow.GetInfo(ctx).WorkflowExecution.RunID
	err := creditClient.InitiateCreditReview(wfID, runID, *order)
	return err
}
