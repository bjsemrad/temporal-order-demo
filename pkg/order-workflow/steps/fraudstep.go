package orderworkflowstep

import (
	"temporal-order-demo/pkg/order"
	fraudactivity "temporal-order-demo/pkg/order-activities/fraud"
	orderworkflowqueues "temporal-order-demo/pkg/order-workflow/queues"
	orderworkflowutils "temporal-order-demo/pkg/order-workflow/utils"
	"temporal-order-demo/pkg/services/fraud"

	"go.temporal.io/sdk/workflow"
)

// TODO: Turn this into a sub-workflow
func DoFraudCheck(ctx workflow.Context, input *order.Order) error {
	eventError := orderworkflowutils.EmitOrderStatusEvent(ctx, input, order.PendingFraudReview)

	if eventError != nil {
		return eventError
	}

	var fraudOutput fraud.FraudDecision
	// Execute Fraud Check
	fraudErr := workflow.ExecuteActivity(
		workflow.WithTaskQueue(ctx, orderworkflowqueues.FraudTaskQueueName),
		fraudactivity.CheckOrderFraudulent,
		input).Get(ctx, &fraudOutput)

	if fraudErr != nil {
		return fraudErr
	}

	if fraudOutput.FraudDetected {
		//Emit fraud detected event
		eventError := orderworkflowutils.EmitOrderStatusEvent(ctx, input, order.Fraudlent)
		if eventError != nil {
			return eventError
		}
		//TODO: What do we want to do at this point, cancel or have some intervention
	} else {
		//Emit Fraud Review Complete
		eventError := orderworkflowutils.EmitOrderStatusEvent(ctx, input, order.NoFraudDetected)

		if eventError != nil {
			return eventError
		}
	}
	return nil
}
