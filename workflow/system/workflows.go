package planner

import "go.temporal.io/sdk/workflow"

type OrderHistory struct {
	ID   string
	Logs []string
}

func BestFitAgentsWorkflow(ctx workflow.Context) {
	// Get signal when Agents Active or Deactivate

	// Select the BestFitAgent when order signal comes in

	// Next time partition against zones + adjacents ..
}

// DeliveryOrderWorkflow - tracks history from Customer, Partner, Agent, System
func DeliveryOrderWorkflow(ctx workflow.Context) {
	// New Order coming in ...
	// Check payment capabilities (if not COD)
	// Allocate the next BestFit Agent

	// Deliver
	// Cash Collect ..
}

// Event Broker function that sends signal for new events (Kafka-like)
// SysConsumerWorkflow - Consumes all events; and send to needed; restart every 1k events
