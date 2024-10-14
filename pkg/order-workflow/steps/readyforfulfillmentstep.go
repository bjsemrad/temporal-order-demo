package orderworkflowstep

import (
	"temporal-order-demo/pkg/order"
	eventactivity "temporal-order-demo/pkg/order-activities/event"
	orderworkflowqueues "temporal-order-demo/pkg/order-workflow/queues"
	orderworkflowutils "temporal-order-demo/pkg/order-workflow/utils"
	"time"

	"go.temporal.io/sdk/workflow"
)

func PrepareOrderForFulfillment(ctx workflow.Context, custOrder *order.Order) error {
	//IF we are not in a terminal status send the fulfillment signal
	if !order.TerminalOrderStatus(custOrder.Status) {
		err := orderworkflowutils.EmitOrderStatusEvent(ctx, custOrder, order.ReadyForFullfilment, "Processing Complete, Ready to Fulfill")

		if err != nil {
			return err
		}

		//TODO: Send Out 2nd Event
		wfID := workflow.GetInfo(ctx).WorkflowExecution.ID
		runID := workflow.GetInfo(ctx).WorkflowExecution.RunID

		options := workflow.ActivityOptions{
			StartToCloseTimeout: 5 * time.Minute,
		}
		ctx = workflow.WithActivityOptions(ctx, options)

		var eb *eventactivity.EventBroker
		emitEventErr := workflow.ExecuteActivity(
			workflow.WithTaskQueue(ctx, orderworkflowqueues.EventEmitterTaskQueueName),
			eb.EmitFullfilmentEvent,
			wfID,
			runID,
			custOrder).Get(ctx, nil)

		if emitEventErr != nil {
			return emitEventErr
		}
	}

	return nil
}
