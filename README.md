# Zeus API

Go backend using Fiber, GORM (Postgres), Redis, JWT auth, OTP via Redis, and rate limiting.

## Features
- Phone-based users (phone is username)
- OTP generation and verification stored in Redis
- JWT issuance on OTP verification
- Global rate limiting (per-IP)
- Users list with pagination (protected)
- Swagger UI docs at `/swagger/`

## Getting Started

### Prerequisites
- Go 1.21+
- Docker & Docker Compose

### Environment
Copy and edit env:

```
cp sample.env .env
```

By default, the app loads `.env`. If absent, it falls back to `sample.env`.

### Run Dependencies

```
docker compose up -d
```

### Run Server

```
go run ./cmd/zeus
```

Server listens on `:${APP_PORT}` (default 8080).

## API Usage Examples

### 1) Request OTP (Login)
Requests an OTP for the given phone. In development, the OTP is printed to stdout (not returned in the response).

```
curl -X POST \
  http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"phone": "+15551234567"}'
```
Response:
```
{"sent": true}
```
Check your server logs for a line like: `DEV OTP for +15551234567: 123456`

### 2) Verify OTP (Creates user if needed + returns JWT)
```
curl -X POST \
  http://localhost:8080/api/auth/otp/verify \
  -H 'Content-Type: application/json' \
  -d '{"phone": "+15551234567", "code": "123456"}'
```
Example response:
```
{"token": "<JWT_TOKEN>"}
```

### 3) List Users (Protected, with pagination)
```
TOKEN="<JWT_TOKEN>"
curl -X GET \
  'http://localhost:8080/api/users?page=1&page_size=20' \
  -H "Authorization: Bearer ${TOKEN}"
```

## Swagger
- Open Swagger UI: `http://localhost:8080/swagger/`
- Click Authorize and paste either `Bearer <JWT>` or just `<JWT>`. The server accepts both formats.

If you modify routes/handlers or tags, regenerate docs:
```
go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g ./cmd/zeus/main.go --output ./docs
```

## Notes
- OTPs are not returned in responses in development; they are printed to stdout.
- Replace the OTP printing with an SMS provider for production.
- Postgres and Redis defaults are set via `.env`/`sample.env`.
- CORS is enabled for Swagger and typical API clients; tighten it for production as needed.