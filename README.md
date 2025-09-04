# Zeus API

Go backend using Fiber, GORM (Postgres), Redis, JWT auth, OTP via Redis, and rate limiting.

## Features
- Phone-based users (phone is username)
- OTP generation and verification stored in Redis
- JWT issuance on OTP verification
- Global rate limiting (per-IP)
- Users list with pagination (protected)

## Getting Started

### Prerequisites
- Go 1.21+
- Docker & Docker Compose

### Environment
Copy and edit env:

```
cp .env.example .env
```

### Run Dependencies

```
docker compose up -d
```

### Run Server

```
go run ./cmd/zeus
```

Server listens on `:${APP_PORT}` (default 8080).

### API
- POST `/api/auth/register` {"phone": "<phone>"}
- POST `/api/auth/login` {"phone": "<phone>"} → issues OTP
- POST `/api/auth/otp/request` {"phone": "<phone>"}
- POST `/api/auth/otp/verify` {"phone": "<phone>", "code": "<code>"} → { token }
- GET `/api/users?page=1&page_size=20` with `Authorization: Bearer <token>`

### Notes
- OTPs are returned in responses for development only. Replace with SMS integration for production.
- Postgres and Redis defaults are set via `.env.example`.