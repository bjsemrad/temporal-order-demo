package orderworkflow

import (
	"errors"
	"strings"
	"temporal-order-demo/pkg/order"
	orderworkflowstep "temporal-order-demo/pkg/order-workflow/steps"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ProcessOrder(ctx workflow.Context, input *order.Order) (order.Order, error) { //TODO comeback to the workflow outpu
	ctx = setupDefaultContext(ctx)

	fraudErr := orderworkflowstep.StartFraudCheck(ctx, input)
	var fraudDetectedError *orderworkflowstep.FraudDetectedError
	if fraudErr != nil && !errors.As(fraudErr, &fraudDetectedError) {
		return *input, fraudErr
	}

	//TODO: Approval

	if input.Payment != nil && strings.Trim(input.Payment.AccountNumber, " ") != "" {
		creditReviewErr := orderworkflowstep.StartCreditReview(ctx, input)
		var creditDeniedError *orderworkflowstep.CreditDeniedError
		if creditReviewErr != nil && !errors.As(creditReviewErr, &creditDeniedError) {
			return *input, creditReviewErr
		}
	}
	//TODO: Operational Rule Checks

	//TODO: Team Intervention

	fulfillError := orderworkflowstep.PrepareOrderForFulfillment(ctx, input)
	if fulfillError != nil {
		return *input, fulfillError
	}

	orderworkflowstep.WaitForConfirmedOrder(ctx, input)

	return *input, nil
}

func setupDefaultContext(ctx workflow.Context) workflow.Context {
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

	return workflow.WithActivityOptions(ctx, options)
}
