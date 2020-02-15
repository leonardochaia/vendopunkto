
# Create invoice using public api
curl -X POST http://127.0.0.1:8080/v1/invoices -d '{"amount":100000000000000, "currency":"XMR"}' -H 'Content-Type: application/json'

# Confirm invoice as if you were a plugin with internal api
curl -X POST http://127.0.0.1:9080/v1/invoices/payments/confirm -d '{"txHash":"thehash","confirmations":0, "address":"", "amount":100000000000000}' -H 'Content-Type: application/json'


# Check public info to see if it's in mempool (status=2, confirmations=0)
curl http://127.0.0.1:8080/v1/invoices/<the invoice id>

# Add confirmations to a tx
curl -X POST http://127.0.0.1:9080/v1/invoices/payments/confirm -d '{"txHash":"thehash","confirmations":5, "address":"", "amount":100000000000000}' -H 'Content-Type: application/json'
