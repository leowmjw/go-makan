package customers

import (
	"github.com/google/go-cmp/cmp"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
	"reflect"
	"testing"
)

func TestOrderWorkflow(t *testing.T) {
	type args struct {
		items      []Item
		delDetails DeliveryDetails
	}
	tests := []struct {
		name    string
		args    args
		want    Order
		wantErr bool
	}{
		// TODO: Add test cases.
		{"happy #1", args{
			items: []Item{
				{
					ShortName: "",
					ShortCode: "",
					Quantity:  0,
					Amount:    0,
					Notes:     "",
				},
				{
					ShortName: "",
					ShortCode: "",
					Quantity:  0,
					Amount:    0,
					Notes:     "",
				},
			},
			delDetails: DeliveryDetails{
				CustomerID:       "",
				FullAddress:      nil,
				ContactNumber:    "",
				DeliveryCost:     0,
				DeliveryRange:    0,
				DeliveryEstimate: 0,
			},
		}, Order{
			ID:        "",
			PartnerID: "",
			Items: []Item{
				{
					ShortName: "",
					ShortCode: "",
					Quantity:  0,
					Amount:    0,
					Notes:     "",
				},
			},
			DeliveryDetails: DeliveryDetails{
				CustomerID:       "",
				FullAddress:      nil,
				ContactNumber:    "",
				DeliveryCost:     4,
				DeliveryRange:    10,
				DeliveryEstimate: 0,
			},
			PaymentDetails: PaymentDetails{
				ID:        "",
				ShortName: "COD",
				Itemized: []string{
					"2 x Fish - RM 20.00",
					"1 x Chicken - RM 30.00",
				},
				Total: 50,
			},
			OffersClaim: OffersClaim{
				ID:        "",
				ShortName: "",
			},
		}, true},
	}
	// Setup temporal testsuite
	ts := testsuite.WorkflowTestSuite{}
	//ts.SetContextPropagators()
	// Run in parallel; see if got race conditions ..
	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// new test environment ..
			env := ts.NewTestWorkflowEnvironment()
			env.SetStartWorkflowOptions(client.StartWorkflowOptions{
				ID: "mleow/baba-ang", // this does not work ..
			})

			env.RegisterWorkflow(OrderWorkflow)
			env.ExecuteWorkflow(OrderWorkflow, tt.args.items, tt.args.delDetails)

			// Add item1, add item2, add item3, mod item3, add item 4, remove item4
			// Confirm + pay
			if !env.IsWorkflowCompleted() {
				t.Fatal("WF NOT COmpleted!! Timed out??")
			}
			if wferr := env.GetWorkflowError(); wferr != nil {
				t.Fatal(wferr)
			}
			var got Order
			rerr := env.GetWorkflowResult(&got)
			if rerr != nil {
				t.Fatal(rerr)
			}

			if diff := cmp.Diff(tt.want.Items, got.Items); diff != "" {
				t.Errorf("OrderWorkflow() mismatch (-want +got):\n%s", diff)
			}
			//got, err := OrderWorkflow(tt.args.ctx, tt.args.items, tt.args.delDetails)
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("OrderWorkflow() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("OrderWorkflow() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestWaitOrderWorkflow(t *testing.T) {
	type args struct {
		ctx   workflow.Context
		order Order
	}
	tests := []struct {
		name    string
		args    args
		want    OrderDelivery
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WaitOrderWorkflow(tt.args.ctx, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("WaitOrderWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WaitOrderWorkflow() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostOrderWorkflow(t *testing.T) {
	type args struct {
		ctx   workflow.Context
		order OrderDelivery
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PostOrderWorkflow(tt.args.ctx, tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("PostOrderWorkflow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
