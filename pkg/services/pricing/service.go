package pricing

import (
	"math/rand/v2"
	"temporal-order-demo/pkg/order"
)

type OrderLinePricing struct {
	LinePricing map[string]float64
}
type PriceServiceClient struct {
	//TODO setup credentials to talk to the pricing service
}

func InitializeClient() *PriceServiceClient {
	return &PriceServiceClient{}
}

func (c *PriceServiceClient) PriceOrder(order *order.Order) (OrderLinePricing, error) {
	result := OrderLinePricing{
		LinePricing: make(map[string]float64, 0),
	}
	for _, line := range order.Lines {
		if line.Price == 0 {
			result.LinePricing[line.Product] = 1 + rand.Float64()*(500-1)
		}
	}
	return result, nil
}
