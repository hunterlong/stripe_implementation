# Stripe Implementation
This Stripe implementation will process customer payments and will return all payments made a customer. Server is built in golang and runs on port `8080`. Each response will return a JSON object, status code 200 if it's successful, and 400 if theres an error. If there's an error, the error will be contained within the `error` JSON obect.

### Running
- `go run ./`
> Server runs on port 8080

### Testing
- `go test -v`
> Already includes stripe test keys (would be environment variable in production)

# HTTP Requests

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

# Scalability
- For hundreds of thousands of requests, I'd recommend creating a **AWS Lamba** service to process payments without being responsible for running a server.
- Scale with Docker by increasing the amount of golang server instances running, (round robin with load balancer). Really only effective if there are multiple servers and not running on a single node.

# Limitations
Based on the specs, the current limitation is the processing power of the server and the bandwidth. To truly make this server scaleable, it would have to be ran with multiple instances and/or servers. If this server is getting hundreds of thousands of requests, then I would run it in a simple AWS Lamba service, or create an ECS to auto scale based on need. (requests, CPU usage, RAM, etc)
- Limited on the performance of the Stripe API. Hopefully there's no rate limiting. If the Stripe API is offline for some reason, then this API would only return errors. I do not recommend keeping any payment information (except confirmations ID's) on the server. Leaving all customer information on Stripe is P.C.I compliant. 
- I would not cache any Stripe request. 
- I would add `limit` as a URL query to limit the amount of customer payment records from Stripe. (only show 20 at a time)
- I would also like to add a custom `Error` type that contains better errors messages and status codes based on Stripe response.
