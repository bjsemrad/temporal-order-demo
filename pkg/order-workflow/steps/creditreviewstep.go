package orderworkflowstep

import (
	"log"
	"temporal-order-demo/pkg/order"
	creditreviewactivity "temporal-order-demo/pkg/order-activities/creditreview"
	orderworkflowqueues "temporal-order-demo/pkg/order-workflow/queues"
	orderworkflowutils "temporal-order-demo/pkg/order-workflow/utils"
	"temporal-order-demo/pkg/services/creditreview"
	"time"

	"go.temporal.io/sdk/workflow"
)

type CreditReviewDecision string

const (
	CreditExtended        CreditReviewDecision = "credit-extended"
	CreditExtensionDenied CreditReviewDecision = "credit-extension-denied"

	CreditReviewDecisionChannel = "OrderCreditReviewDecision"
)

type CreditReviewDecisionSignal struct {
	OrderNumber    string
	CreditDecision CreditReviewDecision
	Reviewier      string
	NewLimit       float64
	DecisionDate   time.Time
}

// TODO: Turn this into a sub-workflow
func DoCreditReview(ctx workflow.Context, custOrder *order.Order) error {
	var creditReservation creditreview.CreditReservationResult
	// Execute Fraud Check
	creditReviewErr := workflow.ExecuteActivity(
		workflow.WithTaskQueue(ctx, orderworkflowqueues.CreditReviewTaskQueueName),
		creditreviewactivity.ValidateAndReserveCredit,
		custOrder).Get(ctx, &creditReservation)

	if creditReviewErr != nil {
		return creditReviewErr
	}

	custOrder.RecordCreditReservation(creditReservation.CreditReserved, creditReservation.AvailableCredit)
	//Credit Not Reserved so we need to do a review
	if !creditReservation.CreditReserved {
		eventError := orderworkflowutils.EmitOrderStatusEvent(ctx, custOrder, order.PendingCreditReview, "Not enough credit available for Order")
		if eventError != nil {
			return eventError
		}

		creditReviewErr := workflow.ExecuteActivity(
			workflow.WithTaskQueue(ctx, orderworkflowqueues.CreditReviewTaskQueueName),
			creditreviewactivity.SubmitCreditReview,
			custOrder).Get(ctx, nil)

		if creditReviewErr != nil {
			return creditReviewErr
		}

		var cdSignalInput CreditReviewDecisionSignal
		workflow.GetSignalChannel(ctx, CreditReviewDecisionChannel).Receive(ctx, &cdSignalInput)
		log.Printf("Received signal for credit review decision")

		custOrder.RecordCreditReviewDecision(string(cdSignalInput.CreditDecision), cdSignalInput.Reviewier, cdSignalInput.NewLimit, cdSignalInput.DecisionDate)

		if cdSignalInput.CreditDecision == CreditExtended {
			eventError = orderworkflowutils.EmitOrderStatusEvent(ctx, custOrder, order.CreditReviewApproved, "Credit Increase Approved")
			if eventError != nil {
				return eventError
			}
		} else {
			eventError = orderworkflowutils.EmitOrderStatusEvent(ctx, custOrder, order.CreditReviewDenied, "Credit Increase Denied")
			if eventError != nil {
				return eventError
			}

			eventError = orderworkflowutils.EmitOrderStatusEvent(ctx, custOrder, order.Canceled, "Not enough credit canceling order")
			if eventError != nil {
				return eventError
			}
		}
	}

	return nil
}
