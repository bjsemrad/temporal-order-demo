package order

const OrderIntakeTaskQueueName = "ORDER_SUBMITTED_INTAKET_TASK_QUEUE"

type Order struct {
	OrderNumber string
	Lines       []*OrderLine
}

type OrderLine struct {
	LineNumber int
	Product    string
	Quantity   int
	Price      float64
}

func NewOrder(orderNumber string) *Order {
	return &Order{
		OrderNumber: orderNumber,
		Lines:       []*OrderLine{},
	}
}

func (o *Order) addLine(lineNumber int, product string, quantity int, price float64) {
	line := &OrderLine{
		LineNumber: lineNumber,
		Product:    product,
		Quantity:   quantity,
		Price:      price,
	}
	o.Lines = append(o.Lines, line)
}
