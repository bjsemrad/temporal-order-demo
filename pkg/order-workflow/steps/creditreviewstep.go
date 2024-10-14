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
	CreditDecision CreditReviewDecision
	Reviewier      string
	NewLimit       float64
	DecisionDate   time.Time
}

type CreditDeniedError struct{}

func (m *CreditDeniedError) Error() string {
	return "Credit Extension Declined"
}

// TODO: Turn this into a sub-workflow
func StartCreditReview(ctx workflow.Context, custOrder *order.Order) error {
	creditReservation, creditReviewErr := startValidateAndReserveActivity(ctx, custOrder)
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

		creditReviewErr := startCreditReviewActivity(ctx, custOrder)
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
			return processCreditDenied(ctx, custOrder)
		}
	}

	return nil
}

func processCreditDenied(ctx workflow.Context, custOrder *order.Order) error {
	eventError := orderworkflowutils.EmitOrderStatusEvent(ctx, custOrder, order.CreditReviewDenied, "Credit Increase Denied")
	if eventError != nil {
		return eventError
	}

	eventError = orderworkflowutils.EmitOrderStatusEvent(ctx, custOrder, order.Canceled, "Not enough credit canceling order")
	if eventError != nil {
		return eventError
	}

	return &CreditDeniedError{}
}

func startValidateAndReserveActivity(ctx workflow.Context, custOrder *order.Order) (creditreview.CreditReservationResult, error) {
	var creditReservation creditreview.CreditReservationResult

	options := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	var creditActivity *creditreviewactivity.CreditActivity
	creditReviewErr := workflow.ExecuteActivity(
		workflow.WithTaskQueue(ctx, orderworkflowqueues.CreditReviewTaskQueueName),
		creditActivity.ValidateAndReserveCredit,
		custOrder).Get(ctx, &creditReservation)
	return creditReservation, creditReviewErr
}

func startCreditReviewActivity(ctx workflow.Context, custOrder *order.Order) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	wfID := workflow.GetInfo(ctx).WorkflowExecution.ID
	runID := workflow.GetInfo(ctx).WorkflowExecution.RunID

	var creditActivity *creditreviewactivity.CreditActivity
	creditReviewErr := workflow.ExecuteActivity(
		workflow.WithTaskQueue(ctx, orderworkflowqueues.CreditReviewTaskQueueName),
		creditActivity.SubmitCreditReview,
		custOrder, wfID, runID).Get(ctx, nil)
	return creditReviewErr
}
