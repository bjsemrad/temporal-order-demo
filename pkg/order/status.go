package order

type OrderStatus string

const (
	Submitted          OrderStatus = "submitted"
	PendingFraudReview OrderStatus = "pendingfraud"
	NoFraudDetected    OrderStatus = "no-fraud"
	Fraudlent          OrderStatus = "fraud"
	ApprovalRequired   OrderStatus = "approvalrequired"
	Approved           OrderStatus = "approved"
	Rejected           OrderStatus = "rejected"
	Canceled           OrderStatus = "canceled"
)
