package orderworkflowutils

import (
	"temporal-order-demo/pkg/order"
	eventactivity "temporal-order-demo/pkg/order-activities/event"
	orderworkflowqueues "temporal-order-demo/pkg/order-workflow/queues"

	"go.temporal.io/sdk/workflow"
)

func EmitOrderStatusEvent(ctx workflow.Context, order *order.Order, status order.OrderStatus, reason string) error {
	order.UpdateStatus(status, reason)
	emitEventErr := workflow.ExecuteActivity(
		workflow.WithTaskQueue(ctx, orderworkflowqueues.EventEmitterTaskQueueName),
		eventactivity.EmitEvent,
		order).Get(ctx, nil)

	if emitEventErr != nil {
		return emitEventErr
	}

	return nil
}
