package processing

import (
	"temporal-order-demo/pkg/order"
	"temporal-order-demo/pkg/processing/event"
	pfraud "temporal-order-demo/pkg/processing/fraud"
	processorqueue "temporal-order-demo/pkg/queue"
	sfraud "temporal-order-demo/pkg/services/fraud"
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

	fraudErr := handleFraudCheck(ctx, input)
	if fraudErr != nil {
		return *input, fraudErr
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
	// result := fmt.Sprintf("Order Procesing complete for order: %s)", input.OrderNumber)
	return *input, nil
}

// TODO Turn this into a sub-workflow
func handleFraudCheck(ctx workflow.Context, input *order.Order) error {
	var fraudEventOutput event.EventEmitOutput
	input.UpdateStatus(order.PendingFraudReview)
	//Emit pending fraud review event
	emitEventErr := workflow.ExecuteActivity(
		workflow.WithTaskQueue(ctx, processorqueue.EventEmitterTaskQueueName),
		event.EmitEvent,
		input).Get(ctx, &fraudEventOutput)

	if emitEventErr != nil {
		return emitEventErr
	}

	var fraudOutput sfraud.FraudDecision
	// Execute Fraud Check
	fraudErr := workflow.ExecuteActivity(
		workflow.WithTaskQueue(ctx, processorqueue.FraudTaskQueueName),
		pfraud.CheckOrderFraudulent,
		input).Get(ctx, &fraudOutput)

	if fraudErr != nil {
		return fraudErr
	}

	var eventOutput event.EventEmitOutput

	if fraudOutput.FraudDetected {
		//TODO REVISIT THIS AS WE WILL WANT TO DEAL WITH EVENTS DIFFERENTLY
		//Emit fraud detected event
		input.UpdateStatus(order.Fraudlent)
		emitEventErr := workflow.ExecuteActivity(
			workflow.WithTaskQueue(ctx, processorqueue.EventEmitterTaskQueueName),
			event.EmitEvent,
			input).Get(ctx, &eventOutput)

		if emitEventErr != nil {
			return emitEventErr
		}
	} else {
		//Emit Fraud Review Complete
		input.UpdateStatus(order.NoFraudDetected)
		emitEventErr = workflow.ExecuteActivity(
			workflow.WithTaskQueue(ctx, processorqueue.EventEmitterTaskQueueName),
			event.EmitEvent,
			input).Get(ctx, &eventOutput)

		if emitEventErr != nil {
			return emitEventErr
		}
	}
	return nil
}
