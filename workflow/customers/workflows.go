package customers

import "go.temporal.io/sdk/workflow"

type Item struct {
	ShortName string
	ShortCode string
	Quantity  int
	Amount    int
	Notes     string
}
type Order struct {
	Items []Item
	COD   bool
}

func CustomerWorkflow(ctx workflow.Context) {
	// Sign up

	// Sign off
}

func CustomerOrderWorkflow(ctx workflow.Context) {
	// Choose Order + Details

	// Make payment (if ot COD)

	// Accept Delivery

	// Send Feedback ..
}
