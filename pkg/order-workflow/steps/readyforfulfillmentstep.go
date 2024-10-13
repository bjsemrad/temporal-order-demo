package orderworkflowstep

import (
	"temporal-order-demo/pkg/order"
	orderworkflowutils "temporal-order-demo/pkg/order-workflow/utils"

	"go.temporal.io/sdk/workflow"
)

func PrepareOrderForFulfillment(ctx workflow.Context, custOrder *order.Order) error {
	//IF we are not in a terminal status send the fulfillment signal
	if !order.TerminalOrderStatus(custOrder.Status) {
		return orderworkflowutils.EmitOrderStatusEvent(ctx, custOrder, order.ReadyForFullfilment, "Processing Complete, Ready to Fulfill")
	}

	return nil
}
