package orderworkflowstep

import (
	"temporal-order-demo/pkg/order"
	creditreviewactivity "temporal-order-demo/pkg/order-activities/creditreview"
	orderworkflowqueues "temporal-order-demo/pkg/order-workflow/queues"
	orderworkflowutils "temporal-order-demo/pkg/order-workflow/utils"
	"temporal-order-demo/pkg/services/creditreview"

	"go.temporal.io/sdk/workflow"
)

// TODO: Turn this into a sub-workflow
func DoCreditReview(ctx workflow.Context, input *order.Order) error {
	var creditReservation creditreview.CreditReservationResult
	// Execute Fraud Check
	creditReviewErr := workflow.ExecuteActivity(
		workflow.WithTaskQueue(ctx, orderworkflowqueues.CreditReviewTaskQueueName),
		creditreviewactivity.ValidateAndReserveCredit,
		input).Get(ctx, &creditReservation)

	if creditReviewErr != nil {
		return creditReviewErr
	}
	//TODO: Append info to order

	//Credit Not Reserved so we need to do a review
	if !creditReservation.CreditReserved {
		//TODO: Append info to order
		eventError := orderworkflowutils.EmitOrderStatusEvent(ctx, input, order.PendingCreditReview)
		if eventError != nil {
			return eventError
		}

		creditReviewErr := workflow.ExecuteActivity(
			workflow.WithTaskQueue(ctx, orderworkflowqueues.CreditReviewTaskQueueName),
			creditreviewactivity.SubmitCreditReview,
			input).Get(ctx, nil)

		if creditReviewErr != nil {
			return creditReviewErr
		}

		//TODO: Await Signal
		if true {
			eventError = orderworkflowutils.EmitOrderStatusEvent(ctx, input, order.CreditReviewApproved)
			if eventError != nil {
				return eventError
			}
		} else {
			//TODO: Append Info
			eventError = orderworkflowutils.EmitOrderStatusEvent(ctx, input, order.CreditReviewDenied)
			if eventError != nil {
				return eventError
			}

			//TODO: Append Info
			eventError = orderworkflowutils.EmitOrderStatusEvent(ctx, input, order.Canceled)
			if eventError != nil {
				return eventError
			}
		}
	}

	return nil
}
