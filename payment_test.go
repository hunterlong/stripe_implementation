package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	customerID string
)

func init() {
	stripe.Key = stripePriv
}

func TestNewCustomer(t *testing.T) {
	params := &stripe.CustomerParams{
		Email: stripe.String("info@dev.com"),
	}
	params.SetSource("tok_mastercard")

	cust, err := customer.New(params)
	require.Nil(t, err)

	customerID = cust.ID
	assert.NotEmpty(t, customerID)
	assert.NotEmpty(t, cust.Email)
	t.Log("New Custom ID: ", customerID)
}

func TestPayment(t *testing.T) {

	tester := []testData{
		{
			Url:    "/payment",
			Method: "POST",
			Data: paymentJson{
				AccountId: customerID,
				Amount:    45.33,
			},
			Handler:      paymentsHandler,
			ExpectedCode: 200,
			ExpectedResponseFunc: checkPayment,
		},
		{
			Url:          "/payment",
			Method:       "POST",
			Data:         nil,
			Handler:      paymentsHandler,
			ExpectedCode: 400,
		},
		{
			Url:    "/payment",
			Method: "POST",
			Data: paymentJson{
				AccountId: "",
				Amount:    45.33,
			},
			Handler:      paymentsHandler,
			ExpectedCode: 400,
		},
		{
			Url:    "/payment",
			Method: "POST",
			Data: paymentJson{
				AccountId: "1234",
				Amount:    0,
			},
			Handler:      paymentsHandler,
			ExpectedCode: 400,
		},
	}

	for _, test := range tester {
		httpChainTest(t, test)
	}
}

func TestListPayments(t *testing.T) {
	tester := []testData{
		{
			Url:          "/" + customerID + "/payments",
			RouteMux:     "/{id}/payments",
			Method:       "GET",
			Handler:      accountPaymentsHandler,
			ExpectedCode: 200,
			ExpectedResponseFunc: checkCustomerPayments,
		},
		{
			Url:          "/missingID/payments",
			RouteMux:     "/{id}/payments",
			Method:       "GET",
			Handler:      accountPaymentsHandler,
			ExpectedCode: 400,
		},
	}

	for _, test := range tester {
		httpChainTest(t, test)
	}
}

func (p paymentJson) asBuffer() *bytes.Buffer {
	d, _ := json.Marshal(p)
	return bytes.NewBuffer(d)
}

type testData struct {
	Url                  string
	RouteMux             string
	Method               string
	Data                 interface{}
	Handler              http.HandlerFunc
	ExpectedCode         int
	ExpectedResponseFunc ResponseFunc
}

type ResponseFunc func(response []byte) error

func httpChainTest(t *testing.T, test testData) {
	var d []byte
	if test.Data != nil {
		d, _ = json.Marshal(test.Data)
	}
	data := bytes.NewBuffer(d)

	req, err := http.NewRequest(test.Method, test.Url, data)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	if test.RouteMux != "" {
		router.HandleFunc(test.RouteMux, test.Handler)
	} else {
		router.HandleFunc(test.Url, test.Handler)
	}
	router.ServeHTTP(rr, req)

	if test.ExpectedCode != 0 {
		assert.Equal(t, test.ExpectedCode, rr.Code)
	}

	if test.ExpectedResponseFunc != nil {
		body, err := ioutil.ReadAll(rr.Body)
		require.Nil(t, err)
		assert.Nil(t, test.ExpectedResponseFunc(body))
	}
}

func parseCustomerPayments(data []byte) (custPayments, error) {
	var ch custPayments
	if err := json.Unmarshal(data, &ch); err != nil {
		return custPayments{}, err
	}
	return ch, nil
}

func parsePaymentCharge(data []byte) (paymentCharge, error) {
	var ch paymentCharge
	if err := json.Unmarshal(data, &ch); err != nil {
		return paymentCharge{}, err
	}
	return ch, nil
}

func checkPayment(response []byte) error {
	charge, err := parsePaymentCharge(response)
	if err != nil {
		return err
	}
	if charge.Status != "succeeded" {
		return errors.New("status is not approved: " + charge.Status)
	}
	return nil
}

func checkCustomerPayments(response []byte) error {
	custPayments, err := parseCustomerPayments(response)
	if err != nil {
		return err
	}
	if len(custPayments.Payments) == 0 {
		return errors.New("customer has no payments")
	}
	payment := custPayments.Payments[0]
	if payment.Status != "succeeded" {
		return errors.New("charge does not have status succeeded: " + payment.Status)
	}
	return nil
}
