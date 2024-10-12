package order

import (
	"time"
)

type Order struct {
	Status      OrderStatus
	OrderNumber string
	Lines       []*OrderLine
	Payment     Payment
	LastUpdated time.Time
}

type OrderLine struct {
	LineNumber int
	Product    string
	Quantity   int
	Price      float64
}

type Payment struct {
	CreditCard    string
	AccountNumber string
}

func NewOrder(orderNumber string) *Order {
	return &Order{
		OrderNumber: orderNumber,
		Lines:       []*OrderLine{},
	}
}

func (o *Order) AddLine(lineNumber int, product string, quantity int, price float64) {
	line := &OrderLine{
		LineNumber: lineNumber,
		Product:    product,
		Quantity:   quantity,
		Price:      price,
	}
	o.Lines = append(o.Lines, line)
}

func (o *Order) UpdateStatus(newStatus OrderStatus) {
	o.Status = newStatus
	o.LastUpdated = time.Now()
}

func (o *Order) Total() float64 {
	total := 0.0
	for _, line := range o.Lines {
		total += line.Price
	}
	return total
}
