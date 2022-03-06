package customers

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go.temporal.io/sdk/workflow"
)

type Customer struct {
	ID              string `json:"id,omitempty"`
	LinkedPhoneNo   string `json:"linked_phone_no,omitempty"`
	FullName        string `json:"full_name,omitempty"`
	DisplayName     string `json:"display_name,omitempty"`
	SignedAgreement bool   `json:"signed_agreement,omitempty"`
	SignUpDate      string `json:"sign_up_date,omitempty"`
	TerminationDate string `json:"termination_date,omitempty"`
}

// Will move some below to the common type to be shared ..

type Item struct {
	ShortName string `json:"name,omitempty"`
	ShortCode string `json:"code,omitempty"`
	Quantity  int    `json:"quantity,omitempty"`
	Amount    int    `json:"amount,omitempty"`
	Notes     string `json:"notes,omitempty"` // can use to indicate discounts?
}

type PaymentDetails struct {
	ID        string `json:"id,omitempty"`
	ShortName string `json:"short_name,omitempty"`
	//COD       bool   // Only CC, others Future
	Itemized []string `json:"itemized,omitempty"` // How Total was caluclated for audit.history ..
	Total    int      `json:"total,omitempty"`
}

type OffersClaim struct {
	ID        string `json:"id,omitempty"`
	ShortName string `json:"short_name,omitempty"`
}

type DeliveryDetails struct {
	CustomerID string `json:"customer_id,omitempty"`
	//DeliveryDate  string // Future
	//DeliveryTime  string // Future
	FullAddress      []string `json:"full_address,omitempty"`
	ContactNumber    string   `json:"contact_number,omitempty"`
	DeliveryCost     int      `json:"delivery_cost,omitempty"`
	DeliveryRange    int      `json:"delivery_range,omitempty"`
	DeliveryEstimate int      `json:"delivery_estimate,omitempty"`
}

type Order struct {
	ID              string `json:"id,omitempty"`
	PartnerID       string `json:"partner_id,omitempty"`
	Items           []Item `json:"items,omitempty"`
	DeliveryDetails `json:"delivery_details,omitempty"`
	PaymentDetails  `json:"payment_details,omitempty"`
	OffersClaim     `json:"offers_claim,omitempty"`
}

type OrderDelivery struct {
	OrderID             string `json:"order_id,omitempty"`
	Status              int    `json:"status,omitempty"`
	PartnerAcceptedDate string `json:"partner_accepted_date,omitempty"`
	PartnerPreparedDate string `json:"partner_prepared_date,omitempty"`
	AgentAcceptedDate   string `json:"agent_accepted_date,omitempty"`
	AgentPickedUpDate   string `json:"agent_picked_up_date,omitempty"`
	DeliveryCompleted   string `json:"delivery_completed,omitempty"`
	Rating              int    `json:"rating,omitempty"`
}

type Cart struct {
	CustomerID string `json:"customer_id"`
	PartnerID  string `json:"partner_id"`
	Items      []Item `json:"items,omitempty"`
}

type ShoppingCart struct {
	CustomerID      string `json:"customer_id"`
	PartnerID       string `json:"partner_id"`
	Items           []Item `json:"items,omitempty"`
	DeliveryDetails `json:"delivery_details,omitempty"`
}

type OrderSignal struct {
	Action string `json:"action"`
	Item   `json:"item,omitempty"`
}

// WorkflowID is customerID/partnerID? guarantee one actively running

func CustomerWorkflow(ctx workflow.Context) {
	// Sign up Account / Contract / Sign Agreement
	// Verify Phone; sign agreement ..
	// Reactivate Account ..

	// Terminate Account / Contract

	// Temporary/Permanent Ban
}

func NewOrderWorkflow(isTest bool) func(workflow.Context, ShoppingCart) (Order, error) {
	if isTest {
		return func(ctx workflow.Context, cart ShoppingCart) (Order, error) {
			// Keep track of current Order
			order := Order{}
			// Waiting for this signal from customer
			signalChan := workflow.GetSignalChannel(ctx, "order-action")
			// Loop until get a Order Submitted or nothing in the ShoppingCart
			for {
			orderSignalLoop:
				// Block until PlaceOrder which is order-action -> complete
				var orderSignal *OrderSignal
				more := signalChan.Receive(ctx, &orderSignal)
				if !more {
					fmt.Println("Channel closed ..")
					break
				}
				// TODO: In Future; when getting close to 1K signals receive;
				// 	can kickoff new ChildWorkflow ..
				// DEBUG
				//spew.Dump(orderSignal)
				// Test loop out
				if orderSignal.Action == "complete" {
					fmt.Println("Order done and submitted!!!")
					break
				}

				if orderSignal.Action == "delete" {
					// Look b ySjortCode; and reform the slcie .,.
				}

				if orderSignal.Action == "upsert" {
					// Validate if got Item; if no; do nothing ...
					if orderSignal.Item.Quantity < 1 {
						// Look and remove item
						// if not found ignore??

					} else {
						// Add item to the slice ...
						order.Items = append(order.Items, orderSignal.Item)
					}

					// If remove all and no more members in Order; can break and exist ..
					if len(order.Items) == 0 {
						// Nothing left in order; quit the flow!!
						fmt.Println("NOTHING in Order .. Quitting!!!!")
						break
					}

					goto orderSignalLoop
				}

				// If reach here means no signal recognize; flag it!!!
				fmt.Println("UNKNOWN Signal ACTION: ", orderSignal.Action)
			}

			return order, nil
		}
	}

	return implOrderWorkflow
}

func implOrderWorkflow(ctx workflow.Context, cart ShoppingCart) (Order, error) {
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

// OrderWorkflow will confirm the COD delivery details
func OrderWorkflow(ctx workflow.Context, cart ShoppingCart) (Order, error) {
	// Query will return shoppingCart ..
	// Default date/time is now .. location is auto-filled based on location

	// Split the wfid to get the CustomerID/PartnerID combo
	i := workflow.GetInfo(ctx)
	wfid := i.WorkflowExecution.ID
	fmt.Println("IN OrderWorkflow!!! WFID: ", wfid)
	fmt.Println("CART: ", cart)
	//fmt.Println("CustomerID: ", cart.CustomerID, " PartnerID: ", cart.PartnerID)
	fmt.Println("CustomerOrderWorkflow ID: ", wfid, " ", i.WorkflowExecution.RunID)

	fmt.Println("CART at the start ...")
	spew.Dump(cart)
	// How to figure out it is test ..
	owf := NewOrderWorkflow(true)
	return owf(ctx, cart)
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
