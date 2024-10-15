package orderworkflowutils

import (
	"temporal-order-demo/pkg/order"
	eventactivity "temporal-order-demo/pkg/order-activities/event"
	orderworkflowqueues "temporal-order-demo/pkg/order-workflow/queues"
	orderstatus "temporal-order-demo/pkg/order/status"
	"time"

	"go.temporal.io/sdk/workflow"
)

func EmitOrderStatusEvent(ctx workflow.Context, order *order.Order, status orderstatus.OrderStatusCode, reason string) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	order.UpdateStatus(status, reason)
	var eb *eventactivity.EventBroker
	emitEventErr := workflow.ExecuteActivity(
		workflow.WithTaskQueue(ctx, orderworkflowqueues.EventEmitterTaskQueueName),
		eb.EmitStatusUpdateEvent,
		order).Get(ctx, nil)

	if emitEventErr != nil {
		return emitEventErr
	}

	return nil
}
