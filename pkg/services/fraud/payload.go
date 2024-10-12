package fraud

import "time"

type FraudDecision struct {
	FraudDetected   bool
	RejectionReason string
	CheckDate       time.Time
}
