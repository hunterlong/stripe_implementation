package main

import "errors"

// paymentCharge is the JSON response for a payment
type paymentCharge struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

// paymentJson is the JSON request for charging a customer
type paymentJson struct {
	AccountId string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

// errorJson returns a error
type errorJson struct {
	Error error `json:"error"`
}

// custPayments returns the customer's payment records
type custPayments struct {
	Payments []paymentCharge `json:"payments"`
}

// Validate will validate fields for the incoming payment details
func (p paymentJson) Validate() error {
	if p.AccountId == "" {
		return errors.New("missing customer id")
	}
	if p.Amount == 0 {
		return errors.New("missing amount")
	}
	return nil
}
