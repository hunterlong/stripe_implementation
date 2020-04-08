# Stripe Implementation

### Running
- `go run ./`
> Server runs on port 8080

### Testing
- `go test -v`

### POST `/payments`
```json
{
  "account_id": "cus_12345",
  "amount": "42.99"
}
```
Returns payment confirmation or an error.

### GET `/{customer_id}/payments`
Returns an array of all payments made by a customer.
