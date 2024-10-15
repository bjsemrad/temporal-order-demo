package orderworkflowstep

import (
	"log"
	"temporal-order-demo/pkg/order"
	orderworkflowutils "temporal-order-demo/pkg/order-workflow/utils"
	"temporal-order-demo/pkg/order/status"

	"go.temporal.io/sdk/workflow"
)

const (
	OrderFulfillmentConfirmedChannel = "OrderFulfillmentConfirmed"
)

type OrderConfirmedSignal struct {
	FulfillmentOrderNumber string
}

func WaitForConfirmedOrder(ctx workflow.Context, custOrder *order.Order) error {
	if !orderstatus.TerminalOrderStatus(custOrder.Status.Code) {
		var confirmSignalInput OrderConfirmedSignal
		workflow.GetSignalChannel(ctx, OrderFulfillmentConfirmedChannel).Receive(ctx, &confirmSignalInput)
		log.Printf("Received confirmed order signal")

		custOrder.FullfilmentOrderNumber = confirmSignalInput.FulfillmentOrderNumber
		eventError := orderworkflowutils.EmitOrderStatusEvent(ctx, custOrder, orderstatus.FullfilmentConfirmed, "Order Confirmed for Fullfilment")
		if eventError != nil {
			return eventError
		}
	}

	return nil
}
