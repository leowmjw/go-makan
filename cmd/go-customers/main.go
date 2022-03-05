package main

import (
	"app/workflow/customers"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create the client
	c, err := client.NewClient(client.Options{})
	if err != nil {
		// log failure
		panic(err)
	}
	defer c.Close()

	q := "food.oder"

	// DEBUG only if need to simulate
	//for i := 1; i < 2; i++ {
	//	go func(id int) {
	//		// Start workflow
	//		wid := fmt.Sprintf("mleow-%d", id)
	//		wfo := client.StartWorkflowOptions{
	//			ID:                       wid,
	//			TaskQueue:                onboarding_patient.WorkflowQueue,
	//			WorkflowExecutionTimeout: 2 * time.Minute,
	//			WorkflowRunTimeout:       time.Minute,
	//			WorkflowTaskTimeout:      time.Second,
	//			RetryPolicy: &temporal.RetryPolicy{
	//				InitialInterval: time.Second,
	//				MaximumAttempts: 2,
	//			},
	//		}
	//		wfr, err := c.ExecuteWorkflow(context.Background(), wfo,
	//			bug_wf_retry.RetryWorkflow,
	//			wid,
	//		)
	//		if err != nil {
	//			panic(err)
	//		}
	//		fmt.Println("WFIF: ", wfr.GetID(), "RID: ", wfr.GetRunID())
	//	}(i)
	//}
	// Create the Worker
	w := worker.New(c, q, worker.Options{})
	// Register your functions
	w.RegisterWorkflow(customers.OrderWorkflow)
	//w.RegisterWorkflow(bug_wf_retry.RetryWorkflow)
	//w.RegisterActivity(bug_wf_retry.FailActivities)
	// Run the Worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		// log failure
		panic(err)
	}

}
