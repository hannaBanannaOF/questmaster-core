# Questmaster Core

Questmaster Core is the authenticated backend API for the Questmaster platform. It manages campaigns, character sheets, and campaign invites, exposes HTTP endpoints under `/core/api/v1`, stores its primary data in PostgreSQL, refreshes permission documents in MongoDB, and publishes its gateway route registration to RabbitMQ during startup.

## What this service does

- Creates, lists, inspects, and deletes campaigns
- Updates campaign status with the transitions `DRAFT -> ACTIVE -> PAUSED/ARCHIVED`
- Creates, lists, inspects, updates HP for, and deletes character sheets
- Resolves campaign and character slugs to numeric IDs
- Creates or reuses invite links for campaigns
- Accepts campaign invites by linking a character sheet to a campaign
- Validates Bearer tokens against the auth server JWK set

## Stack

- Go 1.25.6
- Gin
- PostgreSQL
- MongoDB
- RabbitMQ
- JWT/JWK-based authentication

## HTTP API

Base path: `/core/api/v1`

Authentication: every route requires `Authorization: Bearer <token>`

### Campaign routes

- `GET /campaign`
- `POST /campaign`
- `GET /campaign/resolve/:slug`
- `DELETE /campaign/:campaignID`
- `GET /campaign/:campaignID/details`
- `PATCH /campaign/:campaignID/status`
- `POST /campaign/:campaignID/invite`

### Character routes

- `GET /character`
- `POST /character`
- `GET /character/resolve/:slug`
- `GET /character/:characterID/details`
- `PATCH /character/:characterID/hp/current`
- `DELETE /character/:characterID`

### Invite routes

- `GET /invite/:inviteHash`
- `POST /invite/:inviteHash/accept`

Notes:

- `POST /campaign` accepts `name`, `system`, and optional `overview`
- `POST /character` accepts `name`, `system`, and optional `hp`
- `PATCH /campaign/:campaignID/status` accepts `status` with values `DRAFT`, `ACTIVE`, `PAUSED`, or `ARCHIVED`; valid transitions are `DRAFT -> ACTIVE`, `ACTIVE -> PAUSED/ARCHIVED`, and `PAUSED -> ACTIVE/ARCHIVED`
- `POST /campaign/:campaignID/invite` returns an invite hash and reuses the existing invite when one already exists
- `POST /invite/:inviteHash/accept` accepts `character_sheet_id`

## Runtime dependencies

The service expects these systems to be available:

- PostgreSQL for campaigns, character sheets, invites, and sessions
- MongoDB for the `permissions` collection
- An auth server exposing JWKs at `${AUTH_HOST}/realms/${AUTH_REALM}/protocol/openid-connect/certs`
- RabbitMQ for the startup message that registers `/core/api/v1/**` with the gateway

## Configuration

The application reads environment variables at startup:

| Variable | Purpose |
| --- | --- |
| `AUTH_HOST` | Base URL of the auth server |
| `AUTH_REALM` | Realm used to build the JWK URL |
| `DB_URL` | PostgreSQL connection string |
| `MONGODB_URL` | MongoDB connection string |
| `MONGODB_DATABASENAME` | MongoDB database that stores the `permissions` collection |
| `MQTT_HOST` | RabbitMQ host |
| `MQTT_PORT` | RabbitMQ port |
| `MQTT_USER` | RabbitMQ username |
| `MQTT_PASSWORD` | RabbitMQ password |
| `GATEWAY_EXCHANGE` | Exchange name declared during the startup gateway-registration step |
| `GATEWAY_URL` | This service public URL to be included in the gateway registration message |
| `RUN_ADDR` | HTTP listen address. If omitted, the service falls back to `0.0.0.0:8080` |

## Database setup

The repository includes raw SQL migrations in `migrations/`. Apply them in numeric order before starting the service.

Important details:

- `0001_init.up.sql` enables the `unaccent` extension
- `0004_campaign_invite.up.sql` uses `gen_random_uuid()`, so PostgreSQL must have `pgcrypto` available before that migration is applied
- The repo does not include a migration runner; use the SQL files with your preferred migration tool or psql workflow
- A `session` table is created by migration, but session routes are not currently exposed by the HTTP API

## Running locally

### Option 1: debug mode with `.env`

Debug mode loads `../../.env`, so run it from `cmd/app`:

```bash
cd cmd/app
go run . -debug
```

### Option 2: release mode with shell environment variables

From the repository root:

```bash
go run ./cmd/app
```

In release mode the service does not auto-load `.env`, so the environment must already be set in your shell or process manager.

## Docker

Build the image:

```bash
docker build -t questmaster-core .
```

Run the container:

```bash
docker run --rm -p 8080:8080 questmaster-core
```
