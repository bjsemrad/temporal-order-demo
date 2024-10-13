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

func ProcessOrder(ctx workflow.Context, order *order.Order) (order.Order, error) { //TODO comeback to the workflow outpu
	ctx = setupDefaultContext(ctx)

	orderworkflowstep.EnsureOrderIsPriced(ctx, order)

	fraudErr := orderworkflowstep.StartFraudCheck(ctx, order)
	var fraudDetectedError *orderworkflowstep.FraudDetectedError
	if fraudErr != nil && !errors.As(fraudErr, &fraudDetectedError) {
		return *order, fraudErr
	}

	//TODO: Validate Business Rules

	//TODO: Cust Approval

	if order.Payment != nil && strings.Trim(order.Payment.AccountNumber, " ") != "" {
		creditReviewErr := orderworkflowstep.StartCreditReview(ctx, order)
		var creditDeniedError *orderworkflowstep.CreditDeniedError
		if creditReviewErr != nil && !errors.As(creditReviewErr, &creditDeniedError) {
			return *order, creditReviewErr
		}
	}

	//TODO: Apply Settings

	//TODO: Operational Rule Checks

	//TODO: Team Intervention

	fulfillError := orderworkflowstep.PrepareOrderForFulfillment(ctx, order)
	if fulfillError != nil {
		return *order, fulfillError
	}

	orderworkflowstep.WaitForConfirmedOrder(ctx, order)

	return *order, nil
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
