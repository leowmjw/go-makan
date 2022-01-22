package customers

import (
	"fmt"
	"go.temporal.io/sdk/workflow"
)

type Customer struct {
	ID              string
	LinkedPhoneNo   string
	FullName        string
	DisplayName     string
	SignedAgreement bool
	SignUpDate      string
	TerminationDate string
}

// Will move some below to the common type to be shared ..

type Item struct {
	ShortName string
	ShortCode string
	Quantity  int
	Amount    int
	Notes     string // can use to indicate discounts?
}

type PaymentDetails struct {
	ID        string
	ShortName string
	//COD       bool   // Only CC, others Future
	Itemized []string // How Total was caluclated for audit.history ..
	Total    int
}

type OffersClaim struct {
	ID        string
	ShortName string
}

type DeliveryDetails struct {
	CustomerID string
	//DeliveryDate  string // Future
	//DeliveryTime  string // Future
	FullAddress      []string
	ContactNumber    string
	DeliveryCost     int
	DeliveryRange    int
	DeliveryEstimate int
}

type Order struct {
	ID        string
	PartnerID string
	Items     []Item
	DeliveryDetails
	PaymentDetails
	OffersClaim
}

type OrderDelivery struct {
	OrderID             string
	Status              int
	PartnerAcceptedDate string
	PartnerPreparedDate string
	AgentAcceptedDate   string
	AgentPickedUpDate   string
	DeliveryCompleted   string
	Rating              int
}

type ShoppingCart struct {
	CustomerID string
	PartnerID  string
	Items      []Item
	DeliveryDetails
}

type OrderSignal struct {
	Action string
	Item
}

// WorkflowID is customerID/partnerID? guarantee one actively running

func CustomerWorkflow(ctx workflow.Context) {
	// Sign up Account / Contract / Sign Agreement
	// Verify Phone; sign agreement ..
	// Reactivate Account ..

	// Terminate Account / Contract

	// Temporary/Permanent Ban
}

// OrderWorkflow will confirm the COD delivery details
func OrderWorkflow(ctx workflow.Context, cart ShoppingCart) (Order, error) {
	// Query will return shoppingCart ..
	// Default date/time is now .. location is auto-filled based on location

	// Split the wfid to get the CustomerID/PartnerID combo
	i := workflow.GetInfo(ctx)
	wfid := i.WorkflowExecution.ID
	fmt.Println("IN OrderWorkflow!!! WFID: ", wfid)
	fmt.Println("CustomerID: ", cart.CustomerID, " PartnerID: ", cart.PartnerID)
	fmt.Println("CustomerOrderWorkflow ID: ", wfid, " ", i.WorkflowExecution.RunID)

	// Create a new ShoppingCart for use in the future/interruption ..
	// Basic validation? if got no items? partnerID must be valid?
	//cart := ShoppingCart{
	//	PartnerID:       partnerID,
	//	Items:           items,
	//	DeliveryDetails: delDetails,
	//}

	// If fail validation; how?
	//temporal.NewApplicationError()
	// DEBUG
	//spew.Dump(cart)
	// Now block while wating for human signal ..

	// Choose Order + Details from Partner Business
	// or is reloaded from previously active/abndoned
	//receivedPlaceOrder := false
	signalChan := workflow.GetSignalChannel(ctx, "order-action")
	for {
		// Block until PlaceOrder which is order-action -> complete
		var orderSignal *OrderSignal
		more := signalChan.Receive(ctx, &orderSignal)
		if !more {
			fmt.Println("Channel closed ..")
			break
		}
		// DEBUG
		//spew.Dump(orderSignal)
		// Test loop out
		if orderSignal.Action == "complete" {
			break
		}
		//receivedPlaceOrder = true
		//if receivedPlaceOrder {
		//	break
		//}
	}
	// action - AddToBasket
	// action - RemoveFromBasket

	// action - CRUD .. for Items + DeliveryDetails; adjust only quantity
	// Can cancel, drop-off here .. after a short while; or come back later?

	// Choice made but not confirmed yet ..
	// Make payment (if is COD) Success to Escrowo

	// SeeBasket signal calculation and unblocks

	// Recalculate everything; show the new vakues ..
	// Payment must be attached before can kick off

	// action - PlaceOrder to unblock final and complete ..
	// Calculate the time order OrderID
	order := Order{}
	// Persist Order; just pass back ID??
	return order, nil
}

// CompleteOrderWorkflow is accepted ShoppingCart?
func CompleteOrderWorkflow(ctx workflow.Context, order Order) (OrderDelivery, error) {

	// Payment + offers atatched; confirm the rest ..
	// can still go back to order

	// TODO: Future; scheduled delivery; Partner needs to confirm
	// Child workflow of OrderOffer to Partner
	// Order is accepted by Partner

	// Order is Canceled (By Support)
	// Payment collected Reversed for Canceled

	// Child workflow of OrderOffer to Agent
	// Delivery Marked from Agent

	// Payment collected Reversed for Canceled

	// If get signal from Partner in Kitchen; refresh?
	// If take too long; will force reload ??

	// Get signal is in the kitchen

	// Get signal picked up

	// Get signal delivered ..
	// finalize the struct OrderDelivery (from custoer perspective)
	// Send Feedback ..
	orderdel := OrderDelivery{}
	return orderdel, nil
}

// PostOrderWorkflow will handle things like feedback; wrong delivery; support etc..
func PostOrderWorkflow(ctx workflow.Context, order OrderDelivery) error {
	// Emit  signal for Projection .. Report/Admin
	// Get support to chat for this Order (workflowID == OrderID) or customerID/OrderID? or OrderDelivery
	// Collect rating for Partner
	// Collect rating for Agent
	return nil
}
