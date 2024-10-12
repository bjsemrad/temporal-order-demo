package orderworkflow

import (
	"temporal-order-demo/pkg/order"
	orderworkflowstep "temporal-order-demo/pkg/order-workflow/steps"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ProcessOrder(ctx workflow.Context, input *order.Order) (order.Order, error) { //TODO comeback to the workflow outpu

	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    120 * time.Second,
		MaximumAttempts:    500, // 0 is unlimited retries, this will retry if there is no workers as well.
		// NonRetryableErrorTypes: []string{"TODO"},
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: 3 * time.Minute,
		RetryPolicy:         retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	fraudErr := orderworkflowstep.DoFraudCheck(ctx, input)
	//TODO: Credit Review
	//TODO: Approval
	//TODO: Operational Rule Checks
	//TODO: Team Intervention
	if fraudErr != nil {
		return *input, fraudErr
	}

	return *input, nil
}