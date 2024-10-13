package orderworkflow

import (
	"errors"
	"log"
	"strings"
	"temporal-order-demo/pkg/order"
	orderworkflowstep "temporal-order-demo/pkg/order-workflow/steps"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ProcessOrder(ctx workflow.Context, order *order.Order) (order.Order, error) { //TODO comeback to the workflow outpu
	ctx = setupDefaultContext(ctx)

	log.Printf("Start Ensuring Order is Priced: %s.\n", order.OrderNumber)
	orderworkflowstep.EnsureOrderIsPriced(ctx, order)
	log.Printf("Start Fraud Check: %s.\n", order.OrderNumber)
	fraudErr := orderworkflowstep.StartFraudCheck(ctx, order)
	log.Printf("End Fraud Check Step: %s.\n", order.OrderNumber)
	var fraudDetectedError *orderworkflowstep.FraudDetectedError
	if fraudErr != nil && !errors.As(fraudErr, &fraudDetectedError) {
		return *order, fraudErr
	}
	log.Printf("Start Validate Business Rules Check Step: %s.\n", order.OrderNumber)
	//TODO: Validate Business Rules
	log.Printf("End Validate Business Rules Step: %s.\n", order.OrderNumber)
	log.Printf("Start Cust Approval Step: %s.\n", order.OrderNumber)
	//TODO: Cust Approval
	log.Printf("End Cust Approval Step: %s.\n", order.OrderNumber)
	if order.Payment != nil && strings.Trim(order.Payment.AccountNumber, " ") != "" {
		log.Printf("Start Credit Review Step: %s.\n", order.OrderNumber)
		creditReviewErr := orderworkflowstep.StartCreditReview(ctx, order)
		var creditDeniedError *orderworkflowstep.CreditDeniedError
		if creditReviewErr != nil && !errors.As(creditReviewErr, &creditDeniedError) {
			return *order, creditReviewErr
		}
		log.Printf("End Credit Review Step: %s.\n", order.OrderNumber)

	}
	log.Printf("Start Apply Setting Step: %s.\n", order.OrderNumber)
	//TODO: Apply Settings
	log.Printf("End Apply Setting Step: %s.\n", order.OrderNumber)

	log.Printf("Start Ops Rules Step: %s.\n", order.OrderNumber)
	//TODO: Operational Rule Checks
	log.Printf("End Ops Rules Step: %s.\n", order.OrderNumber)
	log.Printf("Start Intervention Step: %s.\n", order.OrderNumber)
	//TODO: Team Intervention
	log.Printf("End Intervention Step: %s.\n", order.OrderNumber)

	log.Printf("Start Prepare For Fulfillment Step: %s.\n", order.OrderNumber)
	fulfillError := orderworkflowstep.PrepareOrderForFulfillment(ctx, order)
	if fulfillError != nil {
		return *order, fulfillError
	}
	log.Printf("End Prepare for Fulfillment Step: %s.\n", order.OrderNumber)

	log.Printf("Start Wait for Confirmed Order: %s.\n", order.OrderNumber)
	orderworkflowstep.WaitForConfirmedOrder(ctx, order)
	log.Printf("Workflow Complete: %s.\n", order.OrderNumber)
	return *order, nil
}

func setupDefaultContext(ctx workflow.Context) workflow.Context {
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    120 * time.Second,
		MaximumAttempts:    0, // 0 is unlimited retries, this will retry if there is no workers as well.
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: 5 * time.Minute,
		RetryPolicy:         retrypolicy,
	}

	return workflow.WithActivityOptions(ctx, options)
}
