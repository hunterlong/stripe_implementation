package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

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

func returnError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errorJson{err})
}

func returnJSON(data interface{}, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
