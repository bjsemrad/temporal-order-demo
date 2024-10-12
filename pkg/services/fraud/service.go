package fraud

import "temporal-order-demo/pkg/order"

type FraudServiceClient struct {
	//TODO setup credentials to talk to the fraud servicea
}

func InitializeClient() *FraudServiceClient {
	return &FraudServiceClient{}
}

func (c *FraudServiceClient) ValidateOrder(order order.Order) (FraudDecision, error) {
	result := FraudDecision{
		FraudDetected:   false,
		RejectionReason: "",
	}
	if len(order.Lines) > 5 {
		result.FraudDetected = true
		result.RejectionReason = "Large Order"
	}
	return result, nil
}
