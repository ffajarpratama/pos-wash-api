# pos-wash-api

Backend API for **POS Wash (Awash)** — a point-of-sale system for laundry outlets: outlet management, services & pricing, customers, orders, payments, and dashboard.

[![Run In Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/30101452-8765e3f3-afcb-45c1-a16d-95271f9a9836?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D30101452-8765e3f3-afcb-45c1-a16d-95271f9a9836%26entityType%3Dcollection%26workspaceId%3D9a1fe7b6-d2c2-4a93-8249-b43452913c2d)

## Tech stack

- **Go 1.21** — no framework, plain `net/http`
- [`chi`](https://github.com/go-chi/chi) — routing
- [`gorm`](https://gorm.io) + Postgres driver — persistence
- [`golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt) — auth
- [`cloudinary-go`](https://github.com/cloudinary/cloudinary-go) — media storage
- [`go-playground/validator`](https://github.com/go-playground/validator) — request validation
- [`viper`](https://github.com/spf13/viper) — config from `.env`
- [`air`](https://github.com/air-verse/air) — hot reload in development

## Features

- **Auth** — register, login, JWT-protected profile
- **Outlets** — create, fetch
- **Service categories** — listing
- **Services** — CRUD for wash services and pricing
- **Customers** — CRUD
- **Perfumes** — catalog listing (add-on for orders)
- **Payment methods** — listing
- **Orders** — create, list, detail, status updates, payment
- **Dashboard** — summary stats, order trend
- **Media** — file upload (Cloudinary)

## Project structure

```plaintext
cmd/            entrypoint + app wiring (config, DB, router, graceful shutdown)
config/         viper config loader (.env)
internal/
  handler/      HTTP handlers, one file per domain
  router/       chi router assembly, middleware, route mounting
  middleware/   JWT auth, logging, panic recovery
  http/
    request/    request DTOs
    response/   response envelope + writers
  model/        gorm models
  repository/   persistence layer (Postgres via gorm)
  usecase/      business logic
pkg/            integrations & helpers (postgres, cloudinary, jwt, hash,
                custom_error, custom_validator, constant, types, util)
```

## Getting started

**Prerequisites:** Go 1.21+, a Postgres database, a Cloudinary account.

1. Copy the env template and fill in your values:

   ```bash
   cp .env.example .env
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Run the server:

   ```bash
   go run ./cmd/main.go
   ```

   Or with hot reload ([air](https://github.com/air-verse/air)):

   ```bash
   air
   ```

Server listens on `APP_PORT` (see `.env`).

## API

All routes are mounted under `/api/v1/pos`. Public: `/auth/register`, `/auth/login`. Everything else requires a `Bearer` JWT from login.

Request/response shapes are documented in the Postman collection linked above.

## Roadmap

This project is under active development. Planned work:

- [ ] Machine management: register wash/dry machines, link to services
- [ ] Machine connection management: pair/connect physical machines to the system
- [ ] Order add-ons: attach perfume to an order, remove it from an order
- [ ] IoT integration: trigger machine operation via MQTT

## License

[MIT](LICENSE)
