package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

// paymentsHandler handles the POST /payments route
func paymentsHandler(w http.ResponseWriter, r *http.Request) {
	var payment paymentJson
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		returnError(err, w)
		return
	}
	defer r.Body.Close()

	if err := payment.Validate(); err != nil {
		returnError(err, w)
		return
	}
	ch, err := chargeCustomer(payment)
	if err != nil {
		returnError(err, w)
		return
	}
	returnJSON(returnPaymentJSON(ch), w)
}

// accountPaymentsHandler handles the GET /{customer_id}/payments route
func accountPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	customerID := vars["id"]
	if customerID == "" {
		returnError(errors.New("missing customer ID"), w)
		return
	}

	payments, err := customerPayments(customerID)
	if err != nil {
		returnError(err, w)
		return
	}
	returnJSON(payments, w)
}

// returnError returns a JSON response with an error (400 status code)
func returnError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errorJson{err})
}

// returnJSON returns a successful JSON response
func returnJSON(data interface{}, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
