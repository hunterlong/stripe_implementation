package main

import (
	"fmt"
	"github.com/stripe/stripe-go"
	"strconv"
)

func returnPaymentJSON(ch *stripe.Charge) paymentCharge {
	return paymentCharge{
		ID:     ch.ID,
		Amount: intToFloat(ch.Amount),
		Status: ch.Status,
	}
}

func intToFloat(amount int64) float64 {
	value := float64(amount) * 0.01
	rr := fmt.Sprintf("%0.2f", value)
	dollar, _ := strconv.ParseFloat(rr, 64)
	return dollar
}

func floatToDollar(amount float64) int64 {
	amount = amount * 100
	val, _ := strconv.ParseInt(fmt.Sprintf("%0.0f", amount), 10, 64)
	return val
}
