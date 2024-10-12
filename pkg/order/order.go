package order

import (
	"time"
)

type Order struct {
	Status           OrderStatus
	OrderNumber      string
	Lines            []*OrderLine
	Payment          *Payment
	LastUpdated      time.Time
	PipelineMetadata *OrderPipelineMetadata
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

type OrderPipelineMetadata struct {
	FraudReview   OrderFraudReview
	CreditReview  OrderCreditReview
	StatusHistory []*OrderStatusHistory
}

type OrderStatusHistory struct {
	Status OrderStatus
	Date   time.Time
}

type OrderFraudReview struct {
	FraudDetected   bool
	RejectionReason string
	DecisionDate    time.Time
}

type OrderCreditReview struct {
	CreditAvailable bool
	AvailableCredit float64
	CreditDecision  string
	Reviewier       string
	DecisionDate    time.Time
}

func NewOrder(orderNumber string) *Order {
	return &Order{
		OrderNumber: orderNumber,
		Lines:       []*OrderLine{},
		Payment:     &Payment{},
		PipelineMetadata: &OrderPipelineMetadata{
			StatusHistory: make([]*OrderStatusHistory, 0),
		},
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

func (o *Order) UpdateStatus(newStatus OrderStatus, reason string) {
	o.recordStatusChange(o.Status, o.LastUpdated)
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

func (o *Order) RecordCreditReservation(creditAvailable bool, availableCredit float64) {
	o.ensureMetadataInitalized()
	o.PipelineMetadata.CreditReview.CreditAvailable = creditAvailable
	o.PipelineMetadata.CreditReview.AvailableCredit = availableCredit
}

func (o *Order) RecordCreditReviewDecision(decision string, reviewer string, reviewDate time.Time) {
	o.ensureMetadataInitalized()
	o.PipelineMetadata.CreditReview.CreditDecision = decision
	o.PipelineMetadata.CreditReview.Reviewier = reviewer
	o.PipelineMetadata.CreditReview.DecisionDate = reviewDate
}

func (o *Order) recordStatusChange(status OrderStatus, statusDate time.Time) {
	o.PipelineMetadata.StatusHistory = append(o.PipelineMetadata.StatusHistory, &OrderStatusHistory{status, statusDate})
}

func (o *Order) ensureMetadataInitalized() {
	if o.PipelineMetadata == nil {
		o.PipelineMetadata = &OrderPipelineMetadata{
			StatusHistory: make([]*OrderStatusHistory, 0),
		}
	}
}
