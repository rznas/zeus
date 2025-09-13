# Zeus API

Go backend using Fiber, GORM (Postgres), Redis, JWT auth, OTP via Redis, and rate limiting.

## Features
- Phone-based users (phone is username)
- OTP generation and verification stored in Redis
- JWT issuance on OTP verification
- Global rate limiting (per-IP)
- **Per-phone OTP rate limiting** - prevents OTP abuse with configurable limits
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

### Configuration

Key environment variables:

```env
# App Configuration
APP_PORT=8080
APP_ENV=development
JWT_SECRET=supersecretjwt
JWT_EXPIRES_MINUTES=60

# Rate Limiting
RATE_LIMIT_PER_MINUTE=60          # Global rate limit per IP
OTP_RATE_LIMIT_PER_MINUTE=3       # OTP requests per phone per minute
OTP_RATE_LIMIT_TIMEOUT_SECONDS=60 # Rate limit window duration
OTP_TTL_SECONDS=300               # OTP expiration time

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=zeus
POSTGRES_USER=zeus
POSTGRES_PASSWORD=zeus

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
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

## API Usage Examples

### 1) Request OTP (Login)
Requests an OTP for the given phone. In development, the OTP is printed to stdout (not returned in the response).

**Rate Limiting**: Each phone number is limited to 3 OTP requests per minute (configurable).

```
curl -X POST \
  http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"phone": "+15551234567"}'
```

**Success Response:**
```
{"sent": true}
```

**Rate Limit Exceeded Response (HTTP 429):**
```
{"error": "rate limit exceeded, please try again later"}
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

## Rate Limiting

The API implements two levels of rate limiting:

### 1. Global Rate Limiting
- **Scope**: Per IP address
- **Limit**: 60 requests per minute (configurable via `RATE_LIMIT_PER_MINUTE`)
- **Applied to**: All API endpoints

### 2. OTP Rate Limiting
- **Scope**: Per phone number
- **Limit**: 3 OTP requests per minute (configurable via `OTP_RATE_LIMIT_PER_MINUTE`)
- **Window**: 60 seconds (configurable via `OTP_RATE_LIMIT_TIMEOUT_SECONDS`)
- **Applied to**: `/api/auth/login` endpoint only
- **Storage**: Redis with automatic expiration
- **Error Response**: HTTP 429 with message "rate limit exceeded, please try again later"

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
- Rate limiting uses Redis for distributed rate limiting across multiple server instances.
- OTP rate limiting prevents abuse while allowing legitimate users to receive codes within reasonable limits.