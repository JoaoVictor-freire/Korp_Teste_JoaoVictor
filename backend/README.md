# Backend Structure

Initial backend scaffold for the technical challenge, organized around two microservices:

- `stock-service`: product registration and stock control.
- `billing-service`: invoice registration and lifecycle.

## Users + Auth (JWT)

This skeleton includes `users` as an entity (register/login) and uses JWT for authentication.
All `/api/v1/*` endpoints (products/invoices) require `Authorization: Bearer <token>`.

## Run

- Stock: `go run ./cmd/stock-service` (default `:8081`)
- Billing: `go run ./cmd/billing-service` (default `:8082`)

## Run Both Services

- `chmod +x ./scripts/dev.sh`
- `./scripts/dev.sh`

## Quick Test (curl)

Health (no auth):

- `curl -s http://localhost:8081/health`
- `curl -s http://localhost:8082/health`

Register/login (stock-service):

- `curl -s -X POST http://localhost:8081/api/v1/auth/register -H 'Content-Type: application/json' -d '{"email":"user_a@example.com","password":"123456"}'`
- `curl -s -X POST http://localhost:8081/api/v1/auth/login -H 'Content-Type: application/json' -d '{"email":"user_a@example.com","password":"123456"}'`

Create/list products (use the token you got above):

- `curl -s -X POST http://localhost:8081/api/v1/products -H 'Content-Type: application/json' -H 'Authorization: Bearer <token>' -d '{"code":"P001","description":"Caneta","stock":10}'`
- `curl -s http://localhost:8081/api/v1/products -H 'Authorization: Bearer <token>'`

Create/list invoices (billing-service, same token/secret):

- `curl -s -X POST http://localhost:8082/api/v1/invoices -H 'Content-Type: application/json' -H 'Authorization: Bearer <token>' -d '{"number":1,"items":[{"product_code":"P001","quantity":2}]}'`
- `curl -s http://localhost:8082/api/v1/invoices -H 'Authorization: Bearer <token>'`

Suggested next steps:

1. Replace the in-memory repositories with a real database adapter.
2. Create an HTTP or event-based integration between billing and stock.
3. Implement the invoice print flow with status transition and stock deduction.
4. Add idempotency, concurrency control, observability, and tests.
