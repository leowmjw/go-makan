package support

import "go.temporal.io/sdk/workflow"

type Agent struct {
	Name         string
	IsSupervisor bool
}

// SupportAgentWorkflow - Lifecycle of Support Agent
func SupportAgentWorkflow(ctx workflow.Context) {
	// On-Boarding

	// Off-Boarding ..
}
