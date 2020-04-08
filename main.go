package main

import (
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"net/http"
)

const (
	stripePub  = "pk_test_dk5I07Jj4rIitwg1nD6KuvPc"
	stripePriv = "sk_test_ZRAkevsDrPLxTg2Jmq2dnndt"
)

func init() {
	stripe.Key = stripePriv
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/payments", paymentsHandler).Methods("POST")
	r.HandleFunc("/{id}/payments", accountPaymentsHandler).Methods("GET")
	http.ListenAndServe(":8080", r)
}
