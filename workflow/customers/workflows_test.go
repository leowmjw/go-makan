package customers

import (
	"go.temporal.io/sdk/workflow"
	"reflect"
	"testing"
)

func TestOrderWorkflow(t *testing.T) {
	type args struct {
		ctx        workflow.Context
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OrderWorkflow(tt.args.ctx, tt.args.items, tt.args.delDetails)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderWorkflow() got = %v, want %v", got, tt.want)
			}
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
