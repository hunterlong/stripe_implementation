package main

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

// chargeCustomer will charge a customer's billing profile with an amount
func chargeCustomer(payment paymentJson) (*stripe.Charge, error) {
	dollar := floatToDollar(payment.Amount)
	chargeParams := &stripe.ChargeParams{
		Amount:   &dollar,
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Customer: &payment.AccountId,
	}
	ch, err := charge.New(chargeParams)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// customerPayments returns a list of payments made by a customer
func customerPayments(customerID string) (custPayments, error) {
	params := &stripe.ChargeListParams{
		Customer: &customerID,
	}
	i := charge.List(params)
	if i.Err() != nil {
		return custPayments{}, i.Err()
	}

	var payments []paymentCharge
	for i.Next() {
		c := i.Charge()
		p := paymentCharge{
			ID:     c.ID,
			Amount: intToFloat(c.Amount),
			Status: c.Status,
		}
		payments = append(payments, p)
	}
	return custPayments{payments}, nil
}
