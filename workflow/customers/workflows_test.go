package customers

import (
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
	"reflect"
	"testing"
	"time"
)

func TestReceiveFirstItemStartOrderWorkflow(t *testing.T) {
	// When receive the first add item via action signal
	//	start a new OrderWorkflow  combo <customerID>/<partnerID>
	//		with ShoppingCart
	// If already exists, continue on; flag but not fatal?
	//	Optional: with the pre-req Location info
	//	Query Shopping Cart should succeed
}

func TestRemovedLastItemRemoveCartOrderWorkflow(t *testing.T) {
	// Given a bunch of ShoppingCarts OrderWorkflow combo <customerID>/<partnerID>
	// When receive the final remove of Item from ShoppingCart,
	//		then => close OrderWorkflow  combo <customerID>/<partnerID>
}

func TestOrderWorkflow(t *testing.T) {
	type args struct {
		cart ShoppingCart
	}
	tests := []struct {
		name    string
		args    args
		want    Order
		wantErr bool
	}{
		{"happy #1", args{
			cart: ShoppingCart{
				CustomerID: "mleow",
				PartnerID:  "baba-ang",
			},
		}, Order{
			ID:        "",
			PartnerID: "",
			Items: []Item{
				{
					ShortName: "FishBurger",
					ShortCode: "FB-01",
					Quantity:  1,
					Amount:    10,
					Notes:     "is goode",
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
			// Should finish iin 1 sec; if forget to send complete signal is done .. is clock time
			env.SetTestTimeout(time.Second)
			// Add item1, add item2, add item3, mod item3, add item 4, remove item4
			env.RegisterDelayedCallback(func() {
				env.SignalWorkflow("order-action", OrderSignal{
					Action: "upsert",
					Item: Item{
						ShortName: "FishBurger",
						ShortCode: "FB-01",
						Quantity:  1,
						Amount:    10,
						Notes:     "is goode",
					},
				})
			}, time.Millisecond)
			env.RegisterDelayedCallback(func() {
				env.SignalWorkflow("order-action", OrderSignal{
					Action: "upsert",
					Item: Item{
						ShortName: "LAMBIE BURGER",
						ShortCode: "LAMB-01",
						Quantity:  1,
						Amount:    10,
						Notes:     "",
					},
				})
			}, time.Millisecond*2)
			env.RegisterDelayedCallback(func() {
				env.SignalWorkflow("order-action", OrderSignal{
					Action: "delete",
					Item: Item{
						ShortCode: "FB-01",
					},
				})
			}, time.Millisecond*3)
			env.RegisterDelayedCallback(func() {
				env.SignalWorkflow("order-action", OrderSignal{
					Action: "upsert",
					Item: Item{
						ShortName: "LAMBIE BURGER",
						ShortCode: "LAMB-01",
						Quantity:  2,
						Amount:    20,
						Notes:     "",
					},
				})
			}, time.Millisecond*4)
			env.RegisterDelayedCallback(func() {
				env.SignalWorkflow("order-action", OrderSignal{
					Action: "upsert",
					Item: Item{
						ShortName: "Veggy Burger",
						ShortCode: "VEG-01",
						Quantity:  1,
						Amount:    5,
						Notes:     "",
					},
				})
			}, time.Millisecond*4000)
			// If don;t have this; is doom!
			env.RegisterDelayedCallback(func() {
				env.SignalWorkflow("order-action", OrderSignal{
					Action: "complete",
				})
			}, time.Millisecond*5000)

			// all done setup; run things
			env.ExecuteWorkflow(OrderWorkflow, tt.args.cart)
			// Confirm + pay
			if !env.IsWorkflowCompleted() {
				t.Fatal("WF NOT COmpleted!! Timed out??")
			}
			if wferr := env.GetWorkflowError(); wferr != nil {

				var appErr *temporal.ApplicationError
				if errors.As(wferr, &appErr) {
					fmt.Println("ERR_TYPE:", appErr.Type())
				}
				t.Fatal(errors.Unwrap(wferr))
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

func TestOrderCompleteModifyWorkflow(t *testing.T) {

	// Create a basic order; see basket
	// Add offers
	// Go back and modify
	// Checkout again
}

func TestCompleteOrderWorkflow(t *testing.T) {
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
			got, err := CompleteOrderWorkflow(tt.args.ctx, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompleteOrderWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompleteOrderWorkflow() got = %v, want %v", got, tt.want)
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
