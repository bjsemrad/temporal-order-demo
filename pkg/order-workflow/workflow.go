package orderworkflow

import (
	"errors"
	"log"
	"strings"
	"temporal-order-demo/pkg/order"
	orderworkflowstep "temporal-order-demo/pkg/order-workflow/steps"
	orderstatus "temporal-order-demo/pkg/order/status"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ProcessOrder(ctx workflow.Context, custOrder *order.Order) (order.Order, error) { //TODO comeback to the workflow outpu
	ctx = setupDefaultContext(ctx)

	log.Printf("Start Ensuring Order is Priced: %s.\n", custOrder.OrderNumber)
	orderworkflowstep.EnsureOrderIsPriced(ctx, custOrder)

	log.Printf("Start Fraud Check: %s.\n", custOrder.OrderNumber)
	fraudErr := orderworkflowstep.StartFraudCheck(ctx, custOrder)
	var fraudDetectedError *orderworkflowstep.FraudDetectedError
	if fraudErr != nil && !errors.As(fraudErr, &fraudDetectedError) {
		return *custOrder, fraudErr
	}
	log.Printf("End Fraud Check Step: %s.\n", custOrder.OrderNumber)

	log.Printf("Start Validate Business Rules Check Step: %s.\n", custOrder.OrderNumber)
	//TODO: Validate Business Rules
	log.Printf("End Validate Business Rules Step: %s.\n", custOrder.OrderNumber)

	log.Printf("Start Cust Approval Step: %s.\n", custOrder.OrderNumber)
	//TODO: Cust Approval
	log.Printf("End Cust Approval Step: %s.\n", custOrder.OrderNumber)

	if custOrder.Payment != nil && strings.Trim(custOrder.Payment.AccountNumber, " ") != "" {
		log.Printf("Start Credit Review Step: %s.\n", custOrder.OrderNumber)
		creditReviewErr := orderworkflowstep.StartCreditReview(ctx, custOrder)
		var creditDeniedError *orderworkflowstep.CreditDeniedError
		if creditReviewErr != nil && !errors.As(creditReviewErr, &creditDeniedError) {
			return *custOrder, creditReviewErr
		}
		log.Printf("End Credit Review Step: %s.\n", custOrder.OrderNumber)

	}

	log.Printf("Start Apply Setting Step: %s.\n", custOrder.OrderNumber)
	settingsErr := orderworkflowstep.ApplySettings(ctx, custOrder)
	if settingsErr != nil {
		return *custOrder, settingsErr
	}
	log.Printf("End Apply Setting Step: %s.\n", custOrder.OrderNumber)

	log.Printf("Start Ops Rules Step: %s.\n", custOrder.OrderNumber)
	//TODO: Operational Rule Checks
	log.Printf("End Ops Rules Step: %s.\n", custOrder.OrderNumber)

	log.Printf("Start Intervention Step: %s.\n", custOrder.OrderNumber)
	//TODO: Team Intervention
	log.Printf("End Intervention Step: %s.\n", custOrder.OrderNumber)

	//IF we are not in a terminal status send the fulfillment signal
	if !orderstatus.TerminalOrderStatus(custOrder.Status.Code) {
		log.Printf("Start Prepare For Fulfillment Step: %s.\n", custOrder.OrderNumber)
		fulfillError := orderworkflowstep.PrepareOrderForFulfillment(ctx, custOrder)
		if fulfillError != nil {
			return *custOrder, fulfillError
		}
		log.Printf("End Prepare for Fulfillment Step: %s.\n", custOrder.OrderNumber)

		log.Printf("Start Wait for Confirmed Order: %s.\n", custOrder.OrderNumber)
		orderworkflowstep.WaitForConfirmedOrder(ctx, custOrder)
	}
	log.Printf("Workflow Complete: %s.\n", custOrder.OrderNumber)
	return *custOrder, nil
}

func setupDefaultContext(ctx workflow.Context) workflow.Context {
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    120 * time.Second,
		MaximumAttempts:    0, // 0 is unlimited retries, this will retry if there is no workers as well.
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions, this will issue a retry after 5 minutes
		StartToCloseTimeout: 5 * time.Minute,
		RetryPolicy:         retrypolicy,
	}

	return workflow.WithActivityOptions(ctx, options)
}
