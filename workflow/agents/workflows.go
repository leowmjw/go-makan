package agents

import (
	"go.temporal.io/sdk/workflow"
)

// AgentWorkflow - Lifecycle: From Start to End
func AgentWorkflow(ctx workflow.Context) {
	workflow.GetInfo(ctx)

	// Gets paid according to scheduled payout
}

// AgentActiveWorkflow - Actively Delivering
func AgentActiveWorkflow(ctx workflow.Context) {

}

func AgentPayout(ctx workflow.Context) {
	// Examine record since last payout
	// Examine cost incurred in same period

	// Transfer payment to selected institution
}
