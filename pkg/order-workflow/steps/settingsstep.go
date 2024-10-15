package orderworkflowstep

import (
	"log"
	"temporal-order-demo/pkg/order"
	settingsactivity "temporal-order-demo/pkg/order-activities/settings"
	orderstatus "temporal-order-demo/pkg/order/status"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ApplySettings(ctx workflow.Context, custOrder *order.Order) error {
	//IF we are not in a terminal status send the fulfillment signal
	if !orderstatus.TerminalOrderStatus(custOrder.Status.Code) {

		log.Printf("Appending Order Settings to Order")

		retrypolicy := &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    120 * time.Second,
			MaximumAttempts:    0, // 0 is unlimited retries, this will retry if there is no workers as well.
		}

		options := workflow.LocalActivityOptions{
			// Timeout options specify when to automatically timeout Activity functions.
			StartToCloseTimeout: 5 * time.Minute,
			RetryPolicy:         retrypolicy,
		}

		localContext := workflow.WithLocalActivityOptions(ctx, options)
		settingsErr := workflow.ExecuteLocalActivity(localContext,
			settingsactivity.ApplySettings,
			custOrder).Get(ctx, nil)

		if settingsErr != nil {
			return settingsErr
		}
	}
	return nil
}
