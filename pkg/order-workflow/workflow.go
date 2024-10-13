package orderworkflow

import (
	"errors"
	"strings"
	"temporal-order-demo/pkg/order"
	orderworkflowstep "temporal-order-demo/pkg/order-workflow/steps"
	orderworkflowutils "temporal-order-demo/pkg/order-workflow/utils"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ProcessOrder(ctx workflow.Context, input *order.Order) (order.Order, error) { //TODO comeback to the workflow outpu

	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    120 * time.Second,
		MaximumAttempts:    0, // 0 is unlimited retries, this will retry if there is no workers as well.
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: 5 * time.Minute,
		RetryPolicy:         retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	fraudErr := orderworkflowstep.DoFraudCheck(ctx, input)
	var fraudDetectedError *orderworkflowstep.FraudDetectedError
	if fraudErr != nil && !errors.As(fraudErr, &fraudDetectedError) {
		return *input, fraudErr
	}

	if input.Payment != nil && strings.Trim(input.Payment.AccountNumber, " ") != "" {
		creditReviewErr := orderworkflowstep.DoCreditReview(ctx, input)
		var creditDeniedError *orderworkflowstep.CreditDeniedError
		if creditReviewErr != nil && !errors.As(creditReviewErr, &creditDeniedError) {
			return *input, creditReviewErr
		}
	}

	//IF we are not in a terminal status send the fulfillment signal
	if !order.TerminalOrderStatus(input.Status) {
		orderworkflowutils.EmitOrderStatusEvent(ctx, input, order.ReadyForFullfilment, "Processing Complete, Ready to Fulfill")
	}
	//TODO: Approval
	//TODO: Operational Rule Checks
	//TODO: Team Intervention

	return *input, nil
}
