package order

type OrderStatus string

const (
	Submitted OrderStatus = "submitted"

	PendingFraudReview OrderStatus = "pending-fraud"
	NoFraudDetected    OrderStatus = "no-fraud"
	Fraudlent          OrderStatus = "fraud"

	PendingCreditReview  OrderStatus = "pending-credit-review"
	CreditReviewApproved OrderStatus = "credit-review-approved"
	CreditReviewDenied   OrderStatus = "credit-review-denied"

	ApprovalRequired OrderStatus = "approvalrequired"
	Approved         OrderStatus = "approved"
	Rejected         OrderStatus = "rejected"
	Canceled         OrderStatus = "canceled"

	ReadyForFullfilment OrderStatus = "ready-for-fulfillment"
)

func TerminalOrderStatus(status OrderStatus) bool {
	switch status {
	case Fraudlent, CreditReviewDenied, Rejected, Canceled:
		return true
	default:
		return false
	}
}
