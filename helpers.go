package main

import (
	"fmt"
	"github.com/stripe/stripe-go"
	"strconv"
)

// returnPaymentJSON returns the payment JSON based off stripe's payment response
func returnPaymentJSON(ch *stripe.Charge) paymentCharge {
	return paymentCharge{
		ID:     ch.ID,
		Amount: intToFloat(ch.Amount),
		Status: ch.Status,
	}
}

// intToFloat converts an int to a float
func intToFloat(amount int64) float64 {
	value := float64(amount) * 0.01
	rr := fmt.Sprintf("%0.2f", value)
	dollar, _ := strconv.ParseFloat(rr, 64)
	return dollar
}

// floatToDollar converts a dollar amount to an int for stripe's dollar value
func floatToDollar(amount float64) int64 {
	amount = amount * 100
	val, _ := strconv.ParseInt(fmt.Sprintf("%0.0f", amount), 10, 64)
	return val
}
