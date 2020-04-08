package main

import "errors"

type paymentCharge struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

type paymentJson struct {
	AccountId string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type errorJson struct {
	Error error `json:"error"`
}

type custPayments struct {
	Payments []paymentCharge `json:"payments"`
}

func (p paymentJson) Validate() error {
	if p.AccountId == "" {
		return errors.New("missing account id")
	}
	if p.Amount == 0 {
		return errors.New("missing amount")
	}
	return nil
}
