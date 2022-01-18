package customers

import (
	"fmt"
	"go.temporal.io/sdk/workflow"
	"strings"
)

// Will move it to the common type to be shared ..

type Item struct {
	ShortName string
	ShortCode string
	Quantity  int
	Amount    int
	Notes     string
}

type PaymentMethod struct {
	ID        string
	ShortName string
	COD       bool // Signal slightly different flow; since need get cash
}

type OffersClaim struct {
	ID        string
	ShortName string
}
type DeliveryDetails struct {
	CustomerID    string
	FullAddress   []string
	ContactNumber string
}

type Order struct {
	ID        string
	PartnerID string
	Items     []Item
	PaymentMethod
	OffersClaim
}

type ShoppingCart struct {
	PartnerID string
	Items     []Item
}

// WorkflowID is customerID/partnerID? guarantee one actively running

func CustomerWorkflow(ctx workflow.Context) {
	// Sign up

	// Sign off
}

func PreOrderWorkflow(ctx workflow.Context) (ShoppingCart, error) {
	// Query will return shoppingCart ..

	// Split the wfid to get the CustomerID/PartnerID combo
	i := workflow.GetInfo(ctx)
	wfid := i.WorkflowExecution.ID
	ids := strings.Split(wfid, "/")
	fmt.Println("CustomerID: ", ids[0], " PartnerID: ", ids[1])
	fmt.Println("CustomerOrderWorkflow ID: ", wfid, " ", i.WorkflowExecution.RunID)
	// Choose Order + Details from Partner Business
	// or is reloaded from previously active/abndoned

	// Can cancel, drop-off here .. after a short while; or come back later?

	// Choice made but not confirmed yet ..
	// Make payment (if is COD) Success to Escrowo
	// Payment must be attached before can kick off

	cart := ShoppingCart{}
	return cart, nil
}

// CustomerOrderWorkflow is accepted ShoppingCart?
func CustomerOrderWorkflow(ctx workflow.Context) (Order, error) {

	// Child workflow of OrderOffer to Partner
	// Order is accepted by Partner

	// Order is Canceled (By Support)
	// Payment collected Reversed for Canceled

	// Child workflow of OrderOffer to Agent
	// Delivery Marked from Agent

	// Payment collected Reversed for Canceled

	// Send Feedback ..
	order := Order{}
	return order, nil
}

// Confirmed Order; track acceptance to delivery completion?
