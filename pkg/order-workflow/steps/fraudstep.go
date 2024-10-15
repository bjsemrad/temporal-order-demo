package orderworkflowstep

import (
	"temporal-order-demo/pkg/order"
	fraudactivity "temporal-order-demo/pkg/order-activities/fraud"
	orderworkflowqueues "temporal-order-demo/pkg/order-workflow/queues"
	orderworkflowutils "temporal-order-demo/pkg/order-workflow/utils"
	orderstatus "temporal-order-demo/pkg/order/status"
	"temporal-order-demo/pkg/services/fraud"

	"go.temporal.io/sdk/workflow"
)

type FraudDetectedError struct{}

func (m *FraudDetectedError) Error() string {
	return "Fraud Detected"
}

// TODO: Turn this into a sub-workflow
func StartFraudCheck(ctx workflow.Context, input *order.Order) error {
	eventError := orderworkflowutils.EmitOrderStatusEvent(ctx, input, orderstatus.PendingFraudReview, "Begin Fraud Review")

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
	input.RecoardFraudReviewDecision(fraudOutput.FraudDetected, fraudOutput.RejectionReason, fraudOutput.CheckDate)
	if fraudOutput.FraudDetected {
		//Emit fraud detected event
		eventError := orderworkflowutils.EmitOrderStatusEvent(ctx, input, orderstatus.Fraudlent, "Fraud Detected")
		if eventError != nil {
			return eventError
		}

		return &FraudDetectedError{}
		//TODO: What do we want to do at this point, cancel or have some intervention
	} else {
		//Emit Fraud Review Complete
		eventError := orderworkflowutils.EmitOrderStatusEvent(ctx, input, orderstatus.NoFraudDetected, "Order deemed not fraudlent")

		if eventError != nil {
			return eventError
		}
	}
	return nil
}
