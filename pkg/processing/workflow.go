package processing

import (
	"fmt"
	"temporal-order-demo/pkg/order"
	pfraud "temporal-order-demo/pkg/processing/fraud"
	sfraud "temporal-order-demo/pkg/services/fraud"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ProcessOrder(ctx workflow.Context, input order.Order) (string, error) {

	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    120 * time.Second,
		MaximumAttempts:    500, // 0 is unlimited retries
		// NonRetryableErrorTypes: []string{"TODO"},
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: 3 * time.Minute,
		RetryPolicy:         retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var fraudOutput sfraud.FraudDecision

	fraudErr := workflow.ExecuteActivity(ctx, pfraud.CheckOrderFraudulent, input).Get(ctx, &fraudOutput)

	if fraudErr != nil {
		return "", fraudErr
	}

	// // Deposit money.
	// var depositOutput string
	//
	// depositErr := workflow.ExecuteActivity(ctx, Deposit, input).Get(ctx, &depositOutput)
	//
	// if depositErr != nil {
	// 	// The deposit failed; put money back in original account.
	//
	// 	var result string
	//
	// 	refundErr := workflow.ExecuteActivity(ctx, Refund, input).Get(ctx, &result)
	//
	// 	if refundErr != nil {
	// 		return "",
	// 			fmt.Errorf("Deposit: failed to deposit money into %v: %v. Money could not be returned to %v: %w",
	// 				input.TargetAccount, depositErr, input.SourceAccount, refundErr)
	// 	}
	//
	// 	return "", fmt.Errorf("Deposit: failed to deposit money into %v: Money returned to %v: %w",
	// 		input.TargetAccount, input.SourceAccount, depositErr)
	// }
	//
	result := fmt.Sprintf("Order Procesing complete for order: %s)", input.OrderNumber)
	return result, nil
}

// @@@SNIPEND
