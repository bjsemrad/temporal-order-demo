package creditreview

import (
	"temporal-order-demo/pkg/order"
)

type CreditReviewClient struct {
	//TODO setup credentials
}

type CreditApprovalDecision struct {
	Approved bool
	Limit    float64
}

func InitializeClient() *CreditReviewClient {
	return &CreditReviewClient{}
}

type CreditReservationResult struct {
	CreditReserved  bool
	AvailableCredit float64
}

func (c *CreditReviewClient) ReserveCredit(order order.Order) (CreditReservationResult, error) {
	if order.Total() > 4000 {
		return CreditReservationResult{
			CreditReserved:  false,
			AvailableCredit: 100.00,
		}, nil
	}
	return CreditReservationResult{
		CreditReserved:  true,
		AvailableCredit: 5000.00,
	}, nil
}

func (c *CreditReviewClient) InitiateCreditReview(order order.Order) error {
	return nil
}
