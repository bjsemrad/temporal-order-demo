package creditcheck

import "temporal-order-demo/pkg/order"

type CreditCheckClient struct {
	//TODO setup credentials
}

type CreditApprovalDecision struct {
	Approved bool
	Limit    float64
}

func InitializeClient() *CreditCheckClient {
	return &CreditCheckClient{}
}

func (c *CreditCheckClient) ProcessOrder(order order.Order) (CreditApprovalDecision, error) {
	result := CreditApprovalDecision{
		Approved: true,
	}

	if order.Total() > 100 {
		result.Approved = false
		result.Limit = 90
	}
	return result, nil
}
