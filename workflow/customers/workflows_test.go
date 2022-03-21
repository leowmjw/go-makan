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
	// Setup temporal testsuite
	ts := testsuite.WorkflowTestSuite{}

	// When receive the first add item via action signal
	//	start a new OrderWorkflow  combo <customerID>/<partnerID>
	//		with ShoppingCart
	// If already exists, continue on; flag but not fatal?
	//	Optional: with the pre-req Location info
	//	Query Shopping Cart should succeed
	// Check the workflow is still running ..

	env := ts.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(OrderWorkflow)

	// To simulate signal starting WF; still need to start WF manually tho .. :P
	env.RegisterDelayedCallback(func() {
		fmt.Println("Send signal to get started hre ,...")
		env.SignalWorkflow("order-action", OrderSignal{
			Action: "upsert",
			Item: Item{
				ShortName: "FishBurger",
				ShortCode: "FB-01",
				Quantity:  3,
				Amount:    10,
				Notes:     "is goode",
			},
		})

		// If put here ; it will fail .. NOPE ..
		//if env.IsWorkflowCompleted() {
		//	t.Fatal("Should be started? NOT completed")
		//}

	}, 0)

	env.RegisterDelayedCallback(func() {
		fmt.Println("Assert workflow has started and is running ..")

		if env.IsWorkflowCompleted() {
			t.Fatal("Should be started? NOT completed")
		}
	}, time.Second*2)

	env.RegisterDelayedCallback(func() {
		// Add, see if the order changes ..
		env.SignalWorkflow("order-action", OrderSignal{
			Action: "upsert",
			Item: Item{
				ShortName: "ChicSpicy",
				ShortCode: "CB-03",
				Quantity:  2,
				Amount:    20,
				Notes:     "luv chix",
			},
		})
	}, time.Second*3)
	env.RegisterDelayedCallback(func() {
		// Cancel does not seem towork??
		//env.CancelWorkflow()
		env.SignalWorkflow("order-action", OrderSignal{
			Action: "complete",
		})

		//err := env.GetWorkflowError()
		//fmt.Println("FINAL ERR: ", err)

		//for {
		//	if env.IsWorkflowCompleted() {
		//		break
		//	}
		//}

	}, time.Second*5)

	//env.RegisterDelayedCallback(func() {
	//	fmt.Println("************** WRAp up .. *")
	//	time.Sleep(time.Millisecond * 1000)
	//	// Won't reach here as the WF will already be completed :(
	//	// Should fail first?
	//	if !env.IsWorkflowCompleted() {
	//		t.Fatal("Should have finished!")
	//	}
	//
	//	// After completed; then get result .
	//	var order Order
	//	wferr := env.GetWorkflowResult(&order)
	//	if wferr != nil {
	//		t.Fatal(wferr)
	//	}
	//	spew.Dump(order)
	//
	//}, time.Millisecond*3001)

	// in test env; the signaltostart does not seem to work ..
	env.ExecuteWorkflow(OrderWorkflow, ShoppingCart{
		Items: []Item{{
			ShortName: "ChicSpicy",
			ShortCode: "CB-03",
			Quantity:  1,
			Amount:    10,
			Notes:     "from prev session",
		}},
	})
	// Above will block until it is completed
	if !env.IsWorkflowCompleted() {
		t.Fatal("Should have finished!")
	}
	err := env.GetWorkflowError()
	if err != nil {
		t.Fatal("WFERR: ", err)
	}

	// After completed; then get result .
	var order Order
	wferr := env.GetWorkflowResult(&order)
	if wferr != nil {
		t.Fatal(wferr)
	}
	got := order.Items
	// DEBUG
	//spew.Dump(got)
	want := []Item{
		{
			ShortName: "ChicSpicy",
			ShortCode: "CB-03",
			Quantity:  2,
			Amount:    20,
			Notes:     "luv chix",
		},
		{
			ShortName: "FishBurger",
			ShortCode: "FB-01",
			Quantity:  3,
			Amount:    10,
			Notes:     "is goode",
		},
	}
	// After some changes; final result should match
	// including ordering ...
	if !cmp.Equal(want, got) {
		t.Fatal("DIFF: ", cmp.Diff(want, got))
	}
	// DEBUG
	//t.Fatal("TODO: Implement ...")
	//if !env.IsWorkflowCompleted() {
	//	t.Fatal("Shoudl block ..")
	//}
}

func TestRemovalUpsertDeleteOrderWorkflow(t *testing.T) {
	// Setup temporal testsuite
	ts := testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(OrderWorkflow)

	// Scenario: Have 3 items; remove from mid using both delete + upsert, complete; should be expected
	// To simulate signal starting WF; still need to start WF manually tho .. :P
	env.RegisterDelayedCallback(func() {
		fmt.Println("Send signal to get started hre ,...")
		env.SignalWorkflow("order-action", OrderSignal{
			Action: "upsert",
			Item: Item{
				ShortName: "FishBurgerSambal",
				ShortCode: "FB-02",
				Quantity:  3,
				Amount:    30,
				Notes:     "is goode",
			},
		})
	}, 0)
	env.RegisterDelayedCallback(func() {
		// Test delete via upsert of <1 items ..
		env.SignalWorkflow("order-action", OrderSignal{
			Action: "delete",
			Item: Item{
				ShortCode: "FB-01",
			},
		})
	}, time.Second*2)

	env.RegisterDelayedCallback(func() {
		// Test delete via upsert of <1 items ..
		env.SignalWorkflow("order-action", OrderSignal{
			Action: "upsert",
			Item: Item{
				ShortCode: "CB-03",
				Quantity:  0,
			},
		})
	}, time.Second*3)
	// Finish the test ..
	env.RegisterDelayedCallback(func() {
		// Cancel does not seem towork??
		//env.CancelWorkflow()
		env.SignalWorkflow("order-action", OrderSignal{
			Action: "complete",
		})
	}, time.Second*5)

	// in test env; the signaltostart does not seem to work ..
	env.ExecuteWorkflow(OrderWorkflow, ShoppingCart{
		Items: []Item{
			{
				ShortName: "ChicNormal",
				ShortCode: "CB-01",
				Quantity:  1,
				Amount:    7,
				Notes:     "i like chix",
			},
			{
				ShortName: "ChicSpicy",
				ShortCode: "CB-03",
				Quantity:  1,
				Amount:    10,
				Notes:     "from prev session",
			},
			{
				ShortName: "FishBurger",
				ShortCode: "FB-01",
				Quantity:  3,
				Amount:    15,
				Notes:     "3x Fishies",
			},
		},
	})

	// Above will block until it is completed
	if !env.IsWorkflowCompleted() {
		t.Fatal("Should have finished!")
	}
	err := env.GetWorkflowError()
	if err != nil {
		t.Fatal("WFERR: ", err)
	}

	// After completed; then get result .
	var order Order
	wferr := env.GetWorkflowResult(&order)
	if wferr != nil {
		t.Fatal(wferr)
	}
	got := order.Items
	// DEBUG
	//spew.Dump(got)
	want := []Item{
		{
			ShortName: "ChicNormal",
			ShortCode: "CB-01",
			Quantity:  1,
			Amount:    7,
			Notes:     "i like chix",
		},
		{
			ShortName: "FishBurgerSambal",
			ShortCode: "FB-02",
			Quantity:  3,
			Amount:    30,
			Notes:     "is goode",
		},
	}
	// After some changes; final result should match
	// including ordering ...
	if !cmp.Equal(want, got) {
		t.Fatal("DIFF: ", cmp.Diff(want, got))
	}
}

func TestRemovedLastItemRemoveCartOrderWorkflow(t *testing.T) {
	// Setup temporal testsuite
	ts := testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(OrderWorkflow)

	// Scenario #2: Have 1 item; remove; should be completed, nothing returned ..
	// To simulate signal starting WF; still need to start WF manually tho .. :P
	// To simulate signal starting WF; still need to start WF manually tho .. :P
	env.RegisterDelayedCallback(func() {
		fmt.Println("Send signal to get started hre ,...")
		env.SignalWorkflow("order-action", OrderSignal{
			Action: "upsert",
			Item: Item{
				ShortName: "FishBurgerSambal",
				ShortCode: "FB-02",
				Quantity:  3,
				Amount:    30,
				Notes:     "is goode",
			},
		})
	}, 0)

	env.RegisterDelayedCallback(func() {
		// Test delete via upsert of <1 items .. it is NON-exist
		env.SignalWorkflow("order-action", OrderSignal{
			Action: "upsert",
			Item: Item{
				ShortCode: "CB-03",
				Quantity:  0,
			},
		})
	}, time.Second*2)

	env.RegisterDelayedCallback(func() {
		// Test delete via upsert of <1 items ..
		env.SignalWorkflow("order-action", OrderSignal{
			Action: "delete",
			Item: Item{
				ShortCode: "FB-02",
			},
		})
	}, time.Second*3)

	// Finish the test ..
	env.RegisterDelayedCallback(func() {
		// Cancel does not seem towork??
		//env.CancelWorkflow()
		//env.SignalWorkflow("order-action", OrderSignal{
		//	Action: "complete",
		//})

		// No need to manually complete; it should show as done!
		// It is checked further down; so might not be needed ..
	}, time.Second*5)

	// Assert the Workflow is completed .. after signal to remove ..

	// in test env; the signaltostart does not seem to work .. start with empty cart ..
	env.ExecuteWorkflow(OrderWorkflow, ShoppingCart{})

	// Above will block until it is completed
	if !env.IsWorkflowCompleted() {
		t.Fatal("Should have finished!")
	}
	err := env.GetWorkflowError()
	if err != nil {
		t.Fatal("WFERR: ", err)
	}

	// After completed; then get result .
	var order Order
	wferr := env.GetWorkflowResult(&order)
	if wferr != nil {
		t.Fatal(wferr)
	}
	got := order.Items
	// DEBUG
	//spew.Dump(got)
	var want []Item
	// After some changes; final result should match
	// including ordering ...
	if !cmp.Equal(want, got) {
		t.Fatal("DIFF: ", cmp.Diff(want, got))
	}

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
				ID: "mleow/baba-ang", // this does not work .. now it does :P
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
